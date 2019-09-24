package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

func TestLunaPoolUpdate(t *testing.T) {
	input := CreateTestInput(t)

	basePool := input.MarketKeeper.GetBasePool(input.Ctx)
	lunaPool := input.MarketKeeper.GetLunaPool(input.Ctx)
	require.Equal(t, basePool, lunaPool)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetLunaPool(input.Ctx, lunaPool.Sub(diff))

	lunaPool = input.MarketKeeper.GetLunaPool(input.Ctx)
	require.Equal(t, basePool.Sub(diff), lunaPool)
}

func TestTerraPoolUpdate(t *testing.T) {
	input := CreateTestInput(t)

	basePool := input.MarketKeeper.GetBasePool(input.Ctx)
	terraPool := input.MarketKeeper.GetTerraPool(input.Ctx)
	require.Equal(t, basePool, terraPool)

	diff := sdk.NewDec(10)
	input.MarketKeeper.SetTerraPool(input.Ctx, terraPool.Sub(diff))

	terraPool = input.MarketKeeper.GetTerraPool(input.Ctx)
	require.Equal(t, basePool.Sub(diff), terraPool)
}

func TestUpdatePools(t *testing.T) {
	input := CreateTestInput(t)

	// oracle price
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	for i := 0; i < 100; i++ {
		delta := sdk.NewDecWithPrec(rand.Int63n(1000), 4)

		supply := input.SupplyKeeper.GetSupply(input.Ctx)
		total := supply.GetTotal()
		issuance := total.AmountOf(core.MicroLunaDenom)
		issuance = sdk.OneDec().Add(delta).MulInt(issuance).TruncateInt() // (1+delta) * issuance

		total = total.Add(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, issuance)))
		supply = supply.SetTotal(total)
		input.SupplyKeeper.SetSupply(input.Ctx, supply)

		input.MarketKeeper.UpdatePools(input.Ctx)
		expectedBasePool := input.MarketKeeper.DailyTerraLiquidityRatio(input.Ctx).MulInt(total.AmountOf(core.MicroLunaDenom))

		require.Equal(t, expectedBasePool, input.MarketKeeper.GetBasePool(input.Ctx))
	}
}

// TestReplenishPools tests that
// each pools move towards base pool
func TestReplenishPools(t *testing.T) {
	input := CreateTestInput(t)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())
	_, err := input.MarketKeeper.UpdatePools(input.Ctx)
	require.NoError(t, err)

	basePool := input.MarketKeeper.GetBasePool(input.Ctx)
	terraPool := input.MarketKeeper.GetTerraPool(input.Ctx)
	require.Equal(t, basePool, terraPool)

	diff := basePool.QuoInt64(core.BlocksPerDay)
	input.MarketKeeper.SetTerraPool(input.Ctx, terraPool.Add(diff))

	input.MarketKeeper.ReplenishPools(input.Ctx)
	terraPool = input.MarketKeeper.GetTerraPool(input.Ctx)
	require.Equal(t, basePool, terraPool)
}
