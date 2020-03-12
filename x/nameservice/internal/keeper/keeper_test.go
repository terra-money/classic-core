package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/x/nameservice/internal/types"
	"testing"
	"time"
)

func TestAuction(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.name")
	name2 := types.Name("wallet2.name")
	nameHash, _ := name.NameHash()
	nameHash2, _ := name2.NameHash()
	endTime := time.Now().UTC()

	// no auction exists
	_, err := input.NameserviceKeeper.GetAuction(input.Ctx, nameHash)
	require.Error(t, err)

	// set auction
	auction := types.NewAuction(types.Name(name), types.AuctionStatusBid, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash, auction)

	// get auction and compare
	auctionFromKeeper, err := input.NameserviceKeeper.GetAuction(input.Ctx, nameHash)
	require.NoError(t, err)
	require.Equal(t, auction, auctionFromKeeper)

	// delete auction and check existence
	input.NameserviceKeeper.DeleteAuction(input.Ctx, nameHash)
	_, err = input.NameserviceKeeper.GetAuction(input.Ctx, nameHash)
	require.Error(t, err)

	// iterate auctions
	auction2 := types.NewAuction(types.Name(name), types.AuctionStatusBid, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash, auction)
	input.NameserviceKeeper.SetAuction(input.Ctx, nameHash2, auction2)

	input.NameserviceKeeper.IterateAuction(input.Ctx, func(h types.NameHash, ac types.Auction) bool {
		if nameHash.Equal(h) {
			require.Equal(t, auction, ac)
		} else if nameHash2.Equal(h) {
			require.Equal(t, auction2, ac)
		} else {
			assert.Fail(t, "unknown name hash", h)
		}

		return false
	})
}

func TestBid(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.test")
	name2 := types.Name("wallet2.test")
	nameHash, _ := name.NameHash()
	nameHash2, _ := name2.NameHash()
	bidAmount := sdk.NewCoin("foo", sdk.NewInt(123))
	bidHash := types.GetBidHash("salt", name, bidAmount, Addrs[0])
	bidHash2 := types.GetBidHash("salt", name2, bidAmount, Addrs[1])

	// no bid exists
	_, err := input.NameserviceKeeper.GetBid(input.Ctx, nameHash, Addrs[0])
	require.Error(t, err)

	// set bid
	bid := types.NewBid(bidHash, bidAmount, Addrs[0])
	input.NameserviceKeeper.SetBid(input.Ctx, nameHash, bid)

	// get bid and compare
	bidFromKeeper, err := input.NameserviceKeeper.GetBid(input.Ctx, nameHash, Addrs[0])
	require.NoError(t, err)
	require.Equal(t, bid, bidFromKeeper)

	// delete bid and check existence
	input.NameserviceKeeper.DeleteBid(input.Ctx, nameHash, Addrs[0])
	_, err = input.NameserviceKeeper.GetBid(input.Ctx, nameHash, Addrs[0])
	require.Error(t, err)

	// iterate bids
	bid2 := types.NewBid(bidHash2, bidAmount, Addrs[1])
	input.NameserviceKeeper.SetBid(input.Ctx, nameHash, bid)
	input.NameserviceKeeper.SetBid(input.Ctx, nameHash2, bid2)

	input.NameserviceKeeper.IterateBid(input.Ctx, nameHash, func(_ types.NameHash, b types.Bid) bool {
		require.Equal(t, bid, b)

		return false
	})

	input.NameserviceKeeper.IterateBid(input.Ctx, types.NameHash{}, func(h types.NameHash, b types.Bid) bool {
		if nameHash.Equal(h) {
			require.Equal(t, bid, b)
		} else if nameHash2.Equal(h) {
			require.Equal(t, bid2, b)
		} else {
			assert.Fail(t, "unknown name hash", h)
		}

		return false
	})
}

func TestRegistry(t *testing.T) {
	input := CreateTestInput(t)

	name := types.Name("wallet.name")
	name2 := types.Name("wallet.name2")
	endTime := time.Now().UTC()
	nameHash, _ := name.NameHash()
	nameHash2, _ := name2.NameHash()

	// no registry exists
	_, err := input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.Error(t, err)

	// set registry and check
	registry := types.NewRegistry(name, Addrs[0], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	registryFromKeeper, err := input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.Equal(t, registry, registryFromKeeper)

	// delete registry
	input.NameserviceKeeper.DeleteRegistry(input.Ctx, nameHash)
	_, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.Error(t, err)

	// iterate registries
	registry2 := types.NewRegistry(name2, Addrs[1], endTime)
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash2, registry2)

	input.NameserviceKeeper.IterateRegistry(input.Ctx, func(h types.NameHash, b types.Registry) bool {
		if nameHash.Equal(h) {
			require.Equal(t, registry, b)
		} else if nameHash2.Equal(h) {
			require.Equal(t, registry2, b)
		} else {
			assert.Fail(t, "unknown name hash", h)
		}

		return false
	})
}

