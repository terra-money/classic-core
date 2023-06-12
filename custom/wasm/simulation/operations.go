package simulation

import (
	"math/rand"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/CosmWasm/wasmd/app/params"
	"github.com/CosmWasm/wasmd/x/wasm/keeper/testdata"
	wasmsim "github.com/CosmWasm/wasmd/x/wasm/simulation"
	"github.com/CosmWasm/wasmd/x/wasm/types"
)

// Simulation operation weights constants
//
//nolint:gosec
const (
	OpWeightMsgStoreCode           = "op_weight_msg_store_code"
	OpWeightMsgInstantiateContract = "op_weight_msg_instantiate_contract"
	OpWeightMsgExecuteContract     = "op_weight_msg_execute_contract"
	OpWeightMsgUpdateAdmin         = "op_weight_msg_update_admin"
	OpWeightMsgClearAdmin          = "op_weight_msg_clear_admin"
	OpWeightMsgMigrateContract     = "op_weight_msg_migrate_contract"
	OpReflectContractPath          = "op_reflect_contract_path"
)

// WasmKeeper is a subset of the wasm keeper used by simulations
type WasmKeeper interface {
	GetParams(ctx sdk.Context) types.Params
	IterateCodeInfos(ctx sdk.Context, cb func(uint64, types.CodeInfo) bool)
	IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, types.ContractInfo) bool)
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	PeekAutoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64
}
type BankKeeper interface {
	simulation.BankKeeper
	IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool
}

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	simstate *module.SimulationState,
	ak types.AccountKeeper,
	bk BankKeeper,
	wasmKeeper WasmKeeper,
) simulation.WeightedOperations {
	var (
		weightMsgStoreCode           int
		weightMsgInstantiateContract int
		weightMsgExecuteContract     int
		weightMsgUpdateAdmin         int
		weightMsgClearAdmin          int
		weightMsgMigrateContract     int
		wasmContractPath             string
	)

	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgStoreCode, &weightMsgStoreCode, nil,
		func(_ *rand.Rand) {
			weightMsgStoreCode = params.DefaultWeightMsgStoreCode
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgInstantiateContract, &weightMsgInstantiateContract, nil,
		func(_ *rand.Rand) {
			weightMsgInstantiateContract = params.DefaultWeightMsgInstantiateContract
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgExecuteContract, &weightMsgInstantiateContract, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteContract = params.DefaultWeightMsgExecuteContract
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgUpdateAdmin, &weightMsgUpdateAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAdmin = params.DefaultWeightMsgUpdateAdmin
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgClearAdmin, &weightMsgClearAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgClearAdmin = params.DefaultWeightMsgClearAdmin
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgMigrateContract, &weightMsgMigrateContract, nil,
		func(_ *rand.Rand) {
			weightMsgMigrateContract = params.DefaultWeightMsgMigrateContract
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpReflectContractPath, &wasmContractPath, nil,
		func(_ *rand.Rand) {
			wasmContractPath = ""
		},
	)

	var wasmBz []byte
	if wasmContractPath == "" {
		wasmBz = testdata.MigrateReflectContractWasm()
	} else {
		var err error
		wasmBz, err = os.ReadFile(wasmContractPath)
		if err != nil {
			panic(err)
		}
	}

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgStoreCode,
			wasmsim.SimulateMsgStoreCode(ak, bk, wasmKeeper, wasmBz, 5_000_000),
		),
		simulation.NewWeightedOperation(
			weightMsgInstantiateContract,
			wasmsim.SimulateMsgInstantiateContract(ak, bk, wasmKeeper, wasmsim.DefaultSimulationCodeIDSelector),
		),
		simulation.NewWeightedOperation(
			weightMsgExecuteContract,
			SimulateMsgExecuteContract(
				ak,
				bk,
				wasmKeeper,
				wasmsim.DefaultSimulationExecuteContractSelector,
				wasmsim.DefaultSimulationExecuteSenderSelector,
				wasmsim.DefaultSimulationExecutePayloader,
			),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateAdmin,
			wasmsim.SimulateMsgUpdateAmin(
				ak,
				bk,
				wasmKeeper,
				wasmsim.DefaultSimulationUpdateAdminContractSelector,
			),
		),
		simulation.NewWeightedOperation(
			weightMsgClearAdmin,
			wasmsim.SimulateMsgClearAdmin(
				ak,
				bk,
				wasmKeeper,
				wasmsim.DefaultSimulationClearAdminContractSelector,
			),
		),
		simulation.NewWeightedOperation(
			weightMsgMigrateContract,
			wasmsim.SimulateMsgMigrateContract(
				ak,
				bk,
				wasmKeeper,
				wasmsim.DefaultSimulationMigrateContractSelector,
				wasmsim.DefaultSimulationMigrateCodeIDSelector,
			),
		),
	}
}

// SimulateMsgExecuteContract create a execute message a reflect contract instance
func SimulateMsgExecuteContract(
	ak types.AccountKeeper,
	bk BankKeeper,
	wasmKeeper WasmKeeper,
	contractSelector wasmsim.MsgExecuteContractSelector,
	senderSelector wasmsim.MsgExecuteSenderSelector,
	payloader wasmsim.MsgExecutePayloader,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr := contractSelector(ctx, wasmKeeper)
		if contractAddr == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgExecuteContract{}.Type(), "no contract instance available"), nil, nil
		}
		simAccount, err := senderSelector(wasmKeeper, ctx, contractAddr, accs)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgExecuteContract{}.Type(), "query contract owner"), nil, err
		}

		deposit := sdk.Coins{}
		spendableCoins := bk.SpendableCoins(ctx, simAccount.Address)
		for _, v := range spendableCoins {
			if bk.IsSendEnabledCoin(ctx, v) {
				deposit = deposit.Add(simtypes.RandSubsetCoins(r, sdk.NewCoins(v))...)
			}
		}
		if deposit.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgExecuteContract{}.Type(), "broke account"), nil, nil
		}
		msg := types.MsgExecuteContract{
			Sender:   simAccount.Address.String(),
			Contract: contractAddr.String(),
			Funds:    deposit,
		}
		if err := payloader(&msg); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgExecuteContract{}.Type(), "contract execute payload"), nil, err
		}

		txCtx := wasmsim.BuildOperationInput(r, app, ctx, &msg, simAccount, ak, bk, deposit)
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
