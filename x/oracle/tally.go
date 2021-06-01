package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/oracle/keeper"
	"github.com/terra-money/core/x/oracle/types"
)

// Tally calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func Tally(ctx sdk.Context, pb types.ExchangeRateBallot, rewardBand sdk.Dec, validatorClaimMap map[string]types.Claim) (weightedMedian sdk.Dec) {
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

			key := vote.Voter.String()
			claim := validatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			validatorClaimMap[key] = claim
		}
	}

	return
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot, k keeper.Keeper) (sdk.Int, bool) {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx), k.StakingKeeper.PowerReduction(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower, !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}

// choose Reference Terra with the highest voter turnout
// If the voting power of the two denominations is the same,
// select reference Terra in alphabetical order.
func pickReferenceTerra(ctx sdk.Context, k keeper.Keeper, voteTargets map[string]sdk.Dec, voteMap map[string]types.ExchangeRateBallot) string {
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
