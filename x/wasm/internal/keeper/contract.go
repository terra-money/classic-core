package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// StoreCode uploads and compiles a WASM contract bytecode, returning a short identifier for the stored code
func (k Keeper) StoreCode(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte) (codeID uint64, err error) {
	wasmCode, err = k.uncompress(ctx, wasmCode)
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrStoreCodeFailed, err.Error())
	}

	codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrStoreCodeFailed, err.Error())
	}

	codeID, err = k.GetLastCodeID(ctx)
	if err != nil {
		return 0, err
	}

	codeID++
	contractInfo := types.NewCodeInfo(codeHash, creator)

	k.SetLastCodeID(ctx, codeID)
	k.SetCodeInfo(ctx, codeID, contractInfo)

	return codeID, nil
}

// InstantiateContract creates an instance of a WASM contract
func (k Keeper) InstantiateContract(ctx sdk.Context, codeID uint64, creator sdk.AccAddress, initMsg []byte, deposit sdk.Coins) (contractAddress sdk.AccAddress, err error) {
	instanceID, err := k.GetLastInstanceID(ctx)
	if err != nil {
		return nil, err
	}

	instanceID++

	// create contract address
	contractAddress = k.generateContractAddress(ctx, codeID, instanceID)
	existingAcct := k.accountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// create contract account
	contractAccount := k.accountKeeper.NewAccountWithAddress(ctx, contractAddress)
	k.accountKeeper.SetAccount(ctx, contractAccount)

	// deposit initial contract funds
	if !deposit.IsZero() {
		err = k.bankKeeper.SendCoins(ctx, creator, contractAddress, deposit)
		if err != nil {
			return
		}
	}

	// get code info
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeInfoKey(codeID))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", codeID)
		return
	}

	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &codeInfo)

	// prepare params for contract instantiate call
	apiParams := types.NewWasmAPIParams(ctx, creator, deposit, contractAddress)

	// create prefixed data store
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), contractStoreKey)

	// instantiate wasm contract
	gas := k.gasForContract(ctx)
	res, err2 := k.wasmer.Instantiate(codeInfo.CodeHash.Bytes(), apiParams, initMsg, contractStore, cosmwasmAPI, k.querier.WithCtx(ctx), gas)
	if err2 != nil {
		err = sdkerrors.Wrap(types.ErrInstantiateFailed, err2.Error())
		return
	}

	k.consumeGas(ctx, res.GasUsed)

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return
	}

	// persist contractInfo
	contractInfo := types.NewContractInfo(codeID, contractAddress, creator, initMsg)

	k.SetLastInstanceID(ctx, instanceID)
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	return contractAddress, nil
}

// ExecuteContract executes the contract instance
func (k Keeper) ExecuteContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) (sdk.Result, error) {
	codeInfo, storePrefix, sdkerr := k.getContractDetails(ctx, contractAddress)
	if sdkerr != nil {
		return sdk.Result{}, sdkerr
	}

	// add more funds
	if !coins.IsZero() {
		sdkerr = k.bankKeeper.SendCoins(ctx, caller, contractAddress, coins)
		if sdkerr != nil {
			return sdk.Result{}, sdkerr
		}
	}

	apiParams := types.NewWasmAPIParams(ctx, caller, coins, contractAddress)

	gas := k.gasForContract(ctx)
	res, err := k.wasmer.Execute(codeInfo.CodeHash.Bytes(), apiParams, msg, storePrefix, cosmwasmAPI, k.querier.WithCtx(ctx), gas)
	if err != nil {
		// TODO: wasmer doesn't return wasm gas used on error. we should consume it (for error on metering failure)
		// Note: OutOfGas panics (from storage) are caught by go-cosmwasm, subtract one more gas to check if
		// this contract died due to gas limit in Storage
		k.consumeGas(ctx, k.GasMultiplier(ctx))
		return sdk.Result{}, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	k.consumeGas(ctx, res.GasUsed)

	sdkerr = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if sdkerr != nil {
		return sdk.Result{}, sdkerr
	}

	return types.ParseResult(res, contractAddress), nil
}

func (k Keeper) gasForContract(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()
	remaining := (meter.Limit() - meter.GasConsumed()) * k.GasMultiplier(ctx)
	if remaining > k.MaxContractGas(ctx) {
		return k.MaxContractGas(ctx)
	}
	return remaining
}

// converts contract gas usage to sdk gas and consumes it
func (k Keeper) consumeGas(ctx sdk.Context, gas uint64) {
	consumed := gas / k.GasMultiplier(ctx)
	ctx.GasMeter().ConsumeGas(consumed, "wasm contract")
}

// generates a contract address from codeID + instanceID
// and increases last instanceID
func (k Keeper) generateContractAddress(ctx sdk.Context, codeID uint64, instanceID uint64) sdk.AccAddress {
	// NOTE: It is possible to get a duplicate address if either codeID or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	contractID := codeID<<32 + instanceID
	return addrFromUint64(contractID)
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}

func (k Keeper) queryToStore(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) (result []byte) {
	if key == nil {
		return result
	}

	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	result = prefixStore.Get(key)

	return
}

func (k Keeper) queryToContract(ctx sdk.Context, contractAddr sdk.AccAddress, queryMsg []byte) ([]byte, error) {
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(k.queryGasLimit))

	codeInfo, contractStorePrefix, err := k.getContractDetails(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	queryResult, gasUsed, err := k.wasmer.Query(codeInfo.CodeHash.Bytes(), queryMsg, contractStorePrefix, cosmwasmAPI, k.querier.WithCtx(ctx), k.gasForContract(ctx))
	if err != nil {
		return nil, err
	}

	k.consumeGas(ctx, gasUsed)
	return queryResult, nil
}

func (k Keeper) getContractDetails(ctx sdk.Context, contractAddress sdk.AccAddress) (codeInfo types.CodeInfo, contractStorePrefix prefix.Store, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetContractInfoKey(contractAddress))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "contract %s", contractAddress)
		return
	}

	var contractInfo types.ContractInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &contractInfo)

	bz = store.Get(types.GetCodeInfoKey(contractInfo.CodeID))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", contractInfo.CodeID)
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &codeInfo)
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStorePrefix = prefix.NewStore(ctx.KVStore(k.storeKey), contractStoreKey)
	return
}
