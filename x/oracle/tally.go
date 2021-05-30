package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/x/oracle/internal/types"
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

func updateWinnerMap(ballotWinningClaims []types.Claim, validVotesCounterMap map[string]int, winnerMap map[string]types.Claim) {
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
func ballotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot, k Keeper) (sdk.Int, bool) {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower, !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}

// choose Reference Terra with the highest voter turnout
// If the voting power of the two denominations is the same,
// select reference Terra in alphabetical order.
func pickReferenceTerra(ctx sdk.Context, k Keeper, voteTargets map[string]sdk.Dec, voteMap map[string]types.ExchangeRateBallot) string {
	largestBallotPower := int64(0)
	referenceTerra := ""

	for denom, ballot := range voteMap {
		// If denom is not in the voteTargets, or the ballot for it has failed, then skip
		// and remove it from voteMap for iteration efficiency
		if _, exists := voteTargets[denom]; !exists {
			delete(voteMap, denom)
			continue
		}

		ballotPower := int64(0)

		// If the ballot is not passed, remove it from the voteTargets array
		// to prevent slashing validators who did valid vote.
		if power, ok := ballotIsPassing(ctx, ballot, k); ok {
			ballotPower = power.Int64()
		} else {
			delete(voteTargets, denom)
			delete(voteMap, denom)
			continue
		}

		if ballotPower > largestBallotPower || largestBallotPower == 0 {
			referenceTerra = denom
			largestBallotPower = ballotPower
		} else if largestBallotPower == ballotPower && referenceTerra > denom {
			referenceTerra = denom
		}
	}

	return referenceTerra
}
