package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	SupplyKeeper types.SupplyKeeper
	marketKeeper types.MarketKeeper

	codespace sdk.CodespaceType
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, supplyKeeper types.SupplyKeeper,
	marketKeeper types.MarketKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramspace.WithKeyTable(types.ParamKeyTable()),
		codespace:  codespace,

		SupplyKeeper: supplyKeeper,
		marketKeeper: marketKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns a codespace of keeper
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Registry Logic

// GetRegistry returns Registry for the given name hash
func (k Keeper) GetRegistry(ctx sdk.Context, nameHash types.NameHash) (registry types.Registry, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetRegistryKey(nameHash))
	if b == nil {
		return types.Registry{}, types.ErrRegistryNotExists(k.codespace)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &registry)
	return
}

// SetRegistry stores Registry for the given name hash
func (k Keeper) SetRegistry(ctx sdk.Context, nameHash types.NameHash, registry types.Registry) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(registry)
	store.Set(types.GetRegistryKey(nameHash), bz)
}

// DeleteRegistry removes Registry from the store
func (k Keeper) DeleteRegistry(ctx sdk.Context, nameHash types.NameHash) {
	ctx.KVStore(k.storeKey).Delete(types.GetRegistryKey(nameHash))
}

// IterateRegistry iterates registry in the store
func (k Keeper) IterateRegistry(ctx sdk.Context, handler func(nameHash types.NameHash, registry types.Registry) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.RegistryPrefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var registry types.Registry
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &registry)
		if handler(iter.Key()[1:], registry) {
			break
		}
	}
}

// GetResolve returns resolved address of the name
func (k Keeper) GetResolve(ctx sdk.Context, nameHash, childNameHash types.NameHash) (address sdk.AccAddress, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetResolveKey(nameHash, childNameHash))
	if b == nil {
		return sdk.AccAddress{}, types.ErrResolveNotExists(k.codespace)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &address)
	return
}

// SetResolve stores resolve entry to store
func (k Keeper) SetResolve(ctx sdk.Context, nameHash, childNameHash types.NameHash, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(address)
	store.Set(types.GetResolveKey(nameHash, childNameHash), bz)
}

// DeleteResolve removes resolve entry from the store
func (k Keeper) DeleteResolve(ctx sdk.Context, nameHash, childNameHash types.NameHash) {
	ctx.KVStore(k.storeKey).Delete(types.GetResolveKey(nameHash, childNameHash))
}

// IterateResolve iterates resolve entries in the store with nameHash (set it []byte{} to iterate whole entries)
func (k Keeper) IterateResolve(ctx sdk.Context, nameHash types.NameHash, handler func(nameHash, childNameHash types.NameHash, address sdk.AccAddress) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetResolveKey(nameHash, []byte{}))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var address sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &address)
		nameHash, childNameHash := types.SplitResolveKey(iter.Key())
		if handler(nameHash, childNameHash, address) {
			break
		}
	}
}

// GetReverseResolve returns resolved second level name hash of the address
func (k Keeper) GetReverseResolve(ctx sdk.Context, address sdk.AccAddress) (nameHash types.NameHash, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetReverseResolveKey(address))
	if b == nil {
		return types.NameHash{}, types.ErrReverseResolveNotExists(k.codespace, address)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &nameHash)
	return
}

// SetReverseResolve stores reverse resolve entry (address => second level name hash) to store
func (k Keeper) SetReverseResolve(ctx sdk.Context, address sdk.AccAddress, nameHash types.NameHash) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(nameHash)
	store.Set(types.GetReverseResolveKey(address), bz)
}

// DeleteReverseResolve removes reverse resolve entry from the store
func (k Keeper) DeleteReverseResolve(ctx sdk.Context, address sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetReverseResolveKey(address))
}

// IterateReverseResolve iterates resolve entries in the store
func (k Keeper) IterateReverseResolve(ctx sdk.Context, handler func(address sdk.AccAddress, nameHash types.NameHash) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ReverseResolvePrefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var nameHash types.NameHash
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &nameHash)
		if handler(iter.Key()[1:], nameHash) {
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Auction Logic

// GetAuction returns Auction for the given name hash
func (k Keeper) GetAuction(ctx sdk.Context, nameHash types.NameHash) (auction types.Auction, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetAuctionKey(nameHash))
	if b == nil {
		return types.Auction{}, types.ErrAuctionNotExists(k.codespace)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &auction)
	return
}

// SetAuction stores Auction for the given name hash
func (k Keeper) SetAuction(ctx sdk.Context, nameHash types.NameHash, auction types.Auction) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(auction)
	store.Set(types.GetAuctionKey(nameHash), bz)
}

// DeleteAuction removes Auction from the store
func (k Keeper) DeleteAuction(ctx sdk.Context, nameHash types.NameHash) {
	ctx.KVStore(k.storeKey).Delete(types.GetAuctionKey(nameHash))
}

// IterateAuction iterates auction in the store
func (k Keeper) IterateAuction(ctx sdk.Context, handler func(nameHash types.NameHash, auction types.Auction) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AuctionPrefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var auction types.Auction
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &auction)
		if handler(iter.Key()[1:], auction) {
			break
		}
	}
}

// GetBid return Bid for the given name hash and bidder address
func (k Keeper) GetBid(ctx sdk.Context, nameHash types.NameHash, bidder sdk.AccAddress) (bid types.Bid, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetBidKey(nameHash, bidder))
	if b == nil {
		return types.Bid{}, types.ErrBidNotExists(k.codespace)
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &bid)
	return
}

// SetBid stores Bid for the given name hash and bidder address
func (k Keeper) SetBid(ctx sdk.Context, nameHash types.NameHash, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(bid)
	store.Set(types.GetBidKey(nameHash, bid.Bidder), bz)
}

// DeleteBid removes Bid from the store
func (k Keeper) DeleteBid(ctx sdk.Context, nameHash types.NameHash, bidder sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetBidKey(nameHash, bidder))
}

// IterateBid iterates auction in the store
func (k Keeper) IterateBid(ctx sdk.Context, nameHash types.NameHash, handler func(nameHash types.NameHash, bid types.Bid) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidKey(nameHash, sdk.AccAddress{}))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &bid)
		nameHash, _ := types.SplitBidKey(iter.Key())
		if handler(nameHash, bid) {
			break
		}
	}
}
