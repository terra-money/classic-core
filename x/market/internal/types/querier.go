package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the oracle Querier
const (
	QuerySwap            = "swap"
	QueryPrevDayIssuance = "lastDayIssuance"
	QueryParameters      = "parameters"
)

// QuerySwapParams for query
// - 'custom/market/swap'
type QuerySwapParams struct {
	OfferCoin sdk.Coin
	AskDenom  string
}

func NewQuerySwapParams(offerCoin sdk.Coin, askDenom string) QuerySwapParams {
	return QuerySwapParams{
		OfferCoin: offerCoin,
		AskDenom:  askDenom,
	}
}
