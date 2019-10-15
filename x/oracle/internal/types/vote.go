package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

// PricePrevote - struct to store a validator's prevote on the price of Luna in the denom asset
type PricePrevote struct {
	Hash        string         `json:"hash"`  // Vote hex hash to protect centralize data source problem
	Denom       string         `json:"denom"` // Ticker name of target fiat currency
	Voter       sdk.ValAddress `json:"voter"` // Voter val address
	SubmitBlock int64          `json:"submit_block"`
}

func NewPricePrevote(hash string, denom string, voter sdk.ValAddress, submitBlock int64) PricePrevote {
	return PricePrevote{
		Hash:        hash,
		Denom:       denom,
		Voter:       voter,
		SubmitBlock: submitBlock,
	}
}

// String implements fmt.Stringer interface
func (pp PricePrevote) String() string {
	return fmt.Sprintf(`PricePrevote
	Hash:    %s, 
	Denom:    %s, 
	Voter:    %s, 
	SubmitBlock:    %d`,
		pp.Hash, pp.Denom, pp.Voter, pp.SubmitBlock)
}

// PricePrevotes is a collection of PreicePrevote
type PricePrevotes []PricePrevote

// String implements fmt.Stringer interface
func (v PricePrevotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// VoteHash computes hash value of PriceVote
func VoteHash(salt string, price sdk.Dec, denom string, voter sdk.ValAddress) ([]byte, error) {
	hash := tmhash.NewTruncated()
	_, err := hash.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", salt, price, denom, voter)))
	bz := hash.Sum(nil)
	return bz, err
}

// PriceVote - struct to store a validator's vote on the price of Luna in the denom asset
type PriceVote struct {
	Price sdk.Dec        `json:"price"` // Price of Luna in target fiat currency
	Denom string         `json:"denom"` // Ticker name of target fiat currency
	Voter sdk.ValAddress `json:"voter"` // voter val address of validator
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(price sdk.Dec, denom string, voter sdk.ValAddress) PriceVote {
	return PriceVote{
		Price: price,
		Denom: denom,
		Voter: voter,
	}
}

func (pv PriceVote) getPower(ctx sdk.Context, sk StakingKeeper) int64 {
	validator := sk.Validator(ctx, pv.Voter)
	if validator == nil {
		return 0
	}

	return validator.GetConsensusPower()
}

// String implements fmt.Stringer interface
func (pv PriceVote) String() string {
	return fmt.Sprintf(`PriceVote
	Denom:    %s, 
	Voter:    %s, 
	Price:    %s`,
		pv.Denom, pv.Voter, pv.Price)
}

// PriceVotes is a collection of PriceVote
type PriceVotes []PriceVote

// String implements fmt.Stringer interface
func (v PriceVotes) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
