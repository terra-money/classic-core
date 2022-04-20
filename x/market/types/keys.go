package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName is the name of the market module
	ModuleName = "market"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the market module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the market module
	QuerierRoute = ModuleName

	// BurnModuleName is special purpose module name to perform burn coins
	// burn address = terra1sk06e3dyexuq4shw77y3dsv480xv42mq73anxu
	BurnModuleName = "burn"
)

// Keys for market store
// Items are stored with the following key: values
//
// - 0x01: sdk.Dec
//
// - 0x02: SeigniorageRoutes
var (
	// Keys for store prefixed
	TerraPoolDeltaKey    = []byte{0x01} // key for terra pool delta which gap between MintPool from BasePool
	SeigniorageRoutesKey = []byte{0x02} // key for SeigniorageRoutes which seigniorage will be routed

	// AlternateCommunityPoolAddress is intended to replace the community pool.
	// Since the community pool does not have its own address, an alternate
	// address is needed to register as a seigniorage route.
	// - https://github.com/cosmos/cosmos-sdk/issues/10811
	//
	// The bech32 encoded form of the alternate address is
	// terra1pf89qgchfytxmd0fvdn3zfdshja4fqtevepu93
	AlternateCommunityPoolAddress = authtypes.NewModuleAddress("community_pool_placeholder")
)
