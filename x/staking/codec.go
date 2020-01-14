package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(staking.MsgCreateValidator{}, "staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(staking.MsgEditValidator{}, "staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(staking.MsgDelegate{}, "staking/MsgDelegate", nil)
	cdc.RegisterConcrete(staking.MsgUndelegate{}, "staking/MsgUndelegate", nil)
	cdc.RegisterConcrete(staking.MsgBeginRedelegate{}, "staking/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
