package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/terra-money/core/x/market/types"
)

func TestMigrate1to2(t *testing.T) {
	input := CreateTestInput(t)

	routes := input.MarketKeeper.GetSeigniorageRoutes(input.Ctx)
	require.Empty(t, routes)

	migrator := NewMigrator(input.MarketKeeper)
	migrator.Migrate1to2(input.Ctx)

	routes = input.MarketKeeper.GetSeigniorageRoutes(input.Ctx)
	require.Equal(t, []types.SeigniorageRoute{
		{
			Address: types.AlternateCommunityPoolAddress.String(),
			Weight:  sdk.NewDecWithPrec(2, 1),
		},
		{
			Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
			Weight:  sdk.NewDecWithPrec(1, 1),
		},
	}, routes)
}
