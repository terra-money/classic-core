package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// Calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func tally(ctx sdk.Context, pb types.PriceBallot, k Keeper) (weightedMedian sdk.Dec, ballotWinners types.ClaimPool) {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	weightedMedian = pb.WeightedMedian(ctx, k.StakingKeeper)
	standardDeviation := pb.StandardDeviation(ctx, k.StakingKeeper)
	rewardSpread := k.RewardBand(ctx).QuoInt64(2)

	if standardDeviation.GT(rewardSpread) {
		rewardSpread = standardDeviation
	}

	for _, vote := range pb {
		// If a validator is not found, then just ignore the vote
		if validator := k.StakingKeeper.Validator(ctx, vote.Voter); validator != nil {
			if vote.Price.GTE(weightedMedian.Sub(rewardSpread)) && vote.Price.LTE(weightedMedian.Add(rewardSpread)) {
				power := validator.GetConsensusPower()

				ballotWinners = append(ballotWinners, types.Claim{
					Recipient: vote.Voter,
					Weight:    power,
				})
			}
		}
	}

	return
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(ctx sdk.Context, ballot types.PriceBallot, k Keeper) bool {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power(ctx, k.StakingKeeper))
	return ballotPower.GTE(thresholdVotes)
}
