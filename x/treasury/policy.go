package treasury

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func updateTaxPolicy(ctx sdk.Context, k Keeper) (newTaxRate sdk.Dec) {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx, util.GetEpoch(ctx))
	inc := params.MiningIncrement
	tlYear := RollingAverageIndicator(ctx, k, params.EpochLong, TRL)
	tlMonth := RollingAverageIndicator(ctx, k, params.EpochShort, TRL)

	// No revenues, hike as much as possible.
	if tlMonth.Equal(sdk.ZeroDec()) {
		newTaxRate = params.TaxPolicy.RateMax
	} else {
		newTaxRate = oldTaxRate.Mul(tlYear.Add(inc.Amount)).Quo(tlMonth)
	}

	newTaxRate = params.TaxPolicy.Clamp(oldTaxRate, newTaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx, newTaxRate)
	return
}

// w(t+1) = w(t)*SMR_target/SMR_rolling(t)
func updateRewardPolicy(ctx sdk.Context, k Keeper) (newRewardWeight sdk.Dec) {
	params := k.GetParams(ctx)

	curEpoch := util.GetEpoch(ctx)
	oldWeight := k.GetRewardWeight(ctx, curEpoch)
	smrTarget := params.SeigniorageBurdenTarget
	smrAvgMonth := RollingAverageIndicator(ctx, k, params.EpochShort, SMR)

	// No revenues; hike as much as possible
	if smrAvgMonth.Equal(sdk.ZeroDec()) {
		newRewardWeight = params.RewardPolicy.RateMax
	} else {
		newRewardWeight = oldWeight.Mul(smrTarget.Quo(smrAvgMonth))
	}

	newRewardWeight = params.RewardPolicy.Clamp(oldWeight, newRewardWeight)

	// Set the new reward weight
	k.SetRewardWeight(ctx, newRewardWeight)
	return
}
