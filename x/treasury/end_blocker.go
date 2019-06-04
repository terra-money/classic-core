package treasury

import (
	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// not enough data collected to update variables
func isProbationPeriod(ctx sdk.Context, k Keeper) bool {

	// Look 1 block into the future ... at the last block of the epoch, trigger
	futureCtx := ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	futureEpoch := util.GetEpoch(futureCtx)

	return futureEpoch.LT(k.GetParams(ctx).WindowProbation)
}

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	if !util.IsPeriodLastBlock(ctx, util.BlocksPerEpoch) {
		return resTags
	}

	if isProbationPeriod(ctx, k) {
		return resTags
	}

	// Update policy weights
	taxRate := updateTaxPolicy(ctx, k)
	rewardWeight := updateRewardPolicy(ctx, k)

	return sdk.NewTags(
		tags.Action, tags.ActionPolicyUpdate,
		tags.Tax, taxRate.String(),
		tags.MinerReward, rewardWeight.String(),
	),
}
