// Package bank contains a forked version of the bank module. It only contains
// a modified message handler to support a very limited form of transfers during
// mainnet launch -- MsgMultiSend messages.
//
// NOTE: This fork should be removed entirely once transfers are enabled and
// the Gaia router should be reset to using the original bank module handler.
package pay

import (
	"terra/types/util"
	"terra/x/treasury"

	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// var (
// 	uatomDenom    = "uatom"
// 	atomsToUatoms = int64(1000000)

// 	// BurnedCoinsAccAddr represents the burn account address used for
// 	// MsgMultiSend message during the period for which transfers are disabled.
// 	// Its Bech32 address is cosmos1x4p90uuy63fqzsheamn48vq88q3eusykf0a69v.
// 	BurnedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("bankBurnedCoins")))
// )

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case bank.MsgSend:
			return handleMsgSend(ctx, k, tk, fk, msg)

		case bank.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, tk, fk, msg)

		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper, msg bank.MsgSend) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Codespace()).Result()
	}

	tags, err := payTax(ctx, k, tk, fk, msg.FromAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	sendTags, err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	tags = tags.AppendTags(sendTags)
	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper, msg bank.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Codespace()).Result()
	}

	tags := sdk.NewTags()
	for _, output := range msg.Outputs {
		taxTags, taxErr := payTax(ctx, k, tk, fk, output.Address, output.Coins)
		if taxErr != nil {
			return taxErr.Result()
		}

		tags = tags.AppendTags(taxTags)
	}

	sendTags, sendErr := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if sendErr != nil {
		return sendErr.Result()
	}

	tags = tags.AppendTags(sendTags)
	return sdk.Result{
		Tags: tags,
	}
}

// payTax charges the stability tax on SendCoin and InputOutputCoins.
func payTax(ctx sdk.Context, bk bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper,
	taxPayer sdk.AccAddress, principal sdk.Coins) (taxTags sdk.Tags, err sdk.Error) {

	taxes := sdk.Coins{}
	for _, coin := range principal {
		taxRate := tk.GetTaxRate(ctx)
		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).RoundInt()
		taxCap := tk.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		taxCoin := sdk.Coins{sdk.NewCoin(coin.Denom, taxDue)}

		_, payTags, err := bk.SubtractCoins(ctx, taxPayer, taxCoin)
		if err != nil {
			return nil, err
		}

		taxTags = taxTags.AppendTags(payTags)
		taxes = append(taxes, sdk.NewCoin(coin.Denom, taxDue))
		fk.AddCollectedFees(ctx, taxCoin)
	}

	tk.AddTaxProceeds(ctx, util.GetEpoch(ctx), taxes)
	return
}
