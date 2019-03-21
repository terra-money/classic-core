package market

import (
	"terra/types/assets"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandlerMsgSwapValidPrice(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerAmt := sdk.NewInt(2)
	offerCoin := sdk.NewCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(4)
	askLunaPrice := sdk.NewDec(8)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)

	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	swapAmount := offerLunaPrice.Quo(askLunaPrice).MulInt(offerAmt).TruncateInt()
	trader := input.accKeeper.GetAccount(input.ctx, addrs[0])
	require.Equal(t, trader.GetCoins().AmountOf(offerCoin.Denom), initAmt.Sub(offerAmt))
	require.Equal(t, trader.GetCoins().AmountOf(askDenom), swapAmount)
}

func TestHandlerMsgSwapNoBalance(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	// Try to swap a coin I don't have at all
	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.CNYDenom, sdk.OneInt()), assets.GBPDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Try to swap a coin I don't have enough of
	msg.OfferCoin = sdk.NewCoin(assets.SDRDenom, initAmt.Add(sdk.OneInt()))
	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapRecursion(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.SDRDenom, sdk.OneInt()), assets.SDRDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapTooSmall(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerAmt := sdk.NewInt(100)
	offerCoin := sdk.NewCoin(assets.SDRDenom, offerAmt)
	askDenom := assets.CNYDenom
	offerLunaPrice := sdk.NewDec(1)
	askLunaPrice := sdk.NewDecWithPrec(1001, 1)

	// Set oracle price
	input.oracleKeeper.SetPrice(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Reset oracle price
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)
	askLunaPrice = sdk.NewDecWithPrec(1000, 1)
	input.oracleKeeper.SetPrice(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)
}
