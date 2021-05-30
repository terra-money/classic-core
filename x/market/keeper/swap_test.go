package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestApplySwapToPool_Mint(t *testing.T) {
	input := CreateTestInput(t)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000))
	askCoin := sdk.NewDecCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	oldMintDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)

	newMintDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	sdrDiff := newMintDelta.Sub(oldMintDelta)
	require.Equal(t, sdk.NewDec(1700), sdrDiff)
}

func TestApplySwapToPool_Burn(t *testing.T) {
	input := CreateTestInput(t)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	askCoin := sdk.NewDecCoin(core.MicroLunaDenom, sdk.NewInt(1000))
	oldBurnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)

	newBurnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	sdrDiff := newBurnPoolDelta.Sub(oldBurnPoolDelta)
	require.Equal(t, sdk.NewDec(1700), sdrDiff)
}

func TestApplySwapToPool_TerraTerra(t *testing.T) {
	input := CreateTestInput(t)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	// TERRA <> TERRA, no pool changes are expected
	offerCoin := sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1700))
	askCoin := sdk.NewDecCoin(core.MicroKRWDenom, sdk.NewInt(3400))
	oldMintPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	oldBurnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	input.MarketKeeper.ApplySwapToPool(input.Ctx, offerCoin, askCoin)

	newMintPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	newBurnPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)

	require.Equal(t, oldMintPoolDelta, newMintPoolDelta)
	require.Equal(t, oldBurnPoolDelta, newBurnPoolDelta)
}

func TestComputeSwap_Mint(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewInt64Coin(core.MicroLunaDenom, rand.Int63()%10000+2)

	basePool := input.MarketKeeper.MintBasePool(input.Ctx)
	poolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	require.True(t, poolDelta.IsZero())

	cp := basePool.Mul(basePool)
	offerPool := basePool.Add(poolDelta)
	askPool := cp.Quo(offerPool)

	// ask = ask_pool - cp / (offer_pool + offer_amount)
	// spread % = (return - ask) / return
	expectedReturnAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(lunaPriceInSDR)
	expectedSpread := expectedReturnAmount.Sub(askPool.Sub(cp.Quo(offerPool.Add(expectedReturnAmount)))).Quo(expectedReturnAmount)

	// without min stability fee
	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, expectedReturnAmount, retCoin.Amount)
	require.Equal(t, expectedSpread, spread)

	// with min stability fee
	params.MinStabilitySpread = sdk.OneDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	retCoin, spread, err = input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, expectedReturnAmount, retCoin.Amount)
	require.Equal(t, sdk.OneDec(), spread)
}

func TestComputeSwap_Burn(t *testing.T) {
	input := CreateTestInput(t)

	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	offerCoin := sdk.NewInt64Coin(core.MicroSDRDenom, rand.Int63()%10000+2)

	basePool := input.MarketKeeper.BurnBasePool(input.Ctx)
	poolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	require.True(t, poolDelta.IsZero())

	cp := basePool.Mul(basePool)
	offerPool := basePool.Add(poolDelta)
	askPool := cp.Quo(offerPool)

	// ask = ask_pool - cp / (offer_pool + offer_amount)
	// spread % = (return - ask) / return
	expectedReturnAmountInSDR := sdk.NewDecFromInt(offerCoin.Amount)
	expectedReturnAmount := expectedReturnAmountInSDR.Quo(lunaPriceInSDR)
	expectedSpread := expectedReturnAmountInSDR.Sub(askPool.Sub(cp.Quo(offerPool.Add(expectedReturnAmountInSDR)))).Quo(expectedReturnAmountInSDR)

	// without min stability fee
	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	require.NoError(t, err)
	require.Equal(t, expectedReturnAmount, retCoin.Amount)
	require.Equal(t, expectedSpread, spread)

	// with min stability fee
	params.MinStabilitySpread = sdk.OneDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	retCoin, spread, err = input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	require.NoError(t, err)
	require.Equal(t, expectedReturnAmount, retCoin.Amount)
	require.Equal(t, sdk.OneDec(), spread)
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

	tobinTax := sdk.NewDecWithPrec(25, 4)
	params := input.MarketKeeper.GetParams(input.Ctx)
	input.MarketKeeper.SetParams(input.Ctx, params)

	illiquidFactor := sdk.NewDec(2)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroSDRDenom, tobinTax)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroMNTDenom, tobinTax.Mul(illiquidFactor))

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	offerCoin := sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)
	_, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroMNTDenom)
	require.NoError(t, err)
	require.Equal(t, tobinTax.Mul(illiquidFactor), spread)
}
