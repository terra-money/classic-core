package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

//------------------------------------
//------------------------------------
//------------------------------------

// Claim interface. A claimant's periodic income
// = total income / sum(weight claims) * weight of claim.
type Claim interface {
	ID() string
	GetWeight() sdk.Dec
	Settle(ctx sdk.Context, bk bank.Keeper, allocation sdk.Coins) sdk.Error
}

// BaseClaim is an interface that directs its rewards to an attached bank account.
type BaseClaim struct {
	Claim
	id        string
	weight    sdk.Dec
	recipient sdk.AccAddress
}

// NewBaseClaim generates a BaseClaim instance.
func NewBaseClaim(weight sdk.Dec, recipient sdk.AccAddress) BaseClaim {
	return BaseClaim{
		id:        GenerateBaseClaimID(recipient),
		weight:    weight,
		recipient: recipient,
	}
}

// GenerateBaseClaimID generates an id for a BaseClaim.
// TODO: come up with an actual id system
func GenerateBaseClaimID(recipient sdk.AccAddress) string {
	return recipient.String()
}

// GetWeight returns the weight for the BaseClaim.
func (bc BaseClaim) GetWeight() sdk.Dec {
	return bc.weight
}

// ID returns the weight for the BaseClaim.
func (bc BaseClaim) ID() string {
	return bc.id
}

// Settle the BaseClaim by adsding {allotcation} alloted coins to the recipient account.
func (bc BaseClaim) Settle(ctx sdk.Context, bk bank.Keeper, allocation sdk.Coins) sdk.Error {
	_, _, err := bk.AddCoins(ctx, bc.recipient, allocation)
	return err
}
