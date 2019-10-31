package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgDelegateConsent{}
	_ sdk.Msg = &MsgPrevote{}
	_ sdk.Msg = &MsgVote{}
)

//-------------------------------------------------
//-------------------------------------------------

// MsgPrevote - struct for prevoting on the Vote.
// The purpose of prevote is to hide vote exchangeRate with hash
// which is formatted as hex string in SHA256("salt:exchangeRate:denom:voter")
type MsgPrevote struct {
	Hash      string         `json:"hash" yaml:"hash"` // hex string
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgPrevote creates a MsgPrevote instance
func NewMsgPrevote(VoteHash string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgPrevote {
	return MsgPrevote{
		Hash:      VoteHash,
		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgPrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgPrevote) Type() string { return "Prevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgPrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgPrevote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgPrevote) ValidateBasic() sdk.Error {

	if bz, err := hex.DecodeString(msg.Hash); len(bz) != tmhash.TruncatedSize || err != nil {
		return ErrInvalidHashLength(DefaultCodespace, len(bz))
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
func (msg MsgPrevote) String() string {
	return fmt.Sprintf(`MsgVote
	hash:         %s,
	feeder:       %s, 
	validator:    %s, 
	denom:        %s`,
		msg.Hash, msg.Feeder, msg.Validator, msg.Denom)
}

// MsgVote - struct for voting on the exchangeRate of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective exchangeRate of Luna in USD is 10.39, that's
// what the exchangeRate field would be, and if 1213.34 for KRW, same.
type MsgVote struct {
	Price     sdk.Dec        `json:"exchangeRate" yaml:"exchangeRate"` // the effective exchangeRate of Luna in {Denom}
	Salt      string         `json:"salt" yaml:"salt"`
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgVote creates a MsgVote instance
func NewMsgVote(exchangeRate sdk.Dec, salt string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgVote {
	return MsgVote{
		Price:     exchangeRate,
		Salt:      salt,
		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgVote) Type() string { return "Vote" }

// GetSignBytes implements sdk.Msg
func (msg MsgVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic implements sdk.Msg
func (msg MsgVote) ValidateBasic() sdk.Error {

	if len(msg.Denom) == 0 {
		return ErrUnknownDenomination(DefaultCodespace, "")
	}

	if msg.Feeder.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.ExchangeRate.LTE(sdk.ZeroDec()) {
		return ErrInvalidPrice(DefaultCodespace, msg.ExchangeRate)
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgVote) String() string {
	return fmt.Sprintf(`MsgVote
	exchangeRate:      %s,
	salt:       %s,
	feeder:     %s, 
	validator:  %s, 
	denom:      %s`,
		msg.ExchangeRate, msg.Salt, msg.Feeder, msg.Validator, msg.Denom)
}

// MsgDelegateConsent - struct for delegating oracle voting rights to another address.
type MsgDelegateConsent struct {
	Operator  sdk.ValAddress `json:"operator" yaml:"operator"`
	Delegatee sdk.AccAddress `json:"delegatee" yaml:"delegatee"`
}

// NewMsgDelegateConsent creates a MsgDelegateConsent instance
func NewMsgDelegateConsent(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) MsgDelegateConsent {
	return MsgDelegateConsent{
		Operator:  operatorAddress,
		Delegatee: feederAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgDelegateConsent) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateConsent) Type() string { return "delegatefeeder" }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateConsent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateConsent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDelegateConsent) ValidateBasic() sdk.Error {
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	if msg.Delegatee.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgDelegateConsent) String() string {
	return fmt.Sprintf(`MsgDelegateConsent
	operator:    %s, 
	delegatee:   %s`,
		msg.Operator, msg.Delegatee)
}
