package market

import (
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperSwapCoinsBasic(t *testing.T) {

	input := createTestInput(t)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit))
	askCoin := sdk.NewCoin(assets.MicroCNYDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))

	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	retCoin, spread, err := input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, askCoin.Denom)
	require.Nil(t, err)
	require.Zero(t, spread, "Spread should be 0 for non luna swaps")

	require.Equal(t, retCoin, askCoin)
}

func TestKeeperSwapCoinsLunaCap(t *testing.T) {

	input := createTestInput(t)

	// Set params
	params := DefaultParams()
	input.marketKeeper.SetParams(input.ctx, params)

	// Set day to 2 and issuance as the same as the day before
	input.mintKeeper.Mint(input.ctx, sdk.AccAddress{}, sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(10^9)))
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerDay + 1)

	lunaCoin := sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))

	// Set exchange rate. Keep it at 1:1 for simplicity
	lnasdrRate := sdk.NewDec(1)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, lnasdrRate)

	maxDelta := params.DailyLunaDeltaCap.MulInt(sdk.NewInt(10 ^ 9)).TruncateInt()

	// Check cap luna -> sdr swap, at the cap. Should succeed
	offerCoin := sdk.NewCoin(assets.MicroLunaDenom, maxDelta)
	retCoin, spread, err := input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, spread, params.MaxSwapSpread)

	// Check cap luna -> sdr swap, 1 coin higher than the cap. Should fail
	offerCoin = sdk.NewCoin(assets.MicroLunaDenom, maxDelta.Add(sdk.OneInt()))
	_, _, err = input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, assets.MicroSDRDenom)
	require.NotNil(t, err)

	// Check cap sdr -> luna swap,at the cap. Should succeed
	offerCoin = sdk.NewCoin(assets.MicroSDRDenom, maxDelta)
	_, spread, err = input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, assets.MicroLunaDenom)
	require.Nil(t, err)
	require.Equal(t, spread, params.MaxSwapSpread)

	// Check cap sdr -> luna swap, 1 coin higher than the cap. Should fail
	offerCoin = sdk.NewCoin(assets.MicroSDRDenom, maxDelta.Add(sdk.OneInt()))
	_, _, err = input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, assets.MicroLunaDenom)
	require.NotNil(t, err)
}

func TestKeeperSwapDecCoins(t *testing.T) {
	input := createTestInput(t)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewDecCoin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit))
	askCoin := sdk.NewDecCoin(assets.MicroCNYDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))

	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	retCoin, err := input.marketKeeper.GetSwapDecCoins(input.ctx, offerCoin, askCoin.Denom)
	require.Nil(t, err)

	require.Equal(t, retCoin, askCoin)
}
