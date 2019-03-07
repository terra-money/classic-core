package treasury

import (
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	settlementPeriod := k.GetParams(ctx).EpochShort
	currentEpoch := util.GetEpoch(ctx)
	if currentEpoch.Mod(settlementPeriod).Equal(sdk.ZeroInt()) {
		resTags = k.settleClaims(ctx)

		tax := updateTaxes(ctx, k)
		k.pk.SetTaxRate(ctx, tax)

		updateTaxCaps(ctx, k)

		// Update taxes
		resTags = resTags.AppendTags(
			sdk.NewTags(
				tags.Action, tags.ActionTaxUpdate,
				tags.Tax, tax.String(),
			),
		)

		// Update mining rewards
		rewardWeight := updateRewardWeight(ctx, k)
		k.SetRewardWeight(ctx, rewardWeight)

		resTags = resTags.AppendTags(
			sdk.NewTags(
				tags.Action, tags.ActionMinerRewardUpdate,
				tags.MinerReward, rewardWeight.String(),
			),
		)
	}

	return
}

func updateTaxes(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	taxOld := k.pk.GetTaxRate(ctx)

	mrlLong := mrl(ctx, k, params.EpochLong)
	mrlShort := mrl(ctx, k, params.EpochShort)
	taxNew := taxOld.Mul(mrlLong).Quo(mrlShort)

	// Clamp within bounds
	if taxNew.GT(params.TaxRateMax) {
		taxNew = params.TaxRateMax
	} else if taxNew.LT(params.TaxRateMin) {
		taxNew = params.TaxRateMin
	}

	return taxNew
}

func updateTaxCaps(ctx sdk.Context, k Keeper) {
	taxProceeds := k.pk.PeekTaxProceeds(ctx, util.GetEpoch(ctx))
	taxCap := k.GetParams(ctx).TaxCap

	for _, coin := range taxProceeds {
		taxCapForDenom, err := k.mk.SwapCoins(ctx, taxCap, coin.Denom)

		// The coin is more valuable than TerraSDR. just set 1 as the cap.
		if err != nil {
			taxCapForDenom.Amount = sdk.OneInt()
		}

		k.pk.SetTaxCap(ctx, coin.Denom, taxCapForDenom.Amount)
	}
}

func updateRewardWeight(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	weightOld := k.GetRewardWeight(ctx)

	mrlLong := mrl(ctx, k, params.EpochLong)
	mrlShort := mrl(ctx, k, params.EpochShort)
	delta := sdk.OneDec().Sub(mrlShort.Quo(mrlLong))

	weightNew := weightOld.Add(delta)

	// Clamp within bounds
	if weightNew.GT(params.RewardMax) {
		weightNew = params.RewardMax
	} else if weightNew.LT(params.RewardMin) {
		weightNew = params.RewardMin
	}

	return weightNew
}

func translateFees(ctx sdk.Context, k Keeper) sdk.Coin {
	feeSum := sdk.NewCoin(assets.SDRDenom, sdk.ZeroInt())

	taxProceeds := k.pk.PeekTaxProceeds(ctx, util.GetEpoch(ctx))
	for _, proceed := range taxProceeds {
		translation, err := k.mk.SwapCoins(ctx, proceed, assets.SDRDenom)
		if err != nil {
			continue
		}

		feeSum.Amount = feeSum.Amount.Add(translation.Amount)
	}

	return feeSum
}

func mrl(ctx sdk.Context, k Keeper, epochs sdk.Int) (res sdk.Dec) {
	epoch := util.GetEpoch(ctx)
	sum := sdk.ZeroDec()
	for i := int64(0); i < epochs.Int64(); i++ {
		epoch := epoch.Sub(sdk.NewInt(int64(i)))

		if epoch.LT(sdk.ZeroInt()) {
			break
		}

		numLuna := k.pk.GetIssuance(ctx, assets.LunaDenom, epoch)
		taxProceeds := translateFees(ctx, k)
		marginalProceeds := sdk.NewDecFromInt(taxProceeds.Amount).QuoInt(numLuna)
		sum = sum.Add(marginalProceeds)
	}

	return sum.QuoInt(epochs)
}

// AddClaim adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) ProcessClaims(ctx sdk.Context, class ClaimClass, rewardees map[string]sdk.Int) {
	for rAddrStr, rewardWeight := range rewardees {
		addr, err := sdk.AccAddressFromBech32(rAddrStr)
		if err != nil {
			continue
		}

		newClaim := NewClaim(class, rewardWeight, addr)
		k.addClaim(ctx, newClaim)
	}
}

func (k Keeper) settleClaimsForClass(ctx sdk.Context, cReward sdk.DecCoins, cWeightSum sdk.Int, cClaims []Claim) (remainder sdk.DecCoins, classTags sdk.Tags) {
	store := ctx.KVStore(k.key)
	for _, claim := range cClaims {
		claimWeightInDec := sdk.NewDecFromInt(claim.weight)
		sumWeightInDec := sdk.NewDecFromInt(cWeightSum)
		claimReward := cReward.MulDec(claimWeightInDec).QuoDec(sumWeightInDec)

		// translate rewards to SDR
		rewardInSDR, err := k.mk.SwapDecCoins(
			ctx,
			sdk.NewDecCoinFromDec(assets.LunaDenom, claimReward.AmountOf(assets.LunaDenom)),
			assets.SDRDenom,
		)
		if err != nil {
			continue
		}
		rewardInSDRInt, dust := rewardInSDR.TruncateDecimal()

		// credit the recipient's account with the reward
		k.pk.AddCoins(ctx, claim.recipient, sdk.Coins{rewardInSDRInt})
		remainder = remainder.Plus(sdk.DecCoins{dust})

		// We are now done with the claim; remove it from the store
		store.Delete(KeyClaim(claim.ID()))

		classTags = classTags.AppendTags(
			sdk.NewTags(
				tags.Action, tags.ActionReward,
				tags.Rewardee, claim.recipient,
				tags.Amount, rewardInSDR,
				tags.Class, claim.class,
			),
		)
	}
	return
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	totalPool := k.dk.GetFeePool(ctx)

	// Pay mining rewards; just burn Luna
	minerRewardWeight := k.GetRewardWeight(ctx)
	minerRewards := totalPool.CommunityPool.MulDec(minerRewardWeight)

	// Compute the size of oracle + budget claim reward pools
	params := k.GetParams(ctx)

	oracleRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.OracleClaimShare)
	oracleReward := totalPool.CommunityPool.MulDec(oracleRewardWeight)
	budgetRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.BudgetClaimShare)
	budgetReward := totalPool.CommunityPool.MulDec(budgetRewardWeight)

	// Sum the total amount of voting power accumulated in claims by class
	oracleClaimWeightSum, oracleClaims := k.sumClaims(ctx, OracleClaimClass)
	budgetClaimWeightSum, budgetClaims := k.sumClaims(ctx, BudgetClaimClass)

	// Reward claims
	oracleRemainder, oracleTags := k.settleClaimsForClass(ctx, oracleReward, oracleClaimWeightSum, oracleClaims)
	budgetRemainder, budgetTags := k.settleClaimsForClass(ctx, budgetReward, budgetClaimWeightSum, budgetClaims)

	settleTags = settleTags.AppendTags(oracleTags)
	settleTags = settleTags.AppendTags(budgetTags)

	// Add the remainder back to the community pool
	totalPool.CommunityPool = oracleRemainder.Plus(budgetRemainder)
	k.dk.SetFeePool(ctx, totalPool)

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, minerRewards,
		tags.Oracle, oracleReward,
		tags.Budget, budgetReward,
	)
}
