package oracle

import (
	"terra/types/assets"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// StoreKey is string representation of the store key for oracle
const StoreKey = "oracle"

// Keeper of the oracle store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	valset     sdk.ValidatorSet
	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, valset sdk.ValidatorSet, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc: cdc,
		key: key,

		valset:     valset,
		paramSpace: paramspace.WithKeyTable(ParamKeyTable()),
	}
}

//-----------------------------------
// Votes logic

func (k Keeper) getVotes(ctx sdk.Context) (votes map[string]PriceBallot) {
	handler := func(vote PriceVote) (stop bool) {
		votes[vote.Denom] = append(votes[vote.Denom], vote)
		return false
	}
	k.iterateVotes(ctx, handler)

	return
}

// Iterate over votes
func (k Keeper) iterateVotes(ctx sdk.Context, handler func(vote PriceVote) (stop bool)) {
	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixVote)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

func (k Keeper) getVote(ctx sdk.Context, denom string, voter sdk.AccAddress) (vote PriceVote) {
	store := ctx.KVStore(k.key)
	b := store.Get(KeyVote(denom, voter))
	if b == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vote)
	return
}

func (k Keeper) addVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(KeyVote(vote.Denom, vote.Voter), bz)
}

func (k Keeper) deleteVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(k.key)
	store.Delete(KeyVote(vote.Denom, vote.Voter))
}

//-----------------------------------
// Drop counter logic

func (k Keeper) setDropCounter(ctx sdk.Context, denom string, counter sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(counter)
	store.Set(KeyDropCounter(denom), bz)
}

func (k Keeper) deleteDropCounter(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.key)
	store.Delete(KeyDropCounter(denom))
}

func (k Keeper) getDropCounter(ctx sdk.Context, denom string) (counter sdk.Int) {
	store := ctx.KVStore(k.key)
	b := store.Get(KeyDropCounter(denom))
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &counter)
	return
}

//-----------------------------------
// Price logic

func (k Keeper) setPrice(ctx sdk.Context, denom string, price sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(KeyPrice(denom), bz)
}

func (k Keeper) deletePrice(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.key)
	store.Delete(KeyPrice(denom))
}

func (k Keeper) GetPrice(ctx sdk.Context, denom string) (price sdk.Dec, err sdk.Error) {
	if denom == assets.LunaDenom {
		return sdk.OneDec(), nil
	}

	store := ctx.KVStore(k.key)
	b := store.Get(KeyPrice(denom))
	if b == nil {
		return sdk.ZeroDec(), ErrUnknownDenomination(DefaultCodespace, denom)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

//______________________________________________________________________
// Params logic

// GetParams get oralce params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.Get(ctx, ParamStoreKeyParams, &params)
	return params
}

// SetParams set oracle params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, ParamStoreKeyParams, &params)
}
