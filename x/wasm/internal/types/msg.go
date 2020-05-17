package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	core "github.com/terra-project/core/types"
)

const (
	// MaxWasmSize 500 KB (hard-cap)
	MaxWasmSize = 500 * 1024
)

// MsgStoreCode - struct for upload contract wasm byte codes
type MsgStoreCode struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode core.HexBytes `json:"wasm_byte_code" yaml:"wasm_byte_code"`
}

// NewMsgStoreCode creates a MsgStoreCode instance
func NewMsgStoreCode(sender sdk.AccAddress, wasmByteCode core.HexBytes) MsgStoreCode {
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
	return "storecode"
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
	if len(msg.WASMByteCode) == 0 {

		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty wasm code")
	}
	if len(msg.WASMByteCode) > MaxWasmSize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wasm code too large")
	}
	return nil
}

// MsgInstantiateContract - struct for instantiate contract from uploaded code
type MsgInstantiateContract struct {
	Sender    sdk.AccAddress `json:"sender" yaml:"sender"`
	CodeID    uint64         `json:"code_id" yaml:"code_id"`
	InitMsg   core.HexBytes  `json:"init_msg" yaml:"init_msg"`
	InitCoins sdk.Coins      `json:"init_coins" yaml:"init_coins"`
}

// NewMsgInstantiateContract creates a MsgInstantiateContract instance
func NewMsgInstantiateContract(sender sdk.AccAddress, codeID uint64, initMsg []byte, initCoins sdk.Coins) MsgInstantiateContract {
	return MsgInstantiateContract{
		Sender:    sender,
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
	return "instantiatecontract"
}

// ValidateBasic implements sdk.Msg
func (msg MsgInstantiateContract) ValidateBasic() error {
	if !msg.InitCoins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.InitCoins.String())
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgExecuteContract - struct for execute instantiated contract with givn inner msg bytes
type MsgExecuteContract struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Msg      core.HexBytes  `json:"msg" yaml:"msg"`
	Coins    sdk.Coins      `json:"coins" yaml:"coins"`
}

// NewMsgExecuteContract creates a NewMsgExecuteContract instance
func NewMsgExecuteContract(sender sdk.AccAddress, contract sdk.AccAddress, msg []byte, coins sdk.Coins) MsgExecuteContract {
	return MsgExecuteContract{
		Sender:   sender,
		Contract: contract,
		Msg:      msg,
		Coins:    coins,
	}
}

// Route implements sdk.Msg
func (msg MsgExecuteContract) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgExecuteContract) Type() string {
	return "executecontract"
}

// ValidateBasic implements sdk.Msg
func (msg MsgExecuteContract) ValidateBasic() error {
	if !msg.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coins.String())
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
