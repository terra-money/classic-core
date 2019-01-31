package treasury

import (
	"terra/types/assets"
	"terra/types/tax"
	"terra/types/util"

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

// Logic for shares
//------------------------------------
//------------------------------------
//------------------------------------

func (k Keeper) GetShare(ctx sdk.Context, shareID string) (res Share, err sdk.Error) {
	store := ctx.KVStore(k.key)
	bz := store.Get(GetShareKey(shareID))
	if bz == nil {
		err = ErrNoShareFound(DefaultCodespace, shareID)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (k Keeper) ResetShares(ctx sdk.Context, shares []Share) sdk.Error {
	// Ensure the weights sum to below 1
	totalWeight := sdk.ZeroDec()
	for _, share := range shares {
		totalWeight.Add(share.GetWeight())
	}
	if totalWeight.GT(sdk.OneDec()) {
		return ErrExcessiveWeight(DefaultCodespace, totalWeight)
	}

	// Clear existing shares
	util.Clear(k.key, ctx, PrefixShare)

	// Set shares to the store
	for _, share := range shares {
		util.Set(k.key, k.cdc, ctx, GetShareKey(share.ID()), share)
	}

	return nil
}

func dividePool(ratio sdk.Dec, pool sdk.Coins) sdk.Coins {
	if len(pool) != 1 {
		return nil
	}

	return sdk.Coins{sdk.NewCoin(pool[0].Denom, ratio.MulInt(pool[0].Amount).TruncateInt())}
}

func (k Keeper) SettleShares(ctx sdk.Context) {
	incomePool := k.getIncomePool(ctx)

	store := ctx.KVStore(k.key)

	shareIter := sdk.KVStorePrefixIterator(store, PrefixShare)
	for ; shareIter.Valid(); shareIter.Next() {
		var share Share
		k.cdc.MustUnmarshalBinaryLengthPrefixed(shareIter.Value(), &share)

		claimIter := sdk.KVStorePrefixIterator(store, GetClaimsForSharePrefix(share.ID()))
		for ; claimIter.Valid(); claimIter.Next() {
			var claim Claim
			k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

			claimWeight := share.GetWeight().Mul(claim.GetWeight())
			claimCoin := dividePool(claimWeight, incomePool)
			claim.Settle(ctx, k.tk, claimCoin)

			store.Delete(claimIter.Key())
		}
		claimIter.Close()

	}
	shareIter.Close()

	// shares := util.Collect(k.key, k.cdc, ctx, PrefixShare)

	// incomePool, err := k.getIncomePool(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// residualPool := incomePool

	// for _, share := range shares {
	// 	share := share.(Share)
	// 	sharePool := dividePool(share.GetWeight(), incomePool)

	// 	claims := util.Collect(k.key, k.cdc, ctx, GetClaimsForSharePrefix(share.ID()))

	// 	totalWeight := sdk.ZeroDec()
	// 	for _, c := range claims {
	// 		c := c.(Claim)
	// 		totalWeight = totalWeight.Add(c.GetWeight())
	// 	}

	// 	// Settle claims with others
	// 	for _, c := range claims {
	// 		c := c.(Claim)
	// 		adjustedWeight := c.GetWeight().Quo(totalWeight)
	// 		claimCoin := dividePool(adjustedWeight, sharePool)
	// 		c.Settle(ctx, k.tk, claimCoin)

	// 		residualPool.Minus(claimCoin)
	// 	}
	// }

	// // Set remaining coins as the remaining income pool
	// util.Set(k.key, k.cdc, ctx, KeyIncomePool, residualPool)
}

// Logic for Income Pool
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

	if incomePool.AmountOf(assets.LunaDenom).GT(sdk.ZeroInt()) {

	}

}

func (k Keeper) getIncomePool(ctx sdk.Context) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyIncomePool)
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

func (k Keeper) AddClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(GetClaimKey(claim.ShareID(), claim.ID()), bz)
}
