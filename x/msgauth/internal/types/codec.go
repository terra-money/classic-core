package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgGrantAuthorization{}, "msgauth/GrantAuthorization", nil)
	cdc.RegisterConcrete(MsgRevokeAuthorization{}, "msgauth/RevokeAuthorization", nil)
	cdc.RegisterConcrete(MsgExecAuthorized{}, "msgauth/ExecAuthorized", nil)
	cdc.RegisterConcrete(SendAuthorization{}, "msgauth/SendAuthorization", nil)
	cdc.RegisterConcrete(GenericAuthorization{}, "msgauth/GenericAuthorization", nil)

	cdc.RegisterInterface((*Authorization)(nil), nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
