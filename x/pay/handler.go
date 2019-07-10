// Package pay contains a forked version of the bank module. It only contains
// a modified message handler to support the payement of stability taxes.
//
// Taxes are of the fomula: min(principal * taxRate, taxCap).
// TaxCap and taxRate are stored by the treasury module.
// Should transactions fail midway, taxes are still paid and non-refundable.
package pay

import (
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/treasury"

	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

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

// Handle MsgPay.
func handleMsgSend(ctx sdk.Context, k bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper, msg bank.MsgSend) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Codespace()).Result()
	}

	taxes, err := payTax(ctx, k, tk, fk, msg.FromAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	log := NewLog()
	log = log.append(LogKeyTax, taxes.String())

	resultTags := sdk.NewTags()
	sendTags, err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	resultTags = resultTags.AppendTags(sendTags)

	return sdk.Result{
		Tags: resultTags,
		Log:  log.String(),
	}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper, msg bank.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(k.Codespace()).Result()
	}

	totalTaxes := sdk.Coins{}
	for _, input := range msg.Inputs {
		taxes, taxErr := payTax(ctx, k, tk, fk, input.Address, input.Coins)
		if taxErr != nil {
			return taxErr.Result()
		}

		totalTaxes = totalTaxes.Add(taxes).Sort()
	}

	log := NewLog()
	log = log.append(LogKeyTax, totalTaxes.String())

	resultTags := sdk.NewTags()
	sendTags, sendErr := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if sendErr != nil {
		return sendErr.Result()
	}

	resultTags = resultTags.AppendTags(sendTags)
	return sdk.Result{
		Tags: resultTags,
		Log:  log.String(),
	}
}

// payTax charges the stability tax on MsgSend and MsgMultiSend.
func payTax(ctx sdk.Context, bk bank.Keeper, tk treasury.Keeper, fk auth.FeeCollectionKeeper,
	taxPayer sdk.AccAddress, principal sdk.Coins) (taxes sdk.Coins, err sdk.Error) {

	taxRate := tk.GetTaxRate(ctx, util.GetEpoch(ctx))

	if taxRate.Equal(sdk.ZeroDec()) {
		return nil, nil
	}

	for _, coin := range principal {
		// no tax fee for uluna
		if coin.Denom == assets.MicroLunaDenom {
			continue
		}

		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()

		// If tax due is greater than the tax cap, cap!
		taxCap := tk.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		if taxDue.Equal(sdk.ZeroInt()) {
			continue
		}

		taxes = taxes.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxDue))).Sort()
	}

	if taxes.Empty() {
		return taxes, nil
	}

	_, _, err = bk.SubtractCoins(ctx, taxPayer, taxes)
	if err != nil {
		return nil, err
	}

	fk.AddCollectedFees(ctx, taxes)
	tk.RecordTaxProceeds(ctx, taxes)
	return
}
