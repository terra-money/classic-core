package wasm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgStoreCode:
			return handleStoreCode(ctx, k, &msg)
		case *MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case MsgInstantiateContract:
			return handleInstantiate(ctx, k, &msg)
		case *MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case MsgExecuteContract:
			return handleExecute(ctx, k, &msg)
		case *MsgExecuteContract:
			return handleExecute(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStoreCode(ctx sdk.Context, k Keeper, msg *MsgStoreCode) sdk.Result {
	codeID, err := k.StoreCode(ctx, msg.Sender, msg.WASMByteCode)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeStoreCode,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
			),
		},
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg *MsgInstantiateContract) sdk.Result {
	contractAddr, err := k.InstantiateContract(ctx, msg.CodeID, msg.Sender, msg.InitMsg, msg.InitCoins)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeInstantiateContract,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr.String()),
			),
		},
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleExecute(ctx sdk.Context, k Keeper, msg *MsgExecuteContract) sdk.Result {
	res, err := k.ExecuteContract(ctx, msg.Contract, msg.Sender, msg.Coins, msg.Msg)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	res.Events = res.Events.AppendEvents(ctx.EventManager().Events())
	return res
}
