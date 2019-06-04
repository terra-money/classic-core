package oracle

import (
	"github.com/terra-project/core/types/util"

	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/tags"
)

// At the end of every VotePeriod, we give out all the market swap fees collected to the
// oracle voters that voted faithfully.
func rewardPrevBallotWinners(ctx sdk.Context, k Keeper) {
	// Sum weight of the claimpool
	prevBallotWeightSum := sdk.ZeroInt()
	k.iterateClaimPool(ctx, func(_ sdk.AccAddress, weight sdk.Int) (stop bool) {
		prevBallotWeightSum = prevBallotWeightSum.Add(weight)
		return false
	})

	if !prevBallotWeightSum.IsZero() {

		accmFeePool := k.GetSwapFeePool(ctx)
		if !accmFeePool.Empty() {

			// Dole out rewards
			var distributedFee sdk.Coins
			k.iterateClaimPool(ctx, func(recipient sdk.AccAddress, weight sdk.Int) (stop bool) {

				rewardCoins := sdk.NewCoins()
				rewardeeVal := k.valset.Validator(ctx, sdk.ValAddress(recipient))
				for _, feeCoin := range accmFeePool {
					rewardAmt := sdk.NewDecCoinFromCoin(feeCoin).Amount.QuoInt(prevBallotWeightSum).MulInt(weight).TruncateInt()
					rewardCoins = rewardCoins.Add(sdk.NewCoins(sdk.NewCoin(feeCoin.Denom, rewardAmt)))
				}

				// In case absence of the validator, we collect the rewards to fee collect keeper
				if rewardeeVal != nil {
					k.dk.AllocateTokensToValidator(ctx, rewardeeVal, sdk.NewDecCoins(rewardCoins))
				} else {
					k.fck.AddCollectedFees(ctx, rewardCoins)
				}

				distributedFee = distributedFee.Add(rewardCoins)

				return false
			})

			// move left fees to fee collect keeper
			leftFee := accmFeePool.Sub(distributedFee)
			if !leftFee.Empty() && leftFee.IsValid() {
				k.fck.AddCollectedFees(ctx, leftFee)
			}

			// Clear swap fee pool
			k.clearSwapFeePool(ctx)
		}

		// Clear claim and fee pool
		k.clearClaimPool(ctx)
	}
}

// Calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
func tally(ctx sdk.Context, k Keeper, pb PriceBallot) sdk.Dec {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	ballotWinners := types.ClaimPool{}
	weightedMedian := pb.weightedMedian(ctx, k.valset)
	rewardSpread := k.GetParams(ctx).OracleRewardBand.QuoInt64(2)

	for _, vote := range pb {
		if vote.Price.GTE(weightedMedian.Sub(rewardSpread)) && vote.Price.LTE(weightedMedian.Add(rewardSpread)) {
			if validator := k.valset.Validator(ctx, vote.Voter); validator != nil {
				bondSize := validator.GetBondedTokens()

				ballotWinners = append(ballotWinners, types.Claim{
					Recipient: sdk.AccAddress(vote.Voter),
					Weight:    bondSize,
				})
			}
		}
	}

	// add claim winners to the store
	k.addClaimPool(ctx, ballotWinners)

	return weightedMedian
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(totalBondedTokens sdk.Int, voteThreshold sdk.Dec, ballotPower sdk.Int) bool {
	thresholdVotes := voteThreshold.MulInt(totalBondedTokens).RoundInt()
	return ballotPower.GTE(thresholdVotes)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	params := k.GetParams(ctx)

	// Not yet time for a tally
	if !util.IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	// Reward previous ballot winners
	rewardPrevBallotWinners(ctx, k)

	actives := k.getActiveDenoms(ctx)
	votes := k.collectVotes(ctx)

	// Clear swap rates
	for _, activeDenom := range actives {
		k.deletePrice(ctx, activeDenom)
	}

	totalBondedTokens := k.valset.TotalBondedTokens(ctx)

	// Iterate through votes and update prices; drop if not enough votes have been achieved.
	for denom, filteredVotes := range votes {
		if ballotIsPassing(totalBondedTokens, params.VoteThreshold, filteredVotes.power(ctx, k.valset)) {

			// Get weighted median prices, and faithful respondants
			mod := tally(ctx, k, filteredVotes)

			// Set price to the store
			k.SetLunaSwapRate(ctx, denom, mod)

			resTags = sdk.NewTags(
				tags.Action, tags.ActionPriceUpdate,
				tags.Denom, denom,
				tags.Price, mod.String(),
			)
		} else {
			resTags = sdk.NewTags(
				tags.Action, tags.ActionTallyDropped,
				tags.Denom, denom,
			)
		}
	}

	// Clear all prevotes
	k.iteratePrevotes(ctx, func(prevote PricePrevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+params.VotePeriod {
			k.deletePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.iterateVotes(ctx, func(vote PriceVote) (stop bool) {
		k.deleteVote(ctx, vote)
		return false
	})

	return
}
