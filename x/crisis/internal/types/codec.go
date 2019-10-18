package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/crisis"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(crisis.MsgVerifyInvariant{}, "crisis/MsgVerifyInvariant", nil)
}

// ModuleCdc is generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
