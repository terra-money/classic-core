package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all oracle state that must be provided at genesis
type GenesisState struct {
	Params                        Params                         `json:"params" yaml:"params"`
	FeederDelegations             map[string]sdk.AccAddress      `json:"feeder_delegations" yaml:"feeder_delegations"`
	ExchangeRates                 map[string]sdk.Dec             `json:"exchange_rates" yaml:"exchange_rates"`
	ExchangeRatePrevotes          []ExchangeRatePrevote          `json:"exchange_rate_prevotes" yaml:"exchange_rate_prevotes"`
	ExchangeRateVotes             []ExchangeRateVote             `json:"exchange_rate_votes" yaml:"exchange_rate_votes"`
	MissCounters                  map[string]int64               `json:"miss_counters" yaml:"miss_counters"`
	AggregateExchangeRatePrevotes []AggregateExchangeRatePrevote `json:"aggregate_exchange_rate_prevotes" yaml:"aggregate_exchange_rate_prevotes"`
	AggregateExchangeRateVotes    []AggregateExchangeRateVote    `json:"aggregate_exchange_rate_votes" yaml:"aggregate_exchange_rate_votes"`
	TobinTaxes                    map[string]sdk.Dec             `json:"tobin_taxes" yaml:"tobin_taxes"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, exchangeRatePrevotes []ExchangeRatePrevote,
	exchangeRateVotes []ExchangeRateVote, rates map[string]sdk.Dec,
	feederDelegations map[string]sdk.AccAddress, missCounters map[string]int64,
	aggregateExchangeRatePrevotes []AggregateExchangeRatePrevote,
	aggregateExchangeRateVotes []AggregateExchangeRateVote,
	TobinTaxes map[string]sdk.Dec,
) GenesisState {

	return GenesisState{
		Params:                        params,
		ExchangeRatePrevotes:          exchangeRatePrevotes,
		ExchangeRateVotes:             exchangeRateVotes,
		ExchangeRates:                 rates,
		FeederDelegations:             feederDelegations,
		MissCounters:                  missCounters,
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevotes,
		AggregateExchangeRateVotes:    aggregateExchangeRateVotes,
		TobinTaxes:                    TobinTaxes,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:                        DefaultParams(),
		ExchangeRatePrevotes:          []ExchangeRatePrevote{},
		ExchangeRateVotes:             []ExchangeRateVote{},
		ExchangeRates:                 make(map[string]sdk.Dec),
		FeederDelegations:             make(map[string]sdk.AccAddress),
		MissCounters:                  make(map[string]int64),
		AggregateExchangeRatePrevotes: []AggregateExchangeRatePrevote{},
		AggregateExchangeRateVotes:    []AggregateExchangeRateVote{},
		TobinTaxes:                    make(map[string]sdk.Dec),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return data.Params.ValidateBasic()
}

// Equal checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}
