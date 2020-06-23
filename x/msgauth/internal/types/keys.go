package types

import (
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
// - 0x01<accAddress_Bytes><accAddress_Bytes><msgType_Bytes>: Authorization
var (
	// Keys for store prefixes
	AuthorizationKey = []byte{0x01} // prefix for each key to a prevote
)

// GetAuthorizationKey - return authorization store key
func GetAuthorizationKey(granteeAddr sdk.AccAddress, granterAddr sdk.AccAddress, msgType string) []byte {
	return append(append(append(AuthorizationKey, granteeAddr.Bytes()...), granterAddr.Bytes()...), []byte(msgType)...)
}

// ExtractAddressesFromAuthorizationKey - split granter & grantee address from the authorization key
func ExtractAddressesFromAuthorizationKey(key []byte) (granteeAddr, granterAddr sdk.AccAddress) {
	granteeAddr = sdk.AccAddress(key[1 : sdk.AddrLen+1])
	granterAddr = sdk.AccAddress(key[sdk.AddrLen+1 : sdk.AddrLen*2+1])
	return
}
