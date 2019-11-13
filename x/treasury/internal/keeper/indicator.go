package keeper

import (
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Computes important economic indicators for the stability of Terra currencies.
//
// Three important concepts:
// - MR: Fees + Seigniorage for a given epoch sums to Mining Rewards
// - SR: Computes the Seigniorage Reward
// - TRL: Computes the Tax Reward per unit Luna
// - SMR: Computes the ratio of seigniorage rewards to overall mining rewards
//
// Rolling averages are also computed for MRL and SMR respectively.
//

// alignCoins align the coins to the given denom through the market swap
func (k Keeper) alignCoins(ctx sdk.Context, coins sdk.DecCoins, denom string) (alignedAmt sdk.Dec) {
	alignedAmt = sdk.ZeroDec()
	for _, coinReward := range coins {
		if coinReward.Denom != denom {
			swappedReward, err := k.marketKeeper.ComputeInternalSwap(ctx, coinReward, denom)
			if err != nil {
				continue
			}
			alignedAmt = alignedAmt.Add(swappedReward.Amount)
		} else {
			alignedAmt = alignedAmt.Add(coinReward.Amount)
		}
	}

	return alignedAmt
}

// UpdateIndicators updates interal indicators
func (k Keeper) UpdateIndicators(ctx sdk.Context) {
	epoch := core.GetEpoch(ctx)
	lunaTotalBondedAmount := k.stakingKeeper.TotalBondedTokens(ctx)

	// Compute Tax Rewards & Tax Rewards per (staked)Luna
	taxRewards := sdk.NewDecCoins(k.PeekTaxProceeds(ctx))
	TR := k.alignCoins(ctx, taxRewards, core.MicroSDRDenom)
	TRL := TR.QuoInt(lunaTotalBondedAmount)

	k.SetTRL(ctx, epoch, TRL)

	// Reset tax proceeds after computing TRL for the next epoch
	k.SetTaxProceeds(ctx, sdk.Coins{})

	// Compute Seigniorage Rewards
	seigniorage := k.PeekEpochSeigniorage(ctx)
	seigniorageRewardsAmt := k.GetRewardWeight(ctx).MulInt(seigniorage)
	seigniorageRewards := sdk.DecCoins{sdk.NewDecCoinFromDec(core.MicroLunaDenom, seigniorageRewardsAmt)}
	SR := k.alignCoins(ctx, seigniorageRewards, core.MicroSDRDenom)

	k.SetSR(ctx, epoch, SR)

	// Compute Mining Rewards
	MR := TR.Add(SR)

	k.SetMR(ctx, epoch, MR)
}

func (k Keeper) loadIndicatorByEpoch(ctx sdk.Context, indicatorPrefix []byte, epoch int64) (indicator sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSubkeyByEpoch(indicatorPrefix, epoch))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &indicator)
	return
}

// sumIndicator returns the sum of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return sumIndicator(currentEpoch)
func (k Keeper) sumIndicator(ctx sdk.Context, epochs int64, indicatorPrefix []byte) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := core.GetEpoch(ctx)

	for i := curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		sum = sum.Add(k.loadIndicatorByEpoch(ctx, indicatorPrefix, i))
	}

	return sum
}

// rollingAverageIndicator returns the rolling average of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return rollingAverageIndicator(currentEpoch)
func (k Keeper) rollingAverageIndicator(ctx sdk.Context, epochs int64, indicatorPrefix []byte) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := core.GetEpoch(ctx)

	var i int64
	for i = curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		sum = sum.Add(k.loadIndicatorByEpoch(ctx, indicatorPrefix, i))
	}

	computedEpochs := curEpoch - i
	if computedEpochs == 0 {
		return sum
	}

	return sum.QuoInt64(computedEpochs)
}
