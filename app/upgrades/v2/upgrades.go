package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV2UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		currentVersions := mm.GetVersionMap()
		// Skip capability upgrade by moving it to the latest version
		fromVM[capabilitytypes.ModuleName] = currentVersions[capabilitytypes.ModuleName]
		// treasury store migration
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
