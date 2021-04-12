package keeper

import (
	"encoding/binary"
	"time"

	"github.com/tendermint/tendermint/crypto"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-project/core/x/wasm/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CompileCode uncompress the wasm code bytes and store the code to local file system
func (k Keeper) CompileCode(ctx sdk.Context, wasmCode []byte) (codeHash []byte, err error) {
	maxContractSize := k.MaxContractSize(ctx)
	if uint64(len(wasmCode)) > maxContractSize {

		return nil, sdkerrors.Wrap(types.ErrStoreCodeFailed, "contract size is too huge")
	}

	wasmCode, err = k.uncompress(wasmCode, maxContractSize)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrStoreCodeFailed, err.Error())
	}

	// consume gas for compile cost
	ctx.GasMeter().ConsumeGas(types.CompileCostPerByte*uint64(len(wasmCode)), "Compiling WASM Bytes Cost")

	codeHash, err = k.wasmer.Create(wasmCode)
	if err != nil {
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
	migratable bool) (sdk.AccAddress, []byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "instantiate")
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: init")

	if uint64(len(initMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "init msg size is too huge")
	}

	instanceID, err := k.GetLastInstanceID(ctx)
	if err != nil {
		return nil, nil, err
	}

	instanceID++

	// create contract address
	contractAddress := k.generateContractAddress(ctx, codeID, instanceID)
	existingAcct := k.accountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// create contract account
	contractAccount := k.accountKeeper.NewAccountWithAddress(ctx, contractAddress)
	k.accountKeeper.SetAccount(ctx, contractAccount)

	// deposit initial contract funds
	if !deposit.IsZero() {
		if err := k.bankKeeper.SendCoins(ctx, creator, contractAddress, deposit); err != nil {
			return nil, nil, err
		}
	}

	// get code info
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeInfoKey(codeID))
	if bz == nil {
		return nil, nil, sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", codeID)
	}

	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)

	// prepare env and info for contract instantiate call
	env := types.NewEnv(ctx, contractAddress)
	info := types.NewInfo(creator, deposit)

	// create prefixed data store
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), contractStoreKey)

	// instantiate wasm contract
	res, gasUsed, err := k.wasmer.Instantiate(
		codeInfo.CodeHash,
		env,
		info,
		initMsg,
		contractStore,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	// consume gas before raise error
	k.consumeGas(ctx, gasUsed, "Contract init")
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}

	// Must store contract info first, so last part can use it
	contractInfo := types.NewContractInfo(codeID, contractAddress, creator, initMsg, migratable)

	k.SetLastInstanceID(ctx, instanceID)
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	// vaildate events is size and parse to sdk events
	events, err := types.ValidateAndParseEvents(contractAddress, k.EventParams(ctx), res.Attributes...)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// validate data size
	if uint64(len(res.Data)) > k.MaxContractDataSize(ctx) {
		return nil, nil, sdkerrors.Wrap(types.ErrExceedMaxContractDataSize, "returned data size is too huge")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	if err := k.dispatchAll(ctx, contractAddress, res.Submessages, res.Messages); err != nil {
		return nil, nil, sdkerrors.Wrap(err, "dispatch")
	}

	return contractAddress, res.Data, nil
}

