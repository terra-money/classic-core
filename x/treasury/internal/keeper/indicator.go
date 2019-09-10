package keeper

import (
	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Computes important economic indicators for the stability of Terra currencies.
//
// Three important concepts:
// - MR: Fees + Seigniorage for a given epoch sums to Mining Rewards
// - MRL: Computes the Mining Reward per unit Luna
// - SMR: Computes the ratio of seigniorage rewards to overall mining rewards
//
// Rolling averages are also computed for MRL and SMR respectively.
//

// TaxRewardsForEpoch returns tax rewards that have been collected in the epoch
func TaxRewardsForEpoch(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	taxRewards := sdk.NewDecCoins(k.PeekTaxProceeds(ctx, epoch))

	taxRewardInMicroSDR := sdk.ZeroDec()
	for _, coinReward := range taxRewards {
		if coinReward.Denom != core.MicroSDRDenom {
			swappedReward, err := k.marketKeeper.ComputeInternalSwap(ctx, coinReward, core.MicroSDRDenom)
			if err != nil {
				continue
			}
			taxRewardInMicroSDR = taxRewardInMicroSDR.Add(swappedReward.Amount)
		} else {
			taxRewardInMicroSDR = taxRewardInMicroSDR.Add(coinReward.Amount)
		}
	}

	return taxRewardInMicroSDR
}

// SeigniorageRewardsForEpoch returns seigniorage rewards for the epoch
func SeigniorageRewardsForEpoch(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	seignioragePool := k.PeekEpochSeigniorage(ctx, epoch)
	rewardAmt := k.GetRewardWeight(ctx, epoch).MulInt(seignioragePool)
	seigniorageReward := sdk.NewDecCoinFromDec(core.MicroLunaDenom, rewardAmt)

	microSDRReward, err := k.marketKeeper.ComputeInternalSwap(ctx, seigniorageReward, core.MicroSDRDenom)
	if err != nil {
		return sdk.ZeroDec()
	}

	return microSDRReward.Amount
}

// MiningRewardForEpoch returns the sum of tax and seigniorage rewards for the epoch
func MiningRewardForEpoch(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	taxRewards := TaxRewardsForEpoch(ctx, k, epoch)
	seigniorageRewards := SeigniorageRewardsForEpoch(ctx, k, epoch)

	return taxRewards.Add(seigniorageRewards)
}

// TRL returns tax rewards / luna / epoch
func TRL(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, TaxRewardsForEpoch)
}

// SRL returns Seigniorage rewards / luna / epoch
func SRL(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, SeigniorageRewardsForEpoch)
}

// MRL returns mining rewards / luna / epoch
func MRL(ctx sdk.Context, k Keeper, epoch int64) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, MiningRewardForEpoch)
}

// UnitLunaIndicator evaluates the indicator function and divides it by the luna supply for the epoch
func UnitLunaIndicator(ctx sdk.Context, k Keeper, epoch int64,
	indicatorFunction func(sdk.Context, Keeper, int64) sdk.Dec) sdk.Dec {
	indicator := indicatorFunction(ctx, k, epoch)
	lunaTotalBondedAmount := k.stakingKeeper.TotalBondedTokens(ctx)

	return indicator.QuoInt(lunaTotalBondedAmount)
}

// SumIndicator returns the sum of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return SumIndicator(currentEpoch)
func SumIndicator(ctx sdk.Context, k Keeper, epochs int64,
	indicatorFunction func(sdk.Context, Keeper, int64) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := core.GetEpoch(ctx)

	for i := curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		val := indicatorFunction(ctx, k, i)
		sum = sum.Add(val)
	}

	return sum
}

// RollingAverageIndicator returns the rolling average of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return RollingAverageIndicator(currentEpoch)
func RollingAverageIndicator(ctx sdk.Context, k Keeper, epochs int64,
	indicatorFunction func(sdk.Context, Keeper, int64) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := core.GetEpoch(ctx)

	var i int64
	for i = curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		val := indicatorFunction(ctx, k, i)
		sum = sum.Add(val)
	}

	computedEpochs := curEpoch - i
	if computedEpochs == 0 {
		return sum
	}

	return sum.QuoInt64(computedEpochs)
}
