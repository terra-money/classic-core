package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	oracleKeeper types.OracleKeeper
	SupplyKeeper types.SupplyKeeper

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, oracleKeeper types.OracleKeeper,
	supplyKeeper types.SupplyKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   paramspace.WithKeyTable(ParamKeyTable()),
		oracleKeeper: oracleKeeper,
		SupplyKeeper: supplyKeeper,
		codespace:    codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns a codespace of keeper
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetLastDayIssuance returns the last day issuance
func (k Keeper) GetLastDayIssuance(ctx sdk.Context) (issuance sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastDayIssuanceKey)
	if bz == nil {
		return sdk.Coins{}
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &issuance)
	return
}

// UpdateLastDayIssuance stores the last day issuance
func (k Keeper) UpdateLastDayIssuance(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(k.storeKey)

	totalCoins := k.SupplyKeeper.GetSupply(ctx).GetTotal()
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalCoins)
	store.Set(types.LastDayIssuanceKey, bz)

	return totalCoins
}

// ComputeLunaDelta returns the issuance change rate of Luna for the day post-swap
func (k Keeper) ComputeLunaDelta(ctx sdk.Context, change sdk.Int) sdk.Dec {
	curDay := ctx.BlockHeight() / core.BlocksPerDay
	if curDay == 0 {
		return sdk.ZeroDec()
	}

	lastDayLunaIssuance := k.GetLastDayIssuance(ctx).AmountOf(core.MicroLunaDenom)
	if lastDayLunaIssuance.IsZero() {
		return sdk.ZeroDec()
	}

	supply := k.SupplyKeeper.GetSupply(ctx)
	lunaIssuance := supply.GetTotal().AmountOf(core.MicroLunaDenom)

	postSwapIssunace := lunaIssuance.Add(change)

	return sdk.NewDecFromInt(postSwapIssunace.Sub(lastDayLunaIssuance)).QuoInt(lastDayLunaIssuance)
}

// ComputeLunaSwapSpread returns a spread, which is initialiy MinSwapSpread and grows linearly to MaxSwapSpread with delta
func (k Keeper) ComputeLunaSwapSpread(ctx sdk.Context, postLunaDelta sdk.Dec) sdk.Dec {
	if postLunaDelta.GTE(k.DailyLunaDeltaCap(ctx)) {
		return k.MaxSwapSpread(ctx)
	}

	// min + (p / l) (max - min); l = dailyDeltaCap, p = postDailyDelta,
	return k.MinSwapSpread(ctx).Add(postLunaDelta.Quo(k.DailyLunaDeltaCap(ctx)).Mul(k.MaxSwapSpread(ctx).Sub(k.MinSwapSpread(ctx))))
}

// GetSwapCoin returns the amount of asked coins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Returns an Error if the swap is recursive, or the coins to be traded are unknown by the oracle, or the amount
// to trade is too small.
// Ignores caps and spreads if isInternal = true.
func (k Keeper) GetSwapCoin(ctx sdk.Context, offerCoin sdk.Coin, askDenom string, isInternal bool) (retCoin sdk.Coin, spread sdk.Dec, err sdk.Error) {
	offerRate, err := k.oracleKeeper.GetLunaPrice(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), types.ErrNoEffectivePrice(types.DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.oracleKeeper.GetLunaPrice(ctx, askDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), types.ErrNoEffectivePrice(types.DefaultCodespace, askDenom)
	}

	retAmount := sdk.NewDecFromInt(offerCoin.Amount).Mul(askRate).Quo(offerRate).TruncateInt()
	if retAmount.Equal(sdk.ZeroInt()) {
		return sdk.Coin{}, sdk.ZeroDec(), types.ErrInsufficientSwapCoins(types.DefaultCodespace, offerCoin.Amount)
	}

	// We only charge spread for NON-INTERNAL swaps involving luna; if not, just pass.
	if isInternal || (offerCoin.Denom != core.MicroLunaDenom && askDenom != core.MicroLunaDenom) {
		return sdk.NewCoin(askDenom, retAmount), sdk.ZeroDec(), nil
	}

	dailyDelta := sdk.ZeroDec()
	if offerCoin.Denom == core.MicroLunaDenom {
		dailyDelta = k.ComputeLunaDelta(ctx, offerCoin.Amount.Neg())
	} else if askDenom == core.MicroLunaDenom {
		dailyDelta = k.ComputeLunaDelta(ctx, retAmount)
	}

	// delta should be positive to apply spread
	dailyDelta = dailyDelta.Abs()
	spread = k.ComputeLunaSwapSpread(ctx, dailyDelta)

	return sdk.NewCoin(askDenom, retAmount), spread, nil
}

// GetSwapDecCoin returns the amount of asked DecCoins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Different from swapcoins, SwapDecCoins does not charge a spread as its use is system internal.
// Similar to SwapCoins, but operates over sdk.DecCoins for convenience and accuracy.
func (k Keeper) GetSwapDecCoin(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error) {
	offerRate, err := k.oracleKeeper.GetLunaPrice(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.DecCoin{}, types.ErrNoEffectivePrice(types.DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.oracleKeeper.GetLunaPrice(ctx, askDenom)
	if err != nil {
		return sdk.DecCoin{}, types.ErrNoEffectivePrice(types.DefaultCodespace, askDenom)
	}

	retAmount := offerCoin.Amount.Mul(askRate).Quo(offerRate)
	if retAmount.LTE(sdk.ZeroDec()) {
		return sdk.DecCoin{}, types.ErrInsufficientSwapCoins(types.DefaultCodespace, offerCoin.Amount.TruncateInt())
	}

	return sdk.NewDecCoinFromDec(askDenom, retAmount), nil
}
