package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/classic-terra/core/x/feeshare/types"
)

var _ types.MsgServer = &Keeper{}

// GetContractAdminOrCreatorAddress ensures the deployer is the contract's admin OR creator if no admin is set for all msg_server feeshare functions.
func (k Keeper) GetContractAdminOrCreatorAddress(ctx sdk.Context, contract sdk.AccAddress, deployer string) (sdk.AccAddress, error) {
	var controllingAccount sdk.AccAddress

	// Ensures deployer String is valid
	_, err := sdk.AccAddressFromBech32(deployer)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid deployer address %s", deployer)
	}

	info, err := k.wasmKeeper.GetContractInfo(ctx, contract)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "contract not found %s", contract)
	}

	if len(info.Admin) == 0 {
		// no admin, see if they are the creator of the contract
		if info.Creator != deployer {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "you are not the creator of this contract %s", info.Creator)
		}

		creatorAddr, err := sdk.AccAddressFromBech32(info.Creator)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address %s", info.Creator)
		}
		controllingAccount = creatorAddr
	} else {
		// Admin is set, so we check if the deployer is the admin
		if info.Admin != deployer {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "you are not an admin of this contract %s", deployer)
		}

		adminAddr, err := sdk.AccAddressFromBech32(info.Admin)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address %s", info.Admin)
		}
		controllingAccount = adminAddr
	}

	return controllingAccount, nil
}

// RegisterFeeShare registers a contract to receive transaction fees
func (k Keeper) RegisterFeeShare(
	goCtx context.Context,
	msg *types.MsgRegisterFeeShare,
) (*types.MsgRegisterFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	// Get Contract
	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}

	// Check if contract is already registered
	if k.IsFeeShareRegistered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrFeeShareAlreadyRegistered, "contract is already registered %s", contract)
	}

	// Check that the person who signed the message is the wasm contract admin, if so return the deployer address
	deployer, err := k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
	if err != nil {
		return nil, err
	}

	// Get the withdraw address of the contract
	withdrawer, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address %s", msg.WithdrawerAddress)
	}

	// prevent storing the same address for deployer and withdrawer
	feeshare := types.NewFeeShare(contract, deployer, withdrawer)
	k.SetFeeShare(ctx, feeshare)
	k.SetDeployerMap(ctx, deployer, contract)
	k.SetWithdrawerMap(ctx, withdrawer, contract)

	k.Logger(ctx).Debug(
		"registering contract for transaction fees",
		"contract", msg.ContractAddress,
		"deployer", msg.DeployerAddress,
		"withdraw", msg.WithdrawerAddress,
	)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegisterFeeShare,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &types.MsgRegisterFeeShareResponse{}, nil
}

// UpdateFeeShare updates the withdraw address of a given FeeShare. If the given
// withdraw address is empty or the same as the deployer address, the withdraw
// address is removed.
func (k Keeper) UpdateFeeShare(
	goCtx context.Context,
	msg *types.MsgUpdateFeeShare,
) (*types.MsgUpdateFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid contract address (%s)", err,
		)
	}

	feeshare, found := k.GetFeeShare(ctx, contract)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrFeeShareContractNotRegistered,
			"contract %s is not registered", msg.ContractAddress,
		)
	}

	// feeshare with the given withdraw address is already registered
	if msg.WithdrawerAddress == feeshare.WithdrawerAddress {
		return nil, sdkerrors.Wrapf(types.ErrFeeShareAlreadyRegistered, "feeshare with withdraw address %s is already registered", msg.WithdrawerAddress)
	}

	// Check that the person who signed the message is the wasm contract admin, if so return the deployer address
	_, err = k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
	if err != nil {
		return nil, err
	}

	withdrawAddr, err := sdk.AccAddressFromBech32(feeshare.WithdrawerAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid withdrawer address (%s)", err,
		)
	}
	newWithdrawAddr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid WithdrawerAddress %s", msg.WithdrawerAddress)
	}

	k.DeleteWithdrawerMap(ctx, withdrawAddr, contract)
	k.SetWithdrawerMap(ctx, newWithdrawAddr, contract)

	// update feeshare
	feeshare.WithdrawerAddress = newWithdrawAddr.String()
	k.SetFeeShare(ctx, feeshare)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateFeeShare,
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &types.MsgUpdateFeeShareResponse{}, nil
}

// CancelFeeShare deletes the FeeShare for a given contract
func (k Keeper) CancelFeeShare(
	goCtx context.Context,
	msg *types.MsgCancelFeeShare,
) (*types.MsgCancelFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}

	fee, found := k.GetFeeShare(ctx, contract)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrFeeShareContractNotRegistered, "contract %s is not registered", msg.ContractAddress)
	}

	// Check that the person who signed the message is the wasm contract admin, if so return the deployer address
	_, err = k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
	if err != nil {
		return nil, err
	}

	k.DeleteFeeShare(ctx, fee)
	k.DeleteDeployerMap(
		ctx,
		fee.GetDeployerAddr(),
		contract,
	)

	withdrawAddr := fee.GetWithdrawerAddr()
	if withdrawAddr != nil {
		k.DeleteWithdrawerMap(
			ctx,
			withdrawAddr,
			contract,
		)
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCancelFeeShare,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress),
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
			),
		},
	)

	return &types.MsgCancelFeeShareResponse{}, nil
}
