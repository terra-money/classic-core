package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	core "github.com/terra-project/core/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStoreCode{}
	_ sdk.Msg = &MsgInstantiateContract{}
	_ sdk.Msg = &MsgExecuteContract{}
	_ sdk.Msg = &MsgMigrateContract{}
	_ sdk.Msg = &MsgUpdateContractOwner{}
)

// wasm message types
const (
	TypeMsgStoreCode           = "store_code"
	TypeMsgInstantiateContract = "instantiate_contract"
	TypeMsgExecuteContract     = "execute_contract"
	TypeMsgMigrateContract     = "migrate_contract"
	TypeMsgUpdateContractOwner = "update_contract_owner"
)

// NewMsgStoreCode creates a MsgStoreCode instance
func NewMsgStoreCode(sender sdk.AccAddress, wasmByteCode core.Base64Bytes) *MsgStoreCode {
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
	trader, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{trader}
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

// NewMsgInstantiateContract creates a MsgInstantiateContract instance
func NewMsgInstantiateContract(owner sdk.AccAddress, codeID uint64, initMsg []byte, initCoins sdk.Coins, migratable bool) *MsgInstantiateContract {
	return &MsgInstantiateContract{
		Owner:      owner.String(),
		CodeID:     codeID,
		InitMsg:    initMsg,
		InitCoins:  initCoins,
		Migratable: migratable,
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
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	if !msg.InitCoins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitCoins.String())
	}

	if uint64(len(msg.InitMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	if json.Valid(msg.InitMsg) {
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
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{owner}
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

	if json.Valid(msg.ExecuteMsg) {
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
func NewMsgMigrateContract(owner, contract sdk.AccAddress, newCodeID uint64, migrateMsg json.RawMessage) *MsgMigrateContract {
	return &MsgMigrateContract{
		Owner:      owner.String(),
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

	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	if uint64(len(msg.MigrateMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	if json.Valid(msg.MigrateMsg) {
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
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{owner}
}

// NewMsgUpdateContractOwner creates a MsgUpdateContractOwner instance
func NewMsgUpdateContractOwner(owner, newOwner, contract sdk.AccAddress) *MsgUpdateContractOwner {
	return &MsgUpdateContractOwner{
		Owner:    owner.String(),
		NewOwner: newOwner.String(),
		Contract: contract.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateContractOwner) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateContractOwner) Type() string {
	return TypeMsgUpdateContractOwner
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateContractOwner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid new owner address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid contract address (%s)", err)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateContractOwner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateContractOwner) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}
