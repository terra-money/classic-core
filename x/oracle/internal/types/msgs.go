package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strings"
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

// MsgExchangeRatePrevote - struct for prevoting on the ExchangeRateVote.
// The purpose of prevote is to hide vote exchange rate with hash
// which is formatted as hex string in SHA256("{salt}:{exchange_rate}:{denom}:{voter}")
type MsgExchangeRatePrevote struct {
	Hash      string         `json:"hash" yaml:"hash"` // hex string
	Denom     string         `json:"denom" yaml:"denom"`
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgExchangeRatePrevote creates a MsgExchangeRatePrevote instance
func NewMsgExchangeRatePrevote(VoteHash string, denom string, feederAddress sdk.AccAddress, valAddress sdk.ValAddress) MsgExchangeRatePrevote {
	return MsgExchangeRatePrevote{
		Hash:      VoteHash,
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
func (msg MsgExchangeRatePrevote) String() string {
	return fmt.Sprintf(`MsgExchangeRateVote
	hash:         %s,
	feeder:       %s, 
	validator:    %s, 
	denom:        %s`,
		msg.Hash, msg.Feeder, msg.Validator, msg.Denom)
}

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
	Hash      string         `json:"hash" yaml:"hash"` // hex string
	Feeder    sdk.AccAddress `json:"feeder" yaml:"feeder"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewMsgAggregateExchangeRatePrevote returns MsgAggregateExchangeRatePrevote instance
func NewMsgAggregateExchangeRatePrevote(hash string, feeder sdk.AccAddress, validator sdk.ValAddress) MsgAggregateExchangeRatePrevote {
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

	if bz, err := hex.DecodeString(msg.Hash); len(bz) != tmhash.TruncatedSize || err != nil {
		return ErrInvalidHashLength(DefaultCodespace, len(bz))
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

	if len(msg.ExchangeRates) == 0 {
		return sdk.ErrUnknownRequest("must provide at least one oracle exchange rate")
	}

	exchangeRates, err := ParseDecCoins(msg.ExchangeRates)
	if err != nil {
		return sdk.ErrInvalidCoins(err.Error())
	}

	for _, exchangeRate := range exchangeRates {
		// Check overflow bit length
		if exchangeRate.Amount.BitLen() > 100+sdk.DecimalPrecisionBits {
			return ErrInvalidExchangeRate(DefaultCodespace, exchangeRate.Amount)
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

// ParseDecCoins DecCoin parser to treat non-positive values as valid
func ParseDecCoins(coinsStr string) (sdk.DecCoins, error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	coins := make(sdk.DecCoins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := sdk.ParseDecCoin(coinStr)
		if err != nil {
			return nil, err
		}

		coins[i] = coin
	}

	// sort coins for determinism
	coins.Sort()

	return coins, nil
}
