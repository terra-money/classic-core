package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// InitGenesis initialize default parameters
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for nameHashStr, auction := range data.Auctions {
		nameHash, err := NameHashFromHexString(nameHashStr)
		if err != nil {
			panic(err)
		}

		keeper.SetAuction(ctx, nameHash, auction)
		if auction.Status == AuctionStatusBid {
			keeper.InsertBidAuctionQueue(ctx, nameHash, auction.EndTime)
		} else if auction.Status == AuctionStatusReveal {
			keeper.InsertRevealAuctionQueue(ctx, nameHash, auction.EndTime)
		} else {
			panic("invalid auction status")
		}
	}

	for nameHashStr, bid := range data.Bids {
		nameHash, err := NameHashFromHexString(nameHashStr)
		if err != nil {
			panic(err)
		}

		keeper.SetBid(ctx, nameHash, bid)
	}

	for nameHashStr, registry := range data.Registries {
		nameHash, err := NameHashFromHexString(nameHashStr)
		if err != nil {
			panic(err)
		}

		keeper.SetRegistry(ctx, nameHash, registry)
		keeper.InsertActiveRegistryQueue(ctx, nameHash, registry.EndTime)
	}

	for concatNameHashesStr, addr := range data.Resolves {
		nameHashes := strings.Split(concatNameHashesStr, ":")
		nameHash, err := NameHashFromHexString(nameHashes[0])
		if err != nil {
			panic(err)
		}

		childNameHash, err := NameHashFromHexString(nameHashes[1])
		if err != nil {
			panic(err)
		}

		keeper.SetResolve(ctx, nameHash, childNameHash, addr)
		keeper.SetReverseResolve(ctx, addr, nameHash)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := DefaultParams()

	auctions := make(map[string]Auction)
	keeper.IterateAuction(ctx, func(nameHash NameHash, auction Auction) bool {
		auctions[nameHash.String()] = auction
		return false
	})

	bids := make(map[string]Bid)
	keeper.IterateBid(ctx, []byte{}, func(nameHash NameHash, bid Bid) bool {
		bids[nameHash.String()] = bid
		return false
	})

	registries := make(map[string]Registry)
	keeper.IterateRegistry(ctx, func(nameHash NameHash, registry Registry) bool {
		registries[nameHash.String()] = registry
		return false
	})

	resolves := make(map[string]sdk.AccAddress)
	keeper.IterateResolve(ctx, []byte{}, func(nameHash, childNameHash NameHash, addr sdk.AccAddress) bool {
		concatNameHashesStr := nameHash.String() + ":" + childNameHash.String()
		resolves[concatNameHashesStr] = addr
		return false
	})

	return NewGenesisState(params, auctions, bids, registries, resolves)
}
