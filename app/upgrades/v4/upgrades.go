package v3

import (
	"github.com/classic-terra/core/app/keepers"
	"github.com/classic-terra/core/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV4UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// to run staking store migration
		stakingKeeper := keepers.StakingKeeper
		stkingMigrator := keeper.NewMigrator(stakingKeeper)
		stkingMigrator.Migrate13to14(ctx)
		// to run wasm store migration
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
