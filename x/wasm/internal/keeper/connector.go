package keeper

import (
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) dispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs []wasmTypes.CosmosMsg) error {
	for _, msg := range msgs {
		if err := k.dispatchMessage(ctx, contractAddr, msg); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) dispatchMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) error {
	msgs, err := k.msgParser.Parse(contractAddr, msg)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if err := k.handleSdkMessage(ctx, contractAddr, msg); err != nil {
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
