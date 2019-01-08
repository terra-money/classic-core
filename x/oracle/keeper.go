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

func (keeper Keeper) getVoteIterator(ctx sdk.Context, denom string) (iterator sdk.Iterator) {
	store := ctx.KVStore(keeper.key)
	iterator = sdk.KVStorePrefixIterator(store, GetVotePrefix(denom))
	return
}

func (keeper Keeper) addVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(GetVoteKey(vote.FeedMsg.Denom, vote.FeedMsg.Feeder), bz)
}

func (keeper Keeper) deleteVote(ctx sdk.Context, vote PriceVote) {
	store := ctx.KVStore(keeper.key)
	store.Delete(GetVoteKey(vote.FeedMsg.Denom, vote.FeedMsg.Feeder))
}

//-----------------------------------
// Price logic

func (keeper Keeper) setPriceTarget(ctx sdk.Context, denom string, targetPrice sdk.Dec) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(targetPrice)
	store.Set(GetTargetPriceKey(denom), bz)
}

func (keeper Keeper) setPriceObserved(ctx sdk.Context, denom string, observedPrice sdk.Dec) {
	store := ctx.KVStore(keeper.key)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(observedPrice)
	store.Set(GetObservedPriceKey(denom), bz)
}

func (keeper Keeper) GetPriceTarget(ctx sdk.Context, denom string) (targetPrice sdk.Dec) {
	if denom == assets.LunaDenom {
		return sdk.OneDec()
	}

	store := ctx.KVStore(keeper.key)
	b := store.Get(GetTargetPriceKey(denom))
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
	b := store.Get(GetObservedPriceKey(denom))
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
