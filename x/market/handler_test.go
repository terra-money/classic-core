package market

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"
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

func TestSwapMsg_FailZeroReturn(t *testing.T) {
	input, h := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.OneDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	input.Ctx = input.Ctx.WithChainID(core.ColumbusChainID).WithBlockHeight(core.SwapDisableForkHeight)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	swapMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err := h(input.Ctx, swapMsg)
	require.Error(t, err)
}

func TestSwapMsg(t *testing.T) {
	input, h := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	beforeTerraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	swapMsg := types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err := h(input.Ctx, swapMsg)
	require.NoError(t, err)

	afterTerraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
	diff := beforeTerraPoolDelta.Sub(afterTerraPoolDelta)

	// calculate estimation
	basePool := input.MarketKeeper.GetParams(input.Ctx).BasePool
	price, _ := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	cp := basePool.Mul(basePool)

	terraPool := basePool.Add(beforeTerraPoolDelta)
	lunaPool := cp.Quo(terraPool)
	estmiatedDiff := terraPool.Sub(cp.Quo(lunaPool.Add(price.MulInt(amt))))
	require.True(t, estmiatedDiff.Sub(diff.Abs()).LTE(sdk.NewDecWithPrec(1, 6)))

	// invalid recursive swap
	swapMsg = types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroLunaDenom)

	_, err = h(input.Ctx, swapMsg)
	require.Error(t, err)

	// valid zero tobin tax test
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, sdk.ZeroDec())
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroSDRDenom, sdk.ZeroDec())
	offerCoin = sdk.NewCoin(core.MicroSDRDenom, amt)
	swapMsg = types.NewMsgSwap(keeper.Addrs[0], offerCoin, core.MicroKRWDenom)
	_, err = h(input.Ctx, swapMsg)
	require.NoError(t, err)
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
