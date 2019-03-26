package treasury

import (
	"terra/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&types.Claim{}, "treasury/Claim", nil)
}

func init() {
	RegisterCodec(cdc)
}
