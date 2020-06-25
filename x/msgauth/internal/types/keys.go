package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "msgauth"

	// StoreKey is the store key string for msgauth
	StoreKey = ModuleName

	// RouterKey is the message route for msgauth
	RouterKey = ModuleName

	// QuerierRoute is the querier route for msgauth
	QuerierRoute = ModuleName
)

// Keys for msgauth store
// Items are stored with the following key: values
//
// - 0x01<accAddress_Bytes><accAddress_Bytes><msgType_Bytes>: Grant
// - 0x02<timestamp_Bytes>: []GGMPair
var (
	// Keys for store prefixes
	GrantKey      = []byte{0x01} // prefix for each key to a prevote
	GrantQueueKey = []byte{0x02} // prefix for the timestamps in grants queue
)

// GetGrantKey - return grant store key
func GetGrantKey(granterAddr sdk.AccAddress, granteeAddr sdk.AccAddress, msgType string) []byte {
	return append(append(append(GrantKey, granterAddr.Bytes()...), granteeAddr.Bytes()...), []byte(msgType)...)
}

// GetGrantTimeKey - return grant queue store key
func GetGrantTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(GrantQueueKey, bz...)
}

// ExtractAddressesFromGrantKey - split granter & grantee address from the authorization key
func ExtractAddressesFromGrantKey(key []byte) (granterAddr, granteeAddr sdk.AccAddress) {
	granterAddr = sdk.AccAddress(key[1 : sdk.AddrLen+1])
	granteeAddr = sdk.AccAddress(key[sdk.AddrLen+1 : sdk.AddrLen*2+1])
	return
}
