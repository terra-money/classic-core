package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExchangeRatePrevote - struct to store a validator's prevote on the rate of Luna in the denom asset
type ExchangeRatePrevote struct {
	Hash        VoteHash       `json:"hash"`  // Vote hex hash to protect centralize data source problem
	Denom       string         `json:"denom"` // Ticker name of target fiat currency
	Voter       sdk.ValAddress `json:"voter"` // Voter val address
	SubmitBlock int64          `json:"submit_block"`
}

// NewExchangeRatePrevote returns ExchangeRatePrevote object
func NewExchangeRatePrevote(hash VoteHash, denom string, voter sdk.ValAddress, submitBlock int64) ExchangeRatePrevote {
	return ExchangeRatePrevote{
		Hash:        hash,
		Denom:       denom,
		Voter:       voter,
		SubmitBlock: submitBlock,
	}
}

// String implements fmt.Stringer interface
func (pp ExchangeRatePrevote) String() string {
	return fmt.Sprintf(`ExchangeRatePrevote
	Hash:        %s, 
	Denom:       %s, 
	Voter:       %s, 
	SubmitBlock: %d`,
		pp.Hash, pp.Denom, pp.Voter, pp.SubmitBlock)
}

// ExchangeRatePrevotes is a collection of ExchangeRatePrevote
type ExchangeRatePrevotes []ExchangeRatePrevote

// String implements fmt.Stringer interface
func (v ExchangeRatePrevotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ExchangeRateVote - struct to store a validator's vote on the rate of Luna in the denom asset
type ExchangeRateVote struct {
	ExchangeRate sdk.Dec        `json:"exchange_rate"` // ExchangeRate of Luna in target fiat currency
	Denom        string         `json:"denom"`         // Ticker name of target fiat currency
	Voter        sdk.ValAddress `json:"voter"`         // voter val address of validator
}

// NewExchangeRateVote creates a ExchangeRateVote instance
func NewExchangeRateVote(rate sdk.Dec, denom string, voter sdk.ValAddress) ExchangeRateVote {
	return ExchangeRateVote{
		ExchangeRate: rate,
		Denom:        denom,
		Voter:        voter,
	}
}

// String implements fmt.Stringer interface
func (pv ExchangeRateVote) String() string {
	return fmt.Sprintf(`ExchangeRateVote
	Denom:           %s, 
	Voter:           %s, 
	ExchangeRate:    %s`,
		pv.Denom, pv.Voter, pv.ExchangeRate)
}

// ExchangeRateVotes is a collection of ExchangeRateVote
type ExchangeRateVotes []ExchangeRateVote

// String implements fmt.Stringer interface
func (v ExchangeRateVotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// AggregateExchangeRatePrevote - struct to store a validator's aggregate prevote on the rate of Luna in the denom asset
type AggregateExchangeRatePrevote struct {
	Hash        AggregateVoteHash `json:"hash"`  // Vote hex hash to protect centralize data source problem
	Voter       sdk.ValAddress    `json:"voter"` // Voter val address
	SubmitBlock int64             `json:"submit_block"`
}

// NewAggregateExchangeRatePrevote returns AggregateExchangeRatePrevote object
func NewAggregateExchangeRatePrevote(hash AggregateVoteHash, voter sdk.ValAddress, submitBlock int64) AggregateExchangeRatePrevote {
	return AggregateExchangeRatePrevote{
		Hash:        hash,
		Voter:       voter,
		SubmitBlock: submitBlock,
	}
}

// String implements fmt.Stringer interface
func (pp AggregateExchangeRatePrevote) String() string {
	return fmt.Sprintf(`AggregateExchangeRatePrevote
	Hash:        %s,  
	Voter:       %s, 
	SubmitBlock: %d`,
		pp.Hash, pp.Voter, pp.SubmitBlock)
}

// ExchangeRateTuple - struct to represent a exchange rate of Luna in the denom asset
type ExchangeRateTuple struct {
	Denom        string  `json:"denom"`
	ExchangeRate sdk.Dec `json:"exchange_rate"`
}

// String implements fmt.Stringer interface
func (tuple ExchangeRateTuple) String() string {
	return fmt.Sprintf(`ExchangeRateTuple
	Denom:        %s,
	ExchangeRate: %s`,
		tuple.Denom, tuple.ExchangeRate.String())
}

// ExchangeRateTuples - array of ExchangeRateTuple
type ExchangeRateTuples []ExchangeRateTuple

// String implements fmt.Stringer interface
func (tuples ExchangeRateTuples) String() (out string) {
	for _, tuple := range tuples {
		out += tuple.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ParseExchangeRateTuples ExchangeRateTuple parser
func ParseExchangeRateTuples(tuplesStr string) (ExchangeRateTuples, error) {
	tuplesStr = strings.TrimSpace(tuplesStr)
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	tupleStrs := strings.Split(tuplesStr, ",")
	tuples := make(ExchangeRateTuples, len(tupleStrs))
	duplicateCheckMap := make(map[string]bool)
	for i, tupleStr := range tupleStrs {
		decCoin, err := sdk.ParseDecCoin(tupleStr)
		if err != nil {
			return nil, err
		}

		tuples[i] = ExchangeRateTuple{
			Denom:        decCoin.Denom,
			ExchangeRate: decCoin.Amount,
		}

		if _, ok := duplicateCheckMap[decCoin.Denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", decCoin.Denom)
		}

		duplicateCheckMap[decCoin.Denom] = true
	}

	return tuples, nil
}

// AggregateExchangeRateVote - struct to store a validator's aggregate vote on the rate of Luna in the denom asset
type AggregateExchangeRateVote struct {
	ExchangeRateTuples ExchangeRateTuples `json:"exchange_rate_tuples"` // ExchangeRates of Luna in target fiat currencies
	Voter              sdk.ValAddress     `json:"voter"`                // voter val address of validator
}

// NewAggregateExchangeRateVote creates a AggregateExchangeRateVote instance
func NewAggregateExchangeRateVote(tuples ExchangeRateTuples, voter sdk.ValAddress) AggregateExchangeRateVote {
	return AggregateExchangeRateVote{
		ExchangeRateTuples: tuples,
		Voter:              voter,
	}
}

// String implements fmt.Stringer interface
func (pv AggregateExchangeRateVote) String() string {
	return fmt.Sprintf(`AggregateExchangeRateVote
	ExchangeRateTuples:    %s,
	Voter:                 %s`,
		pv.ExchangeRateTuples, pv.Voter)
}
