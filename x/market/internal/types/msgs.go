package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgSwap{}
)

//--------------------------------------------------------
//--------------------------------------------------------

// MsgSwap contains a swap request
type MsgSwap struct {
	Trader    sdk.AccAddress `json:"trader" yaml:"trader"`         // Address of the trader
	OfferCoin sdk.Coin       `json:"offer_coin" yaml:"offer_coin"` // Coin being offered
	AskDenom  string         `json:"ask_denom" yaml:"ask_denom"`   // Denom of the coin to swap to
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trader}
}

// ValidateBasic Implements Msg
func (msg MsgSwap) ValidateBasic() error {
	if len(msg.Trader) == 0 {
		return sdkerrors.ErrInvalidAddress
	}

	if msg.OfferCoin.Amount.LTE(sdk.ZeroInt()) || msg.OfferCoin.Amount.BigInt().BitLen() > 100 {
		return sdkerrors.Wrap(ErrInvalidOfferCoin, msg.OfferCoin.Amount.String())
	}

	if msg.OfferCoin.Denom == msg.AskDenom {
		return sdkerrors.Wrap(ErrRecursiveSwap, msg.AskDenom)
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgSwap) String() string {
	return fmt.Sprintf(`MsgSwap
	trader:    %s, 
	offer:     %s, 
	ask:       %s`,
		msg.Trader, msg.OfferCoin, msg.AskDenom)
}
