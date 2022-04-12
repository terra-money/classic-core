package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/types"
)

func TestBurnCoinsFromBurnAccount(t *testing.T) {
	input := CreateTestInput(t)

	burnAddress := input.AccountKeeper.GetModuleAddress(types.BurnModuleName)
	coins := input.BankKeeper.GetAllBalances(input.Ctx, burnAddress)
	require.Equal(t, InitCoins, coins)

	input.MarketKeeper.BurnCoinsFromBurnAccount(input.Ctx)
	coins = input.BankKeeper.GetAllBalances(input.Ctx, burnAddress)
	require.True(t, coins.IsZero())
}
