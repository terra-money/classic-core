package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the oracle store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	valset sdk.ValidatorSet
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, valset sdk.ValidatorSet) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,

		valset: valset,
	}
}

func (keeper Keeper) SetWhitelist(ctx sdk.Context, whitelist Whitelist) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(whitelist)
	store.Set(KeyWhitelist, bz)
}

//nolint
func (keeper Keeper) WhitelistContains(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(KeyWhitelist)
	if bz == nil {
		return false
	}
	wl := Whitelist{}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &wl)

	contains := false
	for _, wDenom := range wl {
		if wDenom == denom {
			contains = true
			break
		}
	}

	return contains
}

//nolint
func (keeper Keeper) GetWhitelist(ctx sdk.Context) (res Whitelist) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(KeyWhitelist)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (keeper Keeper) SetVotePeriod(ctx sdk.Context, votePeriod sdk.Int) {
	store := ctx.KVStore(keeper.key)

	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(votePeriod)
	store.Set(KeyVotePeriod, bz)
}

//nolint
func (keeper Keeper) GetVotePeriod(ctx sdk.Context) (res sdk.Int) {
	store := ctx.KVStore(keeper.key)

	bz := store.Get(KeyVotePeriod)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)

	return
}

func (keeper Keeper) SetThreshold(ctx sdk.Context, threshold sdk.Dec) {
	store := ctx.KVStore(keeper.key)

	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(threshold)
	store.Set(KeyThreshold, bz)
}

//nolint
func (keeper Keeper) GetThreshold(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(KeyThreshold)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (keeper Keeper) AddVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)
	key := GetVoteKey(vote.FeedMsg.Denom, vote.FeedMsg.Feeder)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(key, bz)
}

func (keeper Keeper) GetAllVotes(ctx sdk.Context, denom string) (res PriceVotes) {
	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, GetVotePrefix(denom))
	for ; iter.Valid(); iter.Next() {
		var pv PriceVote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &pv)
		res = append(res, pv)
	}
	iter.Close()
	return
}

func (keeper Keeper) ClearVotes(ctx sdk.Context) {
	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixVote)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
	iter.Close()
}

func (keeper Keeper) SetElect(ctx sdk.Context, priceVote PriceVote) {
	store := ctx.KVStore(keeper.key)
	key := GetElectKey(priceVote.FeedMsg.Denom)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(priceVote)
	store.Set(key, bz)
}

//nolint
func (keeper Keeper) GetElect(ctx sdk.Context, denom string) (res PriceVote) {
	store := ctx.KVStore(keeper.key)
	key := GetElectKey(denom)
	bz := store.Get(key)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}
