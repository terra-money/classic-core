package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/internal/types"
)

// CompileCode uncompress the wasm code bytes and store the code to local file system
func (k Keeper) CompileCode(ctx sdk.Context, wasmCode []byte) (codeHash []byte, err error) {
	if uint64(len(wasmCode)) > k.MaxContractSize(ctx) {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(types.ErrInternal, "contract size is too huge")
		}

		return nil, sdkerrors.Wrap(types.ErrStoreCodeFailed, "contract size is too huge")
	}

	wasmCode, err = k.uncompress(ctx, wasmCode)
	if err != nil {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(types.ErrInternal, err.Error())
		}

		return nil, sdkerrors.Wrap(types.ErrStoreCodeFailed, err.Error())
	}

	// consume gas for compile cost
	ctx.GasMeter().ConsumeGas(types.CompileCostPerByte*uint64(len(wasmCode)), "Compiling WASM Bytes Cost")

	codeHash, err = k.wasmer.Create(wasmCode)
	if err != nil {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(types.ErrInternal, err.Error())
		}

		return nil, sdkerrors.Wrap(types.ErrStoreCodeFailed, err.Error())
	}

	return
}

// StoreCode uploads and compiles a WASM contract bytecode, returning a short identifier for the stored code
func (k Keeper) StoreCode(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte) (codeID uint64, err error) {
	codeHash, err := k.CompileCode(ctx, wasmCode)
	if err != nil {
		return 0, err
	}

	codeID, err = k.GetLastCodeID(ctx)
	if err != nil {
		return 0, err
	}

	codeID++
	codeInfo := types.NewCodeInfo(codeID, codeHash, creator)

	k.SetLastCodeID(ctx, codeID)
	k.SetCodeInfo(ctx, codeID, codeInfo)

	return codeID, nil
}

// InstantiateContract creates an instance of a WASM contract
func (k Keeper) InstantiateContract(
	ctx sdk.Context,
	codeID uint64,
	creator sdk.AccAddress,
	initMsg []byte,
	deposit sdk.Coins,
	migratable bool) (contractAddress sdk.AccAddress, err error) {
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: init")

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
	res, gasUsed, err := k.wasmer.Instantiate(
		codeInfo.CodeHash.Bytes(),
		apiParams,
		initMsg,
		contractStore,
		k.getCosmwamAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	// consume gas before raise error
	k.consumeGas(ctx, gasUsed, "Contract init")
	if err != nil {
		err = sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
		return
	}

	// Must store contract info first, so last part can use it
	contractInfo := types.NewContractInfo(codeID, contractAddress, creator, initMsg, migratable)

	k.SetLastInstanceID(ctx, instanceID)
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	// check contract creator address is in whitelist
	if _, ok := k.loggingWhitelist[creator.String()]; ok || k.wasmConfig.LoggingAll() {
		events := types.ParseEvents(res.Log, contractAddress)
		ctx.EventManager().EmitEvents(events)
		if ok && !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
			// If a contract is created from whitelist,
			// add the contract to whitelist.
			// It can be canceled due to transaction failure,
			// but that is tiny cost so ignore that case.
			contractAddr := contractAddress.String()
			k.loggingWhitelist[contractAddr] = true
			k.wasmConfig.ContractLoggingWhitelist += "," + contractAddr

			// store updated config to local wasm.toml
			k.StoreConfig()
		}
	}

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return
	}

	return contractAddress, nil
}

// ExecuteContract executes the contract instance
func (k Keeper) ExecuteContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, exeMsg []byte, coins sdk.Coins) ([]byte, error) {
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: execute")

	if uint64(len(exeMsg)) > k.MaxContractMsgSize(ctx) {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(types.ErrInstantiateFailed, "execute msg size is too huge")
		}

		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, "execute msg size is too huge")
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
	res, gasUsed, err := k.wasmer.Execute(
		codeInfo.CodeHash.Bytes(),
		apiParams,
		exeMsg,
		storePrefix,
		k.getCosmwamAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Execution")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// emit all events from the contract in the logging whitelist
	if _, ok := k.loggingWhitelist[contractAddress.String()]; k.wasmConfig.LoggingAll() || ok {
		events := types.ParseEvents(res.Log, contractAddress)
		ctx.EventManager().EmitEvents(events)
	}

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// MigrateContract allows to upgrade a contract to a new code with data migration.
func (k Keeper) MigrateContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newCodeID uint64, migrateMsg []byte) ([]byte, error) {
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: migrate")

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

	res, gasUsed, err := k.wasmer.Migrate(
		newCodeInfo.CodeHash.Bytes(),
		params,
		migrateMsg,
		&prefixStore,
		k.getCosmwamAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Migration")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// emit all events from the contract in the logging whitelist
	if _, ok := k.loggingWhitelist[contractAddress.String()]; k.wasmConfig.LoggingAll() || ok {
		events := types.ParseEvents(res.Log, contractAddress)
		ctx.EventManager().EmitEvents(events)
	}

	contractInfo.CodeID = newCodeID
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	if err := k.dispatchMessages(ctx, contractAddress, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return res.Data, nil
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
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: query")

	codeInfo, contractStorePrefix, err := k.getContractDetails(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	queryResult, gasUsed, err := k.wasmer.Query(
		codeInfo.CodeHash.Bytes(),
		queryMsg,
		contractStorePrefix,
		k.getCosmwamAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Query")
	if err != nil {
		err = sdkerrors.Wrap(types.ErrContractQueryFailed, err.Error())
	}

	return queryResult, err
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
