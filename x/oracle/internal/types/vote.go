package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Prevote - struct to store a validator's prevote on the exchangeRate of Luna in the denom asset
type Prevote struct {
	Hash        string         `json:"hash"`  // Vote hex hash to protect centralize data source problem
	Denom       string         `json:"denom"` // Ticker name of target fiat currency
	Voter       sdk.ValAddress `json:"voter"` // Voter val address
	SubmitBlock int64          `json:"submit_block"`
}

// NewPrevote returns Prevote object
func NewPrevote(hash string, denom string, voter sdk.ValAddress, submitBlock int64) Prevote {
	return Prevote{
		Hash:        hash,
		Denom:       denom,
		Voter:       voter,
		SubmitBlock: submitBlock,
	}
}

// String implements fmt.Stringer interface
func (pv Prevote) String() string {
	return fmt.Sprintf(`Prevote
	Hash:    %s, 
	Denom:    %s, 
	Voter:    %s, 
	SubmitBlock:    %d`,
		pv.Hash, pv.Denom, pv.Voter, pv.SubmitBlock)
}

// Prevotes is a collection of PreicePrevote
type Prevotes []Prevote

// String implements fmt.Stringer interface
func (v Prevotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// VoteHash computes hash value of Vote
func VoteHash(salt string, exchangeRate sdk.Dec, denom string, voter sdk.ValAddress) ([]byte, error) {
	hash := tmhash.NewTruncated()
	_, err := hash.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", salt, exchangeRate, denom, voter)))
	bz := hash.Sum(nil)
	return bz, err
}

// Vote - struct to store a validator's vote on the exchangeRate of Luna in the denom asset
type Vote struct {
	ExchangeRate sdk.Dec        `json:"exchange_rate"` // Price of Luna in target fiat currency
	Denom        string         `json:"denom"`         // Ticker name of target fiat currency
	Voter        sdk.ValAddress `json:"voter"`         // voter val address of validator
}

// NewVote creates a Vote instance
func NewVote(exchangeRate sdk.Dec, denom string, voter sdk.ValAddress) Vote {
	return Vote{
		ExchangeRate: exchangeRate,
		Denom:        denom,
		Voter:        voter,
	}
}

func (v Vote) getPower(ctx sdk.Context, sk StakingKeeper) int64 {
	validator := sk.Validator(ctx, v.Voter)
	if validator == nil {
		return 0
	}

	return validator.GetConsensusPower()
}

// String implements fmt.Stringer interface
func (v Vote) String() string {
	return fmt.Sprintf(`Vote
	Denom:    %s, 
	Voter:    %s, 
	ExchangeRate:    %s`,
		v.Denom, v.Voter, v.ExchangeRate)
}

// Votes is a collection of Vote
type Votes []Vote

// String implements fmt.Stringer interface
func (v Votes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
