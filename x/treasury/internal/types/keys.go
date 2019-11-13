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
// - 0x01: sdk.Dec
//
// - 0x02: sdk.Dec
//
// - 0x03<denom_Bytes>: sdk.Int
//
// - 0x04: sdk.Coins
//
// - 0x05: sdk.Coins
//
// - 0x06<epoch_Bytes>: sdk.Dec
//
// - 0x07<epoch_Bytes>: sdk.Dec
//
// - 0x08<epoch_Bytes>: sdk.Dec
var (
	// Keys for store prefixes
	TaxRateKey              = []byte{0x01} // a key for a tax-rate
	RewardWeightKey         = []byte{0x02} // a key for a reward-weight
	TaxCapKey               = []byte{0x03} // prefix for each key to a tax-cap
	TaxProceedsKey          = []byte{0x04} // a key for a tax-proceeds
	EpochInitialIssuanceKey = []byte{0x05} // a key for a initial epoch issuance

	// Keys for store prefixes of internal purpose variables
	MRKey  = []byte{0x06} // prefix for each key to a MR
	SRKey  = []byte{0x07} // prefix for each key to a SR
	TRLKey = []byte{0x08} // prefix for each key to a TRL
)

// GetTaxCapKey - stored by *denom*
func GetTaxCapKey(denom string) []byte {
	return append(TaxCapKey, []byte(denom)...)
}

// GetMRKey - stored by *epoch*
func GetMRKey(epoch int64) []byte {
	return GetSubkeyByEpoch(MRKey, epoch)
}

// GetSRKey - stored by *epoch*
func GetSRKey(epoch int64) []byte {
	return GetSubkeyByEpoch(SRKey, epoch)
}

// GetTRLKey - stored by *epoch*
func GetTRLKey(epoch int64) []byte {
	return GetSubkeyByEpoch(TRLKey, epoch)
}

// GetSubkeyByEpoch - stored by *epoch*
func GetSubkeyByEpoch(prefix []byte, epoch int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))
	return append(prefix, b...)
}
