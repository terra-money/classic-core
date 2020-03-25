package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle"
)

var (
	exchangeRates string
)

// SimulateMsgPrevote generates a MsgPrevote with random values
func SimulateMsgPrevote(k oracle.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		valAddr := sdk.ValAddress(acc.Address)
		hash := oracle.GetVoteHash("1234", sdk.NewDec(1700), core.MicroSDRDenom, valAddr)

		msg := oracle.NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, acc.Address, valAddr)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(oracle.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := oracle.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgVote generates a MsgVote with random values
func SimulateMsgVote(k oracle.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		valAddr := sdk.ValAddress(acc.Address)

		msg := oracle.NewMsgExchangeRateVote(sdk.NewDec(1700), "1234", core.MicroSDRDenom, acc.Address, valAddr)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(oracle.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := oracle.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgDelegateFeedConsent generates a MsgDelegateFeedConsent with random values
func SimulateMsgDelegateFeedConsent(k oracle.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		acc2 := simulation.RandomAcc(r, accs)
		valAddr := sdk.ValAddress(acc.Address)
		msg := oracle.NewMsgDelegateFeedConsent(valAddr, acc2.Address)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(oracle.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := oracle.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgAggregateExchangeRatePrevote generates a MsgAggregateExchangeRatePrevote with random values
func SimulateMsgAggregateExchangeRatePrevote(k oracle.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		valAddr := sdk.ValAddress(acc.Address)
		exchangeRates = fmt.Sprintf("%s%s,%s%s", simulation.RandomDecAmount(r,
			sdk.NewDec(1000000)).String(),
			core.MicroKRWDenom,
			simulation.RandomDecAmount(r, sdk.NewDec(1)).String(),
			core.MicroSDRDenom)

		hash := oracle.GetAggregateVoteHash("1234", exchangeRates, valAddr)

		msg := oracle.NewMsgAggregateExchangeRatePrevote(hash, acc.Address, valAddr)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(oracle.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := oracle.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

// SimulateMsgAggregateExchangeRateVote generates a MsgVote with random values
func SimulateMsgAggregateExchangeRateVote(k oracle.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		valAddr := sdk.ValAddress(acc.Address)

		if len(exchangeRates) == 0 {
			return simulation.NoOpMsg(oracle.ModuleName), nil, nil
		}

		msg := oracle.NewMsgAggregateExchangeRateVote("1234", exchangeRates, acc.Address, valAddr)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(oracle.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := oracle.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			exchangeRates = ""
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
