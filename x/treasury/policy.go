package treasury

import (
	"github.com/terra-project/core/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func updateTaxPolicy(ctx sdk.Context, k Keeper) (newTaxRate sdk.Dec) {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx, util.GetEpoch(ctx))
	inc := params.MiningIncrement
	tlYear := RollingAverageIndicator(ctx, k, params.WindowLong, TRL)
	tlMonth := RollingAverageIndicator(ctx, k, params.WindowShort, TRL)

	// No revenues, hike as much as possible.
	if tlMonth.Equal(sdk.ZeroDec()) {
		newTaxRate = params.TaxPolicy.RateMax
	} else {
		newTaxRate = oldTaxRate.Mul(tlYear.Mul(inc)).Quo(tlMonth)
	}

	newTaxRate = params.TaxPolicy.Clamp(oldTaxRate, newTaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx.WithBlockHeight(ctx.BlockHeight()+1), newTaxRate)
	return
}

// w(t+1) = w(t)*SB_target/SB_rolling(t)
func updateRewardPolicy(ctx sdk.Context, k Keeper) (newRewardWeight sdk.Dec) {
	params := k.GetParams(ctx)

	curEpoch := util.GetEpoch(ctx)
	oldWeight := k.GetRewardWeight(ctx, curEpoch)
	sbTarget := params.SeigniorageBurdenTarget

	seigniorageSum := SumIndicator(ctx, k, params.WindowShort, SeigniorageRewardsForEpoch)
	totalSum := SumIndicator(ctx, k, params.WindowShort, MiningRewardForEpoch)

	// No revenues; hike as much as possible
	if totalSum.Equal(sdk.ZeroDec()) || seigniorageSum.Equal(sdk.ZeroDec()) {
		newRewardWeight = params.RewardPolicy.RateMax
	} else {
		// Seigniorage burden out of total rewards
		sb := seigniorageSum.Quo(totalSum)
		newRewardWeight = oldWeight.Mul(sbTarget.Quo(sb))
	}

	newRewardWeight = params.RewardPolicy.Clamp(oldWeight, newRewardWeight)

	// Set the new reward weight
	k.SetRewardWeight(ctx.WithBlockHeight(ctx.BlockHeight()+1), newRewardWeight)
	return
}
