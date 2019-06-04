package slashing

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(slashing.MsgUnjail{}, "slashing/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
