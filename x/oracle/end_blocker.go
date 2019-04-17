package oracle

import (
	"github.com/terra-project/core/types/util"

	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/tags"
)

// Calculates the median and returns the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median.
func tally(ctx sdk.Context, k Keeper, pb PriceBallot) (weightedMedian sdk.Dec, ballotWinners types.ClaimPool) {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	ballotWinners = types.ClaimPool{}
	weightedMedian = pb.weightedMedian(ctx, k.valset)
	rewardSpread := k.GetParams(ctx).OracleRewardBand.QuoInt64(2)

	for _, vote := range pb {
		if vote.Price.GTE(weightedMedian.Sub(rewardSpread)) && vote.Price.LTE(weightedMedian.Add(rewardSpread)) {
			valAddr := sdk.ValAddress(vote.Voter)
			if validator := k.valset.Validator(ctx, valAddr); validator != nil {
				bondSize := validator.GetBondedTokens()

				ballotWinners = append(ballotWinners, types.Claim{
					Recipient: vote.Voter,
					Weight:    bondSize,
					Class:     types.OracleClaimClass,
				})
			}
		}
	}

	return
}

// Drop the ballot. If the ballot drops params.DropThreshold times sequentially, then blacklist
func dropBallot(ctx sdk.Context, k Keeper, denom string, params Params) sdk.Tags {
	actionTag := tags.ActionTallyDropped

	// Not enough votes received
	dropCounter := k.incrementDropCounter(ctx, denom)
	if dropCounter.GTE(params.DropThreshold) {

		// Too many drops, blacklist currency
		k.deletePrice(ctx, denom)
		k.resetDropCounter(ctx, denom)

		actionTag = tags.ActionBlacklist
	}

	return sdk.NewTags(
		tags.Action, actionTag,
		tags.Denom, denom,
	)
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(totalBondedTokens sdk.Int, voteThreshold sdk.Dec, ballotPower sdk.Int) bool {
	thresholdVotes := voteThreshold.MulInt(totalBondedTokens).RoundInt()
	return ballotPower.GTE(thresholdVotes)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (rewardees types.ClaimPool, resTags sdk.Tags) {
	params := k.GetParams(ctx)

	// Not yet time for a tally
	if !util.IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	actives := k.getActiveDenoms(ctx)
	votes := k.collectVotes(ctx)

	// Iterate through active oracle assets and drop assets that have no votes received.
	for _, activeDenom := range actives {
		if _, found := votes[activeDenom]; !found {
			dropTags := dropBallot(ctx, k, activeDenom, params)
			resTags = resTags.AppendTags(dropTags)
		}
	}

	rewardees = types.ClaimPool{}
	totalBondedTokens := k.valset.TotalBondedTokens(ctx)

	// Iterate through votes and update prices; drop if not enough votes have been achieved.
	for denom, filteredVotes := range votes {
		if ballotIsPassing(totalBondedTokens, params.VoteThreshold, filteredVotes.power(ctx, k.valset)) {

			// Get weighted median prices, and faithful respondants
			mod, ballotWinners := tally(ctx, k, filteredVotes)

			// Append ballot winners for the denom
			rewardees = append(rewardees, ballotWinners...)

			actionTag := tags.ActionPriceUpdate
			if _, err := k.GetLunaSwapRate(ctx, denom); err != nil {
				actionTag = tags.ActionWhitelist
			}

			// Set price to the store
			k.SetLunaSwapRate(ctx, denom, mod)

			// Reset drop counter for the passed ballot
			k.resetDropCounter(ctx, denom)

			resTags = resTags.AppendTags(
				sdk.NewTags(
					tags.Action, actionTag,
					tags.Denom, denom,
					tags.Price, mod.String(),
				),
			)
		} else {
			dropTags := dropBallot(ctx, k, denom, params)
			resTags = resTags.AppendTags(dropTags)
		}

		// Clear all votes
		k.iterateVotes(ctx, func(vote PriceVote) (stop bool) { k.deleteVote(ctx, vote); return false })
	}

	// Sort rewardees before we return
	rewardees.Sort()

	return
}
