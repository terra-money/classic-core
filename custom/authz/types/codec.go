package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &authz.MsgGrant{}, "msgauth/MsgGrantAuthorization")
	legacy.RegisterAminoMsg(cdc, &authz.MsgRevoke{}, "msgauth/MsgRevokeAuthorization")
	legacy.RegisterAminoMsg(cdc, &authz.MsgExec{}, "msgauth/MsgExecAuthorized")

	cdc.RegisterInterface((*authz.Authorization)(nil), nil)
	cdc.RegisterConcrete(&authz.GenericAuthorization{}, "msgauth/GenericAuthorization", nil)
}
