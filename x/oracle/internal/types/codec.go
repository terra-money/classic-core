package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// module codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgPriceVote{}, "oracle/MsgPriceVote", nil)
	cdc.RegisterConcrete(MsgPricePrevote{}, "oracle/MsgPricePrevote", nil)
	cdc.RegisterConcrete(MsgDelegateFeederPermission{}, "oracle/MsgDelegateFeederPermission", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
