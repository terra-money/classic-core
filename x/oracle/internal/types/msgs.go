package types

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgDelegateFeedConsent{}
	_ sdk.Msg = &MsgExchangeRatePrevote{}
	_ sdk.Msg = &MsgExchangeRateVote{}
	_ sdk.Msg = &MsgAggregateExchangeRatePrevote{}
	_ sdk.Msg = &MsgAggregateExchangeRateVote{}
)

//-------------------------------------------------
//-------------------------------------------------

// Deprecated: normal prevote and vote will be deprecated after columbus-4
// MsgExchangeRatePrevote - struct for prevoting on the ExchangeRateVote.
// The purpose of prevote is to hide vote exchange rate with hash
// which is formatted as hex string in SHA256("{salt}:{exchange_rate}:{denom}:{voter}")
type MsgExchangeRatePrevote struct {
	Hash      VoteHash       `json:"hash" yaml:"hash"`
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgExchangeRatePrevote creates a MsgExchangeRatePrevote instance
func NewMsgExchangeRatePrevote(hash VoteHash, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgExchangeRatePrevote {
	return MsgExchangeRatePrevote{
		Hash:      hash,
		Denom:     denom,
		Feeder:    feederAddress,
		Validator: valAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgExchangeRatePrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgExchangeRatePrevote) Type() string { return "exchangerateprevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgExchangeRatePrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgExchangeRatePrevote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgExchangeRatePrevote) ValidateBasic() error {

	if len(msg.Hash) != tmhash.TruncatedSize {
		return ErrInvalidHashLength
	}

	if len(msg.Denom) == 0 {
		return ErrUnknowDenom
	}

	if msg.Feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgExchangeRatePrevote) String() string {
	return fmt.Sprintf(`MsgExchangeRateVote
	hash:         %s,
	feeder:       %s, 
	validator:    %s, 
	denom:        %s`,
		msg.Hash, msg.Feeder, msg.Validator, msg.Denom)
}

// Deprecated: normal prevote and vote will be deprecated after columbus-4
// MsgExchangeRateVote - struct for voting on the exchange rate of Luna denominated in various Terra assets.
// For example, if the validator believes that the effective exchange rate of Luna in USD is 10.39, that's
// what the exchange rate field would be, and if 1213.34 for KRW, same.
type MsgExchangeRateVote struct {
	ExchangeRate sdk.Dec        `json:"exchange_rate" yaml:"exchange_rate"` // the effective rate of Luna in {Denom}
	Salt         string         `json:"salt" yaml:"salt"`
	Denom        string         `json:"denom" yaml:"denom"`
	Feeder       sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator    sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgExchangeRateVote creates a MsgExchangeRateVote instance
func NewMsgExchangeRateVote(rate sdk.Dec, salt string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgExchangeRateVote {
	return MsgExchangeRateVote{
		ExchangeRate: rate,
		Salt:         salt,
		Denom:        denom,
		Feeder:       feederAddress,
		Validator:    valAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgExchangeRateVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgExchangeRateVote) Type() string { return "exchangeratevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgExchangeRateVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgExchangeRateVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic implements sdk.Msg
func (msg MsgExchangeRateVote) ValidateBasic() error {

	if len(msg.Denom) == 0 {
		return ErrUnknowDenom
	}

	if msg.Feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	// Check overflow bit length
	if msg.ExchangeRate.BitLen() > 100+sdk.DecimalPrecisionBits {
		return sdkerrors.Wrap(ErrInvalidExchangeRate, msg.ExchangeRate.String())
	}

	if l := len(msg.Salt); l > 4 || l < 1 {
		return sdkerrors.Wrap(ErrInvalidSaltLength, strconv.FormatInt(int64(l), 10))
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgExchangeRateVote) String() string {
	return fmt.Sprintf(`MsgExchangeRateVote
	exchange_rate:      %s,
	salt:               %s,
	feeder:             %s, 
	validator:          %s, 
	denom:              %s`,
		msg.ExchangeRate, msg.Salt, msg.Feeder, msg.Validator, msg.Denom)
}

// MsgDelegateFeedConsent - struct for delegating oracle voting rights to another address.
type MsgDelegateFeedConsent struct {
	Operator sdk.ValAddress `json:"operator" yaml:"operator"`
	Delegate sdk.AccAddress `json:"delegate" yaml:"delegate"`
}

// NewMsgDelegateFeedConsent creates a MsgDelegateFeedConsent instance
func NewMsgDelegateFeedConsent(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) MsgDelegateFeedConsent {
	return MsgDelegateFeedConsent{
		Operator: operatorAddress,
		Delegate: feederAddress,
	}
}

// Route implements sdk.Msg
func (msg MsgDelegateFeedConsent) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateFeedConsent) Type() string { return "delegatefeeder" }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDelegateFeedConsent) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	if msg.Delegate.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid delegate address")
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgDelegateFeedConsent) String() string {
	return fmt.Sprintf(`MsgDelegateFeedConsent
	operator:    %s, 
	delegate:    %s`,
		msg.Operator, msg.Delegate)
}

// MsgAggregateExchangeRatePrevote - struct for aggregate prevoting on the ExchangeRateVote.
// The purpose of aggregate prevote is to hide vote exchange rates with hash
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")
type MsgAggregateExchangeRatePrevote struct {
	Hash      AggregateVoteHash `json:"hash" yaml:"hash"`
	Feeder    sdk.AccAddress    `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress    `json:"validator" yaml:"validator"`
}

// NewMsgAggregateExchangeRatePrevote returns MsgAggregateExchangeRatePrevote instance
func NewMsgAggregateExchangeRatePrevote(hash AggregateVoteHash, feeder sdk.AccAddress, validator sdk.ValAddress) MsgAggregateExchangeRatePrevote {
	return MsgAggregateExchangeRatePrevote{
		Hash:      hash,
		Feeder:    feeder,
		Validator: validator,
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Type() string { return "aggregateexchangerateprevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) ValidateBasic() error {

	if len(msg.Hash) != tmhash.TruncatedSize {
		return ErrInvalidHashLength
	}

	if msg.Feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgAggregateExchangeRatePrevote) String() string {
	return fmt.Sprintf(`MsgAggregateExchangeRateVote
	hash:         %s,
	feeder:       %s, 
	validator:    %s`,
		msg.Hash, msg.Feeder, msg.Validator)
}

// MsgAggregateExchangeRateVote - struct for voting on the exchange rates of Luna denominated in various Terra assets.
type MsgAggregateExchangeRateVote struct {
	Salt          string         `json:"salt" yaml:"salt"`
	ExchangeRates string         `json:"exchange_rates" yaml:"exchange_rates"` // comma separated dec coins
	Feeder        sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator     sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgAggregateExchangeRateVote returns MsgAggregateExchangeRateVote instance
func NewMsgAggregateExchangeRateVote(salt string, exchangeRates string, feeder sdk.AccAddress, validator sdk.ValAddress) MsgAggregateExchangeRateVote {
	return MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRates,
		Feeder:        feeder,
		Validator:     validator,
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Type() string { return "aggregateexchangeratevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Feeder}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) ValidateBasic() error {

	if msg.Feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	if l := len(msg.ExchangeRates); l == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "must provide at least one oracle exchange rate")
	} else if l > 4096 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "exchange rates string can not exceed 4096 characters")
	}

	exchangeRateTuples, err := ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "failed to parse exchange rates string cause:"+err.Error())
	}

	for _, tuple := range exchangeRateTuples {
		// Check overflow bit length
		if tuple.ExchangeRate.BitLen() > 100+sdk.DecimalPrecisionBits {
			return sdkerrors.Wrap(ErrInvalidExchangeRate, "overflow")
		}
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return sdkerrors.Wrap(ErrInvalidSaltLength, "salt length must be [1, 4]")
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgAggregateExchangeRateVote) String() string {
	return fmt.Sprintf(`MsgAggregateExchangeRateVote
	exchangerate:      %s,
	salt:              %s,
	feeder:            %s, 
	validator:         %s`,
		msg.ExchangeRates, msg.Salt, msg.Feeder, msg.Validator)
}
