package market

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	minReserveRatio = 1.2
	maxReserveRatio = 1.5
	feeUpdatePeriod = 1209600
)

var feeUpdateTimestamp = 0

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

// // NewEndBlocker checks proposals and generates a EndBlocker
// func NewEndBlocker(k Keeper) sdk.EndBlocker {
// 	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
// 		newTags := sdk.NewTags()

// 		return abci.ResponseEndBlock{
// 			Tags: tags,
// 		}
// 	}
// }

func lunaExchangeRate(ctx sdk.Context, k Keeper, denom string) sdk.Dec {
	if denom == "luna" {
		return sdk.OneDec()
	}

	return k.ok.GetElect(ctx, denom).FeedMsg.CurrentPrice
}

// handleVoteMsg handles the logic of a SwapMsg
func handleSwapMsg(ctx sdk.Context, k Keeper, msg SwapMsg) sdk.Result {
	tags := sdk.NewTags()

	// If swap msg for not whitelisted denom
	if !k.ok.WhitelistContains(ctx, msg.OfferCoin.Denom) {
		return ErrUnknownDenomination(DefaultCodespace, msg.OfferCoin.Denom).Result()
	}

	offerRate := lunaExchangeRate(ctx, k, msg.OfferCoin.Denom)
	askRate := lunaExchangeRate(ctx, k, msg.AskDenom)

	retAmount := sdk.NewDecFromInt(msg.OfferCoin.Amount).Mul(offerRate).Quo(askRate).RoundInt()

	if retAmount.Equal(sdk.ZeroInt()) {
		// drop in this scenario
		return ErrInsufficientSwapCoins(DefaultCodespace, msg.OfferCoin.Amount).Result()
	}

	retCoin := sdk.Coin{
		Denom:  msg.AskDenom,
		Amount: retAmount,
	}

	// Reflect the swap in the trader's wallet
	swapTags, swapErr := k.bk.InputOutputCoins(ctx, []bank.Input{bank.NewInput(msg.Trader, sdk.Coins{retCoin})},
		[]bank.Output{bank.NewOutput(msg.Trader, sdk.Coins{msg.OfferCoin})})

	if swapErr != nil {
		return swapErr.Result()
	}

	tags.AppendTags(swapTags)

	// Update coin supplies
	offerCoinSupply := k.GetCoinSupply(ctx, msg.OfferCoin.Denom)
	askCoinSupply := k.GetCoinSupply(ctx, msg.AskDenom)

	k.SetCoinSupply(ctx, msg.OfferCoin.Denom, offerCoinSupply.Sub(msg.OfferCoin.Amount))
	k.SetCoinSupply(ctx, msg.AskDenom, askCoinSupply.Add(retAmount))

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
