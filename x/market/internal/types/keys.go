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
// - 0x03: sdk.Dec
// - 0x04: int64
var (
	//Keys for store prefixed
	BasePoolKey         = []byte{0x01} // key for a Base Pool
	LunaPoolKey         = []byte{0x02} // key for Luna Pool
	TerraPoolKey        = []byte{0x03} // key for Terra Pool
	LastUpdateHeightKey = []byte{0x04} // key for Last Update Height
)
