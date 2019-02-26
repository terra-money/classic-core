package treasury

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&Claim{}, "treasury/Claim", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
