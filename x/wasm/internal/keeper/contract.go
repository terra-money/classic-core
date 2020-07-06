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
	if uint64(len(wasmCode)) > k.MaxContractSize(ctx) {
		return 0, sdkerrors.Wrap(types.ErrStoreCodeFailed, "contract size is too huge")
	}

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
func (k Keeper) InstantiateContract(ctx sdk.Context, codeID uint64, creator sdk.AccAddress, initMsg []byte, deposit sdk.Coins, migratable bool) (contractAddress sdk.AccAddress, err error) {
	if uint64(len(initMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrInstantiateFailed, "init msg size is too huge")
	}

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
	res, gasUsed, err := k.wasmer.Instantiate(codeInfo.CodeHash.Bytes(), apiParams, initMsg, contractStore, cosmwasmAPI, k.querier.WithCtx(ctx), ctx.GasMeter(), gas)

	// consume gas before raise error
	k.consumeGas(ctx, gasUsed)
	if err != nil {
		err = sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
		return
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return
	}

	// persist contractInfo
	contractInfo := types.NewContractInfo(codeID, contractAddress, creator, initMsg, migratable)

	k.SetLastInstanceID(ctx, instanceID)
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	return contractAddress, nil
}

// ExecuteContract executes the contract instance
func (k Keeper) ExecuteContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, exeMsg []byte, coins sdk.Coins) ([]byte, error) {
	if uint64(len(exeMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrInstantiateFailed, "execute msg size is too huge")
	}

	codeInfo, storePrefix, err := k.getContractDetails(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	// add more funds
	if !coins.IsZero() {
		err = k.bankKeeper.SendCoins(ctx, caller, contractAddress, coins)
		if err != nil {
			return nil, err
		}
	}

	apiParams := types.NewWasmAPIParams(ctx, caller, coins, contractAddress)

	gas := k.gasForContract(ctx)
	res, gasUsed, err := k.wasmer.Execute(codeInfo.CodeHash.Bytes(), apiParams, exeMsg, storePrefix, cosmwasmAPI, k.querier.WithCtx(ctx), ctx.GasMeter(), gas)

	k.consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// MigrateContract allows to upgrade a contract to a new code with data migration.
func (k Keeper) MigrateContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newCodeID uint64, migrateMsg []byte) ([]byte, error) {
	if uint64(len(migrateMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrInstantiateFailed, "migrate msg size is too huge")
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	if !contractInfo.Migratable {
		return nil, types.ErrNotMigratable
	}

	if !contractInfo.Owner.Equals(caller) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no permission")
	}

	newCodeInfo, err := k.GetCodeInfo(ctx, newCodeID)
	if err != nil {
		return nil, err
	}

	var noDeposit sdk.Coins
	params := types.NewWasmAPIParams(ctx, caller, noDeposit, contractAddress)

	// prepare necessary meta data
	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	gas := k.gasForContract(ctx)

	res, gasUsed, err := k.wasmer.Migrate(newCodeInfo.CodeHash.Bytes(), params, migrateMsg, &prefixStore, cosmwasmAPI, k.querier.WithCtx(ctx), ctx.GasMeter(), gas)

	k.consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	contractInfo.CodeID = newCodeID
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	if err := k.dispatchMessages(ctx, contractAddress, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return res.Data, nil
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

	queryResult, gasUsed, err := k.wasmer.Query(codeInfo.CodeHash.Bytes(), queryMsg, contractStorePrefix, cosmwasmAPI, k.querier.WithCtx(ctx), ctx.GasMeter(), k.gasForContract(ctx))

	k.consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, err
	}

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
