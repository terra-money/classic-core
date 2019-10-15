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

	actives := k.GetActiveDenoms(ctx)
	votes := k.CollectVotes(ctx)

	// Clear prices
	for _, activeDenom := range actives {
		k.DeletePrice(ctx, activeDenom)
	}

	// Changes whitelist array to map for fast lookup
	whitelistMap := make(map[string]bool)
	for _, denom := range k.Whitelist(ctx) {
		whitelistMap[denom] = true
	}

	// Iterate through votes and update prices; drop if not enough votes have been achieved.
	claimMap := make(map[string]types.Claim)
	for denom, ballot := range votes {

		// Check whitelist; if denom is not exists or exists but the ballot is not passed, then skip
		if _, exists := whitelistMap[denom]; !exists || !ballotIsPassing(ctx, ballot, k) {
			continue
		}

		// Get weighted median prices, and faithful respondants
		mod, ballotWinners := tally(ctx, ballot, k)

		// Collect claims of ballot winners
		for _, winner := range ballotWinners {
			key := winner.Recipient.String()
			claim, exists := claimMap[key]
			if exists {
				claim.Weight += winner.Weight
				claimMap[key] = claim
			} else {
				claimMap[key] = winner
			}
		}

		// Set price to the store
		k.SetLunaPrice(ctx, denom, mod)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypePriceUpdate,
				sdk.NewAttribute(types.AttributeKeyDenom, denom),
				sdk.NewAttribute(types.AttributeKeyPrice, mod.String()),
			),
		)
	}

	// Convert map to array
	var claimPool types.ClaimPool
	for _, claim := range claimMap {
		claimPool = append(claimPool, claim)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(ctx, claimPool)

	// Clear all prevotes
	k.IteratePrevotes(ctx, func(prevote PricePrevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+params.VotePeriod {
			k.DeletePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.IterateVotes(ctx, func(vote PriceVote) (stop bool) {
		k.DeleteVote(ctx, vote)
		return false
	})

	return
}
