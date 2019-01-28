package treasury

import "github.com/cosmos/cosmos-sdk/codec"

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Share)(nil), nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
