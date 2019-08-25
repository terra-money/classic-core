package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

func TestLastDayLunaIssuanceUpdate(t *testing.T) {
	input := CreateTestInput(t)

	issuance := input.MarketKeeper.GetLastDayIssuance(input.Ctx).AmountOf(core.MicroLunaDenom)
	require.True(t, issuance.IsZero())

	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.OneInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.MarketKeeper.UpdateLastDayIssuance(input.Ctx)
	issuance = input.MarketKeeper.GetLastDayIssuance(input.Ctx).AmountOf(core.MicroLunaDenom)
	require.Equal(t, sdk.OneInt(), issuance)
}

func TestComputeLunaDelta(t *testing.T) {
	input := CreateTestInput(t)

	for i := 0; i < 100; i++ {
		expectedDelta := sdk.NewDecWithPrec(rand.Int63n(1000), 3)
		issuance := input.SupplyKeeper.GetSupply(input.Ctx).GetTotal().AmountOf(core.MicroLunaDenom)
		change := expectedDelta.MulInt(issuance).TruncateInt()
		input.MarketKeeper.UpdateLastDayIssuance(input.Ctx)
		delta := input.MarketKeeper.ComputeLunaDelta(input.Ctx.WithBlockHeight(core.BlocksPerDay), change)

		require.Equal(t, expectedDelta, delta)
	}
}

func TestComputeLunaSwapSpread(t *testing.T) {
	input := CreateTestInput(t)

	for i := 0; i < 100; i++ {
		delta := sdk.NewDecWithPrec(rand.Int63n(1000), 3)
		spread := input.MarketKeeper.ComputeLunaSwapSpread(input.Ctx, delta)
		require.True(t, spread.GTE(input.MarketKeeper.MinSwapSpread(input.Ctx)))
		require.True(t, spread.LTE(input.MarketKeeper.MaxSwapSpread(input.Ctx)))
	}

	spread := input.MarketKeeper.ComputeLunaSwapSpread(input.Ctx, sdk.ZeroDec())
	require.Equal(t, input.MarketKeeper.MinSwapSpread(input.Ctx), spread)

	spread = input.MarketKeeper.ComputeLunaSwapSpread(input.Ctx, sdk.OneDec())
	require.Equal(t, input.MarketKeeper.MaxSwapSpread(input.Ctx), spread)
}

func TestGetSwapCoin(t *testing.T) {
	input := CreateTestInput(t)
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	// zero day (min spread)
	for i := 0; i < 100; i++ {
		offerCoin := sdk.NewCoin(core.MicroSDRDenom, lunaPriceInSDR.MulInt64(rand.Int63()+1).TruncateInt())
		retCoin, spread, err := input.MarketKeeper.GetSwapCoin(input.Ctx, offerCoin, core.MicroLunaDenom, false)
		require.NoError(t, err)
		require.Equal(t, input.MarketKeeper.MinSwapSpread(input.Ctx), spread)
		require.Equal(t, sdk.NewDecFromInt(offerCoin.Amount).Quo(lunaPriceInSDR).TruncateInt(), retCoin.Amount)

		retCoin, spread, err = input.MarketKeeper.GetSwapCoin(input.Ctx, offerCoin, core.MicroLunaDenom, true)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroDec(), spread)
		require.Equal(t, sdk.NewDecFromInt(offerCoin.Amount).Quo(lunaPriceInSDR).TruncateInt(), retCoin.Amount)
	}

	offerCoin := sdk.NewCoin(core.MicroSDRDenom, lunaPriceInSDR.QuoInt64(2).TruncateInt())
	_, _, err := input.MarketKeeper.GetSwapCoin(input.Ctx, offerCoin, core.MicroLunaDenom, false)
	require.Error(t, err)
}

func TestGetDecSwapCoin(t *testing.T) {
	input := CreateTestInput(t)
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	// zero day (min spread)
	for i := 0; i < 100; i++ {
		offerCoin := sdk.NewDecCoin(core.MicroSDRDenom, lunaPriceInSDR.MulInt64(rand.Int63()+1).TruncateInt())
		retCoin, err := input.MarketKeeper.GetSwapDecCoin(input.Ctx, offerCoin, core.MicroLunaDenom)
		require.NoError(t, err)
		require.Equal(t, offerCoin.Amount.Quo(lunaPriceInSDR), retCoin.Amount)
	}

	offerCoin := sdk.NewDecCoin(core.MicroSDRDenom, lunaPriceInSDR.QuoInt64(2).TruncateInt())
	_, err := input.MarketKeeper.GetSwapDecCoin(input.Ctx, offerCoin, core.MicroLunaDenom)
	require.Error(t, err)
}
