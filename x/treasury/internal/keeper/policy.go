package keeper

import (
	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateTaxCap updates all denom's tax cap
func (k Keeper) UpdateTaxCap(ctx sdk.Context) sdk.Coins {
	cap := sdk.NewDecCoinFromCoin(k.GetParams(ctx).TaxPolicy.Cap)
	total := k.supplyKeeper.GetSupply(ctx).GetTotal()

	var newCaps sdk.Coins
	for _, coin := range total {
		// ignore uluna tax cap (uluna has no tax); keep sdr tax cap
		if coin.Denom == core.MicroLunaDenom || coin.Denom == core.MicroSDRDenom {
			continue
		}

		newDecCap, err := k.marketKeeper.ComputeInternalSwap(ctx, cap, coin.Denom)
		if err == nil {
			newCap, _ := newDecCap.TruncateDecimal()
			newCaps = append(newCaps, newCap)
			k.SetTaxCap(ctx, newCap.Denom, newCap.Amount)
		}
	}

	return newCaps
}

// UpdateTaxPolicy updates tax-rate with t(t+1) = t(t) * (TL_year(t) + INC) / TL_month(t)
func (k Keeper) UpdateTaxPolicy(ctx sdk.Context) (newTaxRate sdk.Dec) {
	params := k.GetParams(ctx)

	oldTaxRate := k.GetTaxRate(ctx)
	inc := params.MiningIncrement
	tlYear := k.rollingAverageIndicator(ctx, params.WindowLong, TRL)
	tlMonth := k.rollingAverageIndicator(ctx, params.WindowShort, TRL)

	// No revenues, hike as much as possible.
	if tlMonth.Equal(sdk.ZeroDec()) {
		newTaxRate = params.TaxPolicy.RateMax
	} else {
		newTaxRate = oldTaxRate.Mul(tlYear.Mul(inc)).Quo(tlMonth)
	}

	newTaxRate = params.TaxPolicy.Clamp(oldTaxRate, newTaxRate)

	// Set the new tax rate to the store
	k.SetTaxRate(ctx, newTaxRate)
	return
}

// UpdateRewardPolicy updates reward-weight with w(t+1) = w(t)*SB_target/SB_rolling(t)
func (k Keeper) UpdateRewardPolicy(ctx sdk.Context) (newRewardWeight sdk.Dec) {
	params := k.GetParams(ctx)

	oldWeight := k.GetRewardWeight(ctx)
	sbTarget := params.SeigniorageBurdenTarget

	seigniorageSum := k.sumIndicator(ctx, params.WindowShort, SR)
	totalSum := k.sumIndicator(ctx, params.WindowShort, MR)

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
	k.SetRewardWeight(ctx, newRewardWeight)
	return
}
