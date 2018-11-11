package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		case MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked

	tk, ok := k.(TaxKeeper)
	if ok {
		// Pay stability taxes on inputs and outputs
		for i, input := range msg.Inputs {
			stabilityTax := tk.PayTax(ctx, input.Coins)
			msg.Inputs[i].Coins = input.Coins.Minus(stabilityTax)
		}

		for j, output := range msg.Outputs {
			stabilityTax := tk.GetTax(ctx, output.Coins) // Tax already paid in inputs
			msg.Outputs[j].Coins = output.Coins.Minus(stabilityTax)
		}
	}

	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgIssue.
func handleMsgIssue(ctx sdk.Context, k Keeper, msg MsgIssue) sdk.Result {
	panic("not implemented yet")
}
