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
	Weight    int64          `json:"weight"`
	Recipient sdk.ValAddress `json:"recipient"`
}

// NewClaim generates a Claim instance.
func NewClaim(weight int64, recipient sdk.ValAddress) Claim {
	return Claim{
		Weight:    weight,
		Recipient: recipient,
	}
}

// String implements fmt.Stringer interface
func (c Claim) String() string {
	return fmt.Sprintf(`Claim
	Weight: %v
	Recipient: %v`, c.Weight, c.Recipient)
}
