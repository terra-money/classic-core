package pay

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetTaxRate(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
	store.Set(KeyTaxRate, bz)
}

func (k Keeper) GetTaxRate(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyTaxRate)
	if bz == nil {
		res = sdk.NewDecWithPrec(1, 3)
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}

	return
}

func (k Keeper) SetTaxCap(ctx sdk.Context, denom string, cap sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cap)
	store.Set(KeyTaxCap(denom), bz)
}

func (k Keeper) GetTaxCap(ctx sdk.Context, denom string) (res sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyTaxRate)
	if bz == nil {
		res = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

func (k Keeper) recordTaxProceeds(ctx sdk.Context, taxProceeds sdk.Coins) {
	currentEpoch := util.GetEpoch(ctx)
	proceeds := k.PeekTaxProceeds(ctx, currentEpoch)
	proceeds = proceeds.Plus(taxProceeds)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proceeds)
	store.Set(KeyTaxProceeds(currentEpoch), bz)
}

func (k Keeper) PeekTaxProceeds(ctx sdk.Context, epoch sdk.Int) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyTaxProceeds(epoch))
	if bz == nil {
		res = sdk.Coins{}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	}
	return
}

func calculateTaxes(ctx sdk.Context, keeper Keeper, principal sdk.Coins) sdk.Coins {
	taxes := sdk.Coins{}
	for _, coin := range principal {
		taxRate := keeper.GetTaxRate(ctx)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()
		taxCap := keeper.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
	}

	return taxes
}
