package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

func TestTerraPoolDeltaUpdate(t *testing.T) {
	input := CreateTestInput(t)

	terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	require.Equal(t, sdk.ZeroDec(), terraPoolDelta)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, diff)

	terraPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	require.Equal(t, diff, terraPoolDelta)
}

// TestReplenishPools tests that
// each pools move towards base pool
func TestReplenishPools(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	basePool := input.MarketKeeper.BasePool(input.Ctx)
	terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	require.True(t, terraPoolDelta.IsZero())

	diff := basePool.QuoInt64(core.BlocksPerDay)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, diff)

	input.MarketKeeper.ReplenishPools(input.Ctx)

	terraPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	replenishAmt := diff.QuoInt64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
	expectedDelta := diff.Sub(replenishAmt)
	require.Equal(t, expectedDelta, terraPoolDelta)
}
