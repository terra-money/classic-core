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
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-project/core/x/wasm/internal/keeper"
	"github.com/terra-project/core/x/wasm/internal/types"
)

const (
	OpWeightMsgStoreCoce           = "op_weight_msg_store_code"
	OpWeightMsgInstantiateContract = "op_weight_msg_instantiate_contract"
	OpWeightMsgExecuteContract     = "op_weight_msg_execute_contract"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak authkeeper.AccountKeeper, k keeper.Keeper,
) simulation.WeightedOperations {
	var weightMsgStoreCode int
	var weightMsgInstantiateContract int
	var weightMsgExecuteContract int
	appParams.GetOrGenerate(cdc, OpWeightMsgStoreCoce, &weightMsgStoreCode, nil,
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
			SimulateMsgExecuteContract(ak, k),
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

var (
	bob  sdk.AccAddress
	fred sdk.AccAddress
)

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
		simAccount, _ := simulation.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if bob.Empty() {
			bobAcc, _ := simulation.RandomAcc(r, accs)
			bob = bobAcc.Address
		}

		if fred.Empty() {
			fredAcc, _ := simulation.RandomAcc(r, accs)
			fred = fredAcc.Address
		}

		initMsg := initMsg{
			Verifier:    fred.String(),
			Beneficiary: bob.String(),
		}

		initMsgBz, err := json.Marshal(initMsg)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		_, err = k.GetCodeInfo(ctx, 1)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgInstantiateContract(simAccount.Address, 1, initMsgBz, nil, true)

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
func SimulateMsgExecuteContract(ak authkeeper.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {
		if fred.Empty() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// should fred execute the msg
		simAccount, _ := simulation.FindAccount(accs, fred)
		account := ak.GetAccount(ctx, fred)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		contractAddr, _ := sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
		_, err = k.GetContractInfo(ctx, contractAddr)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgExecuteContract(fred, contractAddr, []byte(`{"release": {}}`), nil)

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
