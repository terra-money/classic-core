package market

import (
	"testing"

	"github.com/terra-project/core/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperSwapCoins(t *testing.T) {

	input := createTestInput(t)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit))
	askCoin := sdk.NewCoin(assets.MicroCNYDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))

	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	retCoin, err := input.marketKeeper.GetSwapCoins(input.ctx, offerCoin, askCoin.Denom)
	require.Nil(t, err)

	require.Equal(t, retCoin, askCoin)
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
