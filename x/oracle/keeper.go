package oracle

import (
	"terra/types/assets"
	"terra/types/util"
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

func (keeper Keeper) getVotes(ctx sdk.Context, denom string) (res []PriceVote) {
	votes := util.Collect(
		keeper.key,
		keeper.cdc,
		ctx,
		GetVotePrefix(denom),
	)

	for _, v := range votes {
		res = append(res, v.(PriceVote))
	}

	return
}

func (keeper Keeper) addVote(ctx sdk.Context, vote PriceVote) {
	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		GetVotePrefix(vote.FeedMsg.Denom),
		vote,
	)
}

func (keeper Keeper) clearVotes(ctx sdk.Context, denom string) {
	util.Clear(
		keeper.key,
		ctx,
		GetVotePrefix(denom),
	)
}

func (keeper Keeper) setPriceTarget(ctx sdk.Context, denom string, targetPrice sdk.Dec) {
	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		GetTargetPriceKey(denom),
		targetPrice,
	)
}

func (keeper Keeper) setPriceObserved(ctx sdk.Context, denom string, observedPrice sdk.Dec) {
	util.Set(
		keeper.key,
		keeper.cdc,
		ctx,
		GetObservedPriceKey(denom),
		observedPrice,
	)
}

func (keeper Keeper) GetPriceTarget(ctx sdk.Context, denom string) sdk.Dec {
	if denom == assets.LunaDenom {
		return sdk.OneDec()
	}

	return util.Get(
		keeper.key,
		keeper.cdc,
		ctx,
		GetTargetPriceKey(denom),
	).(sdk.Dec)
}

func (keeper Keeper) GetPriceObserved(ctx sdk.Context, denom string) sdk.Dec {
	if denom == assets.LunaDenom {
		return sdk.OneDec()
	}

	return util.Get(
		keeper.key,
		keeper.cdc,
		ctx,
		GetObservedPriceKey(denom),
	).(sdk.Dec)
}

//______________________________________________________________________

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
