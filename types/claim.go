package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ClaimClass is used to categorize type of Claim
type ClaimClass byte

const (
	// OracleClaimClass represents oracle claim type
	OracleClaimClass ClaimClass = iota
	// BudgetClaimClass represents budget claim type
	BudgetClaimClass
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

func (c Claim) getClassString() string {
	switch c.Class {
	case OracleClaimClass:
		return "oracle"
	case BudgetClaimClass:
		return "budget"
	}
	return "unknown"
}

func (c Claim) String() string {
	return fmt.Sprintf(`Claim
	Class: %v
	Weight: %v
	Recipient: %v`, c.getClassString(), c.Weight, c.Recipient)
}
