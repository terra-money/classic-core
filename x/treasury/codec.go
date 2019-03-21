package treasury

import (
	"terra/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&types.Claim{}, "treasury/Claim", nil)
}

func init() {
	RegisterCodec(cdc)
}
