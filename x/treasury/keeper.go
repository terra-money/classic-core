package treasury

import (
	"terra/types/assets"
	"terra/types/util"
	"terra/x/market"
	"terra/x/pay"
	"terra/x/treasury/tags"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	pk pay.Keeper
	mk market.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	pk pay.Keeper, mk market.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		pk:         pk,
		mk:         mk,
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

// AddClaim adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) addClaim(ctx sdk.Context, claim Claim) {
	store := ctx.KVStore(k.key)
	claimKey := KeyClaim(claim.id)

	if bz := store.Get(claimKey); bz != nil {
		var prevClaim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, prevClaim)

		claim.weight = claim.weight.Add(prevClaim.weight)
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(claim)
	store.Set(KeyClaim(claim.id), bz)
}

// AddClaim adds a funding claim to the treasury. Settled around once a month.
func (k Keeper) ProcessClaims(ctx sdk.Context, class ClaimClass, rewardees map[string]sdk.Dec) {
	for rAddrStr, rewardWeight := range rewardees {
		addr, err := sdk.AccAddressFromBech32(rAddrStr)
		if err != nil {
			continue
		}

		k.addClaim(ctx, NewClaim(class, rewardWeight, addr))
	}
}

func getRewardPools(ctx sdk.Context, k Keeper) (claimPools map[ClaimClass]sdk.Int) {
	totalPool := k.mk.GetSeigniorage(ctx, util.GetEpoch(ctx)).AmountOf(assets.LunaDenom)

	minerWeight := k.GetRewardWeight(ctx)
	claimWeight := sdk.OneDec().Sub(minerWeight)

	claimPools[MinerClaimClass] = minerWeight.MulInt(totalPool).TruncateInt()

	for class, share := range k.GetParams(ctx).ClaimShares {
		claimPools[class] = claimWeight.MulInt(totalPool).Mul(share).TruncateInt()
	}
	return
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	claimPools := getRewardPools(ctx, k)

	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		claim.Settle(ctx, k, sdk.Coins{
			sdk.NewCoin(assets.LunaDenom,
				claim.weight.MulInt(claimPools[claim.class]).TruncateInt(),
			),
		})

		store.Delete(claimIter.Key())
	}
	claimIter.Close()

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, claimPools[MinerClaimClass],
		tags.Oracle, claimPools[OracleClaimClass],
		tags.Budget, claimPools[BudgetClaimClass],
	)
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
