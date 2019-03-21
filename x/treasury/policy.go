package treasury

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func updateTaxPolicy(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx)
	inc := params.MiningIncrement
	tlYear := rollingAverageIndicator(ctx, k, params.EpochLong, trl)
	tlMonth := rollingAverageIndicator(ctx, k, params.EpochShort, trl)

	newTaxRate := oldTaxRate.Mul(tlYear.Add(inc.Amount)).Quo(tlMonth)
	clampedTaxRate := params.TaxPolicy.clamp(oldTaxRate, newTaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx, clampedTaxRate)

	return newTaxRate
}

// w(t+1) = w(t)*SMR_target/SMR_rolling(t)
func updateRewardPolicy(ctx sdk.Context, k Keeper) sdk.Dec {
	params := k.GetParams(ctx)

	curEpoch := util.GetEpoch(ctx)
	prevWeight := k.GetRewardWeight(ctx, curEpoch.Sub(sdk.OneInt()))
	smrTarget := params.SeigniorageBurdenTarget
	smrAvgMonth := rollingAverageIndicator(ctx, k, params.EpochShort, smr)

	newWeight := prevWeight.Mul(smrTarget.Quo(smrAvgMonth))
	clampedWeight := params.RewardPolicy.clamp(prevWeight, newWeight)

	// Set the new reward weight
	k.SetRewardWeight(ctx, clampedWeight)

	return clampedWeight
}
