package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// RouterKey is they name of the treasury module
	RouterKey = "treasury"

	OracleClaimClass = iota
	BudgetClaimClass
)

//------------------------------------
//------------------------------------
//------------------------------------

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	id        string
	class     byte
	weight    sdk.Dec
	recipient sdk.AccAddress
}

// NewClaim generates a Claim instance.
func NewClaim(blockheight int64, class byte, weight sdk.Dec, recipient sdk.AccAddress) Claim {
	return Claim{
		id:        GenerateClaimID(blockheight, class, recipient),
		class:     class,
		weight:    weight,
		recipient: recipient,
	}
}

// GenerateClaimID generates an id for a Claim.
func GenerateClaimID(blockheight int64, class byte, recipient sdk.AccAddress) string {
	return fmt.Sprintf("%d:%d:%s", blockheight, class, recipient.String())
}

// Settle the Claim by adsding {allotcation} alloted coins to the recipient account.
func (c Claim) Settle(ctx sdk.Context, k Keeper, allocation sdk.Coins) sdk.Error {
	_, _, err := k.pk.AddCoins(ctx, c.recipient, allocation)
	return err
}
