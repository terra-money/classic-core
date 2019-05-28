package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgPriceFeed{}, "oracle/MsgPriceFeed", nil)
	cdc.RegisterConcrete(MsgDelegateFeederPermission{}, "oracle/MsgDelegateFeederPermission", nil)

	cdc.RegisterConcrete(&PriceBallot{}, "oracle/PriceBallot", nil)
	cdc.RegisterConcrete(&PriceVote{}, "oracle/PriceVote", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
