package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeTaxRateUpdate defines the type for a TaxRateUpdateProposal
	ProposalTypeTaxRateUpdate = "TaxRateUpdate"

	// ProposalTypeRewardWeightUpdate defines the type for a RewardWeightUpdateProposal
	ProposalTypeRewardWeightUpdate = "RewardWeightUpdate"
)

// Assert TaxRateUpdateProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &TaxRateUpdateProposal{}
var _ govtypes.Content = &RewardWeightUpdateProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeTaxRateUpdate)
	govtypes.RegisterProposalType(ProposalTypeRewardWeightUpdate)
}

// NewTaxRateUpdateProposal creates an TaxRateUpdateProposal.
func NewTaxRateUpdateProposal(title, description string, taxRate sdk.Dec) *TaxRateUpdateProposal {
	return &TaxRateUpdateProposal{title, description, taxRate}
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
func (p *TaxRateUpdateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.TaxRate.IsNegative() || p.TaxRate.GT(sdk.OneDec()) {
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

// NewRewardWeightUpdateProposal creates an RewardWeightUpdateProposal.
func NewRewardWeightUpdateProposal(title, description string, taxRate sdk.Dec) *RewardWeightUpdateProposal {
	return &RewardWeightUpdateProposal{title, description, taxRate}
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
func (p *RewardWeightUpdateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.RewardWeight.IsNegative() || p.RewardWeight.GT(sdk.OneDec()) {
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
