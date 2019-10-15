package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgDelegateFeederPermission{}
	_ sdk.Msg = &MsgPricePrevote{}
	_ sdk.Msg = &MsgPriceVote{}
)

//-------------------------------------------------
//-------------------------------------------------

// MsgPricePrevote - struct for prevoting on the PriceVote.
// The purpose of prevote is to hide vote price with hash
// which is formatted as hex string in SHA256("salt:price:denom:voter")
type MsgPricePrevote struct {
	Hash      string         `json:"hash" yaml:"hash"` // hex string
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgPricePrevote creates a MsgPricePrevote instance
func NewMsgPricePrevote(VoteHash string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgPricePrevote {
	return MsgPricePrevote{
		Hash:      VoteHash,
		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route Implements Msg
func (msg MsgPricePrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgPricePrevote) Type() string { return "priceprevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgPricePrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgPricePrevote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgPricePrevote) ValidateBasic() sdk.Error {

	if bz, err := hex.DecodeString(msg.Hash); len(bz) != tmhash.TruncatedSize || err != nil {
		return ErrInvalidHashLength(DefaultCodespace, len([]byte(msg.Hash)))
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

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgPricePrevote) String() string {
	return fmt.Sprintf(`MsgPriceVote
	hash:     %s,
	feeder:    %s, 
	validator:    %s, 
	denom:     %s`,
		msg.Hash, msg.Feeder, msg.Validator, msg.Denom)
}

// MsgPriceVote - struct for voting on the price of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective price of Luna in USD is 10.39, that's
// what the price field would be, and if 1213.34 for KRW, same.
type MsgPriceVote struct {
	Price     sdk.Dec        `json:"price" yaml:"price"` // the effective price of Luna in {Denom}
	Salt      string         `json:"salt" yaml:"salt"`
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgPriceVote creates a MsgPriceVote instance
func NewMsgPriceVote(price sdk.Dec, salt string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgPriceVote {
	return MsgPriceVote{
		Price:     price,
		Salt:      salt,
		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route Implements Msg
func (msg MsgPriceVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgPriceVote) Type() string { return "pricevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgPriceVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgPriceVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgPriceVote) ValidateBasic() sdk.Error {

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

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgPriceVote) String() string {
	return fmt.Sprintf(`MsgPriceVote
	price:     %s,
	salt:     %s,
	feeder:    %s, 
	validator:    %s, 
	denom:     %s`,
		msg.Price, msg.Salt, msg.Feeder, msg.Validator, msg.Denom)
}

// MsgDelegateFeederPermission - struct for delegating oracle voting rights to another address.
type MsgDelegateFeederPermission struct {
	Operator  sdk.ValAddress `json:"operator" yaml:"operator"`
	Delegatee sdk.AccAddress `json:"delegatee" yaml:"delegatee"`
}

// NewMsgDelegateFeederPermission creates a MsgDelegateFeederPermission instance
func NewMsgDelegateFeederPermission(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) MsgDelegateFeederPermission {
	return MsgDelegateFeederPermission{
		Operator:  operatorAddress,
		Delegatee: feederAddress,
	}
}

// Route Implements Msg
func (msg MsgDelegateFeederPermission) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateFeederPermission) Type() string { return "delegatefeeder" }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeederPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
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

	if msg.Delegatee.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgDelegateFeederPermission) String() string {
	return fmt.Sprintf(`MsgDelegateFeederPermission
	operator:    %s, 
	delegatee:   %s`,
		msg.Operator, msg.Delegatee)
}
