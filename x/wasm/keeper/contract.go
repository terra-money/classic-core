package keeper

import (
	"context"
	"time"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/wasm/types"

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

	ctx.GasMeter().ConsumeGas(types.CompileCosts(len(wasmCode)), "Compiling WASM Bytes Cost")

	codeHash, err = k.wasmVM.Create(wasmCode)
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
	k.Logger(ctx).Debug("storing new contract", "code_id", codeID)

	return codeID, nil
}

// MigrateCode uploads and compiles a WASM contract bytecode for the existing code id.
// After columbus-5 update, all contract code will be removed from the store
// due to in-compatibility between CosmWasm@v0.10.x and CosmWasm@v0.16.x
// The migration can be executed by once after columbus-5 update.
// TODO - remove after columbus-5 update
func (k Keeper) MigrateCode(ctx sdk.Context, codeID uint64, creator sdk.AccAddress, wasmCode []byte) error {
	codeInfo, err := k.GetCodeInfo(ctx, codeID)
	if err != nil {
		return err
	}

	if len(codeInfo.CodeHash) != 0 || codeInfo.Creator != creator.String() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no permission")
	}

	codeHash, err := k.CompileCode(ctx, wasmCode)
	if err != nil {
		return err
	}

	codeInfo.CodeHash = codeHash
	k.SetCodeInfo(ctx, codeID, codeInfo)
	k.Logger(ctx).Debug("storing new contract", "code_id", codeID)

	return nil
}

// InstantiateContract creates an instance of a WASM contract
func (k Keeper) InstantiateContract(
	ctx sdk.Context,
	codeID uint64,
	creator sdk.AccAddress,
	admin sdk.AccAddress,
	initMsg []byte,
	deposit sdk.Coins) (sdk.AccAddress, []byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "instantiate")
	ctx.GasMeter().ConsumeGas(types.RegisterContractCosts(), "Registering contract to the store")
	ctx.GasMeter().ConsumeGas(types.InstantiateContractCosts(len(initMsg)), "Loading CosmWasm module: init")

	if uint64(len(initMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "init msg size is too huge")
	}

	instanceID, err := k.GetLastInstanceID(ctx)
	if err != nil {
		return nil, nil, err
	}

	instanceID++

	// create contract address
	contractAddress := types.GenerateContractAddress(codeID, instanceID)
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
	k.cdc.MustUnmarshal(bz, &codeInfo)

	// prepare env and info for contract instantiate call
	env := types.NewEnv(ctx, contractAddress)
	info := types.NewInfo(creator, deposit)

	// create prefixed data store
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStore := prefix.NewStore(types.KVStore(ctx, k.storeKey), contractStoreKey)

	// instantiate wasm contract
	res, gasUsed, err := k.wasmVM.Instantiate(
		codeInfo.CodeHash,
		env,
		info,
		initMsg,
		contractStore,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getWasmVMGasMeter(ctx),
		k.getWasmVMGasRemaining(ctx),
		types.JSONDeserializationWasmGasCost,
	)

	// add types.GasMultiplier to occur out of gas panic
	k.consumeWasmVMGas(ctx, gasUsed+types.GasMultiplier, "Contract initialize")
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}

	// consume gas for wasm events
	ctx.GasMeter().ConsumeGas(types.EventCosts(res.Attributes, res.Events), "Event Cost")

	// Must store contract info first, so last part can use it
	contractInfo := types.NewContractInfo(codeID, contractAddress, creator, admin, initMsg)

	k.SetLastInstanceID(ctx, instanceID)
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	// parse wasm events to sdk events
	events, err := types.ParseEvents(contractAddress, res.Attributes, res.Events)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	respData := res.Data
	if replyData, err := k.dispatchMessages(ctx, contractAddress, res.Messages...); err != nil {
		return nil, nil, sdkerrors.Wrap(err, "dispatch")
	} else if replyData != nil {
		respData = replyData
	}

	return contractAddress, respData, nil
}

