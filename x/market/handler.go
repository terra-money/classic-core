package market

import (
	"reflect"
	"terra/types/assets"

	"terra/x/market/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			errMsg := "Unrecognized market Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSwap handles the logic of a MsgSwap
func handleMsgSwap(ctx sdk.Context, k Keeper, msg MsgSwap) sdk.Result {

	// Can't swap to the same coin
	if msg.OfferCoin.Denom == msg.AskDenom {
		return ErrRecursiveSwap(DefaultCodespace, msg.AskDenom).Result()
	}

	// Compute exchange rates between the ask and offer
	swapCoin, swapErr := k.SwapCoins(ctx, msg.OfferCoin, msg.AskDenom)
	if swapErr != nil {
		return swapErr.Result()
	}

	// Burn offered coins and subtract from the trader's account
	burnErr := k.mk.Burn(ctx, msg.Trader, msg.OfferCoin)
	if burnErr != nil {
		return burnErr.Result()
	}

	// Record seigniorage if the offered coin is Luna
	if msg.OfferCoin.Denom == assets.MicroLunaDenom {
		k.mk.AddSeigniorage(ctx, msg.OfferCoin.Amount)
	}

	// Mint asked coins and credit Trader's account
	mintErr := k.mk.Mint(ctx, msg.Trader, swapCoin)
	if mintErr != nil {
		return mintErr.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Offer, msg.OfferCoin.String(),
			tags.Ask, swapCoin.String(),
			tags.Trader, msg.Trader.String(),
		),
	}
}
