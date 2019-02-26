package pay

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	KeyTaxRate = []byte("tax_rate")

	taxRateMin = sdk.ZeroDec()
	taxRateMax = sdk.NewDecWithPrec(2, 2) // 2%
)

func (k Keeper) SetTax(ctx sdk.Context, rate sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rate)
	store.Set(KeyTaxRate, bz)
}

func (k Keeper) GetTax(ctx sdk.Context) (res sdk.Dec) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyTaxRate)
	if bz == nil {
		res = sdk.NewDecWithPrec(1, 3)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func calculateTaxes(ctx sdk.Context, keeper Keeper, principal sdk.Coins) sdk.Coins {
	taxes := sdk.Coins{}
	for _, coin := range principal {
		taxRate := keeper.GetTax(ctx)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()

		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
	}

	return taxes
}
