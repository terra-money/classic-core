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

// returns tax rewards that have been collected in the epoch
func taxRewardsForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	taxRewards := sdk.NewDecCoins(k.PeekTaxProceeds(ctx, epoch))

	taxRewardInSDR := sdk.ZeroDec()
	for _, coinReward := range taxRewards {
		coinRewardInSDR, err := k.mk.SwapDecCoins(ctx, coinReward, assets.SDRDenom)
		if err != nil {
			continue
		}

		taxRewardInSDR = taxRewardInSDR.Add(coinRewardInSDR.Amount)
	}

	return taxRewardInSDR
}

func seigniorageRewardsForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	seignioragePool := k.mtk.PeekSeignioragePool(ctx, epoch)
	seigniorageReward := sdk.NewDecCoinFromDec(assets.LunaDenom, k.GetRewardWeight(ctx, epoch).MulInt(seignioragePool))

	sdrReward, err := k.mk.SwapDecCoins(ctx, seigniorageReward, assets.SDRDenom)
	if err != nil {
		return sdk.ZeroDec()
	}

	return sdrReward.Amount
}

func miningRewardForEpoch(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	taxRewards := taxRewardsForEpoch(ctx, k, epoch)
	seigniorageRewards := seigniorageRewardsForEpoch(ctx, k, epoch)

	return taxRewards.Add(seigniorageRewards)
}

func trl(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	return unitLunaIndicator(ctx, k, epoch, taxRewardsForEpoch)
}

func mrl(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	return unitLunaIndicator(ctx, k, epoch, miningRewardForEpoch)
}

func smr(ctx sdk.Context, k Keeper, epoch sdk.Int) sdk.Dec {
	taxRewardAmount := taxRewardsForEpoch(ctx, k, epoch)
	seigniorageRewardAmount := seigniorageRewardsForEpoch(ctx, k, epoch)

	return seigniorageRewardAmount.Quo(taxRewardAmount.Add(seigniorageRewardAmount))
}

func unitLunaIndicator(ctx sdk.Context, k Keeper, epoch sdk.Int, indicatorFunction func(sdk.Context, Keeper, sdk.Int) sdk.Dec) sdk.Dec {
	indicator := indicatorFunction(ctx, k, epoch)
	lunaIssuance := k.mtk.GetIssuance(ctx, assets.LunaDenom, epoch)

	return indicator.QuoInt(lunaIssuance)
}

func rollingAverageIndicator(ctx sdk.Context, k Keeper, epochs sdk.Int,
	indicatorFunction func(sdk.Context, Keeper, sdk.Int) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	i := sdk.ZeroInt()
	curEpoch := util.GetEpoch(ctx)
	for i = curEpoch; i.GTE(sdk.ZeroInt()) && i.GT(curEpoch.Sub(epochs)); i = i.Sub(sdk.OneInt()) {
		val := indicatorFunction(ctx, k, i)
		sum = sum.Add(val)
	}

	if i.Equal(sdk.ZeroInt()) {
		return sum
	}

	return sum.QuoInt(i)
}
