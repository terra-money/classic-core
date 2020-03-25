package keeper

import (
	"fmt"
	"github.com/terra-project/core/x/wasm/internal/types"

	wasmTypes "github.com/confio/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

func (k Keeper) dispatchMessages(ctx sdk.Context, contract exported.Account, msgs []wasmTypes.CosmosMsg) sdk.Error {
	for _, msg := range msgs {
		if err := k.dispatchMessage(ctx, contract, msg); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) dispatchMessage(ctx sdk.Context, contract exported.Account, msg wasmTypes.CosmosMsg) sdk.Error {
	if (msg.Send != nil && msg.Contract != nil) ||
		(msg.Send != nil && msg.Opaque != nil) ||
		(msg.Contract != nil && msg.Opaque != nil) {
		return sdk.ErrInternal("single msg cannot contain multiple msgs")
	}

	// Handle MsgSend
	if msg.Send != nil {
		sendMsg, err := types.ParseMsgSend(msg.Send)
		if err != nil {
			return err
		}

		return k.handleSdkMessage(ctx, contract, sendMsg)
	}

	// Handle MsgExecuteContract
	if msg.Contract != nil {
		targetAddr, stderr := sdk.AccAddressFromBech32(msg.Contract.ContractAddr)
		if stderr != nil {
			return sdk.ErrInvalidAddress(msg.Contract.ContractAddr)
		}

		coins, err := types.ParseToCoins(msg.Contract.Send)
		if err != nil {
			return err
		}

		_, err = k.ExecuteContract(ctx, targetAddr, contract.GetAddress(), coins, []byte(msg.Contract.Msg))
		return err
	}

	if msg.Opaque != nil {
		sdkMsg, err := types.ParseOpaqueMsg(k.cdc, msg.Opaque)
		if err != nil {
			return err
		}

		return k.handleSdkMessage(ctx, contract, sdkMsg)
	}

	return sdk.ErrInternal(fmt.Sprintf("Unknown Msg: %#v", msg))
}

func (k Keeper) handleSdkMessage(ctx sdk.Context, contract exported.Account, msg sdk.Msg) sdk.Error {
	// make sure this account can send it
	contractAddr := contract.GetAddress()
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
	return nil
}
