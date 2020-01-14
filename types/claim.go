package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//------------------------------------
//------------------------------------
//------------------------------------

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Weight    sdk.Int        `json:"weight"`
	Recipient sdk.AccAddress `json:"recipient"`
}

// NewClaim generates a Claim instance.
func NewClaim(weight sdk.Int, recipient sdk.AccAddress) Claim {
	return Claim{
		Weight:    weight,
		Recipient: recipient,
	}
}

func (c Claim) String() string {
	return fmt.Sprintf(`Claim
	Weight: %v
	Recipient: %v`, c.Weight, c.Recipient)
}
