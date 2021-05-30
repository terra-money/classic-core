package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	msgauthexported "github.com/terra-money/core/x/msgauth/exported"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(bank.MsgSend{}, "bank/MsgSend", nil)
	cdc.RegisterConcrete(bank.MsgMultiSend{}, "bank/MsgMultiSend", nil)
}

// ModuleCdc defines module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()

	msgauthexported.RegisterMsgAuthTypeCodec(bank.MsgSend{}, "bank/MsgSend")
}