func TestResolve(t *testing.T) {
	input := CreateTestInput(t)

	name := "chai.terra"
	name2 := "timon.chai.terra"
	name3 := "dokwon.harvest.terra"
	nameHash, childNameHash := types.Name(name).NameHash()
	nameHash2, childNameHash2 := types.Name(name2).NameHash()
	nameHash3, childNameHash3 := types.Name(name3).NameHash()

	require.Equal(t, nameHash, nameHash2)

	// no resolve exists
	_, err := input.NameserviceKeeper.GetResolve(input.Ctx, nameHash, childNameHash)
	require.Error(t, err)

	// set resolve and check
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childNameHash, Addrs[0])
	resolveFromKeeper, err := input.NameserviceKeeper.GetResolve(input.Ctx, nameHash, childNameHash)
	require.NoError(t, err)
	require.Equal(t, Addrs[0], resolveFromKeeper)

	// delete resolve
	input.NameserviceKeeper.DeleteResolve(input.Ctx, nameHash, childNameHash)
	_, err = input.NameserviceKeeper.GetResolve(input.Ctx, nameHash, childNameHash)
	require.Error(t, err)

	// iterate resolves
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childNameHash, Addrs[0])
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash2, childNameHash2, Addrs[1])
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash3, childNameHash3, Addrs[2])

	// iterate with chai prefix
	input.NameserviceKeeper.IterateResolve(input.Ctx, nameHash, func(ph types.NameHash, ch types.NameHash, a sdk.AccAddress) bool {
		require.True(t, nameHash.Equal(ph))

		if childNameHash.Equal(ch) {
			require.Equal(t, Addrs[0], a)
		} else if childNameHash2.Equal(ch) {
			require.Equal(t, Addrs[1], a)
		} else {
			assert.Fail(t, "unknown child name hash", ch)
		}

		return false
	})

	// iterate with empty prefix
	input.NameserviceKeeper.IterateResolve(input.Ctx, []byte{}, func(ph types.NameHash, ch types.NameHash, a sdk.AccAddress) bool {
		if nameHash.Equal(ph) {
			if childNameHash.Equal(ch) {
				require.Equal(t, Addrs[0], a)
			} else if childNameHash2.Equal(ch) {
				require.Equal(t, Addrs[1], a)
			} else {
				assert.Fail(t, "unknown child name hash", ch)
			}
		} else if nameHash3.Equal(ph) {
			require.Equal(t, childNameHash3, ch)
			require.Equal(t, Addrs[2], a)
		} else {
			assert.Fail(t, "unknown parent name hash", ph)
		}

		return false
	})
}

func TestReverseResolve(t *testing.T) {
	input := CreateTestInput(t)

	name := "chai.terra"
	name2 := "dokwon.harvest.terra"
	nameHash, _ := types.Name(name).NameHash()
	nameHash2, _ := types.Name(name2).NameHash()

	// no reverse resolve exists
	_, err := input.NameserviceKeeper.GetReverseResolve(input.Ctx, Addrs[0])
	require.Error(t, err)

	// set reverse resolve and check
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, Addrs[0], nameHash)
	resolveFromKeeper, err := input.NameserviceKeeper.GetReverseResolve(input.Ctx, Addrs[0])
	require.NoError(t, err)
	require.Equal(t, nameHash, resolveFromKeeper)

	// delete reverse resolve
	input.NameserviceKeeper.DeleteReverseResolve(input.Ctx, Addrs[0])
	_, err = input.NameserviceKeeper.GetReverseResolve(input.Ctx, Addrs[0])
	require.Error(t, err)

	// iterate reverse resolves
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, Addrs[0], nameHash)
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, Addrs[1], nameHash2)

	input.NameserviceKeeper.IterateReverseResolve(input.Ctx, func(a sdk.AccAddress, ph types.NameHash) bool {

		if Addrs[0].Equals(a) {
			require.Equal(t, nameHash, ph)
		} else if Addrs[1].Equals(a) {
			require.Equal(t, nameHash2, ph)
		} else {
			assert.Fail(t, "unknown account address", a)
		}

		return false
	})
}
