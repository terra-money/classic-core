package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc module codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgExchangeRateVote{}, "oracle/MsgExchangeRateVote", nil)
	cdc.RegisterConcrete(MsgExchangeRatePrevote{}, "oracle/MsgExchangeRatePrevote", nil)
	cdc.RegisterConcrete(MsgDelegateFeedConsent{}, "oracle/MsgDelegateFeedConsent", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRatePrevote{}, "oracle/MsgAggregateExchangeRatePrevote", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRateVote{}, "oracle/MsgAggregateExchangeRateVote", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
