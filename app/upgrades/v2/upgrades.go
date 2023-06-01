package v2

import (
	"github.com/classic-terra/core/v2/app/keepers"
	"github.com/classic-terra/core/v2/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	feesharetypes "github.com/classic-terra/core/v2/x/feeshare/types"
)

func CreateV2UpgradeHandler(
	mm *module.Manager,
	cfg module.Configurator,
	_ upgrades.BaseAppParamManager,
	appKeepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// treasury store migration

		// set new FeeShare params
		newFeeShareParams := feesharetypes.Params{
			EnableFeeShare:  true,
			DeveloperShares: sdk.NewDecWithPrec(50, 2), // = 50%
			AllowedDenoms:   []string{"uluna"},
		}
		appKeepers.FeeShareKeeper.SetParams(ctx, newFeeShareParams)

		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
