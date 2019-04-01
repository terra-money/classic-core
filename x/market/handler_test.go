package market

import (
	"terra/types/assets"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandlerMsgSwapValidPrice(t *testing.T) {
	input := createTestInput()
	handler := NewHandler(input.marketKeeper)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewCoin(assets.SDRDenom, sdk.NewInt(2))
	askCoin := sdk.NewCoin(assets.CNYDenom, sdk.NewInt(4))

	msg := NewMsgSwap(addrs[0], offerCoin, askCoin.Denom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set offer asset price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, lnasdrRate)

	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Set ask asset price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askCoin.Denom, lnacnyRate)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)

	retAmt := lnacnyRate.Quo(lnasdrRate).MulInt(offerCoin.Amount).TruncateInt()
	trader := input.accKeeper.GetAccount(input.ctx, addrs[0])
	require.Equal(t, trader.GetCoins().AmountOf(offerCoin.Denom), initAmt.Sub(offerCoin.Amount))
	require.Equal(t, trader.GetCoins().AmountOf(askCoin.Denom), retAmt)
}

func TestHandlerMsgSwapNoBalance(t *testing.T) {
	input := createTestInput()
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
	input := createTestInput()
	handler := NewHandler(input.marketKeeper)

	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.SDRDenom, sdk.OneInt()), assets.SDRDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapTooSmall(t *testing.T) {
	input := createTestInput()
	handler := NewHandler(input.marketKeeper)

	offerCoin := sdk.NewCoin(assets.SDRDenom, sdk.NewInt(100))
	askDenom := assets.CNYDenom
	askLunaPrice := sdk.NewDec(1)
	offerLunaPrice := sdk.NewDecWithPrec(1001, 1)

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
}
