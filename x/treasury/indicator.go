package treasury

import (
	"terra/types/assets"
	"terra/types/util"

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
func TaxRewardsForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	taxRewards := sdk.NewDecCoins(k.PeekTaxProceeds(ctx, epoch))

	taxRewardInMicroSDR := sdk.ZeroDec()
	for _, coinReward := range taxRewards {
		if coinReward.Denom != assets.MicroSDRDenom {
			swappedReward, err := k.mk.SwapDecCoins(ctx, coinReward, assets.MicroSDRDenom)
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
func SeigniorageRewardsForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	seignioragePool := k.mtk.PeekSeignioragePool(ctx, epoch)
	rewardAmt := k.GetRewardWeight(ctx, epoch).MulInt(seignioragePool)
	seigniorageReward := sdk.NewDecCoinFromDec(assets.MicroLunaDenom, rewardAmt)

	microSDRReward, err := k.mk.SwapDecCoins(ctx, seigniorageReward, assets.MicroSDRDenom)
	if err != nil {
		return sdk.ZeroDec()
	}

	return microSDRReward.Amount
}

// MiningRewardForEpoch returns the sum of tax and seigniorage rewards for the epoch
func MiningRewardForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	taxRewards := TaxRewardsForEpoch(ctx, k, epoch)
	seigniorageRewards := SeigniorageRewardsForEpoch(ctx, k, epoch)

	return taxRewards.Add(seigniorageRewards)
}

// TRL returns tax rewards / luna / epoch
func TRL(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, TaxRewardsForEpoch)
}

// SRL returns Seigniorage rewards / luna / epoch
func SRL(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, SeigniorageRewardsForEpoch)
}

// MRL returns mining rewards / luna / epoch
func MRL(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	return UnitLunaIndicator(ctx, k, epoch, MiningRewardForEpoch)
}

// SMR returns the share of seigniorage rewards out of overall mining rewards for the epoch
func SMR(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	miningRewardAmount := MiningRewardForEpoch(ctx, k, epoch)
	seigniorageRewardAmount := SeigniorageRewardsForEpoch(ctx, k, epoch)

	if miningRewardAmount.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	return seigniorageRewardAmount.Quo(miningRewardAmount)
}

// UnitLunaIndicator evaluates the indicator function and divides it by the luna supply for the epoch
func UnitLunaIndicator(ctx sdk.Context, k Keeper, epoch sdk.Int,
	indicatorFunction func(sdk.Context, Keeper, sdk.Int) sdk.Dec) sdk.Dec {
	indicator := indicatorFunction(ctx, k, epoch)
	lunaTotalBondedAmount := k.valset.TotalBondedTokens(ctx)

	return indicator.QuoInt(lunaTotalBondedAmount)
}

// RollingAverageIndicator returns the rolling average of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return RollingAverageIndicator(currentEpoch)
func RollingAverageIndicator(ctx sdk.Context, k Keeper, epochs sdk.Int,
	indicatorFunction func(sdk.Context, Keeper, sdk.Int) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	var i sdk.Int
	curEpoch := util.GetEpoch(ctx)
	for i = curEpoch; i.GTE(sdk.ZeroInt()) && i.GT(curEpoch.Sub(epochs)); i = i.Sub(sdk.OneInt()) {
		val := indicatorFunction(ctx, k, i)
		sum = sum.Add(val)
	}

	computedEpochs := curEpoch.Sub(i)
	if computedEpochs.Equal(sdk.ZeroInt()) {
		return sum
	}

	return sum.QuoInt(computedEpochs)
}
