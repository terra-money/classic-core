package treasury

import (
	"terra/types/assets"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Logic for Income Pool
//------------------------------------
//------------------------------------
//------------------------------------

// AddIncome adds income to the treasury module
func (k Keeper) AddIncome(ctx sdk.Context, income sdk.Coin) {

	// If income is Luna, add it to the income pool
	if income.Denom == assets.LunaDenom {
		incomePool := k.getIncomePool(ctx)
		incomePool = incomePool.Add(income.Amount)

		k.setIncomePool(ctx, incomePool)
	} else {
		// Otherwise, burn them.
		issuance := k.GetIssuance(ctx, income.Denom)
		k.SetIssuance(ctx, income.Denom, issuance.Sub(income.Amount))
	}
}

func (k Keeper) getIncomePool(ctx sdk.Context) (res sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyIncomePool)
	if bz == nil {
		panic(nil)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (k Keeper) setIncomePool(ctx sdk.Context, pool sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(KeyIncomePool, bz)
}
