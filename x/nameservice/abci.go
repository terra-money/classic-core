package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {

	blockTime := ctx.BlockTime()
	revealPeriod := k.RevealPeriod(ctx)

	// fetch auctions whose bid periods have ended (are passed the block time)
	k.IterateBidAuctionQueue(ctx, blockTime, func(nameHash types.NameHash, auction Auction) bool {
		// delete from the queue
		k.RemoveFromBidAuctionQueue(ctx, nameHash, auction.EndTime)

		// check at least one bid is exists
		activateAuction := false
		k.IterateBid(ctx, nameHash, func(nameHash NameHash, bid Bid) bool {
			activateAuction = true
			return true
		})

		if activateAuction {
			auction.Status = AuctionStatusReveal
			auction.EndTime = blockTime.Add(revealPeriod)

			k.SetAuction(ctx, nameHash, auction)
			k.InsertRevealAuctionQueue(ctx, nameHash, auction.EndTime)
		} else {
			k.DeleteAuction(ctx, nameHash)
		}

		return false
	})

	renewalInterval := k.RenewalInterval(ctx)
	// fetch auctions whose reveal periods have ended (are passed the block time)
	k.IterateRevealAuctionQueue(ctx, blockTime, func(nameHash NameHash, auction Auction) (stop bool) {
		// delete from the queue
		k.RemoveFromRevealAuctionQueue(ctx, nameHash, auction.EndTime)

		// delete auction
		k.DeleteAuction(ctx, nameHash)

		// make registry
		if !auction.TopBidder.Empty() {
			endTime := blockTime.Add(renewalInterval)
			registry := NewRegistry(auction.Name, auction.TopBidder, endTime)
			k.SetRegistry(ctx, nameHash, registry)
			k.InsertActiveRegistryQueue(ctx, nameHash, endTime)

			// TODO - how to treat the profit
			err := k.SupplyKeeper.BurnCoins(ctx, ModuleName, auction.TopBidAmount)
			if err != nil {
				panic(err)
			}
		}

		// do slash & delete all left bid for the current name auction
		k.IterateBid(ctx, nameHash, func(nameHash NameHash, bid Bid) bool {
			// TODO - determine&apply slash fraction
			err := k.SupplyKeeper.BurnCoins(ctx, ModuleName, sdk.NewCoins(bid.Deposit))
			if err != nil {
				panic(err)
			}

			k.DeleteBid(ctx, nameHash, bid.Bidder)
			return false
		})

		return false
	})

	// fetch registries whose active periods have ended (are passed the block time)
	k.IterateActiveRegistryQueue(ctx, blockTime, func(nameHash NameHash, registry Registry) (stop bool) {
		// delete expired registry
		k.DeleteRegistry(ctx, nameHash)
		k.RemoveFromActiveRegistryQueue(ctx, nameHash, registry.EndTime)

		// delete expired registry's resolve and reverse resolve entries
		k.IterateResolve(ctx, nameHash, func(_, childNameHash NameHash, address sdk.AccAddress) bool {
			k.DeleteResolve(ctx, nameHash, childNameHash)
			k.DeleteReverseResolve(ctx, address)
			return false
		})

		return false
	})
}
