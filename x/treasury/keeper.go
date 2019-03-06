package treasury

import (
	"terra/x/market"
	"terra/x/pay"

	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// StoreKey is string representation of the store key for treasury
const StoreKey = "treasury"

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	pk pay.Keeper
	mk market.Keeper
	dk distribution.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	pk pay.Keeper, mk market.Keeper, dk distribution.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		pk:         pk,
		mk:         mk,
		dk:         dk,
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

// Logic for Claims
//------------------------------------
//------------------------------------
//------------------------------------

func (k Keeper) addClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.key)
	claimKey := KeyClaim(claim.id)

	// If the recipient has an existing claim in the same class, add to the previous claim
	if bz := store.Get(claimKey); bz != nil {
		var prevClaim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, prevClaim)
		claim.weight = claim.weight.Add(prevClaim.weight)
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(claimKey, bz)
}

func (k Keeper) iterateClaims(ctx sdk.Context, handler func(Claim) (stop bool)) {
	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		if handler(claim) {
			break
		}
	}
	claimIter.Close()
}

func (k Keeper) sumClaims(ctx sdk.Context, class ClaimClass) (weightSumForClass sdk.Int, claimsForClass []Claim) {
	k.iterateClaims(ctx, func(claim Claim) (stop bool) {
		if claim.class == class {
			weightSumForClass = weightSumForClass.Add(claim.weight)
			claimsForClass = append(claimsForClass, claim)
		}
		return false
	})
	return
}

func (k Keeper) clearClaims(ctx sdk.Context) {
	store := ctx.KVStore(k.key)
	k.iterateClaims(ctx, func(claim Claim) (stop bool) {
		claimKey := KeyClaim(claim.id)
		store.Delete(claimKey)
		return false
	})
	return
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
