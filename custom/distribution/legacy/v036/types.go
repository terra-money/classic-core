package v036

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v036distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v036"
)

// RegisterLegacyAminoCodec nolint
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(v036distr.CommunityPoolSpendProposal{}, "distribution/CommunityPoolSpendProposal", nil)
}
