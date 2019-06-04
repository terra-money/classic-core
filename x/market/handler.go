package market

import (
	"reflect"

	"github.com/terra-project/core/x/market/tags"

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
	swapCoin, spread, swapErr := k.GetSwapCoin(ctx, msg.OfferCoin, msg.AskDenom, false)
	if swapErr != nil {
		return swapErr.Result()
	}

	// Charge a spread if applicable; distributed to vote winners in the oracle module
	swapFee := sdk.Coin{}
	if spread.IsPositive() {
		swapFeeAmt := spread.MulInt(swapCoin.Amount).TruncateInt()
		if swapFeeAmt.IsPositive() {
			swapFee = sdk.NewCoin(swapCoin.Denom, swapFeeAmt)
			k.ok.AddSwapFeePool(ctx, sdk.NewCoins(swapFee))

			swapCoin = swapCoin.Sub(swapFee)
		}
	}

	// Burn offered coins and subtract from the trader's account
	burnErr := k.mk.Burn(ctx, msg.Trader, msg.OfferCoin)
	if burnErr != nil {
		return burnErr.Result()
	}

	// Mint asked coins and credit Trader's account
	mintErr := k.mk.Mint(ctx, msg.Trader, swapCoin)
	if mintErr != nil {
		return mintErr.Result()
	}

	log := NewLog()
	log = log.append(LogKeySwapCoin, swapCoin.String())
	log = log.append(LogKeySwapFee, swapFee.String())

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Offer, msg.OfferCoin.Denom,
			tags.Trader, msg.Trader.String(),
		),
		Log: log.String(),
	}
}
