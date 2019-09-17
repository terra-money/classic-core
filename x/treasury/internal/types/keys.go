package types

import (
	"encoding/binary"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "treasury"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route for treasury
	RouterKey = ModuleName

	// QuerierRoute is the querier route for treasury
	QuerierRoute = ModuleName
)

// Keys for treasury store
// Items are stored with the following key: values
//
// - 0x01<epoch_Bytes>: sdk.Dec
//
// - 0x02<epoch_Bytes>: sdk.Dec
//
// - 0x03<denom_Bytes>: sdk.Int
//
// - 0x04<epoch_Bytes>: sdk.Coins
//
// - 0x05<epoch_Bytes>: sdk.Coins
var (
	// Keys for store prefixes
	TaxRateKey            = []byte{0x01} // prefix for each key to a tax-rate
	RewardWeightKey       = []byte{0x02} // prefix for each key to a reward-weight
	TaxCapKey             = []byte{0x03} // prefix for each key to a tax-cap
	TaxProceedsKey        = []byte{0x04} // prefix for each key to a tax-proceeds
	HistoricalIssuanceKey = []byte{0x05} // prefix for each key to a historical issuance
)

// GetTaxRateKey - stored by *epoch*
func GetTaxRateKey(epoch int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	return append(TaxRateKey, b...)
}

// GetRewardWeightKey - stored by *epoch*
func GetRewardWeightKey(epoch int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	return append(RewardWeightKey, b...)
}

// GetTaxCapKey - stored by *denom*
func GetTaxCapKey(denom string) []byte {
	return append(TaxCapKey, []byte(denom)...)
}

// GetTaxProceedsKey - stored by *epoch*
func GetTaxProceedsKey(epoch int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	return append(TaxProceedsKey, b...)
}

// GetHistoricalIssuanceKey - stored by *epoch*
func GetHistoricalIssuanceKey(epoch int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	return append(HistoricalIssuanceKey, b...)
}
