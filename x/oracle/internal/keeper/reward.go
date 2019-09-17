package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/oracle/internal/types"
)

// At the end of every VotePeriod, we give out portion of seigniorage reward(reward-weight) to the
// oracle voters that voted faithfully.
func (k Keeper) RewardPrevBallotWinners(ctx sdk.Context) {
	// Sum weight of the claimpool
	prevBallotWeightSum := int64(0)
	k.IterateClaimPool(ctx, func(_ sdk.ValAddress, weight int64) (stop bool) {
		prevBallotWeightSum += weight
		return false
	})

	if prevBallotWeightSum != 0 {
		rewardPool := k.getRewardPool(ctx)
		if !rewardPool.Empty() {
			// if rewardFraction = 1%; 1/100 module balance will be distributed
			rewardFraction := k.RewardFraction(ctx)

			// Dole out rewards
			var distributedReward sdk.Coins
			k.IterateClaimPool(ctx, func(recipient sdk.ValAddress, weight int64) (stop bool) {

				rewardCoins := sdk.NewCoins()
				rewardeeVal := k.StakingKeeper.Validator(ctx, recipient)
				for _, feeCoin := range rewardPool {
					rewardAmt := sdk.NewDecFromInt(feeCoin.Amount).Mul(rewardFraction).QuoInt64(prevBallotWeightSum).MulInt64(weight).TruncateInt()
					rewardCoins = append(rewardCoins, sdk.NewCoin(feeCoin.Denom, rewardAmt))
				}

				// In case absence of the validator, we just skip distribution
				if rewardeeVal != nil {
					k.distrKeeper.AllocateTokensToValidator(ctx, rewardeeVal, sdk.NewDecCoins(rewardCoins))
					distributedReward = distributedReward.Add(rewardCoins)
				}

				return false
			})

			// Move distributed reward to distribution module
			err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
			if err != nil {
				panic(fmt.Sprintf("[oracle] Failed to send coins to distribution module %s", err.Error()))
			}
		}

		// Clear claim and fee pool
		k.clearClaimPool(ctx)
	}
}
