package oracle

import (
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	core "github.com/terra-project/core/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {
	params := k.GetParams(ctx)

	// Not yet time for a tally
	if !core.IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	// Reward previous ballot winners
	k.RewardPrevBallotWinners(ctx)

	actives := k.GetActiveDenoms(ctx)
	votes := k.CollectVotes(ctx)

	// Clear swap rates
	for _, activeDenom := range actives {
		k.DeletePrice(ctx, activeDenom)
	}

	ballotAttendees := make(map[string]bool)
	k.StakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator exported.ValidatorI) (stop bool) {
		key := validator.GetOperator().String()
		ballotAttendees[key] = true
		return false
	})

	// Iterate through votes and update prices; drop if not enough votes have been achieved.
	for denom, ballot := range votes {
		if ballotIsPassing(ctx, ballot, k) {

			// Get weighted median prices, and faithful respondants
			mod, ballotWinners, ballotLosers := tally(ctx, ballot, k)

			for _, loser := range ballotLosers {
				key := loser.String()
				if _, exists := ballotAttendees[key]; exists {
					ballotAttendees[key] = false // invalid vote
				}
			}

			// Add claim winners to the store
			k.AddClaimPool(ctx, ballotWinners)

			// TODO - update tax-cap
			// Set price to the store
			k.SetLunaPrice(ctx, denom, mod)
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.EventTypePriceUpdate,
					sdk.NewAttribute(types.AttributeKeyDenom, denom),
					sdk.NewAttribute(types.AttributeKeyPrice, mod.String()),
				),
			)
		}
	}

	// Update & check slash condition for the ballot losers
	k.HandleBallotAttendees(ctx, ballotAttendees)

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
