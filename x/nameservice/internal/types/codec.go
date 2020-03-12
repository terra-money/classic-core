package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc module codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgOpenAuction{}, "nameservice/MsgOpenAuction", nil)
	cdc.RegisterConcrete(MsgBidAuction{}, "nameservice/MsgBidAuction", nil)
	cdc.RegisterConcrete(MsgRevealBid{}, "nameservice/MsgRevealBid", nil)
	cdc.RegisterConcrete(MsgRenewRegistry{}, "nameservice/MsgRenewRegistry", nil)
	cdc.RegisterConcrete(MsgUpdateOwner{}, "nameservice/MsgUpdateOwner", nil)
	cdc.RegisterConcrete(MsgRegisterSubName{}, "nameservice/MsgRegisterSubName", nil)
	cdc.RegisterConcrete(MsgUnregisterSubName{}, "nameservice/MsgUnregisterSubName", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
