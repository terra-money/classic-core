package market

// nolint
var (
	PrefixCoinSupply = []byte("coinsupply")
)

// GetCoinSupplyKey is in format of PrefixElect||denom
func GetCoinSupplyKey(denom string) []byte {
	return append(PrefixCoinSupply, []byte(denom)...)
}
