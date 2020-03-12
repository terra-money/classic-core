package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// alignCoins aligns the coins to the given denom through the market swap
func (k Keeper) alignCoins(ctx sdk.Context, coins sdk.Coins, denom string) (alignedAmt sdk.Dec, err sdk.Error) {
	decCoins := sdk.NewDecCoins(coins)

	alignedAmt = sdk.ZeroDec()
	for _, coin := range decCoins {
		if coin.Denom != denom {
			swappedReward, err := k.marketKeeper.ComputeInternalSwap(ctx, coin, denom)
			if err != nil {
				return sdk.Dec{}, err
			}
			alignedAmt = alignedAmt.Add(swappedReward.Amount)
		} else {
			alignedAmt = alignedAmt.Add(coin.Amount)
		}
	}

	return
}

// ConvertRenewalFeeToTime returns extended time period with the given fees for the specific name length
func (k Keeper) ConvertRenewalFeeToTime(ctx sdk.Context, fees sdk.Coins, nameLength int) (extendedTime time.Duration, err sdk.Error) {
	renewalFee := k.RenewalFees(ctx).RenewalFeeForLength(nameLength)
	renewalInterval := k.RenewalInterval(ctx)

	feeAmount, err := k.alignCoins(ctx, fees, renewalFee.Denom)
	if err != nil {
		return 0, err
	}

	extendedTime = time.Duration(
		feeAmount.MulInt64(int64(renewalInterval)).QuoInt(renewalFee.Amount).TruncateInt64(),
	)

	return
}
