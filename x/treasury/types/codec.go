package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/market interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&TaxRateUpdateProposal{}, "treasury/TaxRateUpdateProposal", nil)
	cdc.RegisterConcrete(&RewardWeightUpdateProposal{}, "treasury/RewardWeightUpdateProposal", nil)
}

// RegisterInterfaces registers the x/treasury interfaces types with the interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&TaxRateUpdateProposal{},
		&RewardWeightUpdateProposal{},
	)
}

func init() {
	govtypes.RegisterProposalTypeCodec(&TaxRateUpdateProposal{}, "treasury/TaxRateUpdateProposal")
	govtypes.RegisterProposalTypeCodec(&RewardWeightUpdateProposal{}, "treasury/RewardWeightUpdateProposal")
}
