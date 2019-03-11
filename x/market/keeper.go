package market

import (
	"terra/x/oracle"
	"terra/x/pay"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

// StoreKey is string representation of the store key for market
const StoreKey = "market"

//nolint
type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	ok oracle.Keeper // Read terra & luna prices
	pk pay.Keeper
	dk distribution.Keeper
}

func NewKeeper(ok oracle.Keeper, pk pay.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		ok: ok,
		pk: pk,
		dk: dk,
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

func (k Keeper) SwapDecCoins(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error) {
	offerRate, err := k.ok.GetPrice(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.DecCoin{}, ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.ok.GetPrice(ctx, askDenom)
	if err != nil {
		return sdk.DecCoin{}, ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := offerCoin.Amount.Mul(offerRate).Quo(askRate)
	retCoin := sdk.NewDecCoinFromDec(askDenom, retAmount)
	return retCoin, nil
}
