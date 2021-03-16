package types

import (
	"encoding/json"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Routes of pre-determined wasm querier
const (
	WasmMsgParserRouteBank    = "bank"
	WasmMsgParserRouteStaking = "staking"
	WasmMsgParserRouteMarket  = "market"
	WasmMsgParserRouteWasm    = "wasm"
)

// WasmMsgParserInterface - msg parsers of each module
type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

// WasmCustomMsg - wasm custom msg parser
type WasmCustomMsg struct {
	Route   string          `json:"route"`
	MsgData json.RawMessage `json:"msg_data"`
}

// MsgParser - holds multiple module msg parsers
type MsgParser map[string]WasmMsgParserInterface

// NewModuleMsgParser returns wasm msg parser
func NewModuleMsgParser() MsgParser {
	return make(MsgParser)
}

// Parse convert Wasm raw msg to chain msg
func (p MsgParser) Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Bank != nil:
		if parser, ok := p[WasmMsgParserRouteBank]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteBank)
	case msg.Custom != nil:
		var customMsg WasmCustomMsg
		err := json.Unmarshal(msg.Custom, &customMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if parser, ok := p[customMsg.Route]; ok {
			return parser.ParseCustom(contractAddr, customMsg.MsgData)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, customMsg.Route)
	case msg.Staking != nil:
		if parser, ok := p[WasmMsgParserRouteStaking]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteStaking)
	case msg.Wasm != nil:
		if parser, ok := p[WasmMsgParserRouteWasm]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteWasm)
	}
	return nil, sdkerrors.Wrap(ErrInvalidMsg, "failed to parse empty msg")
}
