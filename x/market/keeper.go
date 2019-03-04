package market

import (
	"terra/types/util"
	"terra/x/oracle"
	"terra/x/pay"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	ok oracle.Keeper // Read terra & luna prices
	pk pay.Keeper
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(ok oracle.Keeper, pk pay.Keeper) Keeper {
	return Keeper{
		ok: ok,
		pk: pk,
	}
}

func (k Keeper) SwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (sdk.Coin, sdk.Error) {
	offerRate, err := k.ok.GetPrice(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.ok.GetPrice(ctx, askDenom)
	if err != nil {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(offerRate).Quo(askRate).RoundInt()
	if retAmount.Equal(sdk.ZeroInt()) {
		// drop in this scenario
		return sdk.Coin{}, ErrInsufficientSwapCoins(DefaultCodespace, offerCoin.Amount)
	}

	retCoin := sdk.Coin{Denom: askDenom, Amount: retAmount}
	return retCoin, nil
}

func (k Keeper) recordSeigniorage(ctx sdk.Context, seigniorage sdk.Coins) {
	currentEpoch := util.GetEpoch(ctx)
	pool := k.GetSeigniorage(ctx, currentEpoch)
	pool = pool.Plus(seigniorage)

	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(seigniorage)
	store.Set(KeySeigniorage(currentEpoch), bz)
}

func (k Keeper) GetSeigniorage(ctx sdk.Context, epoch sdk.Int) (res sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(KeySeigniorage(epoch))
	if bz == nil {
		res = sdk.Coins{}
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return
}
