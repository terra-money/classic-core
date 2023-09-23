package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlockForks is intended to be ran in a chain upgrade.
func BeginBlockForks(ctx sdk.Context, app *TerraApp) {
	for _, fork := range Forks {
		if ctx.BlockHeight() == fork.UpgradeHeight {
			ctx.Logger().Info(fmt.Sprintf("applying fork %s at height %d", fork.UpgradeName, fork.UpgradeHeight))

			fork.BeginForkLogic(ctx, app.AppKeepers, app.mm)
			return
		}
	}
}
