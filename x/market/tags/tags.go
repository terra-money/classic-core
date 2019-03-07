package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Market tags
var (
	ActionSwap = []byte("swap-coins")

	Action      = sdk.TagAction
	OfferDenom  = "offer-denom"
	OfferAmount = "offer-amount"
	AskDenom    = "ask-denom"
	AskAmount   = "ask-amount"
	Trader      = "trader"
)
