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

	// BurnModuleName is special purpose module name to perform burn coins
	// burn address = terra1sk06e3dyexuq4shw77y3dsv480xv42mq73anxu
	BurnModuleName = "burn"
)

// Keys for market store
// Items are stored with the following key: values
//
// - 0x01: sdk.Dec
var (
	// Keys for store prefixed
	TerraPoolDeltaKey = []byte{0x01} // key for terra pool delta which gap between MintPool from BasePool
)
