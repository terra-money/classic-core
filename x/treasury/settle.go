package treasury

import (
	"terra/types/assets"
	"terra/x/treasury/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Logic for Claims
//------------------------------------
//------------------------------------
//------------------------------------

// RequestFunds immediately settles requested amount of coins to the requester's wallet.
// Used when delayed settlement through adding claims is not appropriate.
func (k Keeper) RequestFunds(ctx sdk.Context, req sdk.Coin, requester sdk.AccAddress) (resTags sdk.Tags, err sdk.Error) {
	reqCoins := sdk.Coins{req}
	_, resTags, err = k.pk.AddCoins(ctx, requester, reqCoins)
	if err != nil {
		return
	}

	// Mint new tokens to satisfy the request
	issuance := k.GetIssuance(ctx, req.Denom)
	issuance = issuance.Add(req.Amount)
	k.SetIssuance(ctx, req.Denom, issuance)

	return
}

// AddClaim adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) AddClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(KeyClaim(claim.id), bz)
}

// rewardMiners indirectly rewards Terra miners by burning a portion of the treasury's income every month.
// Only burns until params.LunaTargetIssuance is reached, and no more.
// Returns excess rewards back to the caller.
func rewardMiners(ctx sdk.Context, k Keeper, reward sdk.Int) (remainder sdk.Int) {
	lunaIssuance := k.GetIssuance(ctx, assets.LunaDenom)
	lunaTargetIssuance := k.GetParams(ctx).LunaTargetIssuance
	excessLunaIssuance := lunaIssuance.Sub(lunaTargetIssuance)

	if excessLunaIssuance.GT(reward) {
		lunaIssuance = lunaIssuance.Sub(reward)
		k.SetIssuance(ctx, assets.LunaDenom, lunaIssuance)
	} else {
		k.SetIssuance(ctx, assets.LunaDenom, lunaTargetIssuance)
		remainder = reward.Sub(excessLunaIssuance)
	}

	return
}

// normalizeClaims normalizes weights of the claims over all the claims registered to the treasury
func normalizeClaims(k Keeper, ctx sdk.Context, scale map[byte]sdk.Dec) {
	store := ctx.KVStore(k.key)

	weightSumMap := map[byte]sdk.Dec{}
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		weightSumMap[claim.class] = weightSumMap[claim.class].Add(claim.weight)
	}
	claimIter.Close()

	claimIter = sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		claim.weight = claim.weight.Quo(weightSumMap[claim.class]).Mul(scale[claim.class])
		k.AddClaim(ctx, claim)
	}
	claimIter.Close()
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	incomePool := k.getIncomePool(ctx)
	if incomePool.Equal(sdk.ZeroInt()) {
		return
	}

	// Reward rest to the miners
	claimPool := rewardMiners(ctx, k, incomePool)

	normalizeClaims(k, ctx, map[byte]sdk.Dec{
		OracleClaimClass: sdk.NewDecWithPrec(1, 1),
		BudgetClaimClass: sdk.NewDecWithPrec(9, 1),
	})

	claimed := sdk.ZeroInt()

	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		claimSize := claim.weight.MulInt(claimPool).TruncateInt()
		claimCoin := sdk.NewCoin(assets.LunaDenom, claimSize)
		claim.Settle(ctx, k, sdk.Coins{claimCoin})

		claimed = claimed.Add(claimSize)

		store.Delete(claimIter.Key())
	}
	claimIter.Close()

	settleTags = settleTags.AppendTags(
		sdk.NewTags(
			tags.Action, tags.ActionSettle,
			tags.Amount, incomePool,
			tags.MinerReward, k.GetRewardWeight(ctx).Bytes(),
			tags.Oracle, sdk.OneDec().Sub(k.GetRewardWeight(ctx)).Mul(sdk.NewDecWithPrec(1, 1)).Bytes(),
			tags.Budget, sdk.OneDec().Sub(k.GetRewardWeight(ctx)).Mul(sdk.NewDecWithPrec(9, 1)).Bytes(),
		),
	)

	dust := claimPool.Sub(claimed)

	// Reset income pool
	k.setIncomePool(ctx, dust)

	return
}
