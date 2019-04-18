package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceVote - struct to store a validator's vote on the price of Luna in the denom asset
type PriceVote struct {
	Price sdk.Dec        `json:"price"` // Price of Luna in target fiat currency
	Denom string         `json:"denom"` // Ticker name of target fiat currency
	Voter sdk.AccAddress `json:"voter"` // account address of validator
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(price sdk.Dec, denom string, voter sdk.AccAddress) PriceVote {
	return PriceVote{
		Price: price,
		Denom: denom,
		Voter: voter,
	}
}

func (pv PriceVote) getPower(ctx sdk.Context, valset sdk.ValidatorSet) (sdk.Int, sdk.Error) {
	valAddr := sdk.ValAddress(pv.Voter)
	if validator := valset.Validator(ctx, valAddr); validator != nil {
		return validator.GetBondedTokens(), nil
	}

	return sdk.ZeroInt(), ErrVoterNotValidator(DefaultCodespace, pv.Voter)
}

// String implements fmt.Stringer
func (pv PriceVote) String() string {
	return fmt.Sprintf(`PriceVote
	Denom:    %s, 
	Voter:    %s, 
	Price:    %s`,
		pv.Denom, pv.Voter, pv.Price)
}
