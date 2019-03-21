package treasury

import (
	"terra/types"
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
		k.SetTaxRate(ctx, tax)

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

	taxOld := k.GetTaxRate(ctx)

	mrlLong := mrl(ctx, k, params.EpochLong)
	mrlShort := mrl(ctx, k, params.EpochShort)

	if mrlShort.IsZero() {
		mrlShort = sdk.OneDec()
	}

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
	taxProceeds := k.PeekTaxProceeds(ctx, util.GetEpoch(ctx))
	taxCap := k.GetParams(ctx).TaxCap

	for _, coin := range taxProceeds {
		taxCapForDenom, err := k.mk.SwapCoins(ctx, taxCap, coin.Denom)

		// The coin is more valuable than TerraSDR. just set 1 as the cap.
		if err != nil {
			taxCapForDenom.Amount = sdk.OneInt()
		}

		k.SetTaxCap(ctx, coin.Denom, taxCapForDenom.Amount)
	}
}

func updateRewardWeight(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	weightOld := k.GetRewardWeight(ctx)

	mrlLong := mrl(ctx, k, params.EpochLong)
	mrlShort := mrl(ctx, k, params.EpochShort)

	if mrlLong.IsZero() {
		mrlLong = sdk.OneDec()
	}

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

	taxProceeds := k.PeekTaxProceeds(ctx, util.GetEpoch(ctx))
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

		numLuna := k.mtk.GetIssuance(ctx, assets.LunaDenom, epoch)
		taxProceeds := translateFees(ctx, k)
		marginalProceeds := sdk.NewDecFromInt(taxProceeds.Amount).QuoInt(numLuna)
		sum = sum.Add(marginalProceeds)
	}

	return sum.QuoInt(epochs)
}

// AddClaim adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) ProcessClaims(ctx sdk.Context, claims []types.Claim) {
	for _, claim := range claims {
		k.addClaim(ctx, claim)
	}
}

func (k Keeper) settleClaimsForClass(ctx sdk.Context, cReward sdk.DecCoins, cWeightSum sdk.Int, cClaims []types.Claim) (classTags sdk.Tags) {
	store := ctx.KVStore(k.key)
	for _, claim := range cClaims {
		claimWeightInDec := sdk.NewDecFromInt(claim.Weight)
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
		rewardInSDRInt, _ := rewardInSDR.TruncateDecimal()

		// credit the recipient's account with the reward
		k.mtk.Mint(ctx, claim.Recipient, rewardInSDRInt)

		// We are now done with the claim; remove it from the store
		store.Delete(KeyClaim(claim.ID()))

		classTags = classTags.AppendTags(
			sdk.NewTags(
				tags.Action, tags.ActionReward,
				tags.Rewardee, claim.Recipient,
				tags.Amount, rewardInSDR,
				tags.Class, claim.Class,
			),
		)
	}
	return
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	seignioragePool := k.mtk.ClaimSeigniorage(ctx)
	decPool := sdk.DecCoins{sdk.NewDecCoin(assets.LunaDenom, seignioragePool)}

	// Pay mining rewards; just burn Luna
	minerRewardWeight := k.GetRewardWeight(ctx)
	minerRewards := decPool.MulDec(minerRewardWeight)

	// Compute the size of oracle + budget claim reward pools
	params := k.GetParams(ctx)

	oracleRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.OracleClaimShare)
	oracleReward := decPool.MulDec(oracleRewardWeight)
	budgetRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.BudgetClaimShare)
	budgetReward := decPool.MulDec(budgetRewardWeight)

	// Sum the total amount of voting power accumulated in claims by class
	oracleClaimWeightSum, oracleClaims := k.sumClaims(ctx, types.OracleClaimClass)
	budgetClaimWeightSum, budgetClaims := k.sumClaims(ctx, types.BudgetClaimClass)

	// Reward claims
	oracleTags := k.settleClaimsForClass(ctx, oracleReward, oracleClaimWeightSum, oracleClaims)
	settleTags = settleTags.AppendTags(oracleTags)
	budgetTags := k.settleClaimsForClass(ctx, budgetReward, budgetClaimWeightSum, budgetClaims)
	settleTags = settleTags.AppendTags(budgetTags)

	// Burn luna
	k.mtk.ChangeIssuance(ctx, assets.LunaDenom, seignioragePool.Neg())

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, minerRewards.String(),
		tags.Oracle, oracleReward.String(),
		tags.Budget, budgetReward.String(),
	)
}
