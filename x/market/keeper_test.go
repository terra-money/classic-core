package market

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperSwapCoins(t *testing.T) {

	input := createTestInput(t)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewCoin(assets.SDRDenom, sdk.NewInt(2))
	askCoin := sdk.NewCoin(assets.CNYDenom, sdk.NewInt(4))

	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	retCoin, err := input.marketKeeper.SwapCoins(input.ctx, offerCoin, askCoin.Denom)
	require.Nil(t, err)

	require.Equal(t, retCoin, askCoin)
}

func TestKeeperSwapDecCoins(t *testing.T) {
	input := createTestInput(t)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewDecCoin(assets.SDRDenom, sdk.NewInt(2))
	askCoin := sdk.NewDecCoin(assets.CNYDenom, sdk.NewInt(4))

	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	retCoin, err := input.marketKeeper.SwapDecCoins(input.ctx, offerCoin, askCoin.Denom)
	require.Nil(t, err)

	require.Equal(t, retCoin, askCoin)
}
