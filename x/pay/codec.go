// Pay TODO - mandatory update

package pay

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(bank.MsgSend{}, "pay/MsgSend", nil)
	cdc.RegisterConcrete(bank.MsgMultiSend{}, "pay/MsgMultiSend", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
