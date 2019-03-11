package budget

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(SubmitProgramMsg{}, "budget/SubmitProgramMsg", nil)
	cdc.RegisterConcrete(WithdrawProgramMsg{}, "budget/WithdrawProgramMsg", nil)
	cdc.RegisterConcrete(VoteMsg{}, "budget/VoteMsg", nil)

	cdc.RegisterConcrete(&Program{}, "budget/Program", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
