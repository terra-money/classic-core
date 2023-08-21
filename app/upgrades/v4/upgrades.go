package v4

import (
	"github.com/classic-terra/core/v2/app/keepers"
	"github.com/classic-terra/core/v2/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV4UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	_ *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Migrate13to14 migrates from version v0.45.13 to v0.45.14.
		// Only for this particular version, which do not use the version of module.
		// stakingMigrator.Migrate13to14(ctx)

		// to run wasm store migration
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
