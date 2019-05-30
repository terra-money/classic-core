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

	minerScale = sdk.OneDec().Sub(oracleScale.MulInt(oracleSum)).Sub(budgetScale.MulInt(budgetSum))
	return
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	curEpoch := util.GetEpoch(ctx)
	store := ctx.KVStore(k.key)

	// Convert seigniorage to TerraSDR for rewards
	seigPool := k.mtk.PeekSeignioragePool(ctx, curEpoch)
	rewardPool, err := k.mk.GetSwapDecCoins(ctx, sdk.NewDecCoin(assets.MicroLunaDenom, seigPool), assets.MicroSDRDenom)
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
		case types.BudgetClaimClass:
			budgetSumWeight = budgetSumWeight.Add(claim.Weight)
		}
		return false
	})

	// Need to scale weights in claims by dividing class shares and total amount of weights
	minerScale, oracleScale, budgetScale := getScales(ctx, k, oracleSumWeight, budgetSumWeight)

	// Settle and delete all claims from the store
	k.IterateClaims(ctx, func(claim types.Claim) (stop bool) {

		var err error
		if claim.Class == types.OracleClaimClass {
			rewardAmt := rewardPool.Amount.Mul(oracleScale).MulInt(claim.Weight).TruncateInt()
			rewardeeVallidator := k.valset.Validator(ctx, sdk.ValAddress(claim.Recipient))
			rewardCoins := sdk.NewCoins(sdk.NewCoin(assets.MicroSDRDenom, rewardAmt))

			// In case absence of the validator, we collect the rewards to fee collect keeper
			if rewardeeVallidator != nil {
				k.dk.AllocateTokensToValidator(ctx, rewardeeVallidator, sdk.NewDecCoins(rewardCoins))
			} else {
				k.fck.AddCollectedFees(ctx, rewardCoins)
			}

			// Minted amount is goes to validator outstanding pool or fee pool, so no repcipient is specified
			err = k.mtk.Mint(ctx, sdk.AccAddress{}, sdk.NewCoin(assets.MicroSDRDenom, rewardAmt))
		} else {
			rewardAmt := rewardPool.Amount.Mul(budgetScale).MulInt(claim.Weight).TruncateInt()

			// Credit the recipient's account with the reward
			err = k.mtk.Mint(ctx, claim.Recipient, sdk.NewCoin(assets.MicroSDRDenom, rewardAmt))
		}

		if err != nil {
			return false
		}

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