// ExecuteContract executes the contract instance
func (k Keeper) ExecuteContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	sender sdk.AccAddress,
	execMsg []byte,
	coins sdk.Coins) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "execute")
	ctx.GasMeter().ConsumeGas(types.InstantiateContractCosts(len(execMsg)), "Loading CosmWasm module: execute")

	if uint64(len(execMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "execute msg size is too huge")
	}

	codeInfo, storePrefix, err := k.getContractDetails(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	// add more funds
	if !coins.IsZero() {
		err = k.bankKeeper.SendCoins(ctx, sender, contractAddress, coins)
		if err != nil {
			return nil, err
		}
	}

	env := types.NewEnv(ctx, contractAddress)
	info := types.NewInfo(sender, coins)
	res, gasUsed, err := k.wasmVM.Execute(
		codeInfo.CodeHash,
		env,
		info,
		execMsg,
		storePrefix,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getWasmVMGasMeter(ctx),
		k.getWasmVMGasRemaining(ctx),
		types.JSONDeserializationWasmGasCost,
	)

	// add types.GasMultiplier to occur out of gas panic
	k.consumeWasmVMGas(ctx, gasUsed+types.GasMultiplier, "Contract Execution")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	// consume gas for wasm events
	ctx.GasMeter().ConsumeGas(types.EventCosts(res.Attributes, res.Events), "Event Cost")

	// parse wasm events to sdk events
	events, err := types.ParseEvents(contractAddress, res.Attributes, res.Events)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	respData := res.Data
	if replyData, err := k.dispatchMessages(ctx, contractAddress, res.Messages...); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	} else if replyData != nil {
		respData = replyData
	}

	return respData, nil
}

// MigrateContract allows to upgrade a contract to a new code with data migration.
func (k Keeper) MigrateContract(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	sender sdk.AccAddress,
	newCodeID uint64,
	migrateMsg []byte) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "migrate")
	ctx.GasMeter().ConsumeGas(types.InstantiateContractCosts(len(migrateMsg)), "Loading CosmWasm module: migrate")

	if uint64(len(migrateMsg)) > k.MaxContractMsgSize(ctx) {
		return nil, sdkerrors.Wrap(types.ErrExceedMaxContractMsgSize, "migrate msg size is too huge")
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	if contractInfo.Admin == "" {
		return nil, types.ErrNotMigratable
	}

	if contractInfo.Admin != sender.String() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "no permission")
	}

	newCodeInfo, err := k.GetCodeInfo(ctx, newCodeID)
	if err != nil {
		return nil, err
	}

	env := types.NewEnv(ctx, contractAddress)

	// prepare necessary meta data
	prefixStoreKey := types.GetContractStoreKey(contractAddress)
	prefixStore := prefix.NewStore(types.KVStore(ctx, k.storeKey), prefixStoreKey)

	res, gasUsed, err := k.wasmVM.Migrate(
		newCodeInfo.CodeHash,
		env,
		migrateMsg,
		prefixStore,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getWasmVMGasMeter(ctx),
		k.getWasmVMGasRemaining(ctx),
		types.JSONDeserializationWasmGasCost,
	)

	// add types.GasMultiplier to occur out of gas panic
	k.consumeWasmVMGas(ctx, gasUsed+types.GasMultiplier, "Contract Migration")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// consume gas for wasm events
	ctx.GasMeter().ConsumeGas(types.EventCosts(res.Attributes, res.Events), "Event Cost")

	// parse wasm events to sdk events
	events, err := types.ParseEvents(contractAddress, res.Attributes, res.Events)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	contractInfo.CodeID = newCodeID
	k.SetContractInfo(ctx, contractAddress, contractInfo)

	// dispatch submessages and messages
	respData := res.Data
	if replyData, err := k.dispatchMessages(ctx, contractAddress, res.Messages...); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	} else if replyData != nil {
		respData = replyData
	}

	return respData, nil
}

