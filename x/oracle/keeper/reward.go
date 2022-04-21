package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/types"
)

// RewardBallotWinners implements
// at the end of every VotePeriod, give out a portion of spread fees collected in the oracle reward pool
//  to the oracle voters that voted faithfully.
func (k Keeper) RewardBallotWinners(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	voteTargets map[string]sdk.Dec,
	ballotWinners map[string]types.Claim,
) {
	// softfork for reward distribution
	if (ctx.ChainID() == core.ColumbusChainID && ctx.BlockHeight() < int64(5_100_000)) ||
		(ctx.ChainID() == core.BombayChainID && ctx.BlockHeight() < int64(6_200_000)) {
		k.RewardBallotWinnersLegacy(ctx, votePeriod, rewardDistributionWindow, ballotWinners)
		return
	}

	rewardDenoms := make([]string, len(voteTargets)+1)
	rewardDenoms[0] = core.MicroLunaDenom

	i := 1
	for denom := range voteTargets {
		rewardDenoms[i] = denom
		i++
	}

	// Sum weight of the claims
	ballotPowerSum := int64(0)
	for _, winner := range ballotWinners {
		ballotPowerSum += winner.Weight
	}

	// Exit if the ballot is empty
	if ballotPowerSum == 0 {
		return
	}
	
	// The Reward distributionRatio = votePeriod/rewardDistributionWindow
	distributionRatio := sdk.NewDec(votePeriod).QuoInt64(rewardDistributionWindow)

	var periodRewards sdk.DecCoins
	for _, denom := range rewardDenoms {
		rewardPool := k.GetRewardPool(ctx, denom)

		// return if there's no rewards to give out
		if rewardPool.IsZero() {
			continue
		}

		periodRewards = periodRewards.Add(sdk.NewDecCoinFromDec(
			denom,
			sdk.NewDecFromInt(rewardPool.Amount).Mul(distributionRatio),
		))
	}

	// Dole out rewards
	var distributedReward sdk.Coins
	for _, winner := range ballotWinners {
		receiverVal := k.StakingKeeper.Validator(ctx, winner.Recipient)

		// Reflects contribution
		rewardCoins, _ := periodRewards.MulDec(sdk.NewDec(winner.Weight).QuoInt64(ballotPowerSum)).TruncateDecimal()

		// In case absence of the validator, we just skip distribution
		if receiverVal != nil && !rewardCoins.IsZero() {
			k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			distributedReward = distributedReward.Add(rewardCoins...)
		}
	}

	// Move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Sprintf("[oracle] Failed to send coins to distribution module %s", err.Error()))
	}

}

// RewardBallotWinnersLegacy implements
// at the end of every VotePeriod, we give out portion of seigniorage reward(reward-weight) to the
// oracle voters that voted faithfully.
func (k Keeper) RewardBallotWinnersLegacy(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	ballotWinners map[string]types.Claim) {
	// Sum weight of the claims
	ballotPowerSum := int64(0)
	for _, winner := range ballotWinners {
		ballotPowerSum += winner.Weight
	}

	// Exit if the ballot is empty
	if ballotPowerSum == 0 {
		return
	}

	rewardPool := k.GetRewardPoolLegacy(ctx)

	// return if there's no rewards to give out
	if rewardPool.IsZero() {
		return
	}

	// rewardCoin  = oraclePool * VotePeriod / RewardDistributionWindow
	periodRewards := sdk.NewDecFromInt(rewardPool.AmountOf(core.MicroLunaDenom)).
		MulInt64(votePeriod).QuoInt64(rewardDistributionWindow)

	// Dole out rewards
	var distributedReward sdk.Coins
	for _, winner := range ballotWinners {
		rewardCoins := sdk.NewCoins()
		receiverVal := k.StakingKeeper.Validator(ctx, winner.Recipient)

		// Reflects contribution
		rewardAmt := periodRewards.QuoInt64(ballotPowerSum).MulInt64(winner.Weight).TruncateInt()
		rewardCoins = append(rewardCoins, sdk.NewCoin(core.MicroLunaDenom, rewardAmt))

		// In case absence of the validator, we just skip distribution
		if receiverVal != nil && !rewardCoins.IsZero() {
			k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			distributedReward = distributedReward.Add(rewardCoins...)
		}
	}

	// Move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Sprintf("[oracle] Failed to send coins to distribution module %s", err.Error()))
	}

}
