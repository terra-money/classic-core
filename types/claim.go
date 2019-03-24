package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClaimClass byte

const (
	OracleClaimClass ClaimClass = iota
	BudgetClaimClass ClaimClass = iota
)

//------------------------------------
//------------------------------------
//------------------------------------

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Class     ClaimClass     `json:"class"`
	Weight    sdk.Int        `json:"weight"`
	Recipient sdk.AccAddress `json:"recipient"`
}

// NewClaim generates a Claim instance.
func NewClaim(class ClaimClass, weight sdk.Int, recipient sdk.AccAddress) Claim {
	return Claim{
		Class:     class,
		Weight:    weight,
		Recipient: recipient,
	}
}

// ID returns the id of the claim
func (c Claim) ID() string {
	return fmt.Sprintf("%d:%s", c.Class, c.Recipient.String())
}

func (c Claim) String() string {
	return fmt.Sprintf("Claim{class: %v, weight: %v, recipient: %v}",
		c.Class, c.Weight, c.Recipient)
}
