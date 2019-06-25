package update

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/terra-project/core/update/plan"
	"github.com/terra-project/core/x/market"
	"github.com/terra-project/core/x/oracle"
)

// EndBlocker
func EndBlocker(
	ctx sdk.Context,
	accountKeeper auth.AccountKeeper,
	oracleKeeper oracle.Keeper,
	marketKeeper market.Keeper) (resTags sdk.Tags) {

	if ctx.ChainID() == "columbus-2" {
		updated := plan.Update230000(ctx, accountKeeper, oracleKeeper, marketKeeper)
		if updated {
			resTags.AppendTag(plan.TagUpdate230000, "updated")
		}
	}

	return
}
