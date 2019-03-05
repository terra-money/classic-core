package treasury

import (
	"terra/types/assets"
	"terra/x/market"
	"terra/x/pay"
	"terra/x/treasury/tags"

	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// StoreKey is string representation of the store key for treasury
const StoreKey = "treasury"

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	pk pay.Keeper
	mk market.Keeper
	dk distribution.Keeper

	paramSpace params.Subspace
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	pk pay.Keeper, mk market.Keeper, dk distribution.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		key:        key,
		cdc:        cdc,
		pk:         pk,
		mk:         mk,
		dk:         dk,
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
func (k Keeper) ProcessClaims(ctx sdk.Context, class ClaimClass, rewardees map[string]sdk.Int) {
	store := ctx.KVStore(k.key)

	for rAddrStr, rewardWeight := range rewardees {
		addr, err := sdk.AccAddressFromBech32(rAddrStr)
		if err != nil {
			continue
		}

		newClaim := NewClaim(class, sdk.NewDecFromInt(rewardWeight), addr)
		claimKey := KeyClaim(newClaim.id)

		// If the recipient has an existing claim in the same class, add to the previous claim
		if bz := store.Get(claimKey); bz != nil {
			var prevClaim Claim
			k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, prevClaim)

			newClaim.weight = newClaim.weight.Add(prevClaim.weight)
		}

		bz := k.cdc.MustMarshalBinaryLengthPrefixed(newClaim)
		store.Set(claimKey, bz)
	}
}

// settleClaims distributes the current treasury to the registered claims, and deletes all claims from the store.
func (k Keeper) settleClaims(ctx sdk.Context) (settleTags sdk.Tags) {
	totalPool := k.dk.GetFeePool(ctx)

	// Pay mining rewards; just burn Luna
	minerRewardWeight := k.GetRewardWeight(ctx)
	minerRewards := totalPool.CommunityPool.MulDec(minerRewardWeight)

	// Compute the size of oracle + budget claim reward pools
	params := k.GetParams(ctx)
	oracleRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.OracleClaimShare)
	oracleReward := totalPool.CommunityPool.MulDec(oracleRewardWeight)
	budgetRewardWeight := sdk.OneDec().Sub(minerRewardWeight).Mul(params.BudgetClaimShare)
	budgetReward := totalPool.CommunityPool.MulDec(budgetRewardWeight)

	// Sum the total amount of voting power accumulated in claims by class
	oracleVotingPowerSum := sdk.ZeroDec()
	budgetVotingPowerSum := sdk.ZeroDec()
	store := ctx.KVStore(k.key)
	claimIter := sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		switch claim.class {
		case OracleClaimClass:
			oracleVotingPowerSum = oracleVotingPowerSum.Add(claim.weight)
		case BudgetClaimClass:
			budgetVotingPowerSum = budgetVotingPowerSum.Add(claim.weight)
		}
	}
	claimIter.Close()

	// Reward claims
	remainder := sdk.DecCoins{}
	claimIter = sdk.KVStorePrefixIterator(store, PrefixClaim)
	for ; claimIter.Valid(); claimIter.Next() {
		var claim Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(claimIter.Value(), &claim)

		var claimReward sdk.DecCoins
		switch claim.class {
		case OracleClaimClass:
			claimReward = oracleReward.MulDec(claim.weight).QuoDec(oracleVotingPowerSum)
		case BudgetClaimClass:
			claimReward = budgetReward.MulDec(claim.weight).QuoDec(budgetVotingPowerSum)
		}

		// translate rewards to SDR
		rewardInSDR, err := k.mk.SwapDecCoins(
			ctx,
			sdk.NewDecCoinFromDec(assets.LunaDenom, claimReward.AmountOf(assets.LunaDenom)),
			assets.SDRDenom,
		)
		if err != nil {
			continue
		}
		rewardInSDRInt, dust := rewardInSDR.TruncateDecimal()

		// credit the recipient's account with the reward
		k.pk.AddCoins(ctx, claim.recipient, sdk.Coins{rewardInSDRInt})
		remainder = remainder.Plus(sdk.DecCoins{dust})

		store.Delete(claimIter.Key())
	}
	claimIter.Close()

	totalPool.CommunityPool = remainder
	k.dk.SetFeePool(ctx, totalPool)

	return sdk.NewTags(
		tags.Action, tags.ActionSettle,
		tags.MinerReward, minerRewards,
		tags.Oracle, oracleReward,
		tags.Budget, budgetReward,
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
