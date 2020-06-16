package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines internal Module Codec
var ModuleCdc = codec.New()

// RegisterCodec concretes types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSwap{}, "market/MsgSwap", nil)
	cdc.RegisterConcrete(MsgSwapSend{}, "market/MsgSwapSend", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
