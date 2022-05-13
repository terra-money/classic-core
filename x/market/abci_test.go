package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	terraDelta := sdk.NewDecWithPrec(17987573223725367, 3)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, terraDelta)

	for i := 0; i < 100; i++ {
		terraDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)

		poolRecoveryPeriod := int64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
		terraRegressionAmt := terraDelta.QuoInt64(poolRecoveryPeriod)

		EndBlocker(input.Ctx, input.MarketKeeper)

		terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
		require.Equal(t, terraDelta.Sub(terraRegressionAmt), terraPoolDelta)
	}
}

func TestSetMinSpread(t *testing.T) {
	input := keeper.CreateTestInput(t)
	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	input.Ctx = input.Ctx.WithChainID(core.ColumbusChainID).WithBlockHeight(core.SwapDisableForkHeight)
	BeginBlocker(input.Ctx, input.MarketKeeper)

	require.Equal(t, sdk.OneDec(), input.MarketKeeper.MinStabilitySpread(input.Ctx))
}
