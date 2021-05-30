package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/core/x/wasm/internal/keeper"
	"github.com/terra-money/core/x/wasm/internal/types"
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
	appParams simulation.AppParams, cdc *codec.Codec,
	ak authkeeper.AccountKeeper, bk bank.Keeper, k keeper.Keeper,
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
			SimulateMsgStoreCode(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgInstantiateContract,
			SimulateMsgInstantiateContract(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgExecuteContract,
			SimulateMsgExecuteContract(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgMigrateContract,
			SimulateMsgMigrateContract(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUpdateContractOwner,
			SimulateMsgUpdateContractOwner(ak, k),
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
func SimulateMsgStoreCode(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		_, err := k.GetCodeInfo(ctx, 2)
		if err == nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		if testContract == nil {
			loadContract()
		}

		simAccount, _ := simulation.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgStoreCode(simAccount.Address, testContract)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
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
func SimulateMsgInstantiateContract(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		bobAcc, _ := simulation.RandomAcc(r, accs)
		fredAcc, _ := simulation.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, fredAcc.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		initMsg := initMsg{
			Verifier:    fredAcc.Address.String(),
			Beneficiary: bobAcc.Address.String(),
		}

		initMsgBz, err := json.Marshal(initMsg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		_, err = k.GetCodeInfo(ctx, 1)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgInstantiateContract(fredAcc.Address, 1, initMsgBz, nil, simulation.RandIntBetween(r, 1, 2) == 1)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			fredAcc.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			if strings.Contains(err.Error(), "insufficient fee") {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}

			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// nolint: funlen
func SimulateMsgExecuteContract(ak authkeeper.AccountKeeper, bk bank.Keeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if !bk.GetSendEnabled(ctx) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// should owner execute the msg
		simAccount, _ := simulation.FindAccount(accs, info.Owner)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		spendableCoins = spendableCoins.Sub(fees)
		spendableCoins = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, spendableCoins.AmountOf(sdk.DefaultBondDenom)))

		msg := types.NewMsgExecuteContract(simAccount.Address, contractAddr, []byte(`{"release": {}}`), simulation.RandSubsetCoins(r, spendableCoins))
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			if strings.Contains(err.Error(), "insufficient fee") {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}

			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// nolint: funlen
func SimulateMsgMigrateContract(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		targetCodeID := 1
		_, err = k.GetCodeInfo(ctx, 2)
		if err == nil && info.CodeID == 1 {
			targetCodeID = 2
		}

		// should owner execute the msg
		simAccount, _ := simulation.FindAccount(accs, info.Owner)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		spendableCoins = spendableCoins.Sub(fees)

		migData := map[string]interface{}{
			"verifier": info.Owner.String(),
		}
		migDataBz, err := json.Marshal(migData)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgMigrateContract(simAccount.Address, contractAddr, uint64(targetCodeID), migDataBz)
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		return simulation.NewOperationMsg(msg, true, ""), nil, err
	}
}

// nolint: funlen
func SimulateMsgUpdateContractOwner(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		contractAddr, _ := sdk.AccAddressFromBech32("cosmos1hqrdl6wstt8qzshwc6mrumpjk9338k0lr4dqxd")
		info, err := k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// should owner execute the msg
		simAccount, _ := simulation.FindAccount(accs, info.Owner)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		newOwnerAccount, _ := simulation.RandomAcc(r, accs)
		if simAccount.Address.Equals(newOwnerAccount.Address) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgUpdateContractOwner(simAccount.Address, newOwnerAccount.Address, contractAddr)
		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			fmt.Println(err.Error())
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, err
	}
}
