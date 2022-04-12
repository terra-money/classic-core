package market

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestBurnAddress(t *testing.T) {
	input := keeper.CreateTestInput(t)

	burnAddress := input.AccountKeeper.GetModuleAddress(types.BurnModuleName)
	require.Equal(t, keeper.InitCoins, input.BankKeeper.GetAllBalances(input.Ctx, burnAddress))

	EndBlocker(input.Ctx, input.MarketKeeper)
	require.True(t, input.BankKeeper.GetAllBalances(input.Ctx, burnAddress).IsZero())
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
