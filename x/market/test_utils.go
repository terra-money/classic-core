// nolint:deadcode unused DONTCOVER
package market

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/keeper"
)

var (
	uSDRAmt    = sdk.NewInt(1005 * core.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10)

	randomPrice = sdk.NewDec(1700)
)

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.MarketKeeper.GetParams(input.Ctx)
	input.MarketKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, randomPrice)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, randomPrice)
	h := NewHandler(input.MarketKeeper)

	return input, h
}
