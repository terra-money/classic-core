package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/terra-project/core/x/wasm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the wasm MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	codeID, err := k.Keeper.StoreCode(ctx, senderAddr, msg.WASMByteCode)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeStoreCode,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &types.MsgStoreCodeResponse{
		CodeID: codeID,
	}, nil
}

func (k msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ownerAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	contractAddr, err := k.Keeper.InstantiateContract(ctx, msg.CodeID, ownerAddr, msg.InitMsg, msg.InitCoins, msg.Migratable)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeInstantiateContract,
				sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
			),
		},
	)

	return &types.MsgInstantiateContractResponse{
		ContractAddress: contractAddr.String(),
	}, nil
}

func (k msgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}

	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	data, err := k.Keeper.ExecuteContract(ctx, contractAddr, senderAddr, msg.ExecuteMsg, msg.Coins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeExecuteContract,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			),
		},
	)

	return &types.MsgExecuteContractResponse{
		Data: data,
	}, nil
}

func (k msgServer) MigrateContract(goCtx context.Context, msg *types.MsgMigrateContract) (*types.MsgMigrateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}

	ownerAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	data, err := k.Keeper.MigrateContract(ctx, contractAddr, ownerAddr, msg.NewCodeID, msg.MigrateMsg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeMigrateContract,
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.NewCodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
			),
		},
	)

	return &types.MsgMigrateContractResponse{
		Data: data,
	}, nil
}

func (k msgServer) UpdateContractOwner(goCtx context.Context, msg *types.MsgUpdateContractOwner) (*types.MsgUpdateContractOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return nil, err
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	if contractInfo.Owner != msg.Owner {
		return nil, sdkerrors.ErrUnauthorized
	}

	contractInfo.Owner = msg.NewOwner
	k.SetContractInfo(ctx, contractAddr, contractInfo)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateContractOwner,
				sdk.NewAttribute(types.AttributeKeyOwner, msg.NewOwner),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &types.MsgUpdateContractOwnerResponse{}, nil
}
