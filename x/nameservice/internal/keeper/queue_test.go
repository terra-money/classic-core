package keeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/x/nameservice/internal/types"
	"testing"
	"time"
)

func TestActiveRegistryQueue(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.name")
	nameHash, _ := name.NameHash()
	beforeEndTime := time.Now().UTC()
	endTime := beforeEndTime.Add(time.Hour).UTC()

	// queue must be empty
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Registry) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// set registry and insert queue
	registry := types.NewRegistry(name, Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.InsertActiveRegistryQueue(input.Ctx, nameHash, endTime)

	// check with same endTime
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime, func(_ types.NameHash, r types.Registry) (stop bool) {
		require.Equal(t, registry, r)
		return false
	})

	// check with later endTime
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime.Add(time.Hour), func(_ types.NameHash, r types.Registry) (stop bool) {
		require.Equal(t, registry, r)
		return false
	})

	// check with before endTime
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, beforeEndTime, func(_ types.NameHash, _ types.Registry) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// remove registry from the store
	input.NameserviceKeeper.DeleteRegistry(input.Ctx, nameHash)
	require.Panics(t, func() {
		input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Registry) (stop bool) {
			assert.Fail(t, "panic should be occurs before entering here")
			return false
		})
	})

	// remove from queue
	input.NameserviceKeeper.RemoveFromActiveRegistryQueue(input.Ctx, nameHash, endTime)
	input.NameserviceKeeper.IterateActiveRegistryQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Registry) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})
}

func TestBidAuctionQueue(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.name")
	nameHash, _ := name.NameHash()
	beforeEndTime := time.Now().UTC()
	endTime := beforeEndTime.Add(time.Hour).UTC()

	// queue must be empty
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// set auction and insert queue
	auction := types.NewAuction(name, types.AuctionStatusBid, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash, auction)
	input.NameserviceKeeper.InsertBidAuctionQueue(input.Ctx, nameHash, endTime)

	// check with same endTime
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, a types.Auction) (stop bool) {
		require.Equal(t, auction, a)
		return false
	})

	// check with later endTime
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, endTime.Add(time.Hour), func(_ types.NameHash, a types.Auction) (stop bool) {
		require.Equal(t, auction, a)
		return false
	})

	// check with before endTime
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, beforeEndTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// remove auction from the store
	input.NameserviceKeeper.DeleteAuction(input.Ctx, nameHash)
	require.Panics(t, func() {
		input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
			assert.Fail(t, "panic should be occurs before entering here")
			return false
		})
	})

	// remove from queue
	input.NameserviceKeeper.RemoveFromBidAuctionQueue(input.Ctx, nameHash, endTime)
	input.NameserviceKeeper.IterateBidAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})
}
func TestRevealAuctionQueue(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.name")
	nameHash, _ := name.NameHash()
	beforeEndTime := time.Now().UTC()
	endTime := beforeEndTime.Add(time.Hour).UTC()

	// queue must be empty
	input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// set auction and insert queue
	auction := types.NewAuction(name, types.AuctionStatusReveal, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash, auction)
	input.NameserviceKeeper.InsertRevealAuctionQueue(input.Ctx, nameHash, endTime)

	// check with same endTime
	input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, a types.Auction) (stop bool) {
		require.Equal(t, auction, a)
		return false
	})

	// check with later endTime
	input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, endTime.Add(time.Hour), func(_ types.NameHash, a types.Auction) (stop bool) {
		require.Equal(t, auction, a)
		return false
	})

	// check with before endTime
	input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, beforeEndTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})

	// remove auction from the store
	input.NameserviceKeeper.DeleteAuction(input.Ctx, nameHash)
	require.Panics(t, func() {
		input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
			assert.Fail(t, "panic should be occurs before entering here")
			return false
		})
	})

	// remove from queue
	input.NameserviceKeeper.RemoveFromRevealAuctionQueue(input.Ctx, nameHash, endTime)
	input.NameserviceKeeper.IterateRevealAuctionQueue(input.Ctx, endTime, func(_ types.NameHash, _ types.Auction) (stop bool) {
		assert.Fail(t, "queue must be empty")
		return false
	})
}
