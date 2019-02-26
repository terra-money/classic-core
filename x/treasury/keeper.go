package treasury

import (
	"terra/x/pay"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	pk pay.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	pk pay.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		pk:         pk,
		paramSpace: paramspace.WithKeyTable(ParamKeyTable()),
	}
}

// SetRewardWeight sets the ratio of the treasury that goes to mining rewards, i.e.
// supply of Luna that is burned.
func (k Keeper) SetRewardWeight(ctx sdk.Context, weight sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(weight)
	store.Set(KeyRewardWeight, bz)
}

// GetRewardWeight returns the mining reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyRewardWeight)
	if bz == nil {
		panic(nil)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// GetIssuance gets the current supply of a coin with denom
func (k Keeper) GetIssuance(ctx sdk.Context, denom string) (res sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyIssuance(denom))
	if bz == nil {
		panic(nil)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// SetIssuance sets the issuance of the coin
func (k Keeper) SetIssuance(ctx sdk.Context, denom string, issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(KeyIssuance(denom), bz)
}

//______________________________________________________________________
// Params logic

// GetParams get treasury params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var params Params
	k.paramSpace.Get(ctx, ParamStoreKeyParams, &params)
	return params
}

// SetParams set treasury params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, ParamStoreKeyParams, &params)
}
