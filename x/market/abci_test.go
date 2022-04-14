package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestBurnAddress(t *testing.T) {
	input := keeper.CreateTestInput(t)

	burnAddress := input.AccountKeeper.GetModuleAddress(types.BurnModuleName)
	require.Equal(t, keeper.InitCoins, input.BankKeeper.GetAllBalances(input.Ctx, burnAddress))

	EndBlocker(input.Ctx, input.MarketKeeper)
	require.True(t, input.BankKeeper.GetAllBalances(input.Ctx, burnAddress).IsZero())
}

func TestSettleSeigniorage(t *testing.T) {
	input := keeper.CreateTestInput(t)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(
		sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(1000000)),
	)

	err := keeper.FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	EndBlocker(input.Ctx, input.MarketKeeper)

	balances := input.BankKeeper.GetAllBalances(input.Ctx, moduleAddr)
	require.Empty(t, balances)

	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	input.MarketKeeper.SetSeigniorageRoutes(input.Ctx, []types.SeigniorageRoute{
		{
			Address: types.AlternateCommunityPoolAddress.String(),
			Weight:  sdk.NewDecWithPrec(2, 1),
		},
		{
			Address: feeCollectorAddr.String(),
			Weight:  sdk.NewDecWithPrec(1, 1),
		},
	})

	err = keeper.FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	EndBlocker(input.Ctx, input.MarketKeeper)

	balances = input.BankKeeper.GetAllBalances(input.Ctx, moduleAddr)
	require.Empty(t, balances)

	// only Luna will be distributed
	coins = sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, coins.AmountOf(core.MicroLunaDenom)))

	feePool := input.DistrKeeper.GetFeePool(input.Ctx)
	balances, _ = feePool.CommunityPool.TruncateDecimal()
	receivedCoins, _ := sdk.NewDecCoinsFromCoins(coins...).MulDec(sdk.NewDecWithPrec(2, 1)).TruncateDecimal()
	require.Equal(t, balances, receivedCoins)

	balances = input.BankKeeper.GetAllBalances(input.Ctx, feeCollectorAddr)
	receivedCoins, _ = sdk.NewDecCoinsFromCoins(coins...).MulDec(sdk.NewDecWithPrec(1, 1)).TruncateDecimal()
	require.Equal(t, balances, receivedCoins)
}

func TestReplenishPools(t *testing.T) {
	input := keeper.CreateTestInput(t)

	terraDelta := sdk.NewDecWithPrec(17987573223725367, 3)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, terraDelta)

	for i := 0; i < 100; i++ {
		terraDelta = input.MarketKeeper.GetTerraPoolDelta(input.Ctx)

		poolRecoveryPeriod := int64(input.MarketKeeper.PoolRecoveryPeriod(input.Ctx))
		terraRegressionAmt := terraDelta.QuoInt64(poolRecoveryPeriod)

		EndBlocker(input.Ctx, input.MarketKeeper)

		terraPoolDelta := input.MarketKeeper.GetTerraPoolDelta(input.Ctx)
		require.Equal(t, terraDelta.Sub(terraRegressionAmt), terraPoolDelta)
	}
}
