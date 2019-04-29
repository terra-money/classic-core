// Pay TODO - mandatory update

package pay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
)

// MsgPay - high level transaction of the pay module
type MsgPay struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Coins       sdk.Coins      `json:"coins"`
}

var _ sdk.Msg = MsgPay{}

// NewMsgPay creates a MsgPay instance
func NewMsgPay(fromAddress sdk.AccAddress, toAddress sdk.AccAddress, coins sdk.Coins) MsgPay {
	return MsgPay{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Coins:       coins,
	}
}

// Route Implements Msg
func (msg MsgPay) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgPay) Type() string { return "pay" }

// GetSignBytes Implements Msg
func (msg MsgPay) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg
func (msg MsgPay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// ValidateBasic Implements Msg
func (msg MsgPay) ValidateBasic() sdk.Error {

	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing payer address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Coins.IsValid() {
		return sdk.ErrInvalidCoins("pay amount is invalid: " + msg.Coins.String())
	}
	if !msg.Coins.IsAllPositive() {
		return sdk.ErrInsufficientCoins("pay amount must be positive")
	}

	return nil
}

// MsgMultiPay is copied type of bank module from cosmos-sdk
type MsgMultiPay bank.MsgMultiSend

var _ sdk.Msg = MsgMultiPay{}

// NewMsgMultiPay - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiPay(in []bank.Input, out []bank.Output) MsgMultiPay {
	return MsgMultiPay{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiPay) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgMultiPay) Type() string { return "multipay" }

// GetSignBytes Implements Msg
func (msg MsgMultiPay) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg
func (msg MsgMultiPay) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// ValidateBasic Implements Msg
func (msg MsgMultiPay) ValidateBasic() sdk.Error {

	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return bank.ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	if len(msg.Outputs) == 0 {
		return bank.ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}

	return bank.ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}
