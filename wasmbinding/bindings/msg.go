package bindings

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TerraMsg struct {
	Swap     *Swap     `json:"swap,omitempty"`
	SwapSend *SwapSend `json:"swap_send,omitempty"`
}

type Swap struct {
	OfferCoin sdk.Coin `json:"offer_coin"`
	AskDenom  string   `json:"ask_denom"`
}

type SwapSend struct {
	ToAddress string   `json:"to_address"`
	OfferCoin sdk.Coin `json:"offer_coin"`
	AskDenom  string   `json:"ask_denom"`
}
