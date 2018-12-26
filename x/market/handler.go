package market

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SwapMsg:
			return handleSwapMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized swap Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleSwapMsg handles the logic of a SwapMsg
func handleSwapMsg(ctx sdk.Context, k Keeper, msg SwapMsg) sdk.Result {
	tags := sdk.NewTags()

	retCoin, err := k.SwapCoins(ctx, msg.OfferCoin, msg.AskDenom)
	if err != nil {
		return err.Result()
	}

	// Reflect the swap in the trader's wallet
	swapTags, swapErr := k.bk.InputOutputCoins(ctx, []bank.Input{bank.NewInput(msg.Trader, sdk.Coins{retCoin})},
		[]bank.Output{bank.NewOutput(msg.Trader, sdk.Coins{msg.OfferCoin})})

	if swapErr != nil {
		return swapErr.Result()
	}

	tags.AppendTags(swapTags)

	// Pay gains to the treasury
	k.tk.PayMintIncome(ctx, sdk.Coins{msg.OfferCoin})

	tags.AppendTags(
		sdk.NewTags(
			"action", []byte("swap"),
			"offer", []byte(msg.OfferCoin.String()),
			"ask", []byte(retCoin.String()),
			"trader", msg.Trader.Bytes(),
		),
	)

	return sdk.Result{
		Tags: tags,
	}
}
