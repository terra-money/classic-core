package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	mintDelta := sdk.NewDecWithPrec(17987573223725367, 3)
	burnDelta := sdk.NewDecWithPrec(23984923849121938, 3)
	input.MarketKeeper.SetMintPoolDelta(input.Ctx, mintDelta)
	input.MarketKeeper.SetBurnPoolDelta(input.Ctx, burnDelta)

	for i := 0; i < 100; i++ {
		mintDelta = input.MarketKeeper.GetMintPoolDelta(input.Ctx)
		burnDelta = input.MarketKeeper.GetBurnPoolDelta(input.Ctx)

		poolRecoveryPeriod := int64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
		mintRegressionAmt := mintDelta.QuoInt64(poolRecoveryPeriod)
		burnRegressionAmt := burnDelta.QuoInt64(poolRecoveryPeriod)

		EndBlocker(input.Ctx, input.MarketKeeper)

		mintPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
		burnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
		require.Equal(t, mintDelta.Sub(mintRegressionAmt), mintPoolDelta)
		require.Equal(t, burnDelta.Sub(burnRegressionAmt), burnPoolDelta)
	}
}
