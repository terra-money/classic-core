package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterFeeShare{}
	_ sdk.Msg = &MsgCancelFeeShare{}
	_ sdk.Msg = &MsgUpdateFeeShare{}
)

const (
	TypeMsgRegisterFeeShare = "register_feeshare"
	TypeMsgCancelFeeShare   = "cancel_feeshare"
	TypeMsgUpdateFeeShare   = "update_feeshare"
)

// NewMsgRegisterFeeShare creates new instance of MsgRegisterFeeShare
func NewMsgRegisterFeeShare(
	contract sdk.Address,
	deployer,
	withdrawer sdk.AccAddress,
) *MsgRegisterFeeShare {
	withdrawerAddress := ""
	if withdrawer != nil {
		withdrawerAddress = withdrawer.String()
	}

	return &MsgRegisterFeeShare{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdrawerAddress,
	}
}

// Route returns the name of the module
func (msg MsgRegisterFeeShare) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterFeeShare) Type() string { return TypeMsgRegisterFeeShare }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterFeeShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if msg.WithdrawerAddress != "" {
		if _, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress); err != nil {
			return sdkerrors.Wrapf(err, "invalid withdraw address %s", msg.WithdrawerAddress)
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterFeeShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterFeeShare) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{from}
}

// NewMsgCancelFeeShare creates new instance of MsgCancelFeeShare.
func NewMsgCancelFeeShare(
	contract sdk.Address,
	deployer sdk.AccAddress,
) *MsgCancelFeeShare {
	return &MsgCancelFeeShare{
		ContractAddress: contract.String(),
		DeployerAddress: deployer.String(),
	}
}

// Route returns the message route for a MsgCancelFeeShare.
func (msg MsgCancelFeeShare) Route() string { return RouterKey }

// Type returns the message type for a MsgCancelFeeShare.
func (msg MsgCancelFeeShare) Type() string { return TypeMsgCancelFeeShare }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelFeeShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelFeeShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCancelFeeShare) GetSigners() []sdk.AccAddress {
	funder, _ := sdk.AccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{funder}
}

// NewMsgUpdateFeeShare creates new instance of MsgUpdateFeeShare
func NewMsgUpdateFeeShare(
	contract sdk.Address,
	deployer,
	withdraw sdk.AccAddress,
) *MsgUpdateFeeShare {
	return &MsgUpdateFeeShare{
		ContractAddress:   contract.String(),
		DeployerAddress:   deployer.String(),
		WithdrawerAddress: withdraw.String(),
	}
}

// Route returns the name of the module
func (msg MsgUpdateFeeShare) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateFeeShare) Type() string { return TypeMsgUpdateFeeShare }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateFeeShare) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DeployerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.DeployerAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.ContractAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid contract address %s", msg.ContractAddress)
	}

	if _, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress); err != nil {
		return sdkerrors.Wrapf(err, "invalid withdraw address %s", msg.WithdrawerAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateFeeShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateFeeShare) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.DeployerAddress)
	return []sdk.AccAddress{from}
}
