package pay

import (
	"terra/types/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// GetIssuance fetches the total issuance count of the coin matching {denom}. If the {epoch} applies
// to a previous period, fetches the last stored snapshot issuance of the coin. For virgin calls,
// iterates through the accountkeeper and computes the genesis issuance.
func (k Keeper) GetIssuance(ctx sdk.Context, denom string, epoch sdk.Int) (issuance sdk.Int) {
	store := ctx.KVStore(k.key)
	bz := store.Get(keyIssuance(denom, util.GetEpoch(ctx)))
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
