package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*authz.Authorization)(nil), nil)

	cdc.RegisterConcrete(&authz.MsgGrant{}, "msgauth/MsgGrantAuthorization", nil)
	cdc.RegisterConcrete(&authz.MsgRevoke{}, "msgauth/MsgRevokeAuthorization", nil)
	cdc.RegisterConcrete(&authz.MsgExec{}, "msgauth/MsgExecAuthorized", nil)
	cdc.RegisterConcrete(&authz.GenericAuthorization{}, "msgauth/GenericAuthorization", nil)
	cdc.RegisterConcrete(&banktypes.SendAuthorization{}, "msgauth/SendAuthorization", nil)
	cdc.RegisterConcrete(&stakingtypes.StakeAuthorization{}, "msgauth/StakeAuthorization", nil)
}
