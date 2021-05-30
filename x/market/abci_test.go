package market

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/x/market/internal/keeper"
)

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	delta := sdk.NewDecWithPrec(17987573223725367, 3)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, delta)

	for i := 0; i < 100; i++ {
		delta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
		regressionAmt := delta.QuoInt64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))

		EndBlocker(input.Ctx, input.MarketKeeper)

		terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
		require.Equal(t, delta.Sub(regressionAmt), terraPoolDelta)
	}
}
