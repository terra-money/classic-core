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

	// Build valid votes counter and winner map over all validators in active set
	validVotesCounterMap := make(map[string]int64)
	winnerMap := make(map[string]types.Claim)
	k.StakingKeeper.IterateValidators(ctx, func(index int64, validator exported.ValidatorI) bool {

		// Exclude not bonded vaildator or jailed validators from tallying
		if validator.IsBonded() && !validator.IsJailed() {
			valAddr := validator.GetOperator()
			validVotesCounterMap[valAddr.String()] = int64(0)
			winnerMap[valAddr.String()] = types.NewClaim(0, valAddr)
		}

		return false
	})

	whitelist := make(map[string]bool)
	for _, denom := range k.Whitelist(ctx) {
		whitelist[denom] = true
	}

	// Clear exchange rates
	for denom := range whitelist {
		k.DeleteLunaExchangeRate(ctx, denom)
	}

	// Organize votes to ballot by denom
	// NOTE: **Filter out inative or jailed validators**
	// NOTE: **Make abstain votes to have zero vote power**
	voteMap := k.OrganizeBallotByDenom(ctx)

	// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
	for denom, ballot := range voteMap {

		// If denom is not in the whitelist, or the ballot for it has failed, then skip
		if _, exists := whitelist[denom]; !exists {
			continue
		}

		// If the ballot is not passed, then remove it from the whitelist array
		// to prevent slashing validators who did valid vote.
		if !ballotIsPassing(ctx, ballot, k) {
			delete(whitelist, denom)
			continue
		}

		// Get weighted median exchange rates, and faithful respondants
		ballotMedian, ballotWinningClaims := tally(ctx, ballot, params.RewardBand)

		// Set the exchange rate
		k.SetLunaExchangeRate(ctx, denom, ballotMedian)

		// Collect claims of ballot winners
		for _, ballotWinningClaim := range ballotWinningClaims {
			key := ballotWinningClaim.Recipient.String()

			// Update claim
			prevClaim := winnerMap[key]
			prevClaim.Weight += ballotWinningClaim.Weight
			winnerMap[key] = prevClaim

			// Increase valid votes counter
			validVotesCounterMap[key]++
		}

		// Emit abci events
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeExchangeRateUpdate,
				sdk.NewAttribute(types.AttributeKeyDenom, denom),
				sdk.NewAttribute(types.AttributeKeyExchangeRate, ballotMedian.String()),
			),
		)
	}

	//---------------------------
	// Do miss counting & slashing

	whitelistLen := int64(len(whitelist))
	for operatorBechAddr, count := range validVotesCounterMap {
		// Skip abstain & valid voters
		if count == whitelistLen {
			continue
		}

		// Increase miss counter
		operator, _ := sdk.ValAddressFromBech32(operatorBechAddr) // error never occur
		k.SetMissCounter(ctx, operator, k.GetMissCounter(ctx, operator)+1)
	}

	// Do slash who did miss voting over threshold and
	// reset miss counters of all validators at the last block of slash window
	if core.IsPeriodLastBlock(ctx, params.VotePeriod*params.SlashWindow) {
		SlashAndResetMissCounters(ctx, k)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(ctx, winnerMap)

	// Clear the ballot
	clearBallots(k, ctx, params)

	return
}

// clearBallots clears all tallied prevotes and votes from the store
func clearBallots(k Keeper, ctx sdk.Context, params Params) {
	// Clear all prevotes
	k.IterateExchangeRatePrevotes(ctx, func(prevote types.ExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+params.VotePeriod {
			k.DeleteExchangeRatePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.IterateExchangeRateVotes(ctx, func(vote types.ExchangeRateVote) (stop bool) {
		k.DeleteExchangeRateVote(ctx, vote)
		return false
	})
}
