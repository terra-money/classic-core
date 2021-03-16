package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the oracle module
	ModuleName = "oracle"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the oracle module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the oracle module
	QuerierRoute = ModuleName
)

// Keys for oracle store
// Items are stored with the following key: values
//
// - 0x01<denom_Bytes>: sdk.Dec
//
// - 0x02<valAddress_Bytes>: accAddress
//
// - 0x03<valAddress_Bytes>: int64
//
// - 0x04<valAddress_Bytes>: AggregateExchangeRatePrevote
//
// - 0x05<valAddress_Bytes>: AggregateExchangeRateVote
//
// - 0x06<denom_Bytes>: sdk.Dec
var (
	// Keys for store prefixes
	ExchangeRateKey                 = []byte{0x01} // prefix for each key to a rate
	FeederDelegationKey             = []byte{0x02} // prefix for each key to a feeder delegation
	MissCounterKey                  = []byte{0x03} // prefix for each key to a miss counter
	AggregateExchangeRatePrevoteKey = []byte{0x04} // prefix for each key to a aggregate prevote
	AggregateExchangeRateVoteKey    = []byte{0x05} // prefix for each key to a aggregate vote
	TobinTaxKey                     = []byte{0x06} // prefix for each key to a tobin tax
)

// GetExchangeRateKey - stored by *denom*
func GetExchangeRateKey(denom string) []byte {
	return append(ExchangeRateKey, []byte(denom)...)
}

// GetFeederDelegationKey - stored by *Validator* address
func GetFeederDelegationKey(v sdk.ValAddress) []byte {
	return append(FeederDelegationKey, v.Bytes()...)
}

// GetMissCounterKey - stored by *Validator* address
func GetMissCounterKey(v sdk.ValAddress) []byte {
	return append(MissCounterKey, v.Bytes()...)
}

// GetAggregateExchangeRatePrevoteKey - stored by *Validator* address
func GetAggregateExchangeRatePrevoteKey(v sdk.ValAddress) []byte {
	return append(AggregateExchangeRatePrevoteKey, v.Bytes()...)
}

// GetAggregateExchangeRateVoteKey - stored by *Validator* address
func GetAggregateExchangeRateVoteKey(v sdk.ValAddress) []byte {
	return append(AggregateExchangeRateVoteKey, v.Bytes()...)
}

// GetTobinTaxKey - stored by *denom* bytes
func GetTobinTaxKey(d string) []byte {
	return append(TobinTaxKey, []byte(d)...)
}

// ExtractDenomFromTobinTaxKey - split denom from the tobin tax key
func ExtractDenomFromTobinTaxKey(key []byte) (denom string) {
	denom = string(key[1:])
	return
}
