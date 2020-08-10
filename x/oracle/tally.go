package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// Calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func tally(ctx sdk.Context, pb types.ExchangeRateBallot, rewardBand sdk.Dec) (weightedMedian sdk.Dec, ballotWinners []types.Claim) {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	weightedMedian = pb.WeightedMedian()
	standardDeviation := pb.StandardDeviation()
	rewardSpread := weightedMedian.Mul(rewardBand.QuoInt64(2))

	if standardDeviation.GT(rewardSpread) {
		rewardSpread = standardDeviation
	}

	for _, vote := range pb {
		// Filter ballot winners & abstain voters
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread))) ||
			!vote.ExchangeRate.IsPositive() {

			// Abstain votes have zero vote power
			ballotWinners = append(ballotWinners, types.Claim{
				Recipient: vote.Voter,
				Weight:    vote.Power,
			})
		}

	}

	return
}

func handleBallotWinner(ballotWinningClaims []types.Claim, validVotesCounterMap map[string]int, winnerMap map[string]types.Claim) {
	// Collect claims of ballot winners
	for _, ballotWinningClaim := range ballotWinningClaims {

		// NOTE: we directly stringify byte to string to prevent unnecessary bech32fy works
		key := string(ballotWinningClaim.Recipient)

		// Update claim
		prevClaim := winnerMap[key]
		prevClaim.Weight += ballotWinningClaim.Weight
		winnerMap[key] = prevClaim

		// Increase valid votes counter
		validVotesCounterMap[key]++
	}
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot, k Keeper) bool {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power())
	return !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}
