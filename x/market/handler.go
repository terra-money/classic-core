package market

import (
	"reflect"

	"terra/x/market/tags"

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
			errMsg := "Unrecognized market Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func recordSeigniorage(ctx sdk.Context, k Keeper, seigniorage sdk.Coin) {
	feePool := k.dk.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Plus(sdk.DecCoins{sdk.NewDecCoinFromCoin(seigniorage)})
	k.dk.SetFeePool(ctx, feePool)
}

// handleSwapMsg handles the logic of a SwapMsg
func handleSwapMsg(ctx sdk.Context, k Keeper, msg SwapMsg) sdk.Result {
	swapCoin, swapErr := k.SwapCoins(ctx, msg.OfferCoin, msg.AskDenom)
	if swapErr != nil {
		return swapErr.Result()
	}

	input := bank.Input{Address: msg.Trader, Coins: sdk.Coins{swapCoin}}
	output := bank.Output{Address: msg.Trader, Coins: sdk.Coins{msg.OfferCoin}}

	// Record seigniorage
	recordSeigniorage(ctx, k, swapCoin)

	reqTags, reqErr := k.pk.InputOutputCoins(ctx, []bank.Input{input}, []bank.Output{output})
	if reqErr != nil {
		return reqErr.Result()
	}

	reqTags = reqTags.AppendTags(
		sdk.NewTags(
			sdk.TagAction, tags.ActionSwap,
			tags.OfferDenom, []byte(msg.OfferCoin.Denom),
			tags.OfferAmount, msg.OfferCoin.Amount,
			tags.AskDenom, []byte(swapCoin.Denom),
			tags.AskAmount, swapCoin.Amount,
			tags.Trader, msg.Trader.Bytes(),
		),
	)

	return sdk.Result{
		Tags: reqTags,
	}
}
