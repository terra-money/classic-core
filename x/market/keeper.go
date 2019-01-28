package market

import (
	"terra/x/oracle"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
type Keeper struct {
	storeKey  sdk.StoreKey      // Key to our module's store
	codespace sdk.CodespaceType // Reserves space for error codes
	cdc       *codec.Codec      // Codec to encore/decode structs
	ok        oracle.Keeper     // Read terra & luna prices
	tk        treasury.Keeper   // Pay mint revenues to the treasury
	bk        bank.Keeper
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

func (k Keeper) SwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (sdk.Coin, sdk.Error) {
	// If swap msg for not whitelisted denom
	if !whitelistContains(ctx, k, offerCoin.Denom) {
		return sdk.Coin{}, ErrUnknownDenomination(DefaultCodespace, offerCoin.Denom)
	}

	offerRate, tErr := k.ok.GetPriceTarget(ctx, offerCoin.Denom)
	if tErr != nil {
		panic(tErr)
	}

	askRate, oErr := k.ok.GetPriceObserved(ctx, askDenom)
	if oErr != nil {
		panic(oErr)
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
