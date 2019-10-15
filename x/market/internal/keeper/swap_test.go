package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

func TestApplySwapToPool(t *testing.T) {
	input := CreateTestInput(t)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000))
	askCoin := sdk.NewDecCoin(core.MicroSDRDenom, sdk.NewInt(1700))

	oldSDRPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)
	newSDRPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)

	sdrDiff := newSDRPoolDelta.Sub(oldSDRPoolDelta)

	require.Equal(t, sdk.NewDec(-1700), sdrDiff)
}

func TestComputeSwap(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

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
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

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