// ExecuteContract executes the contract instance
func (k Keeper) ExecuteContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	caller sdk.AccAddress,
	exeMsg []byte,
	coins sdk.Coins) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "execute")
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: execute")

	if uint64(len(exeMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "execute msg size is too huge")
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

	env := types.NewEnv(ctx, contractAddress)
	info := types.NewInfo(caller, coins)
	res, gasUsed, err := k.wasmer.Execute(
		codeInfo.CodeHash,
		env,
		info,
		exeMsg,
		storePrefix,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Execution")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// vaildate events is size and parse to sdk events
	events, err := types.ValidateAndParseEvents(contractAddress, k.EventParams(ctx), res.Attributes...)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// validate data size
	if uint64(len(res.Data)) > k.MaxContractDataSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractDataSize, "returned data size is too huge")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	if err := k.dispatchAll(ctx, contractAddress, res.Submessages, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return res.Data, nil
}

// MigrateContract allows to upgrade a contract to a new code with data migration.
func (k Keeper) MigrateContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	caller sdk.AccAddress,
	newCodeID uint64,
	migrateMsg []byte) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "migrate")
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: migrate")

	if uint64(len(migrateMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "migrate msg size is too huge")
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	if !contractInfo.Migratable {
		return nil, types.ErrNotMigratable
	}

	if contractInfo.Owner != caller.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no permission")
	}

	newCodeInfo, err := k.GetCodeInfo(ctx, newCodeID)
	if err != nil {
		return nil, err
	}

	env := types.NewEnv(ctx, contractAddress)

	// prepare necessary meta data
	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	res, gasUsed, err := k.wasmer.Migrate(
		newCodeInfo.CodeHash,
		env,
		migrateMsg,
		prefixStore,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Migration")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// vaildate events is size and parse to sdk events
	events, err := types.ValidateAndParseEvents(contractAddress, k.EventParams(ctx), res.Attributes...)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// validate data size
	if uint64(len(res.Data)) > k.MaxContractDataSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractDataSize, "returned data size is too huge")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	contractInfo.CodeID = newCodeID
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	// dispatch submessages and messages
	if err := k.dispatchAll(ctx, contractAddress, res.Submessages, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return res.Data, nil
}

// reply is only called from keeper internal functions
// (dispatchSubmessages) after processing the submessages
func (k Keeper) reply(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	reply wasmvmtypes.Reply) error {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "reply")
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: reply")

	eventParams := k.EventParams(ctx)
	codeInfo, storePrefix, err := k.getContractDetails(ctx, contractAddress)

	env := types.NewEnv(ctx, contractAddress)

	// to prevent passing too huge events to wasmvm
	// cap the reply.Events length to eventParams.MaxAttributeNum
	if reply.Result.Ok != nil && uint64(len(reply.Result.Ok.Events)) > eventParams.MaxAttributeNum {
		reply.Result.Ok.Events = reply.Result.Ok.Events[:eventParams.MaxAttributeNum]
	}

	res, gasUsed, err := k.wasmer.Reply(
		codeInfo.CodeHash,
		env,
		reply,
		storePrefix,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getGasMeter(ctx),
		k.getGasRemaining(ctx),
	)

	k.consumeGas(ctx, gasUsed, "Contract Reply")
	if err != nil {
		return sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// vaildate events is size and parse to sdk events
	events, err := types.ValidateAndParseEvents(contractAddress, eventParams, res.Attributes...)
	if err != nil {
		return sdkerrors.Wrap(err, "event validation failed")
	}

	// validate data size
	if uint64(len(res.Data)) > k.MaxContractDataSize(ctx) {
		return sdkerrors.Wrap(types.ErrExceedMaxContractDataSize, "returned data size is too huge")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	if err := k.dispatchAll(ctx, contractAddress, res.Submessages, res.Messages); err != nil {
		return sdkerrors.Wrap(err, "dispatch")
	}

	return nil
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

func (k Keeper) queryToStore(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "query-raw")
	if key == nil {
		return nil
	}

	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	return prefixStore.Get(key)
}

func (k Keeper) queryToContract(ctx sdk.Context, contractAddress sdk.AccAddress, queryMsg []byte) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "query-smart")
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: query")

	codeInfo, contractStorePrefix, err := k.getContractDetails(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	env := types.NewEnv(ctx, contractAddress)
	queryResult, gasUsed, err := k.wasmer.Query(
		codeInfo.CodeHash,
		env,
		queryMsg,
		contractStorePrefix,
		k.getCosmWasmAPI(ctx),
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
	k.cdc.MustUnmarshalBinaryBare(bz, &contractInfo)

	bz = store.Get(types.GetCodeInfoKey(contractInfo.CodeID))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", contractInfo.CodeID)
		return
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStorePrefix = prefix.NewStore(ctx.KVStore(k.storeKey), contractStoreKey)
	return
}
