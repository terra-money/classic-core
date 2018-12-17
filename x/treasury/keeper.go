package treasury

import (
	"terra/types/tax"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the treasury store
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	tk tax.TaxKeeper
}

// NewKeeper constructs a new keeper
func NewKeeper(key sdk.StoreKey, cdc *codec.Codec, taxKeeper tax.TaxKeeper) Keeper {
	return Keeper{
		key: key,
		cdc: cdc,
		tk:  taxKeeper,
	}
}

func (keeper Keeper) CollectRevenues(ctx sdk.Context, revenue sdk.Coin) {
	debtRatio := keeper.tk.GetDebtRatio(ctx)
	budgetSubsidyRate := sdk.OneDec().Sub(debtRatio)

	subsidy := (sdk.NewDecFromInt(revenue.Amount).Mul(budgetSubsidyRate)).RoundInt()

	subsidyPool := keeper.GetSubsidyPool(ctx)
	subsidyPool = append(subsidyPool, sdk.NewCoin(revenue.Denom, subsidy))

	keeper.SetSubsidyPool(ctx, subsidyPool)
	// Let the rest of Coins burn
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
