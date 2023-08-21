package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&AddBurnTaxExemptionAddressProposal{}, "treasury/AddBurnTaxExemptionAddressProposal", nil)
	cdc.RegisterConcrete(&RemoveBurnTaxExemptionAddressProposal{}, "treasury/RemoveBurnTaxExemptionAddressProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddBurnTaxExemptionAddressProposal{},
		&RemoveBurnTaxExemptionAddressProposal{},
	)
}
