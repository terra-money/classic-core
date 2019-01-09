package oracle

import (
	"terra/types/assets"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the oracle store
type Keeper struct {
	key        sdk.StoreKey
	cdc        *codec.Codec
	tk         treasury.Keeper
	valset     sdk.ValidatorSet
	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, tk treasury.Keeper, valset sdk.ValidatorSet, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		tk:         tk,
		valset:     valset,
		paramSpace: paramspace.WithTypeTable(ParamTypeTable()),
	}
}

//-----------------------------------
// Votes logic

func (keeper Keeper) getVotes(ctx sdk.Context, denom string) (votes []PriceVote) {
	handler := func(vote PriceVote) (stop bool) {
		votes = append(votes, vote)
		return false
	}
	keeper.iterateVotes(ctx, denom, handler)

	return
}

// Iterate over votes
func (keeper Keeper) iterateVotes(ctx sdk.Context, denom string, handler func(vote PriceVote) (stop bool)) {
	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, PrefixVote(denom))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote PriceVote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

func (keeper Keeper) addVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(KeyVote(vote.FeedMsg.Denom, vote.FeedMsg.Feeder), bz)
}

func (keeper Keeper) deleteVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)
	store.Delete(KeyVote(vote.FeedMsg.Denom, vote.FeedMsg.Feeder))
}

//-----------------------------------
// Price logic

func (keeper Keeper) setPriceTarget(ctx sdk.Context, denom string, targetPrice sdk.Dec) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(targetPrice)
	store.Set(KeyTargetPrice(denom), bz)
}

func (keeper Keeper) setPriceObserved(ctx sdk.Context, denom string, observedPrice sdk.Dec) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(observedPrice)
	store.Set(KeyObservedPrice(denom), bz)
}

func (keeper Keeper) GetPriceTarget(ctx sdk.Context, denom string) (targetPrice sdk.Dec) {
	if denom == assets.LunaDenom {
		return sdk.OneDec()
	}

	store := ctx.KVStore(keeper.key)
	b := store.Get(KeyTargetPrice(denom))
	if b == nil {
		return sdk.ZeroDec()
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(b, &targetPrice)
	return
}

func (keeper Keeper) GetPriceObserved(ctx sdk.Context, denom string) (observedPrice sdk.Dec) {
	if denom == assets.LunaDenom {
		return sdk.OneDec()
	}

	store := ctx.KVStore(keeper.key)
	b := store.Get(KeyObservedPrice(denom))
	if b == nil {
		return sdk.ZeroDec()
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(b, &observedPrice)
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