// reply is only called from keeper internal functions
// (dispatchSubmessages) after processing the submessages
func (k Keeper) reply(
	ctx sdk.Context,
	contractAddress sdk.AccAddress,
	reply wasmvmtypes.Reply) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "wasm", "contract", "reply")
	ctx.GasMeter().ConsumeGas(types.ReplyCosts(reply), "Loading CosmWasm module: reply")

	codeInfo, storePrefix, err := k.getContractDetails(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	env := types.NewEnv(ctx, contractAddress)
	res, gasUsed, err := k.wasmVM.Reply(
		codeInfo.CodeHash,
		env,
		reply,
		storePrefix,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getWasmVMGasMeter(ctx),
		k.getWasmVMGasRemaining(ctx),
		types.JSONDeserializationWasmGasCost,
	)

	// add types.GasMultiplier to occur out of gas panic
	k.consumeWasmVMGas(ctx, gasUsed+types.GasMultiplier, "Contract Reply")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrReplyFailed, err.Error())
	}

	// consume gas for wasm events
	ctx.GasMeter().ConsumeGas(types.EventCosts(res.Attributes, res.Events), "Event Cost")

	// parse wasm events to sdk events
	events, err := types.ParseEvents(contractAddress, res.Attributes, res.Events)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "event validation failed")
	}

	// emit events
	ctx.EventManager().EmitEvents(events)

	// dispatch submessages and messages
	respData := res.Data
	if replyData, err := k.dispatchMessages(ctx, contractAddress, res.Messages...); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	} else if replyData != nil {
		respData = replyData
	}

	return respData, nil
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
	ctx.GasMeter().ConsumeGas(types.InstantiateContractCosts(len(queryMsg)), "Loading CosmWasm module: query")

	codeInfo, contractStorePrefix, err := k.getContractDetails(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	env := types.NewEnv(ctx, contractAddress)

	// assert and increase query depth
	ctx, err = assertAndIncreaseQueryDepth(ctx)
	if err != nil {
		return nil, err
	}

	queryResult, gasUsed, err := k.wasmVM.Query(
		codeInfo.CodeHash,
		env,
		queryMsg,
		contractStorePrefix,
		k.getCosmWasmAPI(ctx),
		k.querier.WithCtx(ctx),
		k.getWasmVMGasMeter(ctx),
		k.getWasmVMGasRemaining(ctx),
		types.JSONDeserializationWasmGasCost,
	)

	// add types.GasMultiplier to occur out of gas panic
	k.consumeWasmVMGas(ctx, gasUsed+types.GasMultiplier, "Contract Query")
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrContractQueryFailed, err.Error())
	}

	return queryResult, err
}

func assertAndIncreaseQueryDepth(ctx sdk.Context) (sdk.Context, error) {
	var queryDepth uint8
	if depth := ctx.Context().Value(types.WasmVMQueryDepthContextKey); depth != nil {
		queryDepth = depth.(uint8)
	} else {
		queryDepth = 1
	}

	if queryDepth > types.ContractMaxQueryDepth {
		return ctx, types.ErrExceedMaxQueryDepth
	}

	// set next query depth
	ctx = ctx.WithContext(context.WithValue(ctx.Context(), types.WasmVMQueryDepthContextKey, queryDepth+1))

	return ctx, nil
}

func (k Keeper) getContractDetails(ctx sdk.Context, contractAddress sdk.AccAddress) (codeInfo types.CodeInfo, contractStorePrefix prefix.Store, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetContractInfoKey(contractAddress))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "contract %s", contractAddress)
		return
	}

	var contractInfo types.ContractInfo
	k.cdc.MustUnmarshal(bz, &contractInfo)

	bz = store.Get(types.GetCodeInfoKey(contractInfo.CodeID))
	if bz == nil {
		err = sdkerrors.Wrapf(types.ErrNotFound, "codeID %d", contractInfo.CodeID)
		return
	}

	k.cdc.MustUnmarshal(bz, &codeInfo)
	contractStoreKey := types.GetContractStoreKey(contractAddress)
	contractStorePrefix = prefix.NewStore(types.KVStore(ctx, k.storeKey), contractStoreKey)
	return
}
