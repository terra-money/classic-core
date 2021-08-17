package v05

import (
	"sort"

	v04oracle "github.com/terra-money/core/x/oracle/legacy/v04"
	v05oracle "github.com/terra-money/core/x/oracle/types"
)

// Migrate accepts exported v0.4 x/oracle and
// migrates it to v0.5 x/oracle genesis state. The migration includes:
//
// - Remove ExchangeRatePrevote & ExchangeRateVote from x/oracle genesis state.
// - Re-encode in v0.5 GenesisState.
func Migrate(
	oracleGenState v04oracle.GenesisState,
) *v05oracle.GenesisState {
	aggregateExchangeRatePrevote := make([]v05oracle.AggregateExchangeRatePrevote, len(oracleGenState.AggregateExchangeRatePrevotes))
	for i, prevote := range oracleGenState.AggregateExchangeRatePrevotes {
		aggregateExchangeRatePrevote[i] = v05oracle.AggregateExchangeRatePrevote{
			Hash:        prevote.Hash.String(),
			Voter:       prevote.Voter.String(),
			SubmitBlock: uint64(prevote.SubmitBlock),
		}
	}

	aggregateExchangeRateVote := make([]v05oracle.AggregateExchangeRateVote, len(oracleGenState.AggregateExchangeRateVotes))
	for i, prevote := range oracleGenState.AggregateExchangeRateVotes {
		exchangeRateTuples := make([]v05oracle.ExchangeRateTuple, len(prevote.ExchangeRateTuples))
		for j, tuple := range prevote.ExchangeRateTuples {
			exchangeRateTuples[j] = v05oracle.ExchangeRateTuple{
				Denom:        tuple.Denom,
				ExchangeRate: tuple.ExchangeRate,
			}
		}
		aggregateExchangeRateVote[i] = v05oracle.AggregateExchangeRateVote{
			ExchangeRateTuples: exchangeRateTuples,
			Voter:              prevote.Voter.String(),
		}
	}

	whitelist := make([]v05oracle.Denom, len(oracleGenState.Params.Whitelist))
	for i, denom := range oracleGenState.Params.Whitelist {
		whitelist[i] = v05oracle.Denom{
			Name:     denom.Name,
			TobinTax: denom.TobinTax,
		}
	}

	// Note that the four following `for` loop over a map's keys, so are not
	// deterministic.
	i := 0
	missCounters := make([]v05oracle.MissCounter, len(oracleGenState.MissCounters))
	for validatorAddress, missCounter := range oracleGenState.MissCounters {
		missCounters[i] = v05oracle.MissCounter{
			ValidatorAddress: validatorAddress,
			MissCounter:      uint64(missCounter),
		}

		i++
	}

	i = 0
	feederDelegations := make([]v05oracle.FeederDelegation, len(oracleGenState.FeederDelegations))
	for validatorAddress, feederAddress := range oracleGenState.FeederDelegations {
		feederDelegations[i] = v05oracle.FeederDelegation{
			ValidatorAddress: validatorAddress,
			FeederAddress:    feederAddress.String(),
		}

		i++
	}

	i = 0
	exchangeRates := make([]v05oracle.ExchangeRateTuple, len(oracleGenState.ExchangeRates))
	for denom, exchangeRate := range oracleGenState.ExchangeRates {
		exchangeRates[i] = v05oracle.ExchangeRateTuple{
			Denom:        denom,
			ExchangeRate: exchangeRate,
		}

		i++
	}

	i = 0
	tobinTaxes := make([]v05oracle.TobinTax, len(oracleGenState.TobinTaxes))
	for denom, tobinTax := range oracleGenState.TobinTaxes {
		tobinTaxes[i] = v05oracle.TobinTax{
			Denom:    denom,
			TobinTax: tobinTax,
		}

		i++
	}

	// We sort these four arrays by validator address and denom, so that we get determinstic states.
	sort.Slice(missCounters, func(i, j int) bool { return missCounters[i].ValidatorAddress < missCounters[j].ValidatorAddress })
	sort.Slice(feederDelegations, func(i, j int) bool {
		return feederDelegations[i].ValidatorAddress < feederDelegations[j].ValidatorAddress
	})
	sort.Slice(exchangeRates, func(i, j int) bool { return exchangeRates[i].Denom < exchangeRates[j].Denom })
	sort.Slice(tobinTaxes, func(i, j int) bool { return tobinTaxes[i].Denom < tobinTaxes[j].Denom })

	return &v05oracle.GenesisState{
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevote,
		AggregateExchangeRateVotes:    aggregateExchangeRateVote,
		MissCounters:                  missCounters,
		ExchangeRates:                 exchangeRates,
		FeederDelegations:             feederDelegations,
		TobinTaxes:                    tobinTaxes,
		Params: v05oracle.Params{
			VotePeriod:               uint64(oracleGenState.Params.VotePeriod),
			VoteThreshold:            oracleGenState.Params.VoteThreshold,
			RewardBand:               oracleGenState.Params.RewardBand,
			RewardDistributionWindow: uint64(oracleGenState.Params.RewardDistributionWindow),
			SlashFraction:            oracleGenState.Params.SlashFraction,
			SlashWindow:              uint64(oracleGenState.Params.SlashWindow),
			MinValidPerWindow:        oracleGenState.Params.MinValidPerWindow,
			Whitelist:                whitelist,
		},
	}
}
