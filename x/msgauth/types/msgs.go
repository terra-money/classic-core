package types

import (
	fmt "fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/gogo/protobuf/proto"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgGrantAuthorization{}
	_ sdk.Msg = &MsgRevokeAuthorization{}
	_ sdk.Msg = &MsgExecAuthorized{}
)

var _ codectypes.UnpackInterfacesMessage = MsgGrantAuthorization{}

// msgauth message types
const (
	TypeMsgGrantAuthorization  = "grant_authorization"
	TypeMsgRevokeAuthorization = "revoke_authorization"
	TypeMsgExecAuthorized      = "exec_authorized"
)

// NewMsgGrantAuthorization returns new MsgGrantAuthorization instance
func NewMsgGrantAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, authorization AuthorizationI, period time.Duration) (*MsgGrantAuthorization, error) {
	msg, ok := authorization.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%T does not implement proto.Message", authorization)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &MsgGrantAuthorization{
		Granter:       granter.String(),
		Grantee:       grantee.String(),
		Authorization: any,
		Period:        period,
	}, nil
}

// Route implements sdk.Msg
func (msg MsgGrantAuthorization) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgGrantAuthorization) Type() string { return TypeMsgGrantAuthorization }

// GetSigners implements sdk.Msg
func (msg MsgGrantAuthorization) GetSigners() []sdk.AccAddress {
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{granter}
}

// GetSignBytes implements sdk.Msg
func (msg MsgGrantAuthorization) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic implements sdk.Msg
func (msg MsgGrantAuthorization) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid granter address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address (%s)", err)
	}

	if msg.Granter == msg.Grantee {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "can not be grantee == granter")
	}

	if msg.Period <= 0 {
		return ErrInvalidPeriod
	}

	authorization := msg.GetAuthorization()
	if authorization == nil {
		return sdkerrors.Wrap(ErrInvalidAuthorization, "missing authorization")
	}

	if !IsGrantableMsgType(authorization.MsgType()) {
		return sdkerrors.Wrap(ErrInvalidMsgType, authorization.MsgType())
	}

	return nil
}

// GetAuthorization returns the grant Authorization
func (msg MsgGrantAuthorization) GetAuthorization() AuthorizationI {
	return msg.Authorization.GetCachedValue().(AuthorizationI)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgGrantAuthorization) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var authorization AuthorizationI
	return unpacker.UnpackAny(msg.Authorization, &authorization)
}

// NewMsgRevokeAuthorization returns new MsgRevokeAuthorization instance
func NewMsgRevokeAuthorization(granter sdk.AccAddress, grantee sdk.AccAddress, authorizationMsgType string) *MsgRevokeAuthorization {
	return &MsgRevokeAuthorization{
		Granter:              granter.String(),
		Grantee:              grantee.String(),
		AuthorizationMsgType: authorizationMsgType,
	}
}

// Route implements sdk.Msg
func (msg MsgRevokeAuthorization) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRevokeAuthorization) Type() string { return TypeMsgRevokeAuthorization }

// GetSigners implements sdk.Msg
func (msg MsgRevokeAuthorization) GetSigners() []sdk.AccAddress {
	granter, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{granter}
}

// GetSignBytes implements sdk.Msg
func (msg MsgRevokeAuthorization) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic implements sdk.Msg
func (msg MsgRevokeAuthorization) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Granter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid granter address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address (%s)", err)
	}

	return nil
}

// NewMsgExecAuthorized returns new MsgExecAuthorized instance
func NewMsgExecAuthorized(grantee sdk.AccAddress, msgs []sdk.Msg) (*MsgExecAuthorized, error) {
	anys := make([]*codectypes.Any, len(msgs))

	for i, msg := range msgs {
		var err error
		anys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, fmt.Errorf("%T does not implement proto.Message", msg)
		}
	}

	return &MsgExecAuthorized{
		Grantee: grantee.String(),
		Msgs:    anys,
	}, nil
}

// Route implements sdk.Msg
func (msg MsgExecAuthorized) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgExecAuthorized) Type() string { return TypeMsgExecAuthorized }

// GetSigners implements sdk.Msg
func (msg MsgExecAuthorized) GetSigners() []sdk.AccAddress {
	grantee, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{grantee}
}

// GetSignBytes implements sdk.Msg
func (msg MsgExecAuthorized) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic implements sdk.Msg
func (msg MsgExecAuthorized) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid grantee address (%s)", err)
	}

	if len(msg.Msgs) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cannot execute empty msgs")
	}

	sdkMsgs := msg.GetMsgs()
	for _, sdkMsg := range sdkMsgs {
		if sdkMsg == nil {
			return sdkerrors.Wrap(ErrInvalidMsg, "missing msg")
		}
	}

	return nil
}

// GetMsgs returns the grant Authorization
func (msg MsgExecAuthorized) GetMsgs() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, msg := range msg.Msgs {
		sdkMsgs = append(sdkMsgs, msg.GetCachedValue().(sdk.Msg))
	}

	return sdkMsgs
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgExecAuthorized) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, msg := range msg.Msgs {
		var sdkMsg sdk.Msg
		err := unpacker.UnpackAny(msg, &sdkMsg)

		if err != nil {
			return err
		}
	}

	return nil
}
