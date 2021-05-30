package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	govtypes "github.com/terra-money/core/custom/gov/types"
)

// RegisterLegacyAminoCodec registers all necessary param module types with a given LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&proposal.ParameterChangeProposal{}, "params/ParameterChangeProposal", nil)
}

func init() {
	govtypes.RegisterProposalTypeCodec(&proposal.ParameterChangeProposal{}, "params/ParameterChangeProposal")
}
