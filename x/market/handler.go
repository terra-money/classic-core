package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-project/core/x/market/internal/types"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// handleMsgSwap handles the logic of a MsgSwap
func handleMsgSwap(ctx sdk.Context, k Keeper, ms MsgSwap) (*sdk.Result, error) {

	// Can't swap to the same coin
	if ms.OfferCoin.Denom == ms.AskDenom {
		return nil, sdkerrors.Wrap(ErrRecursiveSwap, ms.AskDenom)
	}

	// Compute exchange rates between the ask and offer
	swapCoin, spread, swapErr := k.ComputeSwap(ctx, ms.OfferCoin, ms.AskDenom)
	if swapErr != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "ComputeSwap failed: %s", swapErr.Error())
	}

	// Update pool delta
	deltaUpdateErr := k.ApplySwapToPool(ctx, ms.OfferCoin, swapCoin)
	if deltaUpdateErr != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "ApplySwapToPool failed: %s", deltaUpdateErr.Error())
	}

	// Send offer coins to module account
	offerCoins := sdk.NewCoins(ms.OfferCoin)
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ms.Trader, ModuleName, offerCoins)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "SendCoinsFromAccountToModule failed: %s", err.Error())
	}

	// Charge a spread if applicable; distributed to vote winners in the oracle module
	var swapFee sdk.DecCoin
	if spread.IsPositive() {
		swapFeeAmt := spread.Mul(swapCoin.Amount)
		if swapFeeAmt.IsPositive() {
			swapFee = sdk.NewDecCoinFromDec(swapCoin.Denom, swapFeeAmt)
			swapCoin = swapCoin.Sub(swapFee)
		}
	}

	// Burn offered coins and subtract from the trader's account
	burnErr := k.SupplyKeeper.BurnCoins(ctx, ModuleName, offerCoins)
	if burnErr != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "BurnCoins failed: %s", burnErr.Error())
	}

	// Mint asked coins and credit Trader's account
	retCoin, decimalCoin := swapCoin.TruncateDecimal()
	swapFee = swapFee.Add(decimalCoin) // add truncated decimalCoin to swapFee
	swapCoins := sdk.NewCoins(retCoin)
	mintErr := k.SupplyKeeper.MintCoins(ctx, ModuleName, swapCoins)
	if mintErr != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "MintCoins failed: %s", mintErr.Error())
	}

	sendErr := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, ms.Trader, swapCoins)
	if sendErr != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPanic, "SendCoinsFromModuleToAccount failed: %s", sendErr.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventSwap,
			sdk.NewAttribute(types.AttributeKeyOffer, ms.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTrader, ms.Trader.String()),
			sdk.NewAttribute(types.AttributeKeySwapCoin, retCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFee, swapFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
