package market

import (
	"terra/x/oracle"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
type Keeper struct {
	storeKey sdk.StoreKey    // Key to our module's store
	ok       oracle.Keeper   // Read terra & luna prices
	tk       treasury.Keeper // Pay mint revenues to the treasury
	bk       bank.Keeper
}

// NewKeeper crates a new keeper with write and read access
func NewKeeper(
	ok oracle.Keeper,
	tk treasury.Keeper,
	bk bank.Keeper,
) Keeper {
	return Keeper{
		ok: ok,
		tk: tk,
		bk: bk,
	}
}

func whitelistContains(ctx sdk.Context, k Keeper, denom string) bool {
	whitelist := k.ok.GetParams(ctx).Whitelist
	for _, w := range whitelist {
		if w == denom {
			return true
		}
	}
	return false
}

// SwapCoins swaps the offerCoin for the requisite amount of coins of the askDenom
// at the Target exchange rate for both the offered and asked coins.
// Returns an error if the ask is not registered.
func (k Keeper) SwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (sdk.Coin, sdk.Error) {
	// If swap msg for not whitelisted denom
	if !whitelistContains(ctx, k, offerCoin.Denom) {
		return sdk.Coin{}, ErrUnknownDenomination(DefaultCodespace, offerCoin.Denom)
	}

	offerRate := k.ok.GetPriceTarget(ctx, offerCoin.Denom)
	if offerRate.Equal(sdk.ZeroDec()) {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate := k.ok.GetPriceTarget(ctx, askDenom)
	if askRate.Equal(sdk.ZeroDec()) {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(offerRate).Quo(askRate).RoundInt()

	if retAmount.Equal(sdk.ZeroInt()) {
		// drop in this scenario
		return sdk.Coin{}, ErrInsufficientSwapCoins(DefaultCodespace, offerCoin.Amount)
	}

	retCoin := sdk.Coin{
		Denom:  askDenom,
		Amount: retAmount,
	}

	return retCoin, nil
}
