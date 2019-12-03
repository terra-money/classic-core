package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

func TestApplySwapToPool(t *testing.T) {
	input := CreateTestInput(t)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000))
	askCoin := sdk.NewDecCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	oldSDRPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	sdrDiff := newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(-1700), sdrDiff)

	// reverse swap
	offerCoin = sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	askCoin = sdk.NewDecCoin(core.MicroLunaDenom, sdk.NewInt(1000))
	oldSDRPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	sdrDiff = newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(1700), sdrDiff)

	// TERRA <> TERRA, no pool changes are expected
	offerCoin = sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	askCoin = sdk.NewDecCoin(core.MicroKRWDenom, sdk.NewInt(3400))
	oldSDRPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	sdrDiff = newSDRPoolDelta.Sub(oldSDRPoolDelta)
	require.Equal(t, sdk.NewDec(0), sdrDiff)
}
func TestComputeSwap(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	for i := 0; i < 100; i++ {
		swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
		offerCoin := sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)
		retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)

		require.NoError(t, err)
		require.True(t, spread.GTE(input.MarketKeeper.MinSpread(input.Ctx)))
		require.Equal(t, sdk.NewDecFromInt(offerCoin.Amount).Quo(lunaPriceInSDR), retCoin.Amount)
	}

	offerCoin := sdk.NewCoin(core.MicroSDRDenom, lunaPriceInSDR.QuoInt64(2).TruncateInt())
	_, _, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	require.Error(t, err)
}

func TestComputeInternalSwap(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	for i := 0; i < 100; i++ {
		offerCoin := sdk.NewDecCoin(core.MicroSDRDenom, lunaPriceInSDR.MulInt64(rand.Int63()+1).TruncateInt())
		retCoin, err := input.MarketKeeper.ComputeInternalSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
		require.NoError(t, err)
		require.Equal(t, offerCoin.Amount.Quo(lunaPriceInSDR), retCoin.Amount)
	}

	offerCoin := sdk.NewDecCoin(core.MicroSDRDenom, lunaPriceInSDR.QuoInt64(2).TruncateInt())
	_, err := input.MarketKeeper.ComputeInternalSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	require.Error(t, err)
}

func TestIlliquidTobinTaxListParams(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	lunaPriceInMNT := sdk.NewDecWithPrec(7652, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroMNTDenom, lunaPriceInMNT)

	// Case 1: tobin tax 2% due to umnt denom
	params := input.MarketKeeper.GetParams(input.Ctx)
	params.TobinTax = sdk.NewDecWithPrec(25, 4)
	params.IlliquidTobinTaxList = types.TobinTaxList{
		types.TobinTax{
			Denom:   core.MicroSDRDenom,
			TaxRate: sdk.NewDecWithPrec(25, 4),
		},
		types.TobinTax{
			Denom:   core.MicroMNTDenom,
			TaxRate: sdk.NewDecWithPrec(2, 2),
		},
	}
	input.MarketKeeper.SetParams(input.Ctx, params)

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	offerCoin := sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)
	_, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroMNTDenom)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecWithPrec(2, 2), spread)

	// Case 2: tobin tax 5% due to default
	params.TobinTax = sdk.NewDecWithPrec(5, 2)
	input.MarketKeeper.SetParams(input.Ctx, params)

	swapAmountInSDR = lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	offerCoin = sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)
	_, spread, err = input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroMNTDenom)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecWithPrec(5, 2), spread)

}
