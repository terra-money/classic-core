package keeper

import (
	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) dispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs []wasmTypes.CosmosMsg) sdk.Error {
	for _, msg := range msgs {
		if err := k.dispatchMessage(ctx, contractAddr, msg); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) dispatchMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) sdk.Error {
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

func (k Keeper) handleSdkMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg sdk.Msg) sdk.Error {
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return sdk.ErrUnauthorized("contract doesn't have permission")
		}
	}

	// find the handler and execute it
	h := k.router.Route(msg.Route())
	if h == nil {
		return sdk.ErrUnknownRequest(msg.Route())
	}

	res := h(ctx, msg)
	if !res.IsOK() {
		return sdk.NewError(res.Codespace, res.Code, res.Log)
	}

	// redispatch all events, (type sdk.EventTypeMessage will be filtered out in the handler)
	ctx.EventManager().EmitEvents(res.Events)

	return nil
}
