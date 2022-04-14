package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestSettleSeigniorage(t *testing.T) {
	input := CreateTestInput(t)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	coins := sdk.NewCoins(
		sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000000)),
		sdk.NewCoin(core.MicroUSDDenom, sdk.NewInt(1000000)),
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

	input.MarketKeeper.SettleSeigniorage(input.Ctx)

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
