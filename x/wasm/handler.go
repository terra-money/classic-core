package wasm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgStoreCode:
			return handleStoreCode(ctx, k, &msg)
		case MsgInstantiateContract:
			return handleInstantiate(ctx, k, &msg)
		case MsgExecuteContract:
			return handleExecute(ctx, k, &msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

func handleStoreCode(ctx sdk.Context, k Keeper, msg *MsgStoreCode) (*sdk.Result, error) {
	codeID, err := k.StoreCode(ctx, msg.Sender, msg.WASMByteCode)
	if err != nil {
		return nil, err
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

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg *MsgInstantiateContract) (*sdk.Result, error) {
	contractAddr, err := k.InstantiateContract(ctx, msg.CodeID, msg.Sender, msg.InitMsg, msg.InitCoins)
	if err != nil {
		return nil, err
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

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleExecute(ctx sdk.Context, k Keeper, msg *MsgExecuteContract) (*sdk.Result, error) {
	res, err := k.ExecuteContract(ctx, msg.Contract, msg.Sender, msg.Msg, msg.Coins)
	if err != nil {
		return nil, err
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
	return &res, nil
}

// filterMessageEvents returns the same events with all of type == EventTypeMessage removed.
// this is so only our top-level message event comes through
func filterMessageEvents(manager *sdk.EventManager) sdk.Events {
	events := manager.Events()
	res := make([]sdk.Event, 0, len(events)+1)
	for _, e := range events {
		if e.Type != sdk.EventTypeMessage {
			res = append(res, e)
		}
	}
	return res
}
