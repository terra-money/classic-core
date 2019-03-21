package treasury

import (
	"terra/types"
	"terra/types/util"
	"terra/x/market"
	"terra/x/mint"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the treasury store
type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey

	mtk mint.Keeper
	mk  market.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey,
	mtk mint.Keeper, mk market.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		mtk:        mtk,
		mk:         mk,
		paramSpace: paramspace.WithKeyTable(paramKeyTable()),
	}
}

//-----------------------------------
// Reward weight logic

// SetRewardWeight sets the ratio of the treasury that goes to mining rewards, i.e.
// supply of Luna that is burned.
func (k Keeper) SetRewardWeight(ctx sdk.Context, weight sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(weight)
	store.Set(keyRewardWeight(util.GetEpoch(ctx)), bz)
}

// GetRewardWeight returns the mining reward weight
func (k Keeper) GetRewardWeight(ctx sdk.Context, epoch sdk.Int) (rewardWeight sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyRewardWeight(epoch))
	if bz == nil {
		rewardWeight = k.GetParams(ctx).RewardPolicy.RateMin
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewardWeight)
	}
	return
}

//-----------------------------------
// Claims logic

func (k Keeper) addClaim(ctx sdk.Context, claim types.Claim) {
	store := ctx.KVStore(k.key)
	claimKey := keyClaim(claim.ID())

	// If the recipient has an existing claim in the same class, add to the previous claim
	if bz := store.Get(claimKey); bz != nil {
		var prevClaim types.Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, prevClaim)
		claim.Weight = claim.Weight.Add(prevClaim.Weight)
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(claimKey, bz)
}

func (k Keeper) iterateClaims(ctx sdk.Context, handler func(types.Claim) (stop bool)) {
	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)

	defer claimIter.Close()
	for ; claimIter.Valid(); claimIter.Next() {
		var claim types.Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		if handler(claim) {
			break
		}
	}
}

//-----------------------------------
// Params logic

// GetParams get treasury params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var resultParams Params
	k.paramSpace.Get(ctx, paramStoreKeyParams, &resultParams)
	return resultParams
}

// SetParams set treasury params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, paramStoreKeyParams, &params)
}

//-----------------------------------
// Tax logic

// SetTaxRate sets the tax rate; called from the treasury.
func (k Keeper) SetTaxRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
	store.Set(keyTaxRate, bz)
}

// GetTaxRate gets the tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context) (rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxRate)
	if bz == nil {
		rate = k.GetParams(ctx).TaxPolicy.RateMin
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rate)
	}
	return
}

// SetTaxCap sets the Tax Cap. Denominated in integer units of the reference {denom}
func (k Keeper) SetTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cap)
	store.Set(keyTaxCap(denom), bz)
}

// GetTaxCap gets the Tax Cap. Denominated in integer units of the reference {denom}
func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxCap(denom))
	if bz == nil {

		// Tax cap does not exist for the asset; compute it by
		// comparing it with the tax cap for TerraSDR
		sdrTaxCap := k.GetParams(ctx).TaxPolicy.Cap
		taxCap, err := k.mk.SwapCoins(ctx, sdrTaxCap, denom)

		// The coin is more valuable than TerraSDR. just set 1 as the cap.
		if err != nil {
			taxCap.Amount = sdk.OneInt()
		}

		k.SetTaxCap(ctx, denom, taxCap.Amount)
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &taxCap)
	}
	return
}

// RecordTaxProceeds add tax proceeds that have been added this epoch
func (k Keeper) RecordTaxProceeds(ctx sdk.Context, epoch sdk.Int, delta sdk.Coins) {
	proceeds := k.PeekTaxProceeds(ctx, epoch)
	proceeds = proceeds.Plus(delta)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proceeds)
	store.Set(keyTaxProceeds(epoch), bz)
}

// PeekTaxProceeds peeks the total amount of taxes that have been collected in the given epoch.
func (k Keeper) PeekTaxProceeds(ctx sdk.Context, epoch sdk.Int) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxProceeds(epoch))
	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}
