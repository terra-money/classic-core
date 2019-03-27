package market

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//--------------------------------------------------------
//--------------------------------------------------------

// MsgSwap contains a swap request
type MsgSwap struct {
	Trader    sdk.AccAddress `json:"trader"`     // Address of the trader
	OfferCoin sdk.Coin       `json:"offer_coin"` // Coin being offered
	AskDenom  string         `json:"ask_denom"`  // Denom of the coin to swap to
}

// NewMsgSwap creates a MsgSwap instance
func NewMsgSwap(traderAddress sdk.AccAddress, offerCoin sdk.Coin, askCoin string) MsgSwap {
	return MsgSwap{
		Trader:    traderAddress,
		OfferCoin: offerCoin,
		AskDenom:  askCoin,
	}
}

// Route Implements Msg
func (msg MsgSwap) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgSwap) Type() string { return "swap" }

// GetSignBytes Implements Msg
func (msg MsgSwap) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners Implements Msg
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trader}
}

// ValidateBasic Implements Msg
func (msg MsgSwap) ValidateBasic() sdk.Error {
	if len(msg.Trader) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Trader.String())
	}

	if msg.OfferCoin.Amount.LT(sdk.ZeroInt()) {
		return ErrInsufficientSwapCoins(DefaultCodespace, msg.OfferCoin.Amount)
	}

	if msg.OfferCoin.Denom == msg.AskDenom {
		return ErrRecursiveSwap(DefaultCodespace, msg.AskDenom)
	}

	return nil
}

// String Implements Msg
func (msg MsgSwap) String() string {
	return fmt.Sprintf(`MsgSwap
	trader:    %s, 
	offer:     %s, 
	ask:       %s`,
		msg.Trader, msg.OfferCoin, msg.AskDenom)
}
