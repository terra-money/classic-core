package treasury

import (
	"terra/types/assets"
	"terra/types/tax"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	tk tax.Keeper
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	taxKeeper tax.Keeper) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		tk:  taxKeeper,
	}
}

// Logic for Income
//------------------------------------
//------------------------------------
//------------------------------------

// AddIncome adds income to the treasury module
func (k Keeper) AddIncome(ctx sdk.Context, income sdk.Coins) {
	incomePool := k.getIncomePool(ctx)
	incomePool = incomePool.Plus(income)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(incomePool)
	store.Set(KeyIncomePool, bz)

	// If greater than some threshold
	threshold := sdk.NewInt(assets.LunaTargetIssuance).Div(sdk.NewInt(100))
	if incomePool.AmountOf(assets.LunaDenom).GT(threshold) {
		k.SettleClaims(ctx)
	}
}

func (k Keeper) getIncomePool(ctx sdk.Context) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyIncomePool)
	if bz == nil {
		res = sdk.Coins{}
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (k Keeper) resetIncomePool(ctx sdk.Context) {
	store := ctx.KVStore(k.key)
	store.Delete(KeyIncomePool)
}

// Logic for Claims
//------------------------------------
//------------------------------------
//------------------------------------

// AddClaim adds a claim to the store. Claims will be reset at SettleShares.
func (k Keeper) AddClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(GetClaimKey(claim.ID()), bz)

	curTally := k.getClaimsTally(ctx)
	newTally := curTally.Add(claim.GetWeight())
	bz2 := k.cdc.MustMarshalBinaryLengthPrefixed(newTally)
	store.Set(KeyClaimsTally, bz2)
}

// Total weights of all the claims in the store
func (k Keeper) getClaimsTally(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyClaimsTally)
	if bz == nil {
		res = sdk.ZeroDec()
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

// Sets the weight of all the claims in the store
func (k Keeper) setClaimTally(ctx sdk.Context, tally sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(tally)
	store.Set(KeyClaimsTally, bz)
}

func dividePool(ratio sdk.Dec, pool sdk.Coins) sdk.Coins {
	if len(pool) != 1 {
		return nil
	}

	return sdk.Coins{sdk.NewCoin(pool[0].Denom, ratio.MulInt(pool[0].Amount).TruncateInt())}
}

// SettleClaims settles and pays out claims in the store.
func (k Keeper) SettleClaims(ctx sdk.Context) {
	incomePool := k.getIncomePool(ctx)
	claimsTally := k.getClaimsTally(ctx)

	// Convert the entire income pool for Terra tokens
	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		shareWeight := sdk.OneDec().Sub(k.tk.GetDebtRatio(ctx))
		claimWeight := shareWeight.Mul(claim.GetWeight()).Quo(claimsTally)
		claimCoin := dividePool(claimWeight, incomePool)
		claim.Settle(ctx, k.tk, claimCoin)

		store.Delete(claimIter.Key())
	}
	claimIter.Close()

	// reset income pool and set the claim tally
	k.resetIncomePool(ctx)
	k.setClaimTally(ctx, sdk.ZeroDec())
}
