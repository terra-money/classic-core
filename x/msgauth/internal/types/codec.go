package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ModuleCdc defines internal Module Codec
var ModuleCdc = codec.New()

// RegisterCodec concretes types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgGrantAuthorization{}, "msgauth/MsgGrantAuthorization", nil)
	cdc.RegisterConcrete(MsgRevokeAuthorization{}, "msgauth/MsgRevokeAuthorization", nil)
	cdc.RegisterConcrete(MsgExecAuthorized{}, "msgauth/MsgExecAuthorized", nil)
	cdc.RegisterConcrete(SendAuthorization{}, "msgauth/SendAuthorization", nil)
	cdc.RegisterConcrete(GenericAuthorization{}, "msgauth/GenericAuthorization", nil)

	cdc.RegisterInterface((*Authorization)(nil), nil)
}

// RegisterMsgAuthTypeCodec registers an external msg type defined
// in another module for the internal ModuleCdc. This allows the MsgExecAuthorized
// to be correctly Amino encoded and decoded.
func RegisterMsgAuthTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// Need interface to register codec for other module
func init() {
	sdk.RegisterCodec(ModuleCdc)
	RegisterCodec(ModuleCdc)
}
