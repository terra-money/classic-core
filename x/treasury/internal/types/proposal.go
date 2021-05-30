package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-money/core/x/gov"
)

const (
	// ProposalTypeTaxRateUpdate defines the type for a TaxRateUpdateProposal
	ProposalTypeTaxRateUpdate = "TaxRateUpdate"

	// ProposalTypeRewardWeightUpdate defines the type for a RewardWeightUpdateProposal
	ProposalTypeRewardWeightUpdate = "RewardWeightUpdate"
)

// Assert TaxRateUpdateProposal implements govtypes.Content at compile-time
var _ gov.Content = TaxRateUpdateProposal{}
var _ gov.Content = RewardWeightUpdateProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeTaxRateUpdate)
	gov.RegisterProposalType(ProposalTypeRewardWeightUpdate)
}

// TaxRateUpdateProposal updates treasury tax-rate
type TaxRateUpdateProposal struct {
	Title       string  `json:"title" yaml:"title"`             // Title of the Proposal
	Description string  `json:"description" yaml:"description"` // Description of the Proposal
	TaxRate     sdk.Dec `json:"tax_rate" yaml:"tax_rate"`       // target TaxRate
}

// NewTaxRateUpdateProposal creates an TaxRateUpdateProposal.
func NewTaxRateUpdateProposal(title, description string, taxRate sdk.Dec) TaxRateUpdateProposal {
	return TaxRateUpdateProposal{title, description, taxRate}
}

// GetTitle returns the title of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of an TaxRateUpdateProposal.
func (TaxRateUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of an TaxRateUpdateProposal.
func (p TaxRateUpdateProposal) ProposalType() string { return ProposalTypeTaxRateUpdate }

// ValidateBasic runs basic stateless validity checks
func (p TaxRateUpdateProposal) ValidateBasic() error {
	err := gov.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if !p.TaxRate.IsPositive() || p.TaxRate.GT(sdk.OneDec()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid tax-rate: "+p.TaxRate.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p TaxRateUpdateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:        %s
  Description:  %s
  TaxRate:      %s
`, p.Title, p.Description, p.TaxRate))
	return b.String()
}

// RewardWeightUpdateProposal update treasury tax-rate
type RewardWeightUpdateProposal struct {
	Title        string  `json:"title" yaml:"title"`                 // Title of the Proposal
	Description  string  `json:"description" yaml:"description"`     // Description of the Proposal
	RewardWeight sdk.Dec `json:"reward_weight" yaml:"reward_weight"` // target RewardWeight
}

// NewRewardWeightUpdateProposal creates an RewardWeightUpdateProposal.
func NewRewardWeightUpdateProposal(title, description string, RewardWeight sdk.Dec) RewardWeightUpdateProposal {
	return RewardWeightUpdateProposal{title, description, RewardWeight}
}

// GetTitle returns the title of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of an RewardWeightUpdateProposal.
func (RewardWeightUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of an RewardWeightUpdateProposal.
func (p RewardWeightUpdateProposal) ProposalType() string { return ProposalTypeRewardWeightUpdate }

// ValidateBasic runs basic stateless validity checks
func (p RewardWeightUpdateProposal) ValidateBasic() error {
	err := gov.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if !p.RewardWeight.IsPositive() || p.RewardWeight.GT(sdk.OneDec()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid reward-weight: "+p.RewardWeight.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p RewardWeightUpdateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:        %s
  Description:  %s
  RewardWeight: %s
`, p.Title, p.Description, p.RewardWeight))
	return b.String()
}
