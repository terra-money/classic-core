package keeper

import (
	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetEpoch returns current epoch of (current block height + cumulated block height of past chains)
func (k Keeper) GetEpoch(ctx sdk.Context) int64 {
	return ctx.BlockHeight() / int64(core.BlocksPerWeek)
}

//
// Computes important economic indicators for the stability of Terra currencies.
//
// Three important concepts:
// - MR: Fees + Seigniorage for a given epoch sums to Mining Rewards
// - SR: Computes the Seigniorage Reward
// - TR: Computes the Tax Reward
// - TSL: Total Staked Luna
// - TRL: Computes the Tax Reward per unit Luna (TR/TSL)

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

// UpdateIndicators updates internal indicators
func (k Keeper) UpdateIndicators(ctx sdk.Context) {
	epoch := k.GetEpoch(ctx)

	// Compute Total Staked Luna (TSL)
	totalStakedLuna := k.stakingKeeper.TotalBondedTokens(ctx)

	k.SetTSL(ctx, epoch, totalStakedLuna)

	// Compute Tax Rewards (TR)
	taxRewards := sdk.NewDecCoinsFromCoins(k.PeekEpochTaxProceeds(ctx)...)
	TR := k.alignCoins(ctx, taxRewards, core.MicroSDRDenom)

	k.SetTR(ctx, epoch, TR)

	// Reset tax proceeds after computing TRL for the next epoch
	k.SetEpochTaxProceeds(ctx, sdk.Coins{})

	// Compute Seigniorage Rewards (SR)
	seigniorage := k.PeekEpochSeigniorage(ctx)
	seigniorageRewardsAmt := k.GetRewardWeight(ctx).MulInt(seigniorage)
	seigniorageRewards := sdk.DecCoins{sdk.NewDecCoinFromDec(core.MicroLunaDenom, seigniorageRewardsAmt)}
	SR := k.alignCoins(ctx, seigniorageRewards, core.MicroSDRDenom)

	k.SetSR(ctx, epoch, SR)
}

// TRL returns Tax Rewards per Luna for the epoch
func TRL(ctx sdk.Context, epoch int64, k Keeper) sdk.Dec {
	tr := k.GetTR(ctx, epoch)
	tsl := k.GetTSL(ctx, epoch)

	// division by zero protection
	if tr.IsZero() || tsl.IsZero() {
		return sdk.ZeroDec()
	}

	return tr.QuoInt(tsl)
}

// SR returns Seigniorage Rewards for the epoch
func SR(ctx sdk.Context, epoch int64, k Keeper) sdk.Dec {
	return k.GetSR(ctx, epoch)
}

// MR returns Mining Rewards = Seigniorage Rewards + Tax Rates for the epoch
func MR(ctx sdk.Context, epoch int64, k Keeper) sdk.Dec {
	return k.GetTR(ctx, epoch).Add(k.GetSR(ctx, epoch))
}

// sumIndicator returns the sum of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return sumIndicator(currentEpoch)
func (k Keeper) sumIndicator(ctx sdk.Context, epochs int64,
	indicator func(ctx sdk.Context, epoch int64, k Keeper) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := k.GetEpoch(ctx)

	for i := curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		val := indicator(ctx, i, k)
		sum = sum.Add(val)
	}

	return sum
}

// rollingAverageIndicator returns the rolling average of the indicator over several epochs.
// If current epoch < epochs, we return the best we can and return rollingAverageIndicator(currentEpoch)
func (k Keeper) rollingAverageIndicator(ctx sdk.Context, epochs int64,
	indicator func(ctx sdk.Context, epoch int64, k Keeper) sdk.Dec) sdk.Dec {
	sum := sdk.ZeroDec()
	curEpoch := k.GetEpoch(ctx)

	var i int64
	for i = curEpoch; i >= 0 && i > (curEpoch-epochs); i-- {
		val := indicator(ctx, i, k)
		sum = sum.Add(val)
	}

	computedEpochs := curEpoch - i
	if computedEpochs == 0 {
		return sum
	}

	return sum.QuoInt64(computedEpochs)
}
