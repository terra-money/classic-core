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
	validVotesCounterMap := make(map[string]int)
	winnerMap := make(map[string]types.Claim)
	k.StakingKeeper.IterateValidators(ctx, func(_ int64, validator exported.ValidatorI) bool {

		// Exclude not bonded validator or jailed validators from tallying
		if validator.IsBonded() && !validator.IsJailed() {

			// NOTE: we directly stringify byte to string to prevent unnecessary bech32fy works
			valAddr := validator.GetOperator()
			validVotesCounterMap[string(valAddr)] = 0
			winnerMap[string(valAddr)] = types.NewClaim(0, valAddr)
		}

		return false
	})

	// Denom-TobinTax map
	voteTargets := make(map[string]sdk.Dec)
	k.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) bool {
		voteTargets[denom] = tobinTax
		return false
	})

	// Clear all exchange rates
	k.IterateLunaExchangeRates(ctx, func(denom string, _ sdk.Dec) (stop bool) {
		k.DeleteLunaExchangeRate(ctx, denom)
		return false
	})

	// Organize votes to ballot by denom
	// NOTE: **Filter out inactive or jailed validators**
	// NOTE: **Make abstain votes to have zero vote power**
	voteMap := k.OrganizeBallotByDenom(ctx)

	var referenceTerra string
	LargestBallotPower := int64(0)
	voteMapRT := make(map[string]sdk.Dec)

	// choose Reference Terra with the highest voter turnout
	// If the voting power of the two denominations is the same, select reference Terra in alphabetical order.
	for denom, ballot := range voteMap {
		// If denom is not in the voteTargets, or the ballot for it has failed, then skip
		// and remove it from voteMap for iteration efficiency
		if _, exists := voteTargets[denom]; !exists {
			delete(voteMap, denom)
			continue
		}
		ballotPower := ballot.Power()

		// If the ballot is not passed, remove it from the voteTargets array
		// to prevent slashing validators who did valid vote.
		if !ballotIsPassing(ctx, ballot, k) {
			delete(voteTargets, denom)
			delete(voteMap, denom)
			continue
		}

		if ballotPower > LargestBallotPower || LargestBallotPower == 0 {
			referenceTerra = denom
			LargestBallotPower = ballotPower
		} else if LargestBallotPower == ballotPower && referenceTerra > denom {
			referenceTerra = denom
		}
	}

	if referenceTerra != "" {
		// make voteMap of Reference Terra to calculate cross exchange rates
		ballotRT, _ := voteMap[referenceTerra]
		for _, vote := range ballotRT {
			if vote.ExchangeRate.IsPositive() {
				voteMapRT[string(vote.Voter)] = vote.ExchangeRate
			}
		}

		// Get weighted median exchange rates of Reference Terra, and faithful respondants
		ballotMedianRT, ballotWinningClaimsRT := tally(ctx, ballotRT, params.RewardBand)

		// Set the exchange rate, emit ABCI event
		k.SetLunaExchangeRateWithEvent(ctx, referenceTerra, ballotMedianRT)

		// Update winnerMap, validVotesCounterMap using ballotWinningClaims of Reference Terra ballot
		updateWinnerMap(ballotWinningClaimsRT, validVotesCounterMap, winnerMap)

		// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
		for denom, ballot := range voteMap {
			if denom == referenceTerra {
				continue
			}
			// Ballot based cross exchange rates
			var cerBallot types.ExchangeRateBallot
			for _, vote := range ballot {
				if exchangeRateRT, ok := voteMapRT[string(vote.Voter)]; ok {
					vote.ExchangeRate = exchangeRateRT.Quo(vote.ExchangeRate)
					cerBallot = append(cerBallot, vote)
				}
			}
			// Get weighted median of cross exchange rates
			cerMedian, ballotWinningClaims := tally(ctx, cerBallot, params.RewardBand)

			// Update winnerMap, validVotesCounterMap using ballotWinningClaims of cross exchange rate ballot
			updateWinnerMap(ballotWinningClaims, validVotesCounterMap, winnerMap)

			// Transform into the original form uluna/stablecoin
			exchangeRateByRT := ballotMedianRT.Quo(cerMedian)

			// Set the exchange rate, emit ABCI event
			k.SetLunaExchangeRateWithEvent(ctx, denom, exchangeRateByRT)
		}
	}
	//---------------------------
	// Do miss counting & slashing
	voteTargetsLen := len(voteTargets)
	for operatorAddrByteStr, count := range validVotesCounterMap {
		// Skip abstain & valid voters
		if count == voteTargetsLen {
			continue
		}

		// Increase miss counter
		operator := sdk.ValAddress(operatorAddrByteStr) // error never occur
		k.SetMissCounter(ctx, operator, k.GetMissCounter(ctx, operator)+1)
	}

	// Do slash who did miss voting over threshold and
	// reset miss counters of all validators at the last block of slash window
	if core.IsPeriodLastBlock(ctx, params.SlashWindow) {
		SlashAndResetMissCounters(ctx, k)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(ctx, winnerMap)

	// Clear the ballot
	clearBallots(ctx, k, params.VotePeriod)

	// Update vote targets and tobin tax
	applyWhitelist(ctx, k, params.Whitelist, voteTargets)

	return
}

// clearBallots clears all tallied prevotes and votes from the store
func clearBallots(ctx sdk.Context, k Keeper, votePeriod int64) {
	// Clear all prevotes
	k.IterateExchangeRatePrevotes(ctx, func(prevote types.ExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+votePeriod {
			k.DeleteExchangeRatePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.IterateExchangeRateVotes(ctx, func(vote types.ExchangeRateVote) (stop bool) {
		k.DeleteExchangeRateVote(ctx, vote)
		return false
	})

	// Clear all aggregate prevotes
	k.IterateAggregateExchangeRatePrevotes(ctx, func(aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > aggregatePrevote.SubmitBlock+votePeriod {
			k.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)
		}

		return false
	})

	// Clear all aggregate votes
	k.IterateAggregateExchangeRateVotes(ctx, func(vote types.AggregateExchangeRateVote) (stop bool) {
		k.DeleteAggregateExchangeRateVote(ctx, vote)
		return false
	})
}

// applyWhitelist update vote target denom list and set tobin tax with params whitelist
func applyWhitelist(ctx sdk.Context, k Keeper, whitelist types.DenomList, voteTargets map[string]sdk.Dec) {

	// check is there any update in whitelist params
	updateRequired := false
	if len(voteTargets) != len(whitelist) {
		updateRequired = true
	} else {
		for _, item := range whitelist {
			if tobinTax, ok := voteTargets[item.Name]; !ok || !tobinTax.Equal(item.TobinTax) {
				updateRequired = true
				break
			}
		}
	}

	if updateRequired {
		k.ClearTobinTaxes(ctx)

		for _, item := range whitelist {
			k.SetTobinTax(ctx, item.Name, item.TobinTax)
		}
	}
}
