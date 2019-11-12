package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ExchangeRatePrevote - struct to store a validator's prevote on the rate of Luna in the denom asset
type ExchangeRatePrevote struct {
	Hash        string         `json:"hash"`  // Vote hex hash to protect centralize data source problem
	Denom       string         `json:"denom"` // Ticker name of target fiat currency
	Voter       sdk.ValAddress `json:"voter"` // Voter val address
	SubmitBlock int64          `json:"submit_block"`
}

// NewExchangeRatePrevote returns ExchangeRatePrevote object
func NewExchangeRatePrevote(hash string, denom string, voter sdk.ValAddress, submitBlock int64) ExchangeRatePrevote {
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
	Hash:    %s, 
	Denom:    %s, 
	Voter:    %s, 
	SubmitBlock:    %d`,
		pp.Hash, pp.Denom, pp.Voter, pp.SubmitBlock)
}

// ExchangeRatePrevotes is a collection of PreicePrevote
type ExchangeRatePrevotes []ExchangeRatePrevote

// String implements fmt.Stringer interface
func (v ExchangeRatePrevotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// VoteHash computes hash value of ExchangeRateVote
func VoteHash(salt string, rate sdk.Dec, denom string, voter sdk.ValAddress) ([]byte, error) {
	hash := tmhash.NewTruncated()
	_, err := hash.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", salt, rate, denom, voter)))
	bz := hash.Sum(nil)
	return bz, err
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

func (pv ExchangeRateVote) getPower(ctx sdk.Context, powerMap map[string]int64) int64 {
	if power, ok := powerMap[pv.Voter.String()]; ok {
		return power
	}

	return 0
}

// String implements fmt.Stringer interface
func (pv ExchangeRateVote) String() string {
	return fmt.Sprintf(`ExchangeRateVote
	Denom:    %s, 
	Voter:    %s, 
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
