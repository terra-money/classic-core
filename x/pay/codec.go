// Pay TODO - mandatory update

package pay

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgPay{}, "pay/MsgPay", nil)
	cdc.RegisterConcrete(MsgMultiPay{}, "pay/MsgMultiPay", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
