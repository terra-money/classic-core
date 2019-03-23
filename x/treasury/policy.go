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
	tlYear := RollingAverageIndicator(ctx, k, params.EpochLong, TRL)
	tlMonth := RollingAverageIndicator(ctx, k, params.EpochShort, TRL)

	newTaxRate := oldTaxRate.Mul(tlYear.Add(inc.Amount)).Quo(tlMonth)
	clampedTaxRate := params.TaxPolicy.Clamp(oldTaxRate, newTaxRate)

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
	smrAvgMonth := RollingAverageIndicator(ctx, k, params.EpochShort, SMR)

	newWeight := prevWeight.Mul(smrTarget.Quo(smrAvgMonth))
	clampedWeight := params.RewardPolicy.Clamp(prevWeight, newWeight)

	// Set the new reward weight
	k.SetRewardWeight(ctx, clampedWeight)

	return clampedWeight
}
