package simulation

// DONTCOVER

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

	"github.com/terra-project/core/x/wasm/keeper"
	"github.com/terra-project/core/x/wasm/types"
)

const (
	OpWeightMsgStoreCode           = "op_weight_msg_store_code"
	OpWeightMsgInstantiateContract = "op_weight_msg_instantiate_contract"
	OpWeightMsgExecuteContract     = "op_weight_msg_execute_contract"
	OpWeightMsgMigrateContract     = "op_weight_msg_migrate_contract"
	OpWeightMsgUpdateContractOwner = "op_weight_msg_update_contract_owner"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgStoreCode int
	var weightMsgInstantiateContract int
	var weightMsgExecuteContract int
	var weightMsgMigrateContract int
	var weightMsgUpdateContractOwner int
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

	appParams.GetOrGenerate(cdc, OpWeightMsgUpdateContractOwner, &weightMsgUpdateContractOwner, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateContractOwner = 3
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgStoreCode,
			SimulateMsgStoreCode(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgInstantiateContract,
			SimulateMsgInstantiateContract(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgExecuteContract,
			SimulateMsgExecuteContract(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMigrateContract,
			SimulateMsgMigrateContract(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateContractOwner,
			SimulateMsgUpdateContractOwner(ak, bk, k),
		),
	}
}

func mustLoad(path string) []byte {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

var testContract []byte

// nolint: funlen
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

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

type initMsg struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

// nolint: funlen
func SimulateMsgInstantiateContract(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
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

		msg := types.NewMsgInstantiateContract(fredAcc.Address, 1, initMsgBz, nil, simtypes.RandIntBetween(r, 1, 2) == 1)

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

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// nolint: funlen
func SimulateMsgExecuteContract(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "contract not exists yet"), nil, nil
		}

		// should owner execute the msg
		ownerAddr, _ := sdk.AccAddressFromBech32(info.Owner)
		simAccount, _ := simtypes.FindAccount(accs, ownerAddr)
		account := ak.GetAccount(ctx, simAccount.Address)

		spendableCoins := bk.SpendableCoins(ctx, simAccount.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "unable to generate fee"), nil, err
		}

		spendableCoins = spendableCoins.Sub(fees)
		spendableCoins = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, spendableCoins.AmountOf(sdk.DefaultBondDenom)))

		if err := bk.SendEnabledCoins(ctx, spendableCoins...); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgExecuteContract, "send not enabled"), nil, nil
		}

		msg := types.NewMsgExecuteContract(simAccount.Address, contractAddr, []byte(`{"release": {}}`), simtypes.RandSubsetCoins(r, spendableCoins))
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

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// nolint: funlen
func SimulateMsgMigrateContract(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "contract not exists yet"), nil, nil
		}

		targetCodeID := 1
		_, err = k.GetCodeInfo(ctx, 2)
		if err == nil && info.CodeID == 1 {
			targetCodeID = 2
		}

		// should owner execute the msg
		ownerAddr, _ := sdk.AccAddressFromBech32(info.Owner)
		simAccount, _ := simtypes.FindAccount(accs, ownerAddr)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := bk.SpendableCoins(ctx, ownerAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgMigrateContract, "unable to generate fee"), nil, err
		}

		spendableCoins = spendableCoins.Sub(fees)

		migData := map[string]interface{}{
			"verifier": info.Owner,
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

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// nolint: funlen
func SimulateMsgUpdateContractOwner(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos1hqrdl6wstt8qzshwc6mrumpjk9338k0lr4dqxd")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractOwner, "contract not exists yet"), nil, nil
		}

		// should owner execute the msg
		ownerAddr, _ := sdk.AccAddressFromBech32(info.Owner)
		simAccount, _ := simtypes.FindAccount(accs, ownerAddr)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := bk.SpendableCoins(ctx, ownerAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractOwner, "unable to generate fee"), nil, err
		}

		newOwnerAccount, _ := simtypes.RandomAcc(r, accs)
		if simAccount.Address.Equals(newOwnerAccount.Address) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateContractOwner, "same account selected"), nil, nil
		}

		msg := types.NewMsgUpdateContractOwner(simAccount.Address, newOwnerAccount.Address, contractAddr)

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

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
