package treasury

import (
	"terra/types"
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// not enough data collected to update variables
func isProbationPeriod(ctx sdk.Context, k Keeper) bool {
	return util.GetEpoch(ctx).LT(k.GetParams(ctx).EpochProbation)
}

// at the block height for a tally
func isEpochLastBlock(ctx sdk.Context, k Keeper) bool {
	settlementPeriod := k.GetParams(ctx).EpochShort

	// Look 1 block into the future ... at the last block of the epoch, trigger
	futureCtx := ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	futureEpoch := util.GetEpoch(futureCtx)
	return futureEpoch.GT(sdk.ZeroInt()) && // Skip the first epoch; need to build up history
		futureEpoch.Mod(settlementPeriod).Equal(sdk.ZeroInt())
}

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	if !isEpochLastBlock(ctx, k) {
		return resTags
	}

	// Settle and clear claims from the store
	resTags = k.settleClaims(ctx)

	if isProbationPeriod(ctx, k) {
		return resTags
	}

	// Update policy weights
	taxRate := updateTaxPolicy(ctx, k)
	rewardWeight := updateRewardPolicy(ctx, k)

	return resTags.AppendTags(
		sdk.NewTags(
			tags.Action, tags.ActionPolicyUpdate,
			tags.Tax, taxRate.String(),
			tags.MinerReward, rewardWeight.String(),
		),
	)
}

// compute scales by which the total reward pool must be
func getScales(ctx sdk.Context, k Keeper, oracleSum, budgetSum sdk.Int) (minerScale, oracleScale, budgetScale sdk.Dec) {
	params := k.GetParams(ctx)
	curEpoch := util.GetEpoch(ctx)

	oracleScale = sdk.ZeroDec()
	budgetScale = sdk.ZeroDec()
	rewardWeight := k.GetRewardWeight(ctx, curEpoch)
	if oracleSum.GT(sdk.ZeroInt()) {
		oracleScale = sdk.OneDec().Sub(rewardWeight).Mul(params.OracleClaimShare).QuoInt(oracleSum)
	}

	if budgetSum.GT(sdk.ZeroInt()) {
		budgetScale = sdk.OneDec().Sub(rewardWeight).Mul(params.BudgetClaimShare).QuoInt(budgetSum)
	}

	minerScale = sdk.OneDec().Sub(oracleScale).Sub(budgetScale)
	return
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	curEpoch := util.GetEpoch(ctx)
	store := ctx.KVStore(k.key)

	// Convert seigniorage to TerraSDR for rewards
	seigPool := k.mtk.PeekSeignioragePool(ctx, curEpoch)
	rewardPool, err := k.mk.SwapDecCoins(ctx, sdk.NewDecCoin(assets.LunaDenom, seigPool), assets.SDRDenom)
	if err != nil {
		return // No or too little seigniorage
	}

	oracleSumWeight := sdk.ZeroInt()
	budgetSumWeight := sdk.ZeroInt()

	// Sum weights by class
	k.IterateClaims(ctx, func(claim types.Claim) (stop bool) {
		switch claim.Class {
		case types.OracleClaimClass:
			oracleSumWeight = oracleSumWeight.Add(claim.Weight)
			break
		case types.BudgetClaimClass:
			budgetSumWeight = budgetSumWeight.Add(claim.Weight)
			break
		}
		return false
	})

	// Need to scale weights in claims by dividing class shares and total amount of weights
	minerScale, oracleScale, budgetScale := getScales(ctx, k, oracleSumWeight, budgetSumWeight)

	// Settle and delete all claims from the store
	k.IterateClaims(ctx, func(claim types.Claim) (stop bool) {
		var rewardAmt sdk.Int
		if claim.Class == types.OracleClaimClass {
			rewardAmt = rewardPool.Amount.Mul(oracleScale).MulInt(claim.Weight).TruncateInt()
		} else {
			rewardAmt = rewardPool.Amount.Mul(budgetScale).MulInt(claim.Weight).TruncateInt()
		}

		// Credit the recipient's account with the reward
		k.mtk.Mint(ctx, claim.Recipient, sdk.NewCoin(assets.SDRDenom, rewardAmt))

		// We are now done with the claim; remove it from the store
		store.Delete(keyClaim(claim.ID()))
		return false
	})

	// Just a rough approximation ... we are leaving some dust by rounding down each claim
	oracleRewards := rewardPool.Amount.Mul(oracleScale).MulInt(oracleSumWeight)
	budgetRewards := rewardPool.Amount.Mul(budgetScale).MulInt(budgetSumWeight)
	minerRewards := rewardPool.Amount.Mul(minerScale)

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, minerRewards.TruncateInt().String(),
		tags.Oracle, oracleRewards.TruncateInt().String(),
		tags.Budget, budgetRewards.TruncateInt().String(),
	)
}
