package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClaimClass byte

const (
	// RouterKey is they name of the treasury module
	RouterKey = "treasury"

	OracleClaimClass ClaimClass = iota
	BudgetClaimClass ClaimClass = iota
)

//------------------------------------
//------------------------------------
//------------------------------------

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	class     ClaimClass     `json:"class"`
	weight    sdk.Int        `json:"weight"`
	recipient sdk.AccAddress `json:"recipient"`
}

// NewClaim generates a Claim instance.
func NewClaim(class ClaimClass, weight sdk.Int, recipient sdk.AccAddress) Claim {
	return Claim{
		class:     class,
		weight:    weight,
		recipient: recipient,
	}
}

// ID returns the id of the claim
func (c Claim) ID() string {
	return fmt.Sprintf("%d:%s", c.class, c.recipient.String())
}

func (c Claim) String() string {
	return fmt.Sprintf("Claim{class: %v, weight: %v, recipient: %v}",
		c.class, c.weight, c.recipient)
}

// Claims is a collection of Claim
type Claims []Claim

func (c Claims) String() (out string) {
	for _, claim := range c {
		out += fmt.Sprintf("\n  %s", claim.String())
	}
	return out
}
