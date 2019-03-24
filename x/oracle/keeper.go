package oracle

import (
	"terra/types/assets"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the oracle store
type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey

	valset     sdk.ValidatorSet
	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, valset sdk.ValidatorSet, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,

		valset:     valset,
		paramSpace: paramspace.WithKeyTable(paramKeyTable()),
	}
}

//-----------------------------------
// Votes logic

// collectVotes collects all oracle votes for the period, categorized by the votes' denom parameter
func (k Keeper) collectVotes(ctx sdk.Context) (votes map[string]PriceBallot) {
	votes = map[string]PriceBallot{}
	handler := func(vote PriceVote) (stop bool) {
		votes[vote.Denom] = append(votes[vote.Denom], vote)
		return false
	}
	k.iterateVotes(ctx, handler)

	return
}

// Iterate over votes in the store
func (k Keeper) iterateVotes(ctx sdk.Context, handler func(vote PriceVote) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefixVote)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Iterate over votes in the store
func (k Keeper) iterateVotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote PriceVote) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Retrieves a vote from the store
func (k Keeper) getVote(ctx sdk.Context, denom string, voter sdk.AccAddress) (vote PriceVote, err sdk.Error) {
	store := ctx.KVStore(k.key)
	b := store.Get(keyVote(denom, voter))
	if b == nil {
		err = ErrNoVote(DefaultCodespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vote)
	return
}

// Add a vote to the store
func (k Keeper) addVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(keyVote(vote.Denom, vote.Voter), bz)
}

// Delete a vote from the store
func (k Keeper) deleteVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(k.key)
	store.Delete(keyVote(vote.Denom, vote.Voter))
}

//-----------------------------------
// Drop counter logic

// Increment drop counter. Called when an oracle vote is illiquid.
func (k Keeper) incrementDropCounter(ctx sdk.Context, denom string) (counter sdk.Int) {
	store := ctx.KVStore(k.key)
	b := store.Get(keyDropCounter(denom))
	if b == nil {
		counter = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &counter)
	}

	// Increment counter
	counter = counter.Add(sdk.OneInt())
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(counter)
	store.Set(keyDropCounter(denom), bz)
	return
}

// resets the drop counter.
func (k Keeper) resetDropCounter(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.key)
	store.Delete(keyDropCounter(denom))
}

//-----------------------------------
// Price logic

// GetLunaSwapRate gets the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) GetLunaSwapRate(ctx sdk.Context, denom string) (price sdk.Dec, err sdk.Error) {
	if denom == assets.LunaDenom {
		return sdk.OneDec(), nil
	}

	store := ctx.KVStore(k.key)
	b := store.Get(keyPrice(denom))
	if b == nil {
		return sdk.ZeroDec(), ErrUnknownDenomination(DefaultCodespace, denom)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

// SetLunaSwapRate sets the consensus exchange rate of Luna denominated in the denom asset to the store.
func (k Keeper) SetLunaSwapRate(ctx sdk.Context, denom string, price sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(keyPrice(denom), bz)
}

// deletePrice deletes the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) deletePrice(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.key)
	store.Delete(keyPrice(denom))
}

//-----------------------------------
// Params logic

// GetParams get oracle params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.Get(ctx, paramStoreKeyParams, &params)
	return params
}

// SetParams set oracle params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, paramStoreKeyParams, &params)
}
