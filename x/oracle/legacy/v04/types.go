// DONTCOVER
// nolint
package v04

import (
	"encoding/hex"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName nolint
	ModuleName = "oracle"
)

type (
	// GenesisState - all oracle state that must be provided at genesis
	GenesisState struct {
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

	// ExchangeRatePrevote - struct to store a validator's prevote on the rate of Luna in the denom asset
	ExchangeRatePrevote struct {
		Hash        VoteHash       `json:"hash"`  // Vote hex hash to protect centralize data source problem
		Denom       string         `json:"denom"` // Ticker name of target fiat currency
		Voter       sdk.ValAddress `json:"voter"` // Voter val address
		SubmitBlock int64          `json:"submit_block"`
	}

	// ExchangeRateVote - struct to store a validator's vote on the rate of Luna in the denom asset
	ExchangeRateVote struct {
		ExchangeRate sdk.Dec        `json:"exchange_rate"` // ExchangeRate of Luna in target fiat currency
		Denom        string         `json:"denom"`         // Ticker name of target fiat currency
		Voter        sdk.ValAddress `json:"voter"`         // voter val address of validator
	}

	// AggregateExchangeRatePrevote - struct to store a validator's aggregate prevote on the rate of Luna in the denom asset
	AggregateExchangeRatePrevote struct {
		Hash        VoteHash       `json:"hash"`  // Vote hex hash to protect centralize data source problem
		Voter       sdk.ValAddress `json:"voter"` // Voter val address
		SubmitBlock int64          `json:"submit_block"`
	}

	// AggregateExchangeRateVote - struct to store a validator's aggregate vote on the rate of Luna in the denom asset
	AggregateExchangeRateVote struct {
		ExchangeRateTuples ExchangeRateTuples `json:"exchange_rate_tuples"` // ExchangeRates of Luna in target fiat currencies
		Voter              sdk.ValAddress     `json:"voter"`                // voter val address of validator
	}

	// ExchangeRateTuple - struct to represent a exchange rate of Luna in the denom asset
	ExchangeRateTuple struct {
		Denom        string  `json:"denom"`
		ExchangeRate sdk.Dec `json:"exchange_rate"`
	}

	// ExchangeRateTuples - array of ExchangeRateTuple
	ExchangeRateTuples []ExchangeRateTuple

	// Params oracle parameters
	Params struct {
		VotePeriod               int64     `json:"vote_period" yaml:"vote_period"`                               // the number of blocks during which voting takes place.
		VoteThreshold            sdk.Dec   `json:"vote_threshold" yaml:"vote_threshold"`                         // the minimum percentage of votes that must be received for a ballot to pass.
		RewardBand               sdk.Dec   `json:"reward_band" yaml:"reward_band"`                               // the ratio of allowable exchange rate error that can be rewarded.
		RewardDistributionWindow int64     `json:"reward_distribution_window" yaml:"reward_distribution_window"` // the number of blocks during which seigniorage reward comes in and then is distributed.
		Whitelist                DenomList `json:"whitelist" yaml:"whitelist"`                                   // the denom list that can be activated,
		SlashFraction            sdk.Dec   `json:"slash_fraction" yaml:"slash_fraction"`                         // the ratio of penalty on bonded tokens
		SlashWindow              int64     `json:"slash_window" yaml:"slash_window"`                             // the number of blocks for slashing tallying
		MinValidPerWindow        sdk.Dec   `json:"min_valid_per_window" yaml:"min_valid_per_window"`             // the ratio of minimum valid oracle votes per slash window to avoid slashing
	}

	// Denom is the object to hold configurations of each denom
	Denom struct {
		Name     string  `json:"name" yaml:"name"`
		TobinTax sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`
	}

	// DenomList is array of Denom
	DenomList []Denom

	// VoteHash is hash value to hide vote exchange rate
	// which is formatted as hex string in SHA256("{salt}:{exchange_rate}:{denom}:{voter}")
	VoteHash []byte
)

// VoteHashFromHexString convert hex string to VoteHash
func VoteHashFromHexString(s string) (VoteHash, error) {
	h, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// String implements fmt.Stringer interface
func (h VoteHash) String() string {
	return hex.EncodeToString(h)
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (h VoteHash) Marshal() ([]byte, error) {
	return h, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (h *VoteHash) Unmarshal(data []byte) error {
	*h = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (h VoteHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (h VoteHash) MarshalYAML() (interface{}, error) {
	return h.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (h *VoteHash) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	h2, err := VoteHashFromHexString(s)
	if err != nil {
		return err
	}

	*h = h2
	return nil
}
