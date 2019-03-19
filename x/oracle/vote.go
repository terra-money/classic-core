package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceVote - struct to store a validator's vote on the price of Luna in the denom asset
type PriceVote struct {
	Price sdk.Dec        `json:"price"` // Price of Luna in target fiat currency
	Denom string         `json:"denom"` // Ticker name of target fiat currency
	Power sdk.Int        `json:"power"` // Total bonded tokens of validator
	Voter sdk.AccAddress `json:"voter"` // account address of validator
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(price sdk.Dec, denom string, power sdk.Int, voter sdk.AccAddress) PriceVote {
	return PriceVote{
		Price: price,
		Denom: denom,
		Power: power,
		Voter: voter,
	}
}

// String implements fmt.Stringer
func (pv PriceVote) String() string {
	return fmt.Sprintf(`PriceVote
	Denom:    %s, 
	Voter:    %s, 
	Power:    %s, 
	Price:    %s`,
		pv.Denom, pv.Voter, pv.Power, pv.Price)
}
