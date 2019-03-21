package market

import (
	"terra/types/assets"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeperSwapCoins(t *testing.T) {

	input := createTestInput(t)

	offerAmt := sdk.NewInt(2)
	offerCoin := sdk.NewCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(4)
	askLunaPrice := sdk.NewDec(8)
	expectedAskCoin := sdk.NewCoin(askDenom, sdk.NewInt(1))

	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	retCoin, err := input.marketKeeper.SwapCoins(input.ctx, offerCoin, askDenom)
	require.Nil(t, err)

	require.Equal(t, retCoin, expectedAskCoin)
}

func TestKeeperSwapDecCoins(t *testing.T) {
	input := createTestInput(t)

	offerAmt := sdk.NewInt(2)
	offerCoin := sdk.NewDecCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(4)
	askLunaPrice := sdk.NewDec(8)
	expectedAskCoin := sdk.NewDecCoin(askDenom, sdk.NewInt(1))

	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	retCoin, err := input.marketKeeper.SwapDecCoins(input.ctx, offerCoin, askDenom)
	require.Nil(t, err)

	require.Equal(t, retCoin, expectedAskCoin)
}
