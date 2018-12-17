package tax

// nolint
var (
	PrefixIssuance = []byte("coinsupply")
)

// GetCoinSupplyKey is in format of PrefixElect||denom
func GetCoinSupplyKey(denom string) []byte {
	return append(PrefixIssuance, []byte(denom)...)
}
