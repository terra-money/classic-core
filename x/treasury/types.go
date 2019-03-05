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
	id        string
	class     ClaimClass
	weight    sdk.Dec
	recipient sdk.AccAddress
}

// NewClaim generates a Claim instance.
func NewClaim(class ClaimClass, weight sdk.Dec, recipient sdk.AccAddress) Claim {
	return Claim{
		id:        GenerateClaimID(class, recipient),
		class:     class,
		weight:    weight,
		recipient: recipient,
	}
}

// GenerateClaimID generates an id for a Claim.
func GenerateClaimID(class ClaimClass, recipient sdk.AccAddress) string {
	return fmt.Sprintf("%d:%s", class, recipient.String())
}
