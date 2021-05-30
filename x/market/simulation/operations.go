package simulation

// DONTCOVER

import (
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle"

	"github.com/terra-money/core/x/market/internal/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgSwap = "op_weight_msg_swap"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak authkeeper.AccountKeeper, ok oracle.Keeper,
) simulation.WeightedOperations {
	var weightMsgSwap int
	appParams.GetOrGenerate(cdc, OpWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = simappparams.DefaultWeightMsgSend
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgSwap,
			SimulateMsgSwap(ak, ok),
		),
	}
}

// SimulateMsgSwap generates a MsgSwap with random values.
// nolint: funlen
func SimulateMsgSwap(ak authkeeper.AccountKeeper, ok oracle.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		var whitelist []string
		ok.IterateLunaExchangeRates(ctx, func(denom string, ex sdk.Dec) bool {
			whitelist = append(whitelist, denom)
			return false
		})

		var offerDenom string
		var askDenom string
		whitelistLen := len(whitelist)
		if whitelistLen == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		if randVal := simulation.RandIntBetween(r, 0, whitelistLen*2); randVal < whitelistLen {
			offerDenom = core.MicroLunaDenom
			askDenom = whitelist[randVal]
		} else {
			offerDenom = whitelist[randVal-whitelistLen]
			askDenom = core.MicroLunaDenom
		}

		amount := simulation.RandomAmount(r, spendableCoins.AmountOf(offerDenom).Sub(fees.AmountOf(offerDenom)))
		if amount.Equal(sdk.ZeroInt()) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgSwap(simAccount.Address, sdk.NewCoin(offerDenom, amount), askDenom)

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

		if err != nil && !strings.Contains(err.Error(), "no price registered") {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgSwapSend generates a MsgSwapSend with random values.
// nolint: funlen
func SimulateMsgSwapSend(ak authkeeper.AccountKeeper, ok oracle.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)
		receiverAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)

		spendableCoins := account.SpendableCoins(ctx.BlockTime())
		fees, err := simulation.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		var whitelist []string
		ok.IterateLunaExchangeRates(ctx, func(denom string, ex sdk.Dec) bool {
			whitelist = append(whitelist, denom)
			return false
		})

		var offerDenom string
		var askDenom string
		whitelistLen := len(whitelist)
		if whitelistLen == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		if randVal := simulation.RandIntBetween(r, 0, whitelistLen*2); randVal < whitelistLen {
			offerDenom = core.MicroLunaDenom
			askDenom = whitelist[randVal]
		} else {
			offerDenom = whitelist[randVal-whitelistLen]
			askDenom = core.MicroLunaDenom
		}

		amount := simulation.RandomAmount(r, spendableCoins.AmountOf(offerDenom).Sub(fees.AmountOf(offerDenom)))
		if amount.Equal(sdk.ZeroInt()) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgSwapSend(simAccount.Address, receiverAccount.Address, sdk.NewCoin(offerDenom, amount), askDenom)

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

		if err != nil && !strings.Contains(err.Error(), "no price registered") {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
