package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddBurnTaxExemptionAddress    = "AddBurnTaxExemptionAddress"
	ProposalTypeRemoveBurnTaxExemptionAddress = "RemoveBurnTaxExemptionAddress"
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeAddBurnTaxExemptionAddress)
	govv1beta1.ModuleCdc.LegacyAmino.RegisterConcrete(&AddBurnTaxExemptionAddressProposal{}, "treasury/AddBurnTaxExemptionAddressProposal", nil)
	govv1beta1.RegisterProposalType(ProposalTypeRemoveBurnTaxExemptionAddress)
	govv1beta1.ModuleCdc.LegacyAmino.RegisterConcrete(&RemoveBurnTaxExemptionAddressProposal{}, "treasury/RemoveBurnTaxExemptionAddressProposal", nil)
}

var (
	_ govv1beta1.Content = &AddBurnTaxExemptionAddressProposal{}
	_ govv1beta1.Content = &RemoveBurnTaxExemptionAddressProposal{}
)

// ======AddBurnTaxExemptionAddressProposal======

func NewAddBurnTaxExemptionAddressProposal(title, description string, addresses []string) govv1beta1.Content {
	return &AddBurnTaxExemptionAddressProposal{
		Title:       title,
		Description: description,
		Addresses:   addresses,
	}
}

func (p *AddBurnTaxExemptionAddressProposal) GetTitle() string { return p.Title }

func (p *AddBurnTaxExemptionAddressProposal) GetDescription() string { return p.Description }

func (p *AddBurnTaxExemptionAddressProposal) ProposalRoute() string { return RouterKey }

func (p *AddBurnTaxExemptionAddressProposal) ProposalType() string {
	return ProposalTypeAddBurnTaxExemptionAddress
}

func (p AddBurnTaxExemptionAddressProposal) String() string {
	return fmt.Sprintf(`AddBurnTaxExemptionAddressProposal:
	Title:       %s
	Description: %s
	Addresses:   %v
  `, p.Title, p.Description, p.Addresses)
}

func (p *AddBurnTaxExemptionAddressProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(p)
	if err != nil {
		return err
	}

	for _, address := range p.Addresses {
		_, err = sdk.AccAddressFromBech32(address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s: %s", err, address)
		}
	}

	return nil
}

// ======RemoveBurnTaxExemptionAddressProposal======

func NewRemoveBurnTaxExemptionAddressProposal(title, description string, addresses []string) govv1beta1.Content {
	return &RemoveBurnTaxExemptionAddressProposal{
		Title:       title,
		Description: description,
		Addresses:   addresses,
	}
}

func (p *RemoveBurnTaxExemptionAddressProposal) GetTitle() string { return p.Title }

func (p *RemoveBurnTaxExemptionAddressProposal) GetDescription() string { return p.Description }

func (p *RemoveBurnTaxExemptionAddressProposal) ProposalRoute() string { return RouterKey }

func (p *RemoveBurnTaxExemptionAddressProposal) ProposalType() string {
	return ProposalTypeRemoveBurnTaxExemptionAddress
}

func (p RemoveBurnTaxExemptionAddressProposal) String() string {
	return fmt.Sprintf(`RemoveBurnTaxExemptionAddressProposal:
	Title:       %s
	Description: %s
	Addresses:   %v
  `, p.Title, p.Description, p.Addresses)
}

func (p *RemoveBurnTaxExemptionAddressProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(p)
	if err != nil {
		return err
	}

	for _, address := range p.Addresses {
		_, err = sdk.AccAddressFromBech32(address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s: %s", err, address)
		}
	}

	return nil
}
