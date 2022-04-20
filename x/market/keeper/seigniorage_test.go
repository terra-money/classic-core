package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func checkBalance(t *testing.T, input TestInput, address sdk.AccAddress, expectedBalance sdk.Coins) {
	balances := input.BankKeeper.GetAllBalances(input.Ctx, address)
	fmt.Println("Expected:", balances, "Actual:", expectedBalance)
	require.Equal(t, expectedBalance, balances)
}

func checkCommunityPoolBalance(t *testing.T, input TestInput, expectedBalance sdk.Coins) {
	feePool := input.DistrKeeper.GetFeePool(input.Ctx)
	balances, _ := feePool.CommunityPool.TruncateDecimal()
	require.Equal(t, expectedBalance, balances)
}

func TestSettleSeigniorage_ZeroBalance(t *testing.T) {
	input := CreateTestInput(t)

	require.NotPanics(t, func() {
		input.MarketKeeper.SettleSeigniorage(input.Ctx)
	})
}

func TestSettleSeigniorage_ZeroSeigniorage(t *testing.T) {
	input := CreateTestInput(t)
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(
		sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(1000000)),
	)

	err := FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	require.NotPanics(t, func() {
		input.MarketKeeper.SettleSeigniorage(input.Ctx)
	})

	// Check whether the coins are burned
	checkBalance(t, input, moduleAddr, sdk.NewCoins())
}

func TestSettleSeigniorage_OnlyLuna(t *testing.T) {
	input := CreateTestInput(t)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(
		sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000000)),
	)

	err := FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	input.MarketKeeper.SettleSeigniorage(input.Ctx)

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

	err = FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	require.NotPanics(t, func() {
		input.MarketKeeper.SettleSeigniorage(input.Ctx)
	})

	checkBalance(t, input, moduleAddr, sdk.NewCoins())
	checkBalance(t, input, feeCollectorAddr, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100000))))
	checkCommunityPoolBalance(t, input, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(200000))))
}

func TestSettleSeigniorage_WithSeigniorage(t *testing.T) {
	input := CreateTestInput(t)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(
		sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(core.MicroKRWDenom, sdk.NewInt(1000000)),
	)

	err := FundAccount(input, moduleAddr, coins)
	require.NoError(t, err)

	// Set seigniorage routes
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	input.MarketKeeper.SetSeigniorageRoutes(input.Ctx, []types.SeigniorageRoute{
		{
			Address: types.AlternateCommunityPoolAddress.String(),
			Weight:  sdk.NewDecWithPrec(2, 1), // 20% to community pool
		},
		{
			Address: feeCollectorAddr.String(),
			Weight:  sdk.NewDecWithPrec(1, 1), // 10% to fee collector
		},
	})

	require.NotPanics(t, func() {
		input.MarketKeeper.SettleSeigniorage(input.Ctx)
	})

	// Check whether the coins are redirected
	// and left coins are burned
	checkBalance(t, input, moduleAddr, sdk.NewCoins())
	checkBalance(t, input, feeCollectorAddr, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(100000))))
	checkCommunityPoolBalance(t, input, sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(200000))))
}
