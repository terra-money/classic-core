package budget

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitProgram{}, "budget/MsgSubmitProgram", nil)
	cdc.RegisterConcrete(MsgWithdrawProgram{}, "budget/MsgWithdrawProgram", nil)
	cdc.RegisterConcrete(MsgVoteProgram{}, "budget/MsgVoteProgram", nil)

	cdc.RegisterConcrete(&Program{}, "budget/Program", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
