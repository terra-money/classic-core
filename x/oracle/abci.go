package oracle

import (
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {
	params := k.GetParams(ctx)

	// Not yet time for a tally
	if !core.IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	whitelist := k.Whitelist(ctx)

	// Clear exchange rates
	for denom := range whitelist {
		k.DeleteLunaExchangeRate(ctx, denom)
	}

	winnerMap := make(map[string]types.Claim)

	// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
	voteMap := OrganizeBallotByDenom(k, ctx)
	for denom, ballot := range voteMap {

		// If denom is not in the whitelist, or the ballot for it has failed, then skip
		if _, exists := whitelist[denom]; !exists || !ballotIsPassing(ctx, ballot, k) {
			continue
		}

		// Get weighted median exchange rates, and faithful respondants
		ballotMedian, ballotWinningClaims := tally(ctx, ballot, k)

		// Set the exchange rate
		k.SetLunaExchangeRate(ctx, denom, ballotMedian)

		// Collect claims of ballot winners
		for _, ballotWinningClaim := range ballotWinningClaims {
			key := ballotWinningClaim.Recipient.String()
			prevClaim, exists := winnerMap[key]
			if !exists {
				winnerMap[key] = ballotWinningClaim
			} else {
				prevClaim.Weight += ballotWinningClaim.Weight
				winnerMap[key] = prevClaim
			}
		}

		// Emit abci events
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeExchangeRateUpdate,
				sdk.NewAttribute(types.AttributeKeyDenom, denom),
				sdk.NewAttribute(types.AttributeKeyExchangeRate, ballotMedian.String()),
			),
		)
	}

	// Convert map to array
	var claimPool types.ClaimPool
	for _, claim := range winnerMap {
		claimPool = append(claimPool, claim)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(ctx, claimPool)

	// Clear the ballot
	clearBallots(k, ctx, params)

	return
}

// clearBallots clears all tallied prevotes and votes from the store
func clearBallots(k Keeper, ctx sdk.Context, params Params) {
	// Clear all prevotes
	k.IteratePrevotes(ctx, func(prevote Prevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+params.VotePeriod {
			k.DeletePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.IterateVotes(ctx, func(vote Vote) (stop bool) {
		k.DeleteVote(ctx, vote)
		return false
	})
}
