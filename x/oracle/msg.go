package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-------------------------------------------------
//-------------------------------------------------

// MsgPriceFeed - struct for voting on the price of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective price of Luna in USD is 10.39, that's
// what the price field would be, and if 1213.34 for KRW, same.
type MsgPriceFeed struct {
	Denom     string         `json:"denom"`
	Price     sdk.Dec        `json:"price"` // the effective price of Luna in {Denom}
	Feeder    sdk.AccAddress `json:"feeder"`
	Validator sdk.ValAddress `json:"validator"`
}

// NewMsgPriceFeed creates a MsgPriceFeed instance
func NewMsgPriceFeed(denom string, price sdk.Dec, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgPriceFeed {
	return MsgPriceFeed{
		Denom:     denom,
		Price:     price,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route Implements Msg
func (msg MsgPriceFeed) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgPriceFeed) Type() string { return "pricefeed" }

// GetSignBytes implements sdk.Msg
func (msg MsgPriceFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgPriceFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgPriceFeed) ValidateBasic() sdk.Error {
	if len(msg.Denom) == 0 {
		return ErrUnknownDenomination(DefaultCodespace, "")
	}

	if msg.Feeder.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Price.LTE(sdk.ZeroDec()) {
		return ErrInvalidPrice(DefaultCodespace, msg.Price)
	}

	return nil
}

// String Implements Msg
func (msg MsgPriceFeed) String() string {
	return fmt.Sprintf(`MsgPriceFeed
	feeder:    %s, 
	validator:    %s, 
	denom:     %s, 
	price:     %s`,
		msg.Feeder, msg.Validator, msg.Denom, msg.Price)
}

// MsgDelegateFeederPermission - struct for delegating oracle voting rights to another address.
type MsgDelegateFeederPermission struct {
	Operator     sdk.ValAddress `json:"operator"`
	FeedDelegate sdk.AccAddress `json:"feed_delegate"`
}

// NewMsgDelegateFeederPermission creates a MsgDelegateFeederPermission instance
func NewMsgDelegateFeederPermission(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) MsgDelegateFeederPermission {
	return MsgDelegateFeederPermission{
		Operator:     operatorAddress,
		FeedDelegate: feederAddress,
	}
}

// Route Implements Msg
func (msg MsgDelegateFeederPermission) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateFeederPermission) Type() string { return "delegatefeeder" }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeederPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateFeederPermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Operator)}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgDelegateFeederPermission) ValidateBasic() sdk.Error {
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	if msg.FeedDelegate.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	return nil
}

// String Implements Msg
func (msg MsgDelegateFeederPermission) String() string {
	return fmt.Sprintf(`MsgDelegateFeederPermission
	operator:    %s, 
	feed_delegate:     %s`,
		msg.Operator, msg.FeedDelegate)
}
