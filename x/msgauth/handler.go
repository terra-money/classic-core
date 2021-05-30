package msgauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-money/core/x/msgauth/internal/types"
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
	expiration := ctx.BlockTime().Add(msg.Period)

	if !k.IsGrantable(msg.Authorization.MsgType()) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMsgType, "Msg %s is not allowed to grant", msg.Authorization.MsgType())
	}

	k.SetGrant(ctx, msg.Granter, msg.Grantee, NewAuthorizationGrant(msg.Authorization, expiration))
	k.InsertGrantQueue(ctx, msg.Granter, msg.Grantee, msg.Authorization.MsgType(), expiration)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventGrantAuthorization,
			sdk.NewAttribute(types.AttributeKeyGrantType, msg.Authorization.MsgType()),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter.String()),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRevokeAuthorization(ctx sdk.Context, msg MsgRevokeAuthorization, k Keeper) (*sdk.Result, error) {
	grant, found := k.GetGrant(ctx, msg.Granter, msg.Grantee, msg.AuthorizationMsgType)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "grant is not found")
	}

	k.RevokeGrant(ctx, msg.Granter, msg.Grantee, msg.AuthorizationMsgType)
	k.RevokeFromGrantQueue(ctx, msg.Granter, msg.Grantee, msg.AuthorizationMsgType, grant.Expiration)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventRevokeAuthorization,
			sdk.NewAttribute(types.AttributeKeyGrantType, msg.AuthorizationMsgType),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter.String()),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgExecAuthorized(ctx sdk.Context, msg MsgExecAuthorized, k Keeper) (*sdk.Result, error) {
	err := k.DispatchActions(ctx, msg.Grantee, msg.Msgs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventExecuteAuthorization,
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
