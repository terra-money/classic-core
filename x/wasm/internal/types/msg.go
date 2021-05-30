package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	core "github.com/terra-money/core/types"
)

// MsgStoreCode - struct for upload contract wasm byte codes
type MsgStoreCode struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode core.Base64Bytes `json:"wasm_byte_code" yaml:"wasm_byte_code"`
}

// NewMsgStoreCode creates a MsgStoreCode instance
func NewMsgStoreCode(sender sdk.AccAddress, wasmByteCode core.Base64Bytes) MsgStoreCode {
	return MsgStoreCode{
		Sender:       sender,
		WASMByteCode: wasmByteCode,
	}
}

// Route implements sdk.Msg
func (msg MsgStoreCode) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgStoreCode) Type() string {
	return "store_code"
}

// GetSignBytes implements sdk.Msg
func (msg MsgStoreCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgStoreCode) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty sender")
	}

	if len(msg.WASMByteCode) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty wasm code")
	}

	if uint64(len(msg.WASMByteCode)) > EnforcedMaxContractSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm code too large")
	}

	return nil
}

// MsgInstantiateContract - struct for instantiate contract from uploaded code
type MsgInstantiateContract struct {
	Owner      sdk.AccAddress   `json:"owner" yaml:"owner"`
	CodeID     uint64           `json:"code_id" yaml:"code_id"`
	InitMsg    core.Base64Bytes `json:"init_msg" yaml:"init_msg"`
	InitCoins  sdk.Coins        `json:"init_coins" yaml:"init_coins"`
	Migratable bool             `json:"migratable" yaml:"migratable"`
}

// NewMsgInstantiateContract creates a MsgInstantiateContract instance
func NewMsgInstantiateContract(owner sdk.AccAddress, codeID uint64, initMsg []byte, initCoins sdk.Coins, migratable bool) MsgInstantiateContract {
	return MsgInstantiateContract{
		Owner:      owner,
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
	return "instantiate_contract"
}

// ValidateBasic implements sdk.Msg
func (msg MsgInstantiateContract) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner")
	}

	if !msg.InitCoins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitCoins.String())
	}

	if uint64(len(msg.InitMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgExecuteContract - struct for execute instantiated contract with givn inner msg bytes
type MsgExecuteContract struct {
	Sender     sdk.AccAddress   `json:"sender" yaml:"sender"`
	Contract   sdk.AccAddress   `json:"contract" yaml:"contract"`
	ExecuteMsg core.Base64Bytes `json:"execute_msg" yaml:"execute_msg"`
	Coins      sdk.Coins        `json:"coins" yaml:"coins"`
}

// NewMsgExecuteContract creates a NewMsgExecuteContract instance
func NewMsgExecuteContract(sender sdk.AccAddress, contract sdk.AccAddress, execMsg []byte, coins sdk.Coins) MsgExecuteContract {
	return MsgExecuteContract{
		Sender:     sender,
		Contract:   contract,
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
	return "execute_contract"
}

// ValidateBasic implements sdk.Msg
func (msg MsgExecuteContract) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender")
	}

	if msg.Contract.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing contract")
	}

	if !msg.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coins.String())
	}

	if uint64(len(msg.ExecuteMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgExecuteContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgMigrateContract - struct to migrate contract to new code
type MsgMigrateContract struct {
	Owner      sdk.AccAddress   `json:"owner" yaml:"owner"`
	Contract   sdk.AccAddress   `json:"contract" yaml:"contract"`
	NewCodeID  uint64           `json:"new_code_id" yaml:"new_code_id"`
	MigrateMsg core.Base64Bytes `json:"migrate_msg" yaml:"migrate_msg"`
}

// NewMsgMigrateContract creates a MsgMigrateContract instance
func NewMsgMigrateContract(owner, contract sdk.AccAddress, newCodeID uint64, migrateMsg core.Base64Bytes) MsgMigrateContract {
	return MsgMigrateContract{
		Owner:      owner,
		Contract:   contract,
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
	return "migrate_contract"
}

// ValidateBasic implements sdk.Msg
func (msg MsgMigrateContract) ValidateBasic() error {
	if msg.NewCodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing new_code_id")
	}

	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}

	if msg.Contract.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing contract address")
	}

	if uint64(len(msg.MigrateMsg)) > EnforcedMaxContractMsgSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm msg byte size is too huge")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgMigrateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgMigrateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgUpdateContractOwner - struct for update contract owner
type MsgUpdateContractOwner struct {
	Owner    sdk.AccAddress `json:"owner" yaml:"owner"`
	NewOwner sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
}

// NewMsgUpdateContractOwner creates a MsgUpdateContractOwner instance
func NewMsgUpdateContractOwner(owner, newOwner, contract sdk.AccAddress) MsgUpdateContractOwner {
	return MsgUpdateContractOwner{
		Owner:    owner,
		NewOwner: newOwner,
		Contract: contract,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateContractOwner) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateContractOwner) Type() string {
	return "update_contract_owner"
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateContractOwner) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty owner")
	}

	if msg.NewOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty new owner")
	}

	if msg.Owner.Equals(msg.NewOwner) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "new owner must be different from the origin owner")
	}

	if msg.Contract.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty contract address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateContractOwner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateContractOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
