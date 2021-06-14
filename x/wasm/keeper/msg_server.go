package keeper

import (
	"context"
	"fmt"

	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) MigrateCode(goCtx context.Context, msg *types.MsgMigrateCode) (*types.MsgMigrateCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.MigrateCode(ctx, msg.CodeID, senderAddr, msg.WASMByteCode)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeMigrateCode,
				sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &types.MsgMigrateCodeResponse{}, nil
}

func (k msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	adminAddr := sdk.AccAddress{}
	if len(msg.Admin) != 0 {
		adminAddr, err = sdk.AccAddressFromBech32(msg.Admin)
		if err != nil {
			return nil, err
		}
	}

	maxGas := k.MaxContractGas(ctx)
	remain := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	if remain > maxGas {
		remain = maxGas
	}

	subCtx := ctx.WithEventManager(sdk.NewEventManager()).WithGasMeter(sdk.NewGasMeter(remain))
	contractAddr, data, err := k.Keeper.InstantiateContract(
		subCtx,
		msg.CodeID,
		senderAddr,
		adminAddr,
		msg.InitMsg,
		msg.InitCoins,
	)
	if err != nil {
		return nil, err
	}

	// consume gas used from wasm execution
	ctx.GasMeter().ConsumeGas(subCtx.GasMeter().GasConsumed(), "wasm vm execute")

	// prepend the event to keep the events order
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeInstantiateContract,
				sdk.NewAttribute(types.AttributeKeyCreator, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
				sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
				sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			),
		}.AppendEvents(subCtx.EventManager().Events()),
	)

	return &types.MsgInstantiateContractResponse{
		ContractAddress: contractAddr.String(),
		Data:            data,
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

	maxGas := k.MaxContractGas(ctx)
	remain := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	if remain > maxGas {
		remain = maxGas
	}

	subCtx := ctx.WithEventManager(sdk.NewEventManager()).WithGasMeter(sdk.NewGasMeter(remain))
	data, err := k.Keeper.ExecuteContract(
		subCtx,
		contractAddr,
		senderAddr,
		msg.ExecuteMsg,
		msg.Coins,
	)
	if err != nil {
		return nil, err
	}

	// consume gas used from wasm execution
	ctx.GasMeter().ConsumeGas(subCtx.GasMeter().GasConsumed(), "wasm vm execute")

	// prepend the event to keep the events order
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
		}.AppendEvents(subCtx.EventManager().Events()),
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

	adminAddr, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, err
	}

	maxGas := k.MaxContractGas(ctx)
	remain := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
	if remain > maxGas {
		remain = maxGas
	}

	subCtx := ctx.WithEventManager(sdk.NewEventManager()).WithGasMeter(sdk.NewGasMeter(remain))
	data, err := k.Keeper.MigrateContract(
		subCtx,
		contractAddr,
		adminAddr,
		msg.NewCodeID,
		msg.MigrateMsg,
	)
	if err != nil {
		return nil, err
	}

	// consume gas used from wasm execution
	ctx.GasMeter().ConsumeGas(subCtx.GasMeter().GasConsumed(), "wasm vm execute")

	// prepend the event to keep the events order
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
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Admin),
			),
		}.AppendEvents(subCtx.EventManager().Events()),
	)

	return &types.MsgMigrateContractResponse{
		Data: data,
	}, nil
}

func (k msgServer) UpdateContractAdmin(goCtx context.Context, msg *types.MsgUpdateContractAdmin) (*types.MsgUpdateContractAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return nil, err
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	if contractInfo.Admin != msg.Admin {
		return nil, sdkerrors.ErrUnauthorized
	}

	contractInfo.Admin = msg.NewAdmin
	k.SetContractInfo(ctx, contractAddr, contractInfo)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateContractAdmin,
				sdk.NewAttribute(types.AttributeKeyAdmin, msg.NewAdmin),
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &types.MsgUpdateContractAdminResponse{}, nil
}

func (k msgServer) ClearContractAdmin(goCtx context.Context, msg *types.MsgClearContractAdmin) (*types.MsgClearContractAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, err
	}

	_, err = sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, err
	}

	contractInfo, err := k.GetContractInfo(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	if contractInfo.Admin != msg.Admin {
		return nil, sdkerrors.ErrUnauthorized
	}

	contractInfo.Admin = ""
	k.SetContractInfo(ctx, contractAddr, contractInfo)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeClearContractAdmin,
				sdk.NewAttribute(types.AttributeKeyContractAddress, msg.Contract),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		},
	)

	return &types.MsgClearContractAdminResponse{}, nil
}
