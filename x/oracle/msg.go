package oracle

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

//-------------------------------------------------
//-------------------------------------------------

// MsgPriceFeed - struct for voting on the price of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective price of Luna in USD is 10.39, that's
// what the price field would be, and if 1213.34 for KRW, same.
// (Hash,Denom,Feeder,Validator) are the contents for prevote of price feed msg,
// in vote period feeder should submit proof price to verify prevote hash
type MsgPriceFeed struct {
	Hash      string         `json:"hash"` // hex string
	Denom     string         `json:"denom"`
	Feeder    sdk.AccAddress `json:"feeder"`
	Validator sdk.ValAddress `json:"validator"`

	Salt  string  `json:"salt"`
	Price sdk.Dec `json:"price"` // the effective price of Luna in {Denom}
}

// NewMsgPriceFeed creates a MsgPriceFeed instance
// price and salt are for prevote hash
func NewMsgPriceFeed(VoteHash string, salt string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress, price sdk.Dec) MsgPriceFeed {
	return MsgPriceFeed{
		Hash: VoteHash,
		Salt: salt,

		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
		Price:     price,
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
	if len(msg.Hash) > 0 {
		if bz, err := hex.DecodeString(msg.Hash); len(bz) != tmhash.TruncatedSize || err != nil {
			return ErrInvalidHashLength(DefaultCodespace, len([]byte(msg.Hash)))
		}
	} else if msg.Price.Equal(sdk.ZeroDec()) {
		return ErrInvalidMsgFormat(DefaultCodespace, "cannot skip both of hash and price")
	}

	if len(msg.Denom) == 0 {
		return ErrUnknownDenomination(DefaultCodespace, "")
	}

	if msg.Feeder.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if !msg.Price.Equal(sdk.ZeroDec()) {
		if len(msg.Salt) < 1 || len(msg.Salt) > 4 {
			return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
		}
	}

	// For initial prevote, the price is not required
	// if msg.Price.LTE(sdk.ZeroDec()) {
	// 	return ErrInvalidPrice(DefaultCodespace, msg.Price)
	// }

	return nil
}

// String Implements Msg
func (msg MsgPriceFeed) String() string {
	return fmt.Sprintf(`MsgPriceFeed
	hash: %s,
	feeder:    %s, 
	validator:    %s, 
	denom:     %s, 
	price:     %s`,
		msg.Hash, msg.Feeder, msg.Validator, msg.Denom, msg.Price)
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
