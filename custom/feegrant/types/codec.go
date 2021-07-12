package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*feegrant.FeeAllowanceI)(nil), nil)

	cdc.RegisterConcrete(&feegrant.MsgGrantAllowance{}, "feegrant/MsgGrantAllowance", nil)
	cdc.RegisterConcrete(&feegrant.MsgRevokeAllowance{}, "feegrant/MsgRevokeAllowance", nil)
	cdc.RegisterConcrete(&feegrant.BasicAllowance{}, "feegrant/BasicAllowance", nil)
	cdc.RegisterConcrete(&feegrant.PeriodicAllowance{}, "feegrant/PeriodicAllowance", nil)
	cdc.RegisterConcrete(&feegrant.AllowedMsgAllowance{}, "feegrant/AllowedMsgAllowance", nil)
}
