package auth

import (
	codec "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	collectedFeesKey = []byte("collectedFees")
	taxRateKey       = []byte("taxRate")
)

var (
	taxMaxRate      = sdk.NewInt(20) // 2%
	taxCeiling      = sdk.NewInt(10) // 10 Terra
	taxRateBase     = 1000
	taxExemptDenoms = []string{"luna"}
)

// FeeCollectionKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeCollectionKeeper struct {

	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
}

//nolint
func NewFeeCollectionKeeper(cdc *codec.Codec, key sdk.StoreKey) FeeCollectionKeeper {
	return FeeCollectionKeeper{
		key: key,
		cdc: cdc,
	}
}

// GetTaxRate retrieves the effective tax rate
func (fck FeeCollectionKeeper) GetTaxRate(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(fck.key)
	bz := store.Get(taxRateKey)
	if bz == nil {
		return sdk.ZeroInt()
	}

	var taxrate sdk.Int
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &taxrate)
	return taxrate
}

// GetTax computes the tax for a given amount of coins to be transferred
func (fck FeeCollectionKeeper) GetTax(ctx sdk.Context, principal sdk.Coins) sdk.Coins {
	taxRate := fck.GetTaxRate(ctx)

	tax := sdk.Coins{}

	// Skip over whitelisted denoms
	for _, tfd := range taxExemptDenoms {
		for _, coin := range principal {
			if coin.Denom != tfd {
				taxCharge := coin.Amount.Mul(taxRate).Div(sdk.NewInt(int64(taxRateBase)))

				// Enforce absolute ceiling on the amount of tax being charged for message
				if taxCharge.GT(taxCeiling) {
					taxCharge = taxCeiling
				}

				tax.Plus(
					sdk.Coins{
						sdk.NewCoin(coin.Denom, taxCharge),
					},
				)
			}
		}
	}

	return tax
}

// SetTaxRate sets the effective tax rate
func (fck FeeCollectionKeeper) SetTaxRate(ctx sdk.Context, taxrate sdk.Int) {
	// Roll out negative tax rates
	if taxrate.LT(sdk.ZeroInt()) {
		return
	}

	// enforce a maximum tax rate
	if taxrate.GT(taxMaxRate) {
		taxrate = taxMaxRate
	}

	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(taxrate)
	store := ctx.KVStore(fck.key)
	store.Set(taxRateKey, bz)
}

// GetCollectedFees retrieves the collected fee pool
func (fck FeeCollectionKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fck.key)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.Coins{}
	}

	feePool := &(sdk.Coins{})
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

// SetCollectedFees sets the collected fee pool
func (fck FeeCollectionKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fck.key)
	store.Set(collectedFeesKey, bz)
}

// AddCollectedFees add to the fee pool
func (fck FeeCollectionKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Plus(coins)
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}

// ClearCollectedFees clear the fee pool
func (fck FeeCollectionKeeper) ClearCollectedFees(ctx sdk.Context) {
	fck.setCollectedFees(ctx, sdk.Coins{})
}
