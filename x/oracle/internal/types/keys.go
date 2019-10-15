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
// - 0x01<denom_Bytes><valAddress_Bytes>: Prevote
//
// - 0x02<denom_Bytes><valAddress_Bytes>: Vote
//
// - 0x03<denom_Bytes>: sdk.Dec
//
// - 0x04<valAddress_Bytes>: accAddress
//
// - 0x05<valAddress_Bytes>: Claim
//
// - 0x06<valAddress_Bytes><period_Bytes>: bool
var (
	// Keys for store prefixes
	PrevoteKey          = []byte{0x01} // prefix for each key to a prevote
	VoteKey             = []byte{0x02} // prefix for each key to a vote
	PriceKey            = []byte{0x03} // prefix for each key to a price
	FeederDelegationKey = []byte{0x04} // prefix for each key to a feeder delegation
)

// GetPrevoteKey - stored by *Validator* address and denom
func GetPrevoteKey(denom string, v sdk.ValAddress) []byte {
	return append(append(PrevoteKey, []byte(denom)...), v.Bytes()...)
}

// GetVoteKey - stored by *Validator* address and denom
func GetVoteKey(denom string, v sdk.ValAddress) []byte {
	return append(append(VoteKey, []byte(denom)...), v.Bytes()...)
}

// GetPriceKey - stored by *denom*
func GetPriceKey(denom string) []byte {
	return append(PriceKey, []byte(denom)...)
}

// GetFeederDelegationKey - stored by *Validator* address
func GetFeederDelegationKey(v sdk.ValAddress) []byte {
	return append(FeederDelegationKey, v.Bytes()...)
}
