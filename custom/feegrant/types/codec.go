package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// RegisterLegacyAminoCodec registers the necessary x/authz interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &feegrant.MsgGrantAllowance{}, "feegrant/MsgGrantAllowance")
	legacy.RegisterAminoMsg(cdc, &feegrant.MsgRevokeAllowance{}, "feegrant/MsgRevokeAllowance")

	cdc.RegisterInterface((*feegrant.FeeAllowanceI)(nil), nil)
	cdc.RegisterConcrete(&feegrant.BasicAllowance{}, "feegrant/BasicAllowance", nil)
	cdc.RegisterConcrete(&feegrant.PeriodicAllowance{}, "feegrant/PeriodicAllowance", nil)
	cdc.RegisterConcrete(&feegrant.AllowedMsgAllowance{}, "feegrant/AllowedMsgAllowance", nil)
}
