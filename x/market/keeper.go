package market

import (
	"terra/x/oracle"
	"terra/x/treasury"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
type Keeper struct {
	storeKey sdk.StoreKey    // Key to our module's store
	ok       oracle.Keeper   // Read terra & luna prices
	tk       treasury.Keeper // Pay mint revenues to the treasury
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(
	ok oracle.Keeper,
	tk treasury.Keeper,
) Keeper {
	return Keeper{
		ok: ok,
		tk: tk,
	}
}

// SwapCoins swaps the offerCoin for the requisite amount of coins of the askDenom
// at the Target exchange rate for both the offered and asked coins.
// Returns an error if the ask is not registered.
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

	return sdk.Coin{Denom: askDenom, Amount: retAmount}, nil
}
