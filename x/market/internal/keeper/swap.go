package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// ApplySwapToPool updates each pool with offerCoin and askCoin taken from swap operation,
// OfferPool = OfferPool + offerAmt (Fills the swap pool with offerAmt)
// AskPool = AskPool - askAmt       (Uses askAmt from the swap pool)
func (k Keeper) ApplySwapToPool(ctx sdk.Context, offerCoin sdk.Coin, askCoin sdk.DecCoin) sdk.Error {
	// No delta update in case Terra to Terra swap
	if offerCoin.Denom != core.MicroLunaDenom && askCoin.Denom != core.MicroLunaDenom {
		return nil
	}

	terraPoolDelta := k.GetTerraPoolDelta(ctx)

	// In case swapping Terra to Luna, the terra swap pool(offer) must be increased and the luna swap pool(ask) must be decreased
	if offerCoin.Denom != core.MicroLunaDenom && askCoin.Denom == core.MicroLunaDenom {
		offerBaseCoin, err := k.ComputeInternalSwap(ctx, sdk.NewDecCoinFromCoin(offerCoin), core.MicroSDRDenom)
		if err != nil {
			return err
		}

		terraPoolDelta = terraPoolDelta.Add(offerBaseCoin.Amount)
	}

	// In case swapping Luna to Terra, the luna swap pool(offer) must be increased and the terra swap pool(ask) must be decreased
	if offerCoin.Denom == core.MicroLunaDenom && askCoin.Denom != core.MicroLunaDenom {
		askBaseCoin, err := k.ComputeInternalSwap(ctx, askCoin, core.MicroSDRDenom)
		if err != nil {
			return err
		}

		terraPoolDelta = terraPoolDelta.Sub(askBaseCoin.Amount)
	}

	k.SetTerraPoolDelta(ctx, terraPoolDelta)

	return nil
}

// ComputeSwap returns the amount of asked coins should be returned for a given offerCoin at the effective
// exchange rate registered with the oracle.
// Returns an Error if the swap is recursive, or the coins to be traded are unknown by the oracle, or the amount
// to trade is too small.
func (k Keeper) ComputeSwap(ctx sdk.Context, offerCoin sdk.Coin, askDenom string) (retDecCoin sdk.DecCoin, spread sdk.Dec, err sdk.Error) {

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

	// Terra->Terra swap
	// Apply only tobin tax without constant product spread
	if offerCoin.Denom != core.MicroLunaDenom && askDenom != core.MicroLunaDenom {
		var tobinTax sdk.Dec
		offerTobinTax, err2 := k.oracleKeeper.GetTobinTax(ctx, offerCoin.Denom)
		if err2 != nil {
			return sdk.DecCoin{}, sdk.Dec{}, err2
		}

		askTobinTax, err2 := k.oracleKeeper.GetTobinTax(ctx, askDenom)
		if err2 != nil {
			return sdk.DecCoin{}, sdk.Dec{}, err2
		}

		// Apply highest tobin tax for the denoms in the swap operation
		if askTobinTax.GT(offerTobinTax) {
			tobinTax = askTobinTax
		} else {
			tobinTax = offerTobinTax
		}

		spread = tobinTax
		return
	}

	basePool := k.BasePool(ctx)
	minSpread := k.MinStabilitySpread(ctx)

	// constant-product, which by construction is square of base(equilibrium) pool
	cp := basePool.Mul(basePool)
	terraPoolDelta := k.GetTerraPoolDelta(ctx)
	terraPool := basePool.Add(terraPoolDelta)
	lunaPool := cp.Quo(terraPool)

	var offerPool sdk.Dec // base denom(usdr) unit
	var askPool sdk.Dec   // base denom(usdr) unit
	if offerCoin.Denom != core.MicroLunaDenom {
		// Terra->Luna swap
		offerPool = terraPool
		askPool = lunaPool
	} else {
		// Luna->Terra swap
		offerPool = lunaPool
		askPool = terraPool
	}

	// Get cp(constant-product) based swap amount
	// askBaseAmount = askPool - cp / (offerPool + offerBaseAmount)
	// askBaseAmount is base denom(usdr) unit
	askBaseAmount := askPool.Sub(cp.Quo(offerPool.Add(baseOfferDecCoin.Amount)))

	// Both baseOffer and baseAsk are usdr units, so spread can be calculated by
	// spread = (baseOfferAmt - baseAskAmt) / baseOfferAmt
	baseOfferAmount := baseOfferDecCoin.Amount
	spread = baseOfferAmount.Sub(askBaseAmount).Quo(baseOfferAmount)

	if spread.LT(minSpread) {
		spread = minSpread
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

	offerRate, err := k.oracleKeeper.GetLunaExchangeRate(ctx, offerCoin.Denom)
	if err != nil {
		return sdk.DecCoin{}, types.ErrNoEffectivePrice(types.DefaultCodespace, offerCoin.Denom)
	}

	askRate, err := k.oracleKeeper.GetLunaExchangeRate(ctx, askDenom)
	if err != nil {
		return sdk.DecCoin{}, types.ErrNoEffectivePrice(types.DefaultCodespace, askDenom)
	}

	retAmount := offerCoin.Amount.Mul(askRate).Quo(offerRate)
	if retAmount.LTE(sdk.ZeroDec()) {
		return sdk.DecCoin{}, types.ErrInvalidOfferCoin(types.DefaultCodespace, offerCoin.Amount.TruncateInt())
	}

	return sdk.NewDecCoinFromDec(askDenom, retAmount), nil
}
