package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"
	oracletypes "github.com/terra-money/core/x/oracle/types"
)

func TestSwapMsg(t *testing.T) {
	input, msgServer := setup(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	params.MinStabilitySpread = sdk.ZeroDec()
	input.MarketKeeper.SetParams(input.Ctx, params)

	beforeTerraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	swapMsg := types.NewMsgSwap(Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err := msgServer.Swap(sdk.WrapSDKContext(input.Ctx), swapMsg)
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
	swapMsg = types.NewMsgSwap(Addrs[0], offerCoin, core.MicroLunaDenom)

	_, err = msgServer.Swap(sdk.WrapSDKContext(input.Ctx), swapMsg)
	require.Error(t, err)

	// valid zero tobin tax test
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, sdk.ZeroDec())
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroSDRDenom, sdk.ZeroDec())
	offerCoin = sdk.NewCoin(core.MicroSDRDenom, amt)
	swapMsg = types.NewMsgSwap(Addrs[0], offerCoin, core.MicroKRWDenom)
	_, err = msgServer.Swap(sdk.WrapSDKContext(input.Ctx), swapMsg)
	require.NoError(t, err)
}

func TestSwapSendMsg(t *testing.T) {
	input, msgServer := setup(t)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroSDRDenom)
	require.NoError(t, err)

	expectedAmt := retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()

	swapSendMsg := types.NewMsgSwapSend(Addrs[0], Addrs[1], offerCoin, core.MicroSDRDenom)
	_, err = msgServer.SwapSend(sdk.WrapSDKContext(input.Ctx), swapSendMsg)
	require.NoError(t, err)

	balance := input.BankKeeper.GetBalance(input.Ctx, Addrs[1], core.MicroSDRDenom)
	require.Equal(t, expectedAmt, balance.Amount)
}

func TestSpreadDistribution(t *testing.T) {
	input, msgServer := setup(t)

	amt := sdk.NewInt(10)
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, amt)
	swapDecCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroSDRDenom)
	require.NoError(t, err)

	feeDecCoin := sdk.NewDecCoinFromDec(swapDecCoin.Denom, spread.Mul(swapDecCoin.Amount))
	swapDecCoin.Amount = swapDecCoin.Amount.Sub(feeDecCoin.Amount)

	_, decimalCoin := swapDecCoin.TruncateDecimal()
	feeDecCoin = feeDecCoin.Add(decimalCoin) // add truncated decimalCoin to swapFee
	feeCoin, _ := feeDecCoin.TruncateDecimal()

	swapMsg := types.NewMsgSwap(Addrs[0], offerCoin, core.MicroSDRDenom)
	_, err = msgServer.Swap(sdk.WrapSDKContext(input.Ctx), swapMsg)
	require.NoError(t, err)

	blockValidationReward := feeCoin.Amount.QuoRaw(2)
	oracleVotingReward := feeCoin.Amount.Sub(blockValidationReward)

	balanceReq := banktypes.QueryBalanceRequest{
		Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
		Denom:   core.MicroSDRDenom,
	}
	balanceRes, err := input.BankKeeper.Balance(sdk.WrapSDKContext(input.Ctx), &balanceReq)
	require.NoError(t, err)
	require.Equal(t, balanceRes.Balance.Amount, blockValidationReward)

	balanceReq = banktypes.QueryBalanceRequest{
		Address: authtypes.NewModuleAddress(oracletypes.ModuleName).String(),
		Denom:   core.MicroSDRDenom,
	}
	balanceRes, err = input.BankKeeper.Balance(sdk.WrapSDKContext(input.Ctx), &balanceReq)
	require.NoError(t, err)
	require.Equal(t, balanceRes.Balance.Amount, oracleVotingReward)
}

var (
	uSDRAmt    = sdk.NewInt(1005 * core.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	randomPrice = sdk.NewDec(1700)
)

func setup(t *testing.T) (TestInput, types.MsgServer) {
	input := CreateTestInput(t)

	params := input.MarketKeeper.GetParams(input.Ctx)
	input.MarketKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, randomPrice)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, randomPrice)
	msgServer := NewMsgServerImpl(input.MarketKeeper)

	return input, msgServer
}
