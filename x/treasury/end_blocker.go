package treasury

import (
	"terra/types/assets"
	"terra/types/util"
	"terra/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called to adjust macro weights (tax, mining reward) and settle outstanding claims.
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {

	settlementPeriod := k.GetParams(ctx).SettlementPeriod
	if sdk.NewInt(ctx.BlockHeight()).Mod(settlementPeriod).Equal(sdk.ZeroInt()) {
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

// // ProcessOracleRewardees process rewardees from the oracle module
// func processClaims(ctx sdk.Context, k Keeper, claimants map[string]sdk.Int) {
// 	for claimantAddr, claimWeight := range claimants {
// 		addr, err := sdk.AccAddressFromBech32(claimantAddr)
// 		if err != nil {
// 			continue
// 		}

// 		oracleRewardClaim := NewClaim(OracleClaimClass, sdk.NewDecFromInt(rewardWeight), addr)
// 		k.AddClaim(ctx, oracleRewardClaim)
// 	}
// }

func updateTaxes(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	taxOld := k.pk.GetTaxRate(ctx)

	mrlLong := k.mrl(ctx, 52)
	mrlShort := k.mrl(ctx, 4)
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
	taxProceeds := k.pk.GetTaxProceeds(ctx, util.GetEpoch(ctx))
	taxCap := k.GetParams(ctx).TaxCap

	for _, coin := range taxProceeds {
		taxCapForDenom, err := k.mk.SwapCoins(ctx, taxCap, coin.Denom)
		if err != nil {
			continue
		}

		k.pk.SetTaxCap(ctx, coin.Denom, taxCapForDenom.Amount)
	}
}

func updateRewardWeight(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	weightOld := k.GetRewardWeight(ctx)

	mrlLong := k.mrl(ctx, 52)
	mrlShort := k.mrl(ctx, 4)
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

	taxProceeds := k.pk.GetTaxProceeds(ctx, util.GetEpoch(ctx))
	for _, proceed := range taxProceeds {
		translation, err := k.mk.SwapCoins(ctx, proceed, assets.SDRDenom)
		if err != nil {
			continue
		}

		feeSum.Amount = feeSum.Amount.Add(translation.Amount)
	}

	return feeSum
}

func (k Keeper) mrl(ctx sdk.Context, epochs int) (res sdk.Dec) {
	epoch := util.GetEpoch(ctx)
	sum := sdk.ZeroDec()
	for i := 0; i < epochs; i++ {
		epoch := epoch.Sub(sdk.NewInt(int64(i)))

		if epoch.LT(sdk.ZeroInt()) {
			break
		}

		numLuna := k.pk.GetIssuance(ctx, assets.LunaDenom, epoch)
		taxProceeds := translateFees(ctx, k)

		marginalProceeds := sdk.NewDecFromInt(taxProceeds.Amount).QuoInt(numLuna)

		sum = sum.Add(marginalProceeds)
	}

	return sum.QuoInt(sdk.NewInt(int64(epochs)))
}
