package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

/////////////////////////////////////////////////////////////////////////
// Queue logic

// InsertActiveRegistryQueue inserts a name hash into the active registry queue at endTime
func (k Keeper) InsertActiveRegistryQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(nameHash)
	store.Set(types.GetActiveRegistryKey(endTime, nameHash), bz)
}

// RemoveFromActiveRegistryQueue removes a name hash from the Active Registry Queue
func (k Keeper) RemoveFromActiveRegistryQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveRegistryKey(endTime, nameHash))
}

// IterateActiveRegistryQueue iterates over the registries in the active registry queue
// and performs a callback function
func (k Keeper) IterateActiveRegistryQueue(ctx sdk.Context, endTime time.Time, cb func(nameHash types.NameHash, registry types.Registry) (stop bool)) {
	iterator := k.ActiveRegistryQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		_, nameHash := types.SplitActiveRegistryKey(iterator.Key())
		registry, err := k.GetRegistry(ctx, nameHash)
		if err != nil {
			panic(err)
		}

		if cb(nameHash, registry) {
			break
		}
	}
}

// InsertBidAuctionQueue inserts a name hash into the bid auction queue at endTime
func (k Keeper) InsertBidAuctionQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(nameHash)
	store.Set(types.GetBidAuctionKey(endTime, nameHash), bz)
}

// RemoveFromBidAuctionQueue removes a name hash from the Bid Auction Queue
func (k Keeper) RemoveFromBidAuctionQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBidAuctionKey(endTime, nameHash))
}

// IterateBidAuctionQueue iterates over the auctions in the bid auction queue
// and performs a callback function
func (k Keeper) IterateBidAuctionQueue(ctx sdk.Context, endTime time.Time, cb func(nameHash types.NameHash, auction types.Auction) (stop bool)) {
	iterator := k.BidAuctionQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		_, nameHash := types.SplitBidAuctionKey(iterator.Key())
		auction, err := k.GetAuction(ctx, nameHash)
		if err != nil {
			panic(err)
		}

		if cb(nameHash, auction) {
			break
		}
	}
}

// InsertRevealAuctionQueue inserts a name hash into the reveal auction queue at endTime
func (k Keeper) InsertRevealAuctionQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(nameHash)
	store.Set(types.GetRevealAuctionKey(endTime, nameHash), bz)
}

// RemoveFromRevealAuctionQueue removes a name hash from the Reveal Auction Queue
func (k Keeper) RemoveFromRevealAuctionQueue(ctx sdk.Context, nameHash types.NameHash, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetRevealAuctionKey(endTime, nameHash))
}

// IterateRevealAuctionQueue iterates over the auctions in the reveal auction queue
// and performs a callback function
func (k Keeper) IterateRevealAuctionQueue(ctx sdk.Context, endTime time.Time, cb func(nameHash types.NameHash, auction types.Auction) (stop bool)) {
	iterator := k.RevealAuctionQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		_, nameHash := types.SplitRevealAuctionKey(iterator.Key())
		auction, err := k.GetAuction(ctx, nameHash)
		if err != nil {
			panic(err)
		}

		if cb(nameHash, auction) {
			break
		}
	}
}

// ActiveRegistryQueueIterator returns an sdk.Iterator for all the name hashes in the Active Queue that expire by endTime
func (k Keeper) ActiveRegistryQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ActiveRegistryQueuePrefixKey, sdk.PrefixEndBytes(types.GetActiveRegistryQueueKey(endTime)))
}

// BidAuctionQueueIterator returns an sdk.Iterator for all the name hashes in the Bid Queue that expire by endTime
func (k Keeper) BidAuctionQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.BidAuctionQueuePrefixKey, sdk.PrefixEndBytes(types.GetBidAuctionQueueKey(endTime)))
}

// RevealAuctionQueueIterator returns an sdk.Iterator for all the name hashes in the Reveal Queue that expire by endTime
func (k Keeper) RevealAuctionQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.RevealAuctionQueuePrefixKey, sdk.PrefixEndBytes(types.GetRevealAuctionQueueKey(endTime)))
}
