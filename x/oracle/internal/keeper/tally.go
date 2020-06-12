package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/oracle/internal/types"
	"sort"
)

// ballot for the asset is passing the threshold amount of voting power
func (k Keeper) BallotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot) bool {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower.GTE(thresholdVotes)
}

// Calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func (k Keeper) Tally(pb types.ExchangeRateBallot, rewardBand sdk.Dec) (weightedMedian sdk.Dec, ballotWinners []types.Claim) {
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

// Calculates the median for cross exchange rate and returns it.
func (k Keeper) TallyCrossRate(ctx sdk.Context, voteMap map[string]types.ExchangeRateBallot, voteTargets map[string]sdk.Dec) (crossExchangeRates types.CrossExchangeRates) {
	crossRateMapByVali := make(map[string]map[string]sdk.Dec)
	crossRateMapByDenom := make(map[string]map[string]types.ExchangeRateBallot)

	// Organize by denom pair for cross rate
	for denom, ballot := range voteMap {
		if _, exists := voteTargets[denom]; !exists {
			continue
		}
		for _, vote := range ballot {
			if vote.ExchangeRate.IsPositive() {
				for k, v := range crossRateMapByVali[vote.Voter.String()]{
					var crossRate sdk.Dec
					denom1, denom2 := types.GetDenomOrderAsc(denom, k)
					if denom > k {
						crossRate = v.Quo(vote.ExchangeRate)
					} else {
						crossRate = vote.ExchangeRate.Quo(v)
					}
					if crossRate.IsPositive() {
						if _, ok := crossRateMapByDenom[denom1]; !ok {
							crossRateMapByDenom[denom1] = make(map[string]types.ExchangeRateBallot)
						}
						crossRateMapByDenom[denom1][denom2] = append(crossRateMapByDenom[denom1][denom2],
							types.NewVoteForTally(
								types.NewExchangeRateVote(crossRate, denom1+"_"+denom2, vote.Voter),
								vote.Power,
							),
						)
					}
				}
				if _, ok := crossRateMapByVali[vote.Voter.String()]; !ok {
					crossRateMapByVali[vote.Voter.String()] = make(map[string]sdk.Dec)
				}
				crossRateMapByVali[vote.Voter.String()][denom] = vote.ExchangeRate  // for only existing flag tmp
			}
		}
	}

	// Get Weighted Median for each denom pair
	for denom1, denom2List := range crossRateMapByDenom {
		for denom2, erb := range denom2List {
			// Check quorum threshold, Get Weighted Median
			if !k.BallotIsPassing(ctx, erb){
				continue
			}
			wm := erb.WeightedMedian()
			if wm.IsPositive(){
				crossExchangeRates = append(crossExchangeRates, types.NewCrossExchangeRate(denom1, denom2, wm))
			}
		}
	}
	return
}