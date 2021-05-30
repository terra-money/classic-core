package keeper

import (
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/terra-money/core/x/auth/ante"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

func (k Keeper) dispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs []wasmTypes.CosmosMsg) error {
	var sdkMsgs []sdk.Msg
	for _, msg := range msgs {

		msgs, err := k.msgParser.Parse(contractAddr, msg)
		if err != nil {
			return err
		}

		sdkMsgs = append(sdkMsgs, msgs...)
	}

	// Charge tax on result msg
	taxes := ante.FilterMsgAndComputeTax(ctx, k.treasuryKeeper, sdkMsgs)
	if !taxes.IsZero() {
		contractAcc := k.accountKeeper.GetAccount(ctx, contractAddr)
		if err := cosmosante.DeductFees(k.supplyKeeper, ctx, contractAcc, taxes); err != nil {
			return err
		}
	}

	for _, sdkMsg := range sdkMsgs {
		if err := k.handleSdkMessage(ctx, contractAddr, sdkMsg); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) handleSdkMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg sdk.Msg) error {
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract doesn't have permission")
		}
	}

	// find the handler and execute it
	h := k.router.Route(ctx, msg.Route())
	if h == nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg.Route())
	}

	res, err := h(ctx, msg)
	if err != nil {
		return err
	}

	// redispatch all events, (type sdk.EventTypeMessage will be filtered out in the handler)
	ctx.EventManager().EmitEvents(res.Events)

	return nil
}
