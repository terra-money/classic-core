package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-money/core/x/market/internal/types"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		case MsgSwapSend:
			return handleMsgSwapSend(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

func handleMsgSwapSend(ctx sdk.Context, k Keeper, mss MsgSwapSend) (*sdk.Result, error) {
	return handleSwapRequest(ctx, k, mss.FromAddress, mss.ToAddress, mss.OfferCoin, mss.AskDenom)
}

func handleMsgSwap(ctx sdk.Context, k Keeper, ms MsgSwap) (*sdk.Result, error) {
	return handleSwapRequest(ctx, k, ms.Trader, ms.Trader, ms.OfferCoin, ms.AskDenom)
}

// handleMsgSwap handles the logic of a MsgSwap
func handleSwapRequest(ctx sdk.Context, k Keeper,
	trader sdk.AccAddress, receiver sdk.AccAddress,
	offerCoin sdk.Coin, askDenom string) (*sdk.Result, error) {
	// Can't swap to the same coin
	if offerCoin.Denom == askDenom {
		return nil, ErrRecursiveSwap
	}

	// Compute exchange rates between the ask and offer
	swapCoin, spread, err := k.ComputeSwap(ctx, offerCoin, askDenom)
	if err != nil {
		return nil, err
	}

	// Charge a spread if applicable; the spread is burned
	var swapFee sdk.DecCoin
	if spread.IsPositive() {
		swapFeeAmt := spread.Mul(swapCoin.Amount)
		if swapFeeAmt.IsPositive() {
			swapFee = sdk.NewDecCoinFromDec(swapCoin.Denom, swapFeeAmt)
			swapCoin = swapCoin.Sub(swapFee)
		}
	} else {
		swapFee = sdk.NewDecCoin(swapCoin.Denom, sdk.ZeroInt())
	}

	// Update pool delta
	err = k.ApplySwapToPool(ctx, offerCoin, swapCoin)
	if err != nil {
		return nil, err
	}

	// Send offer coins to module account
	offerCoins := sdk.NewCoins(offerCoin)
	err = k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, trader, ModuleName, offerCoins)
	if err != nil {
		return nil, err
	}

	// Burn offered coins and subtract from the trader's account
	err = k.SupplyKeeper.BurnCoins(ctx, ModuleName, offerCoins)
	if err != nil {
		return nil, err
	}

	// Mint asked coins and credit Trader's account
	retCoin, decimalCoin := swapCoin.TruncateDecimal()
	swapFee = swapFee.Add(decimalCoin) // add truncated decimalCoin to swapFee
	swapCoins := sdk.NewCoins(retCoin)
	err = k.SupplyKeeper.MintCoins(ctx, ModuleName, swapCoins)
	if err != nil {
		return nil, err
	}

	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, receiver, swapCoins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventSwap,
			sdk.NewAttribute(types.AttributeKeyOffer, offerCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTrader, trader.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, receiver.String()),
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
