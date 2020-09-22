package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgSwap{}
	_ sdk.Msg = &MsgSwapSend{}
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
	if msg.Trader.Empty() {
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

// MsgSwapSend contains a swap request
type MsgSwapSend struct {
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"` // Address of the offer coin payer
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`     // Address of the recipient
	OfferCoin   sdk.Coin       `json:"offer_coin" yaml:"offer_coin"`     // Coin being offered
	AskDenom    string         `json:"ask_denom" yaml:"ask_denom"`       // Denom of the coin to swap to
}

// NewMsgSwapSend conducts market swap and send all the result coins to recipient
func NewMsgSwapSend(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, offerCoin sdk.Coin, askCoin string) MsgSwapSend {
	return MsgSwapSend{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		OfferCoin:   offerCoin,
		AskDenom:    askCoin,
	}
}

// Route Implements Msg
func (msg MsgSwapSend) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgSwapSend) Type() string { return "swapsend" }

// GetSignBytes Implements Msg
func (msg MsgSwapSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg
func (msg MsgSwapSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// ValidateBasic Implements Msg
func (msg MsgSwapSend) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.ErrInvalidAddress
	}

	if msg.ToAddress.Empty() {
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
func (msg MsgSwapSend) String() string {
	return fmt.Sprintf(`MsgSwapSend
	fromAddress:    %s,
	toAddress:      %s, 
	offer:          %s, 
	ask:            %s`,
		msg.FromAddress, msg.ToAddress, msg.OfferCoin, msg.AskDenom)
}
