package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market"
)

// SimulateMsgSwap generates a MsgSwap with random values
func SimulateMsgSwap(k market.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)

		msg := market.NewMsgSwap(acc.Address, sdk.NewInt64Coin(core.MicroLunaDenom, rand.Int63()), core.MicroSDRDenom)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(market.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := market.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
