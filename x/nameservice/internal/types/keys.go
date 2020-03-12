package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"time"
)

const (
	// ModuleName is the name of the nameservice module
	ModuleName = "nameservice"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the nameservice module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the nameservice module
	QuerierRoute = ModuleName
)

// Keys for nameservice store
// Items are stored with the following key: values
//
// - 0x01<NameHash(second_level_name)>: Registry
//
// - 0x02<NameHash(second_level_name)><NameHash(third_level_name)>: sdk.AccAddress
//
// - 0x03<accAddress_Bytes>: NameHash(second_level_name)
//
// - 0x04<endTime_Bytes>: NameHash(second_level_name)
//
// - 0x05<NameHash(second_level_name)>: Auction
//
// - 0x06<NameHash(second_level_name)><accAddress_Bytes>: Bid
//
// - 0x07<endTime_Bytes>: NameHash(second_level_name)
//
// - 0x08<endTime_Bytes>: NameHash(second_level_name)
var (
	// Keys for store prefixes
	RegistryPrefixKey            = []byte{0x01} // prefix for each key to a registry
	ResolvePrefixKey             = []byte{0x02} // prefix for each key to a address
	ReverseResolvePrefixKey      = []byte{0x03} // prefix for each key to a name hash
	ActiveRegistryQueuePrefixKey = []byte{0x04} // prefix for each key to a active period registry for expiration check
	AuctionPrefixKey             = []byte{0x05} // prefix for each key to a auction
	BidPrefixKey                 = []byte{0x06} // prefix for each key to a bid
	BidAuctionQueuePrefixKey     = []byte{0x07} // prefix for each key to a bid period auction
	RevealAuctionQueuePrefixKey  = []byte{0x08} // prefix for each key to a reveal period auction
)

// GetRegistryKey gets a specific registry from the store
func GetRegistryKey(nameHash NameHash) []byte {
	return append(RegistryPrefixKey, nameHash...)
}

// GetResolveKey gets a resolved account address from the store
func GetResolveKey(nameHash, childNameHash NameHash) []byte {
	return append(append(ResolvePrefixKey, nameHash...), childNameHash...)
}

// SplitResolveKey splits the resolve key and returns the parent/child name hash
func SplitResolveKey(key []byte) (nameHash, childNameHash NameHash) {
	nameHash = key[1 : tmhash.TruncatedSize+1]
	childNameHash = key[1+tmhash.TruncatedSize:]
	return
}

// GetReverseResolveKey gets a related second level name hash from the store
func GetReverseResolveKey(accAddress sdk.AccAddress) []byte {
	return append(ReverseResolvePrefixKey, accAddress.Bytes()...)
}

// GetActiveRegistryQueueKey gets the active registry queue key by endTime from the store
func GetActiveRegistryQueueKey(endTime time.Time) []byte {
	return append(ActiveRegistryQueuePrefixKey, sdk.FormatTimeBytes(endTime)...)
}

// GetActiveRegistryKey gets the active registry key from the store
func GetActiveRegistryKey(endTime time.Time, nameHash NameHash) []byte {
	return append(append(ActiveRegistryQueuePrefixKey, sdk.FormatTimeBytes(endTime)...), nameHash...)
}

// SplitActiveRegistryKey split the active registry key and returns the endTime and name hash
func SplitActiveRegistryKey(key []byte) (endTime time.Time, nameHash NameHash) {
	return splitKeyWithTime(key)
}

// GetAuctionKey gets a specific auction from the store
func GetAuctionKey(nameHash NameHash) []byte {
	return append(AuctionPrefixKey, nameHash...)
}

// GetBidKey gets a specific bid from the store
func GetBidKey(nameHash NameHash, bidderAddress sdk.AccAddress) []byte {
	return append(append(BidPrefixKey, nameHash...), bidderAddress.Bytes()...)
}

// SplitBidKey split the bid key and returns the name hash and bidder address
func SplitBidKey(key []byte) (nameHash NameHash, addr sdk.AccAddress) {
	return splitKeyWithAddress(key)
}

// SplitRevealKey split the reveal key and returns the name hash and bidder address
func SplitRevealKey(key []byte) (nameHash NameHash, addr sdk.AccAddress) {
	return splitKeyWithAddress(key)
}

// GetBidAuctionQueueKey gets the bid period auction queue key by endTime
func GetBidAuctionQueueKey(endTime time.Time) []byte {
	return append(BidAuctionQueuePrefixKey, sdk.FormatTimeBytes(endTime)...)
}

// GetBidAuctionKey gets the bid period auction key from the store
func GetBidAuctionKey(endTime time.Time, nameHash NameHash) []byte {
	return append(append(BidAuctionQueuePrefixKey, sdk.FormatTimeBytes(endTime)...), nameHash...)
}

// SplitBidAuctionKey split the bid auction key and returns the endTime and name hash
func SplitBidAuctionKey(key []byte) (endTime time.Time, nameHash NameHash) {
	return splitKeyWithTime(key)
}

// GetRevealAuctionQueueKey gets the bid period auction queue key by endTime
func GetRevealAuctionQueueKey(endTime time.Time) []byte {
	return append(RevealAuctionQueuePrefixKey, sdk.FormatTimeBytes(endTime)...)
}

// GetRevealAuctionKey gets the reveal period auction key from the store
func GetRevealAuctionKey(endTime time.Time, nameHash NameHash) []byte {
	return append(append(RevealAuctionQueuePrefixKey, sdk.FormatTimeBytes(endTime)...), nameHash...)
}

// SplitRevealAuctionKey split the reveal auction key and returns the endTime and name hash
func SplitRevealAuctionKey(key []byte) (endTime time.Time, nameHash NameHash) {
	return splitKeyWithTime(key)
}

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// private function
func splitKeyWithTime(key []byte) (endTime time.Time, nameHash NameHash) {
	if len(key[1:]) != tmhash.TruncatedSize+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}
	nameHash = key[1+lenTime:]
	return
}

func splitKeyWithAddress(key []byte) (nameHash NameHash, addr sdk.AccAddress) {
	if len(key[1:]) != tmhash.TruncatedSize+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}

	nameHash = key[1 : tmhash.TruncatedSize+1]
	addr = sdk.AccAddress(key[tmhash.TruncatedSize+1:])
	return
}
