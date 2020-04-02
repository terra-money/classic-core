package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
func (msg MsgExchangeRatePrevote) ValidateBasic() sdk.Error {

	if len(msg.Hash) != tmhash.TruncatedSize {
		return ErrInvalidHashLength(DefaultCodespace, len(msg.Hash))
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
func (msg MsgExchangeRateVote) ValidateBasic() sdk.Error {

	if len(msg.Denom) == 0 {
		return ErrUnknownDenomination(DefaultCodespace, "")
	}

	if msg.Feeder.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	// Check overflow bit length
	if msg.ExchangeRate.BitLen() > 100+sdk.DecimalPrecisionBits {
		return ErrInvalidExchangeRate(DefaultCodespace, msg.ExchangeRate)
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgExchangeRateVote) String() string {
	return fmt.Sprintf(`MsgExchangeRateVote
	exchangerate:      %s,
	salt:       %s,
	feeder:     %s, 
	validator:  %s, 
	denom:      %s`,
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
func (msg MsgDelegateFeedConsent) ValidateBasic() sdk.Error {
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	if msg.Delegate.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Operator.String())
	}

	return nil
}

// String implements fmt.Stringer interface
func (msg MsgDelegateFeedConsent) String() string {
	return fmt.Sprintf(`MsgDelegateFeedConsent
	operator:    %s, 
	delegate:   %s`,
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
func (msg MsgAggregateExchangeRatePrevote) ValidateBasic() sdk.Error {

	if len(msg.Hash) != tmhash.TruncatedSize {
		return ErrInvalidHashLength(DefaultCodespace, len(msg.Hash))
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
func (msg MsgAggregateExchangeRateVote) ValidateBasic() sdk.Error {

	if msg.Feeder.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Feeder.String())
	}

	if l := len(msg.ExchangeRates); l == 0 {
		return sdk.ErrUnknownRequest("must provide at least one oracle exchange rate")
	} else if l > 4096 {
		return sdk.ErrInternal("exchange rates string can not exceed 512 character")
	}

	exchangeRateTuples, err := ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return sdk.ErrInvalidCoins(err.Error())
	}

	for _, tuple := range exchangeRateTuples {
		// Check overflow bit length
		if tuple.ExchangeRate.BitLen() > 100+sdk.DecimalPrecisionBits {
			return ErrInvalidExchangeRate(DefaultCodespace, tuple.ExchangeRate)
		}
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return ErrInvalidSaltLength(DefaultCodespace, len(msg.Salt))
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
