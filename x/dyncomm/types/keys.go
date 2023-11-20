package types

const (
	// ModuleName defines the module's name.
	ModuleName = "dyncomm"
	// StoreKey is the string store representation
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// MinCommissionRatesPrefix store prefixes
var MinCommissionRatesPrefix = []byte{0x01} // prefix for each MinCommissionRate entry

// MinCommissionRates - stored by *validator addr*
func GetMinCommissionRatesKey(addr string) []byte {
	return append(MinCommissionRatesPrefix, []byte(addr)...)
}
