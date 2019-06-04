package market

import (
	"github.com/terra-project/core/x/mint"
	"github.com/terra-project/core/x/oracle"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

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


// ComputeLunaDelta returns the issuance rate change of Luna for the day post-swap
func (k Keeper) ComputeLunaDelta(ctx sdk.Context, change sdk.Int) sdk.Dec {
	curDay := ctx.BlockHeight() / util.BlocksPerDay

	// Start limits on day 2
	if curDay != 0 {
		curIssuance := k.mk.GetIssuance(ctx, assets.MicroLunaDenom, sdk.NewInt(curDay))
		postSwapIssuance := curIssuance.Add(change)
		prevIssuance := k.mk.GetIssuance(ctx, assets.MicroLunaDenom, sdk.NewInt(curDay-1))

		return sdk.NewDecFromInt(postSwapIssuance.Sub(prevIssuance)).QuoInt(prevIssuance)
	}

	return sdk.ZeroDec()
}

// GetSwapCoins returns the amount of asked coins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Returns an Error if the swap is recursive, or the coins to be traded are unknown by the oracle, or the amount
// to trade is too small.
// Ignores caps and spreads if isInternal = true. 
func (k Keeper) GetSwapCoins(ctx sdk.Context, offerCoin sdk.Coin, askDenom string, isInternal bool) 
(retCoin sdk.Coin, spread sdk.Dec, err sdk.Error) {
	params := k.GetParams(ctx)

	offerRate, err := k.ok.GetLunaSwapRate(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), ErrNoEffectivePrice(DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.ok.GetLunaSwapRate(ctx, askDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), ErrNoEffectivePrice(DefaultCodespace, askDenom)
	}

	retAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(askRate).Quo(offerRate).TruncateInt()
	if retAmount.Equal(sdk.ZeroInt()) {
		return sdk.Coin{}, sdk.ZeroDec(), ErrInsufficientSwapCoins(DefaultCodespace, offerCoin.Amount)
	}

	// We only charge spread for NON-INTERNAL swaps involving luna; if not, just pass. 
	if isInternal || (offerCoin.Denom != assets.MicroLunaDenom && askDenom != assets.MicroLunaDenom) {
		return sdk.NewCoin(askDenom, retAmount), sdk.ZeroDec(), nil
	}

	dailyDelta := sdk.ZeroDec()
	if offerCoin.Denom == assets.MicroLunaDenom {
		dailyDelta = k.ComputeLunaDelta(ctx, offerCoin.Amount.Neg())
	} else if askDenom == assets.MicroLunaDenom {
		dailyDelta = k.ComputeLunaDelta(ctx, retAmount)
	}

	// Do not allow swaps beyond the daily cap
	maxDelta := params.DailyLunaDeltaCap
	if dailyDelta.Abs().GT(maxDelta) {
		return sdk.Coin{}, sdk.ZeroDec(), ErrExceedsDailySwapLimit(DefaultCodespace)
	}

	// Compute a spread, which is at most MinSwapSpread and grows linearly to MaxSwapSpread with delta
	spread = dailyDelta.Quo(maxDelta).Mul(params.MaxSwapSpread.Sub(params.MinSwapSpread))
	
	return sdk.NewCoin(askDenom, retAmount), spread, nil
}

// GetSwapDecCoins returns the amount of asked DecCoins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Different from swapcoins, SwapDecCoins does not charge a spread as its use is system internal. 
// Similar to SwapCoins, but operates over sdk.DecCoins for convenience and accuracy.
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
// Swap fee pool logic

// GetSwapFeePool retrieves the claim pool from the store
func (k Keeper) GetSwapFeePool(ctx sdk.Context) (pool sdk.Coins) {
	store := ctx.KVStore(k.key)
	b := store.Get(keySwapFeePool)
	if b == nil {
		return sdk.Coins{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &pool)
	return
}

// setSwapFeePool sets the claim pool to the store
func (k Keeper) addSwapFeePool(ctx sdk.Context, fees sdk.Coins) {
	pool := k.GetSwapFeePool(ctx)
	pool = pool.Add(fees)
	
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set(keySwapFeePool, bz)
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