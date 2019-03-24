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

<<<<<<< HEAD
	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, offerLunaPrice)

	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)
=======
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	retCoin, err := input.marketKeeper.SwapCoins(input.ctx, offerCoin, askDenom)
	require.Nil(t, err)
>>>>>>> 95f3dd212d1beea3dd2355026d70588ef939b46a

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

<<<<<<< HEAD
	// Set oracle price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Reset oracle price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)
	askLunaPrice = sdk.NewDecWithPrec(1000, 1)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)
=======
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	retCoin, err := input.marketKeeper.SwapDecCoins(input.ctx, offerCoin, askDenom)
	require.Nil(t, err)

	require.Equal(t, retCoin, expectedAskCoin)
>>>>>>> 95f3dd212d1beea3dd2355026d70588ef939b46a
}
