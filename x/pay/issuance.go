package pay

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func (k Keeper) setIssuance(ctx sdk.Context, denom string, issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(issuance)
	store.Set(KeyIssuance(denom, util.GetEpoch(ctx)), bz)
}

func (k Keeper) subtractIssuance(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		issuance := k.GetIssuance(ctx, coin.Denom, util.GetEpoch(ctx))
		issuance = issuance.Sub(coin.Amount)
		k.setIssuance(ctx, coin.Denom, issuance)
	}
}

func (k Keeper) addIssuance(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		issuance := k.GetIssuance(ctx, coin.Denom, util.GetEpoch(ctx))
		issuance = issuance.Add(coin.Amount)
		k.setIssuance(ctx, coin.Denom, issuance)
	}
}

func (k Keeper) GetIssuance(ctx sdk.Context, denom string, epoch sdk.Int) (issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeyIssuance(denom, util.GetEpoch(ctx)))
	if bz == nil {
		if epoch.Equal(sdk.ZeroInt()) {
			countIssuance := func(acc auth.Account) (stop bool) {
				issuance = issuance.Add(acc.GetCoins().AmountOf(denom))
				return false
			}
			k.ak.IterateAccounts(ctx, countIssuance)
			k.setIssuance(ctx, denom, issuance)
		} else {
			issuance = k.GetIssuance(ctx, denom, epoch.Sub(sdk.OneInt()))
		}
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &issuance)
	}

	return
}
