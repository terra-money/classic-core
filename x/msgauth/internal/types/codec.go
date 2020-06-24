package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/bank"
	"github.com/terra-project/core/x/market"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgGrantAuthorization{}, "msgauth/MsgGrantAuthorization", nil)
	cdc.RegisterConcrete(MsgRevokeAuthorization{}, "msgauth/MsgRevokeAuthorization", nil)
	cdc.RegisterConcrete(MsgExecAuthorized{}, "msgauth/MsgExecAuthorized", nil)
	cdc.RegisterConcrete(SendAuthorization{}, "msgauth/SendAuthorization", nil)
	cdc.RegisterConcrete(GenericAuthorization{}, "msgauth/GenericAuthorization", nil)

	cdc.RegisterInterface((*Authorization)(nil), nil)
}

// Need interface to register codec for other module
func init() {
	sdk.RegisterCodec(ModuleCdc)
	bank.RegisterCodec(ModuleCdc)
	market.RegisterCodec(ModuleCdc)
	RegisterCodec(ModuleCdc)
}
