package types

const (
	// ModuleName is the name of the market module
	ModuleName = "market"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the market module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the market module
	QuerierRoute = ModuleName
)

// Keys for market store
// Items are stored with the following key: values
//
// - 0x01: sdk.Dec
// - 0x02: sdk.Dec
var (
	//Keys for store prefixed
	MintPoolDeltaKey = []byte{0x01} // key for mint pool delta which gap between MintPool from BasePool
	BurnPoolDeltaKey = []byte{0x02} // key for burn pool delta which gap between MintPool from BasePool
)
