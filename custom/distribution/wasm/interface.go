package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmMsgParserInterface = MsgParser{}

// MsgParser - wasm msg parser for staking msgs
type MsgParser struct{}

// NewWasmMsgParser returns staking wasm msg parser
func NewWasmMsgParser() MsgParser {
	return MsgParser{}
}

// Parse implements wasm staking msg parser
func (parser MsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmvmtypes.CosmosMsg) (msgs sdk.Msg, err error) {
	msg := wasmMsg.Distribution

	if msg.SetWithdrawAddress != nil {
		rcpt, err := sdk.AccAddressFromBech32(msg.SetWithdrawAddress.Address)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.SetWithdrawAddress.Address)
		}

		cosmosMsg := distrtypes.NewMsgSetWithdrawAddress(
			contractAddr,
			rcpt,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.WithdrawDelegatorReward != nil {
		validator, err := sdk.ValAddressFromBech32(msg.WithdrawDelegatorReward.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.WithdrawDelegatorReward.Validator)
		}

		cosmosMsg := distrtypes.NewMsgWithdrawDelegatorReward(
			contractAddr,
			validator,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Staking")
}

// ParseCustom implements custom parser
func (parser MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) (sdk.Msg, error) {
	return nil, nil
}
