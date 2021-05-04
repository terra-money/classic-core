package market

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/keeper"
	"github.com/terra-project/core/x/market/types"
)

func TestMarketFilters(t *testing.T) {
	input, h := setup(t)

	// Case 1: non-oracle message being sent fails
	bankMsg := banktypes.MsgSend{}
	_, err := h(input.Ctx, &bankMsg)
	require.Error(t, err)

	// Case 2: Normal MsgSwap submission goes through
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	prevoteMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)
}

func TestSwapMsg_Mint(t *testing.T) {
	input, h := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	beforePoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)

	amt := sdk.NewInt(10000)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	swapMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err := h(input.Ctx, swapMsg)
	require.NoError(t, err)

	afterPoolDelta := input.MarketKeeper.GetMintPoolDelta(input.Ctx)
	diff := beforePoolDelta.Sub(afterPoolDelta)

	// calculate estimation
	basePool := input.MarketKeeper.GetParams(input.Ctx).MintBasePool
	price, _ := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	cp := basePool.Mul(basePool)

	offerPool := basePool.Add(beforePoolDelta)
	askPool := cp.Quo(offerPool)
	estmiatedDiff := offerPool.Sub(cp.Quo(askPool.Add(price.MulInt(amt))))
	require.True(t, estmiatedDiff.Sub(diff.Abs()).LTE(sdk.NewDecWithPrec(1, 6)))

	// invalid recursive swap
	swapMsg = types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroLunaDenom)

	_, err = h(input.Ctx, swapMsg)
	require.Error(t, err)
}

func TestSwapMsg_Burn(t *testing.T) {
	input, h := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	beforePoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)

	amt := sdk.NewInt(10000)
	offerCoin := sdk.NewCoin(core.MicroSDRDenom, amt)
	err := input.BankKeeper.AddCoins(input.Ctx, keeper.Addrs[0], sdk.NewCoins(offerCoin))
	require.NoError(t, err)

	supply := input.BankKeeper.GetSupply(input.Ctx)
	supply.SetTotal(supply.GetTotal().Add(offerCoin))
	input.BankKeeper.SetSupply(input.Ctx, supply)

	swapMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroLunaDenom)
	_, err = h(input.Ctx, swapMsg)
	require.NoError(t, err)

	afterPoolDelta := input.MarketKeeper.GetBurnPoolDelta(input.Ctx)
	diff := beforePoolDelta.Sub(afterPoolDelta)

	// calculate estimation
	basePool := input.MarketKeeper.GetParams(input.Ctx).BurnBasePool
	// price, _ := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	cp := basePool.Mul(basePool)

	offerPool := basePool.Add(beforePoolDelta)
	askPool := cp.Quo(offerPool)
	estmiatedDiff := offerPool.Sub(cp.Quo(askPool.Add(sdk.NewDecFromInt(amt))))
	require.True(t, estmiatedDiff.Sub(diff.Abs()).LTE(sdk.NewDecWithPrec(1, 6)))
}

func TestSwapSendMsg(t *testing.T) {
	input, h := setup(t)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroSDRDenom)
	require.NoError(t, err)

	expectedAmt := retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()

	swapSendMsg := types.NewMsgSwapSend(keeper.Addrs[0], keeper.Addrs[1], offerCoin, core.MicroSDRDenom)
	_, err = h(input.Ctx, swapSendMsg)
	require.NoError(t, err)

	balance := input.BankKeeper.GetBalance(input.Ctx, keeper.Addrs[1], core.MicroSDRDenom)
	require.Equal(t, expectedAmt, balance.Amount)
}
