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

func (k Keeper) GetShare(ctx sdk.Context, shareID string) Share {
	return util.Get(k.key, k.cdc, ctx, GetShareKey(shareID)).(Share)
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
	shares := util.Collect(k.key, k.cdc, ctx, PrefixShare)

	incomePool := util.Get(k.key, k.cdc, ctx, KeyIncomePool).(sdk.Coins)
	residualPool := incomePool

	for _, share := range shares {
		share := share.(Share)
		sharePool := dividePool(share.GetWeight(), incomePool)

		claims := util.Collect(k.key, k.cdc, ctx, GetClaimsForSharePrefix(share.ID()))

		totalWeight := sdk.ZeroDec()
		for _, c := range claims {
			c := c.(Claim)
			totalWeight = totalWeight.Add(c.GetWeight())
		}

		// Settle claims with others
		for _, c := range claims {
			c := c.(Claim)
			adjustedWeight := c.GetWeight().Quo(totalWeight)
			claimCoin := dividePool(adjustedWeight, sharePool)
			c.Settle(ctx, k.tk, claimCoin)

			residualPool.Minus(claimCoin)

			util.Delete(k.key, ctx, GetClaimKey(share.ID(), c.ID()))
		}
	}

	// Set remaining coins as the remaining income pool
	util.Set(k.key, k.cdc, ctx, KeyIncomePool, residualPool)
}

// Logic for Income Pool
//------------------------------------
//------------------------------------
//------------------------------------

// AddIncome adds income to the treasury module
func (k Keeper) AddIncome(ctx sdk.Context, income sdk.Coins) sdk.Error {

	taxDenom := income[0].Denom

	// Error if income is not paid in Terra tokens
	if taxDenom != assets.TerraDenom {
		return ErrWrongTaxDenomination(DefaultCodespace, taxDenom)
	}

	incomePool := util.Get(k.key, k.cdc, ctx, KeyIncomePool).(sdk.Coins)
	incomePool = incomePool.Plus(income)

	util.Set(k.key, k.cdc, ctx, KeyIncomePool, incomePool)
	return nil
}

// Logic for Claims
//------------------------------------
//------------------------------------
//------------------------------------

func (k Keeper) AddClaim(ctx sdk.Context, claim Claim) {
	util.Set(k.key, k.cdc, ctx, GetClaimKey(claim.ShareID(), claim.ID()), claim)
}
