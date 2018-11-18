package market

import (
	"encoding/binary"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
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

// NewEndBlocker checks proposals and generates a EndBlocker
func NewEndBlocker(k Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		newTags := sdk.NewTags()

		rp := k.GetReserveParams(ctx)
		curIssuance := k.GetIssuanceMeta()
		curFee := k.GetTerraFee()

		terraSupply := curIssuance.AmountOf("terra")
		lunaSupply := curIssuance.AmountOf("luna")

		terraPrice := k.ok.getTerraPrice()
		lunaPrice := k.ok.getLunaPrice()

		// Time to update fees
		if feeUpdateTimestamp+feeUpdatePeriod < ctx.BlockHeight {
			
			// Overcapitalized
			if terraSupply*terraPrice*maxReserveRatio < (rp.Cap - lunaSupply)*lunaPrice {
				newFee := curFee.Sub(sdk.Rat.SetFloat64(0.0002)) // Iteratively decrease 0.2% in tx fees
				newTags.AppendTag("action", []byte("slashFee"))
				newTags.AppendTag("newFee", newFee.Bytes())
				k.SetTerraFee(ctx, newFee)
			} 

			// Undercapitalized
			else if terraSupply*terraPrice*minReserveRatio > (rp.Cap - lunaSupply)*lunaPrice {
				newFee := curFee.Mul(sdk.Rat.SetInt(2))	// Double transaction fees ... TODO: change to something better
				newTags.AppendTag("action", []byte("hikeFee"))
				newTags.AppendTag("newFee", newFee.Bytes())
				k.SetTerraFee(ctx, newFee)
			}

			// No more collateral ... recap the reserve
			if lunaSupply > rp.Cap {
				rp.Target = rp.Cap
				rp.Cap = rp.Cap.Mul(2)
				k.SetReserveParams(rp)

				newTags.AppendTags(
					"action", []byte("dilute"),
					"cap", reserve.Cap.Bytes(),
				)
			} 
			
			// Reserve is overcollateralized ... need to refund the reserve
			else if lunaSupply < rp.Target {
				// TODO: do seigniorage
				newTags.AppendTags(
					"action", []byte("seigniorage"),
					"cap", (reserve.Target - reserve.Current).Bytes(),
				)
			}

			feeUpdateBlockNum = ctx.BlockHeight
		}

		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// handleVoteMsg handles the logic of a SwapMsg
func handleSwapMsg(ctx sdk.Context, k Keeper, msg SwapMsg) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	var rval sdk.Coin
	switch msg.Coin.Denom {
	case "terra":
		rval = sdk.Coin{Denom: "luna", Amount: msg.Coin.Amount / k.ok.getLunaPrice()}
	case "luna":
		rval = sdk.Coin{Denom: "terra", Amount: msg.Coin.Amount * k.ok.getLunaPrice()}
	default:
		errMsg := "Unrecognized swap Msg type: " + reflect.TypeOf(msg).Name()
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	// Reflect the swap in the trader's wallet
	tags, swaperr := k.bk.InputOutputCoins(ctx, []Input{bank.NewInput(msg.Trader, sdk.Coins{rval})},
		[]Output{bank.NewOutput(msg.Trader, sdk.Coins{rval})})

	// Update the issuance meta with the swap
	curIssuance := k.GetIssuanceMeta()
	curIssuance.Minus(msg.Coin)
	curIssuance.Plus(rval)
	k.SetIssuanceMeta(curIssuance)

	if swaperr != nil {
		return swaperr.Result()
	}

	tags.AppendTags(
		"action", []byte("swap"),
		"subject", []byte(msg.Coin.String),
		"trader", msg.Trader.Bytes(),
	)

	return sdk.Result{
		tags,
	}
}