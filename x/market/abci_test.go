package market

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/x/market/internal/keeper"
)

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	delta := sdk.NewDec(1000000)
	regressionAmt := delta.QuoInt64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, delta)

	EndBlocker(input.Ctx, input.MarketKeeper)

	terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	require.Equal(t, delta.Sub(regressionAmt), terraPoolDelta)
}
