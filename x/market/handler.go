package market

import (
	"reflect"

	"terra/x/market/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SwapMsg:
			return handleSwapMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized market Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleSwapMsg handles the logic of a SwapMsg
func handleSwapMsg(ctx sdk.Context, k Keeper, msg SwapMsg) sdk.Result {
	retCoin, err := k.SwapCoins(ctx, msg.OfferCoin, msg.AskDenom)
	if err != nil {
		return err.Result()
	}

	// Pay gains to the treasury
	k.tk.AddIncome(ctx, msg.OfferCoin)

	reqTags, reqErr := k.tk.RequestFunds(ctx, retCoin, msg.Trader)
	if reqErr != nil {
		return reqErr.Result()
	}

	reqTags = reqTags.AppendTags(
		sdk.NewTags(
			sdk.TagAction, tags.ActionSwap,
			tags.OfferDenom, []byte(msg.OfferCoin.Denom),
			tags.OfferAmount, msg.OfferCoin.Amount,
			tags.AskDenom, []byte(retCoin.Denom),
			tags.AskAmount, retCoin.Amount,
			tags.Trader, msg.Trader.Bytes(),
		),
	)

	return sdk.Result{
		Tags: reqTags,
	}
}
