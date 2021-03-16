package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/terra-project/core/x/msgauth/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) GrantAuthorization(goCtx context.Context, msg *types.MsgGrantAuthorization) (*types.MsgGrantAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	expiration := ctx.BlockTime().Add(msg.Period)

	authorization := msg.GetAuthorization()
	if !types.IsGrantableMsgType(authorization.MsgType()) {
		return nil, types.ErrInvalidMsgType
	}

	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		return nil, err
	}

	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return nil, err
	}

	grant, err := types.NewAuthorizationGrant(authorization, expiration)
	if err != nil {
		return nil, err
	}

	_, found := k.GetGrant(ctx, granter, grantee, authorization.MsgType())
	if found {
		return nil, sdkerrors.Wrap(types.ErrGrantExists, "cannot regrant existing authorization")
	}

	k.SetGrant(ctx, granter, grantee, authorization.MsgType(), grant)
	k.InsertGrantQueue(ctx, granter, grantee, authorization.MsgType(), expiration)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventGrantAuthorization,
			sdk.NewAttribute(types.AttributeKeyGrantType, authorization.MsgType()),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgGrantAuthorizationResponse{}, nil
}

func (k msgServer) RevokeAuthorization(goCtx context.Context, msg *types.MsgRevokeAuthorization) (*types.MsgRevokeAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		return nil, err
	}

	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return nil, err
	}

	grant, found := k.GetGrant(ctx, granter, grantee, msg.AuthorizationMsgType)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "grant is not found")
	}

	k.RevokeGrant(ctx, granter, grantee, msg.AuthorizationMsgType)
	k.RevokeFromGrantQueue(ctx, granter, grantee, msg.AuthorizationMsgType, grant.Expiration)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventRevokeAuthorization,
			sdk.NewAttribute(types.AttributeKeyGrantType, msg.AuthorizationMsgType),
			sdk.NewAttribute(types.AttributeKeyGranterAddress, msg.Granter),
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgRevokeAuthorizationResponse{}, nil
}

func (k msgServer) ExecAuthorized(goCtx context.Context, msg *types.MsgExecAuthorized) (*types.MsgExecAuthorizedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return nil, err
	}

	err = k.DispatchActions(ctx, grantee, msg.GetMsgs())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventExecuteAuthorization,
			sdk.NewAttribute(types.AttributeKeyGranteeAddress, msg.Grantee),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &types.MsgExecAuthorizedResponse{}, nil
}
