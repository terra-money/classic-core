package treasury

import (
	"terra/types/assets"
	"terra/x/market"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	mk market.Keeper
	ok oracle.Keeper
}

var (
	taxMin = sdk.ZeroDec()
	taxMax = sdk.NewDecWithPrec(2, 2) // 2%
)

// GetTaxRate gets the effective stability tax rate. TODO: add substrate
func (keeper Keeper) GetTaxRate(ctx sdk.Context, denom string) sdk.Dec {
	debtRatio := keeper.getDebtRatio(ctx)

	return taxMin.Add(taxMax.Sub(taxMin).Mul(debtRatio))
}

// getDebtRatio returns the ratio of debt and luna total issuance.
func (keeper Keeper) getDebtRatio(ctx sdk.Context) sdk.Dec {

	lunaCurrentIssuance := keeper.mk.GetCoinSupply(ctx, assets.LunaDenom)
	lunaTargetIssuance := keeper.GetLunaTargetIssuance(ctx)

	lunaDebt := lunaCurrentIssuance.Sub(lunaTargetIssuance)

	return sdk.NewDecFromInt(lunaDebt).Quo(sdk.NewDecFromInt(lunaCurrentIssuance))
}

// func (keeper Keeper) CollectTax(tax sdk.Coins) {
// 	debtRatio := keeper.getDebtRatio(ctx)

// }

func (keeper Keeper) CollectRevenues(ctx sdk.Context, revenue sdk.Coin) {
	debtRatio := keeper.getDebtRatio(ctx)
	budgetSubsidyRate := sdk.OneDec().Sub(debtRatio)

	subsidy := (sdk.NewDecFromInt(revenue.Amount).Mul(budgetSubsidyRate)).RoundInt()

	subsidyPool := keeper.GetSubsidyPool(ctx)
	subsidyPool = append(subsidyPool, sdk.NewCoin(revenue.Denom, subsidy))

	keeper.SetSubsidyPool(ctx, subsidyPool)
	// Let the rest of Coins burn
}

func (keeper Keeper) SetLunaTargetIssuance(ctx sdk.Context, genesisIssuance sdk.Int) {
	store := ctx.KVStore(keeper.key)
	key := KeyLunaTargetIssuance
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(genesisIssuance)
	store.Set(key, bz)
}

//nolint
func (keeper Keeper) GetLunaTargetIssuance(ctx sdk.Context) (res sdk.Int) {
	store := ctx.KVStore(keeper.key)
	key := KeyLunaTargetIssuance
	bz := store.Get(key)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}

func (keeper Keeper) SetSubsidyPool(ctx sdk.Context, subsidyPool sdk.Coins) {
	store := ctx.KVStore(keeper.key)
	key := KeySubsidyPool
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(subsidyPool)
	store.Set(key, bz)
}

//nolint
func (keeper Keeper) GetSubsidyPool(ctx sdk.Context) (res sdk.Coins) {
	store := ctx.KVStore(keeper.key)
	key := KeySubsidyPool
	bz := store.Get(key)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}
