package v036

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
)

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*v036gov.Content)(nil), nil)
	cdc.RegisterConcrete(v036gov.TextProposal{}, "gov/TextProposal", nil)
}
