package market

import (
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandlerMsgSwapValidPrice(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	lnasdrRate := sdk.NewDec(4)
	lnacnyRate := sdk.NewDec(8)
	offerCoin := sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(2).MulRaw(assets.MicroUnit))
	askCoin := sdk.NewCoin(assets.MicroCNYDenom, sdk.NewInt(4).MulRaw(assets.MicroUnit))

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
	require.Equal(t, trader.GetCoins().AmountOf(offerCoin.Denom), uSDRAmt.Sub(offerCoin.Amount))
	require.Equal(t, trader.GetCoins().AmountOf(askCoin.Denom), retAmt)
}

func TestHandlerMsgSwapNoBalance(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	// Try to swap a coin I don't have at all
	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.MicroCNYDenom, sdk.OneInt().MulRaw(assets.MicroUnit)), assets.MicroGBPDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Try to swap a coin I don't have enough of
	msg.OfferCoin = sdk.NewCoin(assets.MicroSDRDenom, uSDRAmt.Add(sdk.OneInt().MulRaw(assets.MicroUnit)))
	res = handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapRecursion(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	msg := NewMsgSwap(addrs[0], sdk.NewCoin(assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit)), assets.MicroSDRDenom)
	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)
}

func TestHandlerMsgSwapTooSmall(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerCoin := sdk.NewCoin(assets.MicroSDRDenom, sdk.NewDecFromIntWithPrec(sdk.OneInt(), 4).MulInt64(assets.MicroUnit).TruncateInt())
	askDenom := assets.MicroCNYDenom
	askLunaPrice := sdk.NewDec(1)
	offerLunaPrice := sdk.NewDecWithPrec(1001, 1)

	// Set oracle price
	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, offerLunaPrice)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)

	msg := NewMsgSwap(addrs[0], offerCoin, askDenom)

	res := handler(input.ctx, msg)
	require.False(t, res.IsOK(), "expected failed message execution: %v", res.Log)

	// Reset oracle price
	askLunaPrice = sdk.NewDecWithPrec(1000, 1)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, askDenom, askLunaPrice)

	res = handler(input.ctx, msg)
	require.True(t, res.IsOK(), "expected successful message execution: %v", res.Log)
}

func TestHandlerExceedDailySwapLimit(t *testing.T) {
	input := createTestInput(t)
	handler := NewHandler(input.marketKeeper)

	offerCoin := sdk.NewInt64Coin(assets.MicroSDRDenom, 100)

	// Set oracle price
	offerLunaPrice := sdk.NewDec(1)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, offerCoin.Denom, offerLunaPrice)

	// Day 0 ... trade goes through, even though Luna doesn't even have a stated issuance.
	msg := NewMsgSwap(addrs[0], offerCoin, assets.MicroLunaDenom)
	res := handler(input.ctx, msg)
	require.True(t, res.IsOK())

	// Day 1+ ... Set luna issuance, try to oscillate within the limit, and things should be ok
	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerWeek)
	err := input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewInt64Coin(assets.MicroLunaDenom, 1000000))
	assert.Nil(t, err)
	msg = NewMsgSwap(addrs[0], offerCoin, assets.MicroLunaDenom)
	res = handler(input.ctx, msg)
	require.True(t, res.IsOK())

	// Day 1+ ... Outside of the limit fails
	msg = NewMsgSwap(addrs[0], sdk.NewInt64Coin(assets.MicroLunaDenom, 10005), assets.MicroLunaDenom)
	res = handler(input.ctx, msg)
	require.False(t, res.IsOK())

	// Swapping Terra with each other should be unlimited
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroCNYDenom, sdk.OneDec())
	msg = NewMsgSwap(addrs[1], sdk.NewCoin(assets.MicroSDRDenom, uSDRAmt), assets.MicroCNYDenom) // 1/3 of SDR issuance
	res = handler(input.ctx, msg)
	require.True(t, res.IsOK())
}
