package msgauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-project/core/x/msgauth/internal/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgGrantAuthorization:
			return handleMsgGrantAuthorization(ctx, msg, k)
		case MsgRevokeAuthorization:
			return handleMsgRevokeAuthorization(ctx, msg, k)
		case MsgExecAuthorized:
			return handleMsgExecAuthorized(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized authorization message type: %T", msg)
		}
	}
}

func handleMsgGrantAuthorization(ctx sdk.Context, msg MsgGrantAuthorization, k Keeper) (*sdk.Result, error) {
	k.Grant(ctx, msg.Grantee, msg.Granter, msg.Authorization, msg.Expiration)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventGrantAuthorization,
			sdk.NewAttribute(types.AttributeKeyGrantType, msg.Authorization.MsgType()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter.String()),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRevokeAuthorization(ctx sdk.Context, msg MsgRevokeAuthorization, k Keeper) (*sdk.Result, error) {
	err := k.Revoke(ctx, msg.Grantee, msg.Granter, msg.AuthorizationMsgType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventRevokeAuthorization,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyGrantType, msg.AuthorizationMsgType),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter.String()),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgExecAuthorized(ctx sdk.Context, msg MsgExecAuthorized, k Keeper) (*sdk.Result, error) {
	err := k.DispatchActions(ctx, msg.Grantee, msg.Msgs)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
