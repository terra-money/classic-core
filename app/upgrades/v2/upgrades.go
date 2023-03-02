package v2

import (
	treasurytypes "github.com/classic-terra/core/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV2UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	upgradeKeeper upgradekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// fix empty module map before upgrading modules
		fromVM = mm.GetVersionMap()
		// set treasury module to be 1 so that it gets the new upgrade
		fromVM[treasurytypes.ModuleName] = 1
		upgradeKeeper.SetModuleVersionMap(ctx, fromVM)

		// treasury store migration
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
