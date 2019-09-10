package market

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/keeper"
)

func TestPoolsUpdate(t *testing.T) {
	input := keeper.CreateTestInput(t)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, sdk.OneDec())

	// Update Pools at the non-last block of interval
	input.Ctx = input.Ctx.WithBlockHeight(1)
	issuance := sdk.NewInt(12345)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, issuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	EndBlocker(input.Ctx, input.MarketKeeper)
	basePool := input.MarketKeeper.GetBasePool(input.Ctx)
	require.Equal(t, input.MarketKeeper.DailyTerraLiquidityRatio(input.Ctx).MulInt(issuance), basePool)

	// Update Pools at the last block of interval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch - 1)
	issuance = sdk.NewInt(1000000)
	supply = input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, issuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	EndBlocker(input.Ctx, input.MarketKeeper)
	basePool = input.MarketKeeper.GetBasePool(input.Ctx)
	require.Equal(t, input.MarketKeeper.DailyTerraLiquidityRatio(input.Ctx).MulInt(issuance), basePool)

	// Update Pools at the last block of the another interval
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerEpoch*2 - 1)
	issuance = sdk.NewInt(10000000000)
	supply = input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, issuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	EndBlocker(input.Ctx, input.MarketKeeper)
	basePool = input.MarketKeeper.GetBasePool(input.Ctx)
	require.Equal(t, input.MarketKeeper.DailyTerraLiquidityRatio(input.Ctx).MulInt(issuance), basePool)
}

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	delta := sdk.NewDec(1000000)
	regressionAmt := sdk.NewDec(1)
	basePool := regressionAmt.MulInt64(core.BlocksPerDay)
	input.MarketKeeper.SetBasePool(input.Ctx, basePool)
	input.MarketKeeper.SetTerraPool(input.Ctx, basePool.Sub(delta))

	EndBlocker(input.Ctx, input.MarketKeeper)

	terraPool := input.MarketKeeper.GetTerraPool(input.Ctx)

	require.Equal(t, basePool.Sub(delta).Add(regressionAmt), terraPool)
}
