package v036

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v036params "github.com/cosmos/cosmos-sdk/x/params/legacy/v036"
)

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(v036params.ParameterChangeProposal{}, "params/ParameterChangeProposal", nil)
}
