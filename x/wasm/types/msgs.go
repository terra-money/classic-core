package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStoreCode{}
	_ sdk.Msg = &MsgMigrateCode{}
	_ sdk.Msg = &MsgInstantiateContract{}
	_ sdk.Msg = &MsgExecuteContract{}
	_ sdk.Msg = &MsgMigrateContract{}
	_ sdk.Msg = &MsgUpdateContractAdmin{}
	_ sdk.Msg = &MsgClearContractAdmin{}
)

// wasm message types
const (
	TypeMsgStoreCode           = "store_code"
	TypeMsgMigrateCode         = "migrate_code"
	TypeMsgInstantiateContract = "instantiate_contract"
	TypeMsgExecuteContract     = "execute_contract"
	TypeMsgMigrateContract     = "migrate_contract"
	TypeMsgUpdateContractAdmin = "update_contract_admin"
	TypeMsgClearContractAdmin  = "clear_contract_admin"
)

// NewMsgStoreCode creates a MsgStoreCode instance
func NewMsgStoreCode(sender sdk.AccAddress, wasmByteCode []byte) *MsgStoreCode {
	return &MsgStoreCode{
		Sender:       sender.String(),
		WASMByteCode: wasmByteCode,
	}
}

// Route implements sdk.Msg
func (msg MsgStoreCode) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgStoreCode) Type() string { return TypeMsgStoreCode }

// GetSignBytes Implements Msg
func (msg MsgStoreCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg
func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgStoreCode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if len(msg.WASMByteCode) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty wasm code")
	}

	if uint64(len(msg.WASMByteCode)) > EnforcedMaxContractSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm code too large")
	}

	return nil
}

// NewMsgMigrateCode creates a MsgMigrateCode instance
// TODO - remove after columbus-5 update
func NewMsgMigrateCode(codeID uint64, sender sdk.AccAddress, wasmByteCode []byte) *MsgMigrateCode {
	return &MsgMigrateCode{
		CodeID:       codeID,
		Sender:       sender.String(),
		WASMByteCode: wasmByteCode,
	}
}

// Route implements sdk.Msg
func (msg MsgMigrateCode) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgMigrateCode) Type() string { return TypeMsgMigrateCode }

// GetSignBytes Implements Msg
func (msg MsgMigrateCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg
func (msg MsgMigrateCode) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgMigrateCode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if len(msg.WASMByteCode) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty wasm code")
	}

	if uint64(len(msg.WASMByteCode)) > EnforcedMaxContractSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm code too large")
	}

	return nil
}

// NewMsgInstantiateContract creates a MsgInstantiateContract instance
func NewMsgInstantiateContract(sender, admin sdk.AccAddress, codeID uint64, initMsg []byte, initCoins sdk.Coins) *MsgInstantiateContract {
	var adminAddr string
	if !admin.Empty() {
		adminAddr = admin.String()
	}

	return &MsgInstantiateContract{
		Sender:    sender.String(),
		Admin:     adminAddr,
		CodeID:    codeID,
		InitMsg:   initMsg,
		InitCoins: initCoins,
	}
}

// Route implements sdk.Msg
func (msg MsgInstantiateContract) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgInstantiateContract) Type() string {
	return TypeMsgInstantiateContract
}

// ValidateBasic implements sdk.Msg
func (msg MsgInstantiateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if len(msg.Admin) != 0 {
		_, err := sdk.AccAddressFromBech32(msg.Admin)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid admin address (%s)", err)
		}
	}

	if !msg.InitCoins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitCoins.String())
	}

	if uint64(len(msg.InitMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	if !json.Valid(msg.InitMsg) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte format is invalid json")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

// NewMsgExecuteContract creates a NewMsgExecuteContract instance
func NewMsgExecuteContract(sender sdk.AccAddress, contract sdk.AccAddress, execMsg []byte, coins sdk.Coins) *MsgExecuteContract {
	return &MsgExecuteContract{
		Sender:     sender.String(),
		Contract:   contract.String(),
		ExecuteMsg: execMsg,
		Coins:      coins,
	}
}

// Route implements sdk.Msg
func (msg MsgExecuteContract) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgExecuteContract) Type() string {
	return TypeMsgExecuteContract
}

// ValidateBasic implements sdk.Msg
func (msg MsgExecuteContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	if !msg.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coins.String())
	}

	if uint64(len(msg.ExecuteMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	if !json.Valid(msg.ExecuteMsg) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte format is invalid json")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgExecuteContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender}
}

// NewMsgMigrateContract creates a MsgMigrateContract instance
func NewMsgMigrateContract(admin, contract sdk.AccAddress, newCodeID uint64, migrateMsg json.RawMessage) *MsgMigrateContract {
	return &MsgMigrateContract{
		Admin:      admin.String(),
		Contract:   contract.String(),
		NewCodeID:  newCodeID,
		MigrateMsg: migrateMsg,
	}
}

// Route implements sdk.Msg
func (msg MsgMigrateContract) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgMigrateContract) Type() string {
	return TypeMsgMigrateContract
}

// ValidateBasic implements sdk.Msg
func (msg MsgMigrateContract) ValidateBasic() error {
	if msg.NewCodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing new_code_id")
	}

	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid admin address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	if uint64(len(msg.MigrateMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	if !json.Valid(msg.MigrateMsg) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte format is invalid json")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgMigrateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgMigrateContract) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{admin}
}

// NewMsgUpdateContractAdmin creates a MsgUpdateContractAdmin instance
func NewMsgUpdateContractAdmin(admin, newAdmin, contract sdk.AccAddress) *MsgUpdateContractAdmin {
	return &MsgUpdateContractAdmin{
		Admin:    admin.String(),
		NewAdmin: newAdmin.String(),
		Contract: contract.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateContractAdmin) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateContractAdmin) Type() string {
	return TypeMsgUpdateContractAdmin
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateContractAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid admin address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid new admin address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateContractAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateContractAdmin) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// NewMsgClearContractAdmin creates a MsgClearContractAdmin instance
func NewMsgClearContractAdmin(admin, contract sdk.AccAddress) *MsgClearContractAdmin {
	return &MsgClearContractAdmin{
		Admin:    admin.String(),
		Contract: contract.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgClearContractAdmin) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgClearContractAdmin) Type() string {
	return TypeMsgClearContractAdmin
}

// ValidateBasic implements sdk.Msg
func (msg MsgClearContractAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgClearContractAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgClearContractAdmin) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}
