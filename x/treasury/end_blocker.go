package treasury

import (
	"terra/types"
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// at the block height for a tally
func isAtEpochEnd(ctx sdk.Context, k Keeper) bool {
	settlementPeriod := k.GetParams(ctx).EpochShort

	// Look 1 block into the future ... at the last block of the epoch
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	currentEpoch := util.GetEpoch(ctx)
	return currentEpoch.Mod(settlementPeriod).Equal(sdk.ZeroInt())
}

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	if !isAtEpochEnd(ctx, k) {
		return resTags
	}

	// Settle and clear claims from the store
	resTags = k.settleClaims(ctx)

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

// ProcessClaims adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) ProcessClaims(ctx sdk.Context, claims []types.Claim) {
	for _, claim := range claims {
		k.addClaim(ctx, claim)
	}
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	curEpoch := util.GetEpoch(ctx)
	store := ctx.KVStore(k.key)

	// Convert seigniorage to TerraSDR for rewards
	seigPool := k.mtk.PeekSeignioragePool(ctx, curEpoch)
	rewardPool, err := k.mk.SwapDecCoins(ctx, sdk.NewDecCoin(assets.LunaDenom, seigPool), assets.SDRDenom)
	if err != nil {
		// Bad practice, but if Luna assets can't be converted to SDR, there is something
		// seriously wrong...
		panic(nil)
	}

	// Decrement Luna issuance
	k.mtk.ChangeIssuance(ctx, assets.LunaDenom, seigPool.Neg())

	oracleSumWeight := sdk.ZeroInt()
	budgetSumWeight := sdk.ZeroInt()

	// Sum all weights by class
	k.iterateClaims(ctx, func(claim types.Claim) (stop bool) {
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

	params := k.GetParams(ctx)

	minerScale := k.GetRewardWeight(ctx, curEpoch)
	oracleScale := sdk.OneDec().Sub(minerScale).Mul(params.OracleClaimShare).QuoInt(oracleSumWeight)
	budgetScale := sdk.OneDec().Sub(minerScale).Mul(params.BudgetClaimShare).QuoInt(budgetSumWeight)

	// Settle and delete all claims from the store
	k.iterateClaims(ctx, func(claim types.Claim) (stop bool) {
		var rewardAmt sdk.Int
		if claim.Class == types.OracleClaimClass {
			rewardAmt = rewardPool.Amount.Mul(oracleScale).TruncateInt()
		} else {
			rewardAmt = rewardPool.Amount.Mul(budgetScale).TruncateInt()
		}

		// Credit the recipient's account with the reward
		k.mtk.Mint(ctx, claim.Recipient, sdk.NewCoin(assets.SDRDenom, rewardAmt))

		// We are now done with the claim; remove it from the store
		store.Delete(keyClaim(claim.ID()))
		return false
	})

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, rewardPool.Amount.Mul(minerScale),
		tags.Oracle, rewardPool.Amount.Mul(sdk.OneDec().Sub(minerScale).Mul(params.OracleClaimShare)).String(),
		tags.Budget, rewardPool.Amount.Mul(sdk.OneDec().Sub(minerScale).Mul(params.BudgetClaimShare)).String(),
	)
}
