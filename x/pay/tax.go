package pay

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// payTax charges the stability tax on SendCoin and InputOutputCoins.
func (k Keeper) payTax(ctx sdk.Context, taxPayer sdk.AccAddress, principal sdk.Coins) (taxes sdk.Coins, taxTags sdk.Tags, err sdk.Error) {
	for _, coin := range principal {
		taxRate := k.GetTaxRate(ctx)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()
		taxCap := k.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		taxCoin := sdk.Coins{sdk.NewCoin(coin.Denom, taxDue)}

		_, payTags, err := subtractCoins(ctx, k.ak, taxPayer, taxCoin)
		if err != nil {
			return nil, nil, err
		}

		taxTags = taxTags.AppendTags(payTags)
		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
		k.fk.AddCollectedFees(ctx, taxCoin)
	}

	// Record tax income; can be retrieved by PeekTaxIncome
	currentEpoch := util.GetEpoch(ctx)
	proceeds := k.PeekTaxProceeds(ctx, currentEpoch)
	proceeds = proceeds.Plus(taxes)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(proceeds)
	store.Set(keyTaxProceeds(currentEpoch), bz)

	return
}

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
