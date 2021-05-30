package wasm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-money/core/x/wasm/internal/types"
)

// NewHandler returns a handler for "wasm" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		case MsgMigrateContract:
			return handleMigrate(ctx, k, msg)
		case MsgUpdateContractOwner:
			return handleUpdateContractOwner(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasm message type: %T", msg)
		}
	}
}

func handleStoreCode(ctx sdk.Context, k Keeper, msg MsgStoreCode) (*sdk.Result, error) {
	codeID, err := k.StoreCode(ctx, msg.Sender, msg.WASMByteCode)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeStoreCode,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleInstantiate(ctx sdk.Context, k Keeper, msg MsgInstantiateContract) (*sdk.Result, error) {
	contractAddr, err := k.InstantiateContract(ctx, msg.CodeID, msg.Owner, msg.InitMsg, msg.InitCoins, msg.Migratable)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: filterMessageEvents(ctx.EventManager()).AppendEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeInstantiateContract,
				sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			),
		},
	)}, nil
}

func handleExecute(ctx sdk.Context, k Keeper, msg MsgExecuteContract) (*sdk.Result, error) {
	data, err := k.ExecuteContract(ctx, msg.Contract, msg.Sender, msg.ExecuteMsg, msg.Coins)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: filterMessageEvents(ctx.EventManager()).AppendEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeExecuteContract,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
		},
	), Data: data}, nil
}

func handleMigrate(ctx sdk.Context, k Keeper, msg MsgMigrateContract) (*sdk.Result, error) {
	data, err := k.MigrateContract(ctx, msg.Contract, msg.Owner, msg.NewCodeID, msg.MigrateMsg)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: filterMessageEvents(ctx.EventManager()).AppendEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeMigrateContract,
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.NewCodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			),
		},
	), Data: data}, nil
}

func handleUpdateContractOwner(ctx sdk.Context, k Keeper, msg MsgUpdateContractOwner) (*sdk.Result, error) {
	contractInfo, err := k.GetContractInfo(ctx, msg.Contract)
	if err != nil {
		return nil, err
	}

	if !contractInfo.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.ErrUnauthorized
	}

	contractInfo.Owner = msg.NewOwner
	k.SetContractInfo(ctx, msg.Contract, contractInfo)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateContractOwner,
				sdk.NewAttribute(types.AttributeKeyOwner, msg.NewOwner.String()),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
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
