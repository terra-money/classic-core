package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc module codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgVote{}, "oracle/MsgVote", nil)
	cdc.RegisterConcrete(MsgPrevote{}, "oracle/MsgPrevote", nil)
	cdc.RegisterConcrete(MsgDelegateConsent{}, "oracle/MsgDelegateConsent", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
