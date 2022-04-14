package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/terra-money/core/x/market/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
// Register following initial seigniorage routes
// - community pool 20%
// - fee collector 10%
// TODO - check final percentage before release
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	m.keeper.SetSeigniorageRoutes(ctx, []types.SeigniorageRoute{
		{
			Address: authtypes.NewModuleAddress(types.AlternateCommunityPoolAddress).String(),
			Weight:  sdk.NewDecWithPrec(2, 1),
		},
		{
			Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
			Weight:  sdk.NewDecWithPrec(1, 1),
		},
	})
	return nil
}
