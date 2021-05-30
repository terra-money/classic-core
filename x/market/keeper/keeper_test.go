package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMintPoolDeltaUpdate(t *testing.T) {
	input := CreateTestInput(t)

	terraPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	require.Equal(t, sdk.ZeroDec(), terraPoolDelta)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetMintPoolDelta(input.Ctx, diff)

	terraPoolDelta = input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	require.Equal(t, diff, terraPoolDelta)
}

func TestBurnPoolDeltaUpdate(t *testing.T) {
	input := CreateTestInput(t)

	terraPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	require.Equal(t, sdk.ZeroDec(), terraPoolDelta)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetBurnPoolDelta(input.Ctx, diff)

	terraPoolDelta = input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	require.Equal(t, diff, terraPoolDelta)
}

// TestReplenishPools tests that
// each pools move towards base pool
func TestReplenishPools(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	recoveryPeriod := int64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
	mintBasePool := input.MarketKeeper.MintBasePool(input.Ctx)
	mintPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	require.True(t, mintPoolDelta.IsZero())

	burnBasePool := input.MarketKeeper.BurnBasePool(input.Ctx)
	burnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	require.True(t, burnPoolDelta.IsZero())

	// Positive delta
	mintDiff := mintBasePool.QuoInt64((int64)(core.BlocksPerDay))
	input.MarketKeeper.SetMintPoolDelta(input.Ctx, mintDiff)

	burnDiff := burnBasePool.QuoInt64((int64)(core.BlocksPerDay))
	input.MarketKeeper.SetBurnPoolDelta(input.Ctx, burnDiff)

	input.MarketKeeper.ReplenishPools(input.Ctx)

	mintPoolDelta = input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	mintReplenishAmt := mintDiff.QuoInt64(recoveryPeriod)
	expectedMintDelta := mintDiff.Sub(mintReplenishAmt)
	require.Equal(t, expectedMintDelta, mintPoolDelta)

	burnPoolDelta = input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	burnReplenishAmt := burnDiff.QuoInt64(recoveryPeriod)
	expectedBurnDelta := mintDiff.Sub(burnReplenishAmt)
	require.Equal(t, expectedBurnDelta, burnPoolDelta)
}
