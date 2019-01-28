package treasury

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	// OracleShareID - Share ID for the Oracle module
	OracleShareID = "share_oracle"

	// BudgetShareID - Share ID for the budget module
	BudgetShareID = "share_budget"

	// DebtShareID - Share ID for the Debt module
	DebtShareID = "share_debt"
)

// Share interface. A Share is an abstraction to claim a fixed portion of the Treasury income,
// as denominated by its "Weight." Shareholders need add "claims" indicating to the Treasury
// how their income share should be settled, otherwise the entire income pool ends up in the fee-
// collection pool.
type Share interface {
	ID() string
	GetWeight() sdk.Dec
}

type BaseShare struct {
	Id     string  `json:"id"`
	Weight sdk.Dec `json:"weight"`
}

var _ (Share) = (*BaseShare)(nil)

func NewBaseShare(id string, weight sdk.Dec) BaseShare {
	return BaseShare{
		Id:     id,
		Weight: weight,
	}
}

// GetWeight returns the weight for the BaseShare.
func (bs BaseShare) GetWeight() sdk.Dec {
	return bs.Weight
}

// ID returns the weight for the BaseShare.
func (bs BaseShare) ID() string {
	return bs.Id
}

//------------------------------------
//------------------------------------
//------------------------------------

// Claim interface. A claimant's periodic income
// = total income * share weight / sum(weight claims for share) * weight of claim.
// A share without attached claims has its income directed to the feecollectionpool.
type Claim interface {
	ID() string
	ShareID() string
	GetWeight() sdk.Dec
	Settle(ctx sdk.Context, bk bank.Keeper, allocation sdk.Coins) sdk.Error
}

// BaseClaim is an interface that directs its rewards to an attached bank account.
type BaseClaim struct {
	Claim
	id        string
	shareID   string
	weight    sdk.Dec
	recipient sdk.AccAddress
}

// NewBaseClaim generates a BaseClaim instance.
func NewBaseClaim(shareID string, weight sdk.Dec, recipient sdk.AccAddress) BaseClaim {
	return BaseClaim{
		id:        GenerateBaseClaimID(shareID, recipient),
		shareID:   shareID,
		weight:    weight,
		recipient: recipient,
	}
}

// GenerateBaseClaimID generates an id for a BaseClaim.
func GenerateBaseClaimID(shareID string, recipient sdk.AccAddress) string {
	return fmt.Sprintf("%s:%s", shareID, recipient.String())
}

// GetWeight returns the weight for the BaseClaim.
func (bc BaseClaim) GetWeight() sdk.Dec {
	return bc.weight
}

// ID returns the weight for the BaseClaim.
func (bc BaseClaim) ShareID() string {
	return bc.shareID
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
