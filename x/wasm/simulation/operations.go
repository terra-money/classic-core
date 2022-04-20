package simulation

//DONTCOVER

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"
)

// nolint
const (
	OpWeightMsgStoreCode           = "op_weight_msg_store_code"
	OpWeightMsgInstantiateContract = "op_weight_msg_instantiate_contract"
	OpWeightMsgExecuteContract     = "op_weight_msg_execute_contract"
	OpWeightMsgMigrateContract     = "op_weight_msg_migrate_contract"
	OpWeightMsgUpdateContractAdmin = "op_weight_msg_update_contract_admin"
	OpWeightMsgClearContractAdmin  = "op_weight_msg_update_contract_admin"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	protoCdc *codec.ProtoCodec,
) simulation.WeightedOperations {
	var weightMsgStoreCode int
	var weightMsgInstantiateContract int
	var weightMsgExecuteContract int
	var weightMsgMigrateContract int
	var weightMsgUpdateContractAdmin int
	var weightMsgClearContractAdmin int
	appParams.GetOrGenerate(cdc, OpWeightMsgStoreCode, &weightMsgStoreCode, nil,
		func(_ *rand.Rand) {
			weightMsgStoreCode = 1
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgInstantiateContract, &weightMsgInstantiateContract, nil,
		func(_ *rand.Rand) {
			weightMsgInstantiateContract = 1
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgExecuteContract, &weightMsgExecuteContract, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteContract = 1
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgMigrateContract, &weightMsgMigrateContract, nil,
		func(_ *rand.Rand) {
			weightMsgMigrateContract = 1
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgUpdateContractAdmin, &weightMsgUpdateContractAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateContractAdmin = 3
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgClearContractAdmin, &weightMsgClearContractAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgClearContractAdmin = 1
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgStoreCode,
			SimulateMsgStoreCode(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgInstantiateContract,
			SimulateMsgInstantiateContract(ak, bk, k, protoCdc),
		),
		simulation.NewWeightedOperation(
			weightMsgExecuteContract,
			SimulateMsgExecuteContract(ak, bk, k, protoCdc),
		),
		simulation.NewWeightedOperation(
			weightMsgMigrateContract,
			SimulateMsgMigrateContract(ak, bk, k, protoCdc),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateContractAdmin,
			SimulateMsgUpdateContractAdmin(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgClearContractAdmin,
			SimulateMsgClearContractAdmin(ak, bk, k),
		),
	}
}

// nolint:deadcode,unused
func mustLoad(path string) []byte {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

var testContract []byte

// nolint:funlen
func SimulateMsgStoreCode(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		_, err := k.GetCodeInfo(ctx, 2)
		if err == nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStoreCode, "code already registered"), nil, nil
		}

		if testContract == nil {
			loadContract()
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, simAccount.Address))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgStoreCode, "unable to generate fee"), nil, nil
		}

		msg := types.NewMsgStoreCode(simAccount.Address, testContract)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas*10,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

type initMsg struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
}

// nolint:unused,deadcode
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// nolint:funlen
func SimulateMsgInstantiateContract(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	protoCdc *codec.ProtoCodec) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		bobAcc, _ := simtypes.RandomAcc(r, accs)
		fredAcc, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, fredAcc.Address)
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, fredAcc.Address))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgInstantiateContract, "unable to generate fee"), nil, err
		}

		initMsg := initMsg{
			Verifier:    fredAcc.Address.String(),
			Beneficiary: bobAcc.Address.String(),
		}

		initMsgBz, err := json.Marshal(initMsg)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgInstantiateContract, "failed to marshal json"), nil, err
		}

		_, err = k.GetCodeInfo(ctx, 1)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgInstantiateContract, "code not exists yet"), nil, nil
		}

		msg := types.NewMsgInstantiateContract(fredAcc.Address, fredAcc.Address, 1, initMsgBz, nil)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			fredAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", protoCdc), nil, nil
	}
}

// nolint: funlen
func SimulateMsgExecuteContract(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	protoCdc *codec.ProtoCodec) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "contract not exists yet"), nil, nil
		}

		// should creator execute the msg
		creatorAddr, _ := sdk.AccAddressFromBech32(info.Creator)
		simAccount, _ := simtypes.FindAccount(accs, creatorAddr)
		account := ak.GetAccount(ctx, simAccount.Address)

		spendableCoins := bk.SpendableCoins(ctx, simAccount.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "unable to generate fee"), nil, err
		}

		spendableCoins = spendableCoins.Sub(fees)
		spendableCoins = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, spendableCoins.AmountOf(core.MicroLunaDenom)))
		if spendableCoins.Empty() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "unable to generate deposit"), nil, err
		}

		if err := bk.IsSendEnabledCoins(ctx, spendableCoins...); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "send not enabled"), nil, nil
		}

		msg := types.NewMsgExecuteContract(simAccount.Address, contractAddr, []byte(`{"release": {}}`), spendableCoins)
		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			if strings.Contains(err.Error(), "insufficient fee") {
				return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "ignore tax error"), nil, nil
			}

			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", protoCdc), nil, nil
	}
}

// nolint: funlen
func SimulateMsgMigrateContract(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	protoCdc *codec.ProtoCodec) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "contract not exists yet"), nil, nil
		}

		if len(info.Admin) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "contract has no admin"), nil, nil
		}

		targetCodeID := 1
		_, err = k.GetCodeInfo(ctx, 2)
		if err == nil && info.CodeID == 1 {
			targetCodeID = 2
		}

		// should admin migrate the msg
		adminAddr, _ := sdk.AccAddressFromBech32(info.Admin)
		simAccount, _ := simtypes.FindAccount(accs, adminAddr)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := bk.SpendableCoins(ctx, adminAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "unable to generate fee"), nil, err
		}

		// never used more
		// spendableCoins = spendableCoins.Sub(fees)

		migData := map[string]interface{}{
			"verifier": info.Creator,
		}
		migDataBz, err := json.Marshal(migData)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "unable to marshal json"), nil, err
		}

		msg := types.NewMsgMigrateContract(simAccount.Address, contractAddr, uint64(targetCodeID), migDataBz)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", protoCdc), nil, nil
	}
}

// nolint: funlen
func SimulateMsgUpdateContractAdmin(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractAdmin, "contract not exists yet"), nil, nil
		}

		if len(info.Admin) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractAdmin, "contract has no admin"), nil, nil
		}

		// should admin execute the msg
		adminAddr, _ := sdk.AccAddressFromBech32(info.Admin)
		simAccount, _ := simtypes.FindAccount(accs, adminAddr)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := bk.SpendableCoins(ctx, adminAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractAdmin, "unable to generate fee"), nil, err
		}

		newAdminAccount, _ := simtypes.RandomAcc(r, accs)
		if simAccount.Address.Equals(newAdminAccount.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractAdmin, "same account selected"), nil, nil
		}

		msg := types.NewMsgUpdateContractAdmin(simAccount.Address, newAdminAccount.Address, contractAddr)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}

// nolint: funlen
func SimulateMsgClearContractAdmin(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClearContractAdmin, "contract not exists yet"), nil, nil
		}

		if len(info.Admin) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClearContractAdmin, "contract has no admin"), nil, nil
		}

		// should admin execute the msg
		adminAddr, _ := sdk.AccAddressFromBech32(info.Admin)
		simAccount, _ := simtypes.FindAccount(accs, adminAddr)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := bk.SpendableCoins(ctx, adminAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgClearContractAdmin, "unable to generate fee"), nil, err
		}

		msg := types.NewMsgClearContractAdmin(simAccount.Address, contractAddr)

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "", nil), nil, nil
	}
}
