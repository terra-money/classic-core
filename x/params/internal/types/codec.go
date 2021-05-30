package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/terra-money/core/x/gov"
)

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(params.ParameterChangeProposal{}, "params/ParameterChangeProposal", nil)
}

// ModuleCdc is the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()

	gov.RegisterProposalTypeCodec(params.ParameterChangeProposal{}, "params/ParameterChangeProposal")
}
