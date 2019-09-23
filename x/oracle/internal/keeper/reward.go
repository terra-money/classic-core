package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/oracle/internal/types"
)

// At the end of every VotePeriod, we give out portion of seigniorage reward(reward-weight) to the
// oracle voters that voted faithfully.
func (k Keeper) RewardBallotWinners(ctx sdk.Context, ballotWinners types.ClaimPool) {
	// Sum weight of the claimpool
	prevBallotWeightSum := int64(0)
	for _, winner := range ballotWinners {
		prevBallotWeightSum += winner.Weight
	}

	if prevBallotWeightSum != 0 {
		rewardPool := k.getRewardPool(ctx)
		if !rewardPool.Empty() {
			// In case rewardFraction = 1%; 1/100 module balance will be distributed
			rewardFraction := k.RewardFraction(ctx)

			// Dole out rewards
			var distributedReward sdk.Coins
			for _, winner := range ballotWinners {
				rewardCoins := sdk.NewCoins()
				rewardeeVal := k.StakingKeeper.Validator(ctx, winner.Recipient)
				for _, feeCoin := range rewardPool {
					rewardAmt := sdk.NewDecFromInt(feeCoin.Amount).Mul(rewardFraction).QuoInt64(prevBallotWeightSum).MulInt64(winner.Weight).TruncateInt()
					rewardCoins = append(rewardCoins, sdk.NewCoin(feeCoin.Denom, rewardAmt))
				}

				// In case absence of the validator, we just skip distribution
				if rewardeeVal != nil {
					k.distrKeeper.AllocateTokensToValidator(ctx, rewardeeVal, sdk.NewDecCoins(rewardCoins))
					distributedReward = distributedReward.Add(rewardCoins)
				}
			}

			// Move distributed reward to distribution module
			err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
			if err != nil {
				panic(fmt.Sprintf("[oracle] Failed to send coins to distribution module %s", err.Error()))
			}
		}
	}
}
