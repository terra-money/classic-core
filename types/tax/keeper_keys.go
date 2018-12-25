package tax

// nolint
var (
	PrefixIssuance = []byte("coinsupply")
	PrefixTaxRate  = []byte("taxrate")
)

// GetCoinSupplyKey is in format of PrefixElect||denom
func GetCoinSupplyKey(denom string) []byte {
	return append(PrefixIssuance, []byte(denom)...)
}

// GetTaxRateKey gets the effective tax rate
func GetTaxRateKey() []byte {
	return PrefixTaxRate
}
