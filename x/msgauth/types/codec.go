package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/msgauth interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*AuthorizationI)(nil), nil)

	cdc.RegisterConcrete(&MsgGrantAuthorization{}, "msgauth/MsgGrantAuthorization", nil)
	cdc.RegisterConcrete(&MsgRevokeAuthorization{}, "msgauth/MsgRevokeAuthorization", nil)
	cdc.RegisterConcrete(&MsgExecAuthorized{}, "msgauth/MsgExecAuthorized", nil)
	cdc.RegisterConcrete(&SendAuthorization{}, "msgauth/SendAuthorization", nil)
	cdc.RegisterConcrete(&GenericAuthorization{}, "msgauth/GenericAuthorization", nil)
}

// RegisterInterfaces registers the x/msgauth interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGrantAuthorization{},
		&MsgRevokeAuthorization{},
		&MsgExecAuthorized{},
	)
	registry.RegisterInterface(
		"terra.msgauth.v1beta1.AuthorizationI",
		(*AuthorizationI)(nil),
		&SendAuthorization{},
		&GenericAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterMsgAuthTypeCodec registers an external msg type defined
// in another module for the internal ModuleCdc. This allows the MsgExecAuthorized
// to be correctly Amino encoded and decoded.
func RegisterMsgAuthTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/gov module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/gov and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
