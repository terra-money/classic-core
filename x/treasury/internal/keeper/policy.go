package keeper

import (
	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateTaxCap updates all denom's tax cap
func (k Keeper) UpdateTaxCap(ctx sdk.Context) sdk.Coins {
	cap := k.GetParams(ctx).TaxPolicy.Cap
	total := k.supplyKeeper.GetSupply(ctx).GetTotal()

	var newCaps sdk.Coins
	for _, coin := range total {
		// ignore uluna tax cap (uluna has no tax); keep sdr tax cap
		if coin.Denom == core.MicroLunaDenom || coin.Denom == core.MicroSDRDenom {
			continue
		}

		newDecCap, err := k.marketKeeper.ComputeInternalSwap(ctx, sdk.NewDecCoinFromCoin(cap), coin.Denom)
		if err == nil {
			newCap, _ := newDecCap.TruncateDecimal()
			newCaps = append(newCaps, newCap)
			k.SetTaxCap(ctx, newCap.Denom, newCap.Amount)
		}
	}

	return newCaps
}

// t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func (k Keeper) UpdateTaxPolicy(ctx sdk.Context) (newTaxRate sdk.Dec) {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx, core.GetEpoch(ctx))
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
func (k Keeper) UpdateRewardPolicy(ctx sdk.Context) (newRewardWeight sdk.Dec) {
	params := k.GetParams(ctx)

	curEpoch := core.GetEpoch(ctx)
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
