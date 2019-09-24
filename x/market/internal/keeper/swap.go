package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// ApplySwapToPool updates each pool with offerCoin and askCoin taken from swap operation,
// OfferPool = OfferPool + offerAmt (Fills the swap pool with offerAmt)
// AskPool = AskPool - askAmt       (Uses askAmt from the swap pool)
func (k Keeper) ApplySwapToPool(ctx sdk.Context, offerCoin, askCoin sdk.Coin) sdk.Error {
	// No delta update in case TERRA to TERRA swap
	if offerCoin.Denom != core.MicroLunaDenom && askCoin.Denom != core.MicroLunaDenom {
		return nil
	}

	terraPool := k.GetTerraPool(ctx)
	lunaPool := k.GetLunaPool(ctx)

	offerBaseCoin, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoinFromCoin(offerCoin), core.MicroSDRDenom)
	if err != nil {
		return err
	}

	askBaseCoin, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoinFromCoin(askCoin), core.MicroSDRDenom)
	if err != nil {
		return err
	}

	// In case swapping TERRA to LUNA, the terra swap pool(offer) is increased and the luna swap pool(ask) is decreased
	if offerCoin.Denom != core.MicroLunaDenom && askCoin.Denom == core.MicroLunaDenom {
		terraPool = terraPool.Add(offerBaseCoin.Amount)
		lunaPool = lunaPool.Sub(askBaseCoin.Amount)
	}

	// In case swapping LUNA to TERRA, the luna swap pool(offer) is increased and the terra swap pool(ask) is decreased
	if offerCoin.Denom == core.MicroLunaDenom && askCoin.Denom != core.MicroLunaDenom {
		terraPool = terraPool.Sub(askBaseCoin.Amount)
		lunaPool = lunaPool.Add(offerBaseCoin.Amount)
	}

	k.SetTerraPool(ctx, terraPool)
	k.SetLunaPool(ctx, lunaPool)

	return nil
}

// ComputeSwap returns the amount of asked coins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Returns an Error if the swap is recursive, or the coins to be traded are unknown by the oracle, or the amount
// to trade is too small.
func (k Keeper) ComputeSwap(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (retDecCoin sdk.DecCoin, spread sdk.Dec, err sdk.Error) {

	// BasePool update is delayed, so block swap
	if !k.IsMarketActive(ctx) {
		return sdk.DecCoin{}, sdk.ZeroDec(), types.ErrInactive(k.codespace)
	}

	// Return invalid recursive swap err
	if offerCoin.Denom == askDenom {
		return sdk.DecCoin{}, sdk.ZeroDec(), types.ErrRecursiveSwap(k.codespace, askDenom)
	}

	// Swap offer coin to base denom for simplicity of swap process
	baseOfferDecCoin, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoinFromCoin(offerCoin), core.MicroSDRDenom)
	if err != nil {
		return sdk.DecCoin{}, sdk.Dec{}, err
	}

	// Get swap amount based on the oracle price
	retDecCoin, err = k.ComputeInternalSwap(ctx, baseOfferDecCoin, askDenom)
	if err != nil {
		return sdk.DecCoin{}, sdk.Dec{}, err
	}

	// TERRA->TERRA swap
	// Apply only tobin tax without constant product spread
	if offerCoin.Denom != core.MicroLunaDenom && askDenom != core.MicroLunaDenom {
		spread = k.TobinTax(ctx)
		return
	}

	basePool := k.GetBasePool(ctx)
	minSpread := k.MinSpread(ctx)

	var offerPool sdk.Dec // base denom(usdr) unit
	var askPool sdk.Dec   // base denom(usdr) unit
	if offerCoin.Denom != core.MicroLunaDenom {
		// TERRA->LUNA swap
		offerPool = k.GetTerraPool(ctx)
		askPool = k.GetLunaPool(ctx)
	} else {
		// LUNA->TERRA swap
		offerPool = k.GetLunaPool(ctx)
		askPool = k.GetTerraPool(ctx)
	}

	// constant-product, which by construction is square of base(equilibrium) Terra pool
	cp := basePool.Mul(basePool)

	// Get cp(constant-product) based swap amount
	// askBaseAmount = askPool - cp / (offerPool + offerBaseAmount)
	// askBaseAmount is base denom(usdr) unit
	askBaseAmount := askPool.Sub(cp.Quo(baseOfferDecCoin.Amount.Add(offerPool)))

	// Swap base coin to ask denom
	askDecCoin, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoinFromDec(core.MicroSDRDenom, askBaseAmount), askDenom)
	if err != nil {
		return sdk.DecCoin{}, sdk.ZeroDec(), err
	}

	// spread = max(contant_product_spread + tobin_tax, tobin_tax)
	// contant_product_spread can be negative
	askDecAmount := askDecCoin.Amount
	retDecAmount := retDecCoin.Amount
	spread = retDecAmount.Sub(askDecAmount).Quo(retDecAmount).Add(minSpread)

	if spread.LT(minSpread) {
		spread = minSpread
	}

	if spread.GT(sdk.OneDec()) {
		spread = sdk.OneDec()
	}

	return
}

// ComputeInternalSwap returns the amount of asked DecCoin should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Different from ComputeSwap, ComputeInternalSwap does not charge a spread as its use is system internal.
func (k Keeper) ComputeInternalSwap(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error) {
	if offerCoin.Denom == askDenom {
		return offerCoin, nil
	}

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