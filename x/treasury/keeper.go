package treasury

import (
	"terra/types/util"
	"terra/x/market"

	"github.com/cosmos/cosmos-sdk/x/bank"
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

	bk bank.Keeper
	mk market.Keeper
	dk distribution.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	bk bank.Keeper, mk market.Keeper, dk distribution.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		bk:         bk,
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
	claimKey := KeyClaim(claim.ID())

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
		claimKey := KeyClaim(claim.ID())
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

//______________________________________________________________________
// Issuance logic

// GetIssuance fetches the total issuance count of the coin matching {denom}. If the {epoch} applies
// to a previous period, fetches the last stored snapshot issuance of the coin. For virgin calls,
// iterates through the accountkeeper and computes the genesis issuance.
func (k Keeper) GetIssuance(ctx sdk.Context, denom string, epoch sdk.Int) (issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyIssuance(denom, util.GetEpoch(ctx)))
	if bz == nil {

		// Genesis epoch; nothing exists in store so we must read it
		// from accountkeeper
		if epoch.Equal(sdk.ZeroInt()) {
			// countIssuance := func(acc auth.Account) (stop bool) {
			// 	issuance = issuance.Add(acc.GetCoins().AmountOf(denom))
			// 	return false
			// }
			//k.ak.IterateAccounts(ctx, countIssuance)
			//k.setIssuance(ctx, denom, issuance)
		} else {

			// Fetch the issuance snapshot of the previous epoch
			issuance = k.GetIssuance(ctx, denom, epoch.Sub(sdk.OneInt()))
		}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &issuance)
	}

	return
}

// sets the issuance in the store
func (k Keeper) setIssuance(ctx sdk.Context, denom string, issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(keyIssuance(denom, util.GetEpoch(ctx)), bz)
}

// convinience function. substracts the issuance counter in the store.
func (k Keeper) subtractIssuance(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		issuance := k.GetIssuance(ctx, coin.Denom, util.GetEpoch(ctx))
		issuance = issuance.Sub(coin.Amount)
		k.setIssuance(ctx, coin.Denom, issuance)
	}
}

// convinience function. adds to the issuance counter in the store.
func (k Keeper) addIssuance(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		issuance := k.GetIssuance(ctx, coin.Denom, util.GetEpoch(ctx))
		issuance = issuance.Add(coin.Amount)
		k.setIssuance(ctx, coin.Denom, issuance)
	}
}

//______________________________________________________________________
// Tax logic

// SetTaxRate sets the tax rate; called from the treasury.
func (k Keeper) SetTaxRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
	store.Set(keyTaxRate, bz)
}

// GetTaxRate gets the tax rate
func (k Keeper) GetTaxRate(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxRate)
	if bz == nil {
		res = sdk.ZeroDec()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
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
func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) (res sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyTaxCap(denom))
	if bz == nil {
		res = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

// AddTaxProceeds add tax proceeds that have been added this epoch
func (k Keeper) AddTaxProceeds(ctx sdk.Context, epoch sdk.Int, addition sdk.Coins) {
	proceeds := k.PeekTaxProceeds(ctx, epoch)
	proceeds = proceeds.Plus(addition)

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
