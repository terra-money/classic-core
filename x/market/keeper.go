package market

import (
	"github.com/terra-project/core/x/mint"
	"github.com/terra-project/core/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper holds data structures for the market module
type Keeper struct {
	ok         oracle.Keeper
	mk         mint.Keeper
	paramSpace params.Subspace
}

// NewKeeper creates a new Keeper for the market module
func NewKeeper(ok oracle.Keeper, mk mint.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		ok:         ok,
		mk:         mk,
		paramSpace: paramspace.WithKeyTable(paramKeyTable()),
	}
}

// GetSwapCoins returns the amount of asked coins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Returns an Error if the swap is recursive, or the coins to be traded are unknown by the oracle, or the amount
// to trade is too small.
func (k Keeper) GetSwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (sdk.Coin, sdk.Error) {
	offerRate, err := k.ok.GetLunaSwapRate(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.ok.GetLunaSwapRate(ctx, askDenom)
	if err != nil {
		return sdk.Coin{}, ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(askRate).Quo(offerRate).TruncateInt()
	if retAmount.Equal(sdk.ZeroInt()) {
		return sdk.Coin{}, ErrInsufficientSwapCoins(DefaultCodespace, offerCoin.Amount)
	}

	return sdk.NewCoin(askDenom, retAmount), nil
}

// GetSwapDecCoins returns the amount of asked DecCoins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Similar to GetSwapCoins, but operates over sdk.DecCoins for convinience and accuracy.
func (k Keeper) GetSwapDecCoins(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error) {
	offerRate, err := k.ok.GetLunaSwapRate(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.DecCoin{}, ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.ok.GetLunaSwapRate(ctx, askDenom)
	if err != nil {
		return sdk.DecCoin{}, ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := offerCoin.Amount.Mul(askRate).Quo(offerRate)
	if retAmount.LTE(sdk.ZeroDec()) {
		return sdk.DecCoin{}, ErrInsufficientSwapCoins(DefaultCodespace, offerCoin.Amount.TruncateInt())
	}

	return sdk.NewDecCoinFromDec(askDenom, retAmount), nil
}

//-----------------------------------
// Params logic

// GetParams get budget params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var resultParams Params
	k.paramSpace.Get(ctx, paramStoreKeyParams, &resultParams)
	return resultParams
}

// SetParams set budget params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.Set(ctx, paramStoreKeyParams, &params)
}
