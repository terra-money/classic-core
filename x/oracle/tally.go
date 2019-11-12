package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// Calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func tally(ctx sdk.Context, pb types.ExchangeRateBallot, powerMap map[string]int64, rewardBand sdk.Dec) (weightedMedian sdk.Dec, ballotWinners []types.Claim) {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	weightedMedian = pb.WeightedMedian(ctx, powerMap)
	standardDeviation := pb.StandardDeviation(ctx, powerMap)
	rewardSpread := rewardBand.QuoInt64(2)

	if standardDeviation.GT(rewardSpread) {
		rewardSpread = standardDeviation
	}

	for _, vote := range pb {
		// If a validator is not found, then just ignore the vote
		if power, ok := powerMap[vote.Voter.String()]; ok {
			if vote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) && vote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread)) {
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
func ballotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot, k Keeper, powerMap map[string]int64) bool {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power(ctx, powerMap))
	return ballotPower.GTE(thresholdVotes)
}
