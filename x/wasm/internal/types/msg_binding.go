package types

import (
	"encoding/json"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Routes of pre-determined wasm querier
const (
	WasmMsgParserRouteBank    = "bank"
	WasmMsgParserRouteStaking = "staking"
	WasmMsgParserRouteWasm    = "wasm"
)

// WasmMsgParser - msg parsers of each module
type WasmMsgParser interface {
	Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, sdk.Error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, sdk.Error)
}

// WasmCustomMsg - wasm custom msg parser
type WasmCustomMsg struct {
	Route   string          `json:"route"`
	MsgData json.RawMessage `json:"msg_data"`
}

// MsgParser - holds multiple module msg parsers
type MsgParser map[string]WasmMsgParser

// Parse convert Wasm raw msg to chain msg
func (p MsgParser) Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, sdk.Error) {
	switch {
	case msg.Bank != nil:
		if parser, ok := p[WasmMsgParserRouteBank]; ok {
			return parser.Parse(contractAddr, msg)
		} else {
			return nil, ErrNoRegisteredParser(WasmMsgParserRouteBank)
		}

	case msg.Custom != nil:
		var customMsg WasmCustomMsg
		err := json.Unmarshal(msg.Custom, &customMsg)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}

		if parser, ok := p[customMsg.Route]; ok {
			return parser.ParseCustom(contractAddr, customMsg.MsgData)
		} else {
			return nil, ErrNoRegisteredParser(customMsg.Route)
		}

	case msg.Staking != nil:
		if parser, ok := p[WasmMsgParserRouteStaking]; ok {
			return parser.Parse(contractAddr, msg)
		} else {
			return nil, ErrNoRegisteredParser(WasmMsgParserRouteStaking)
		}
	case msg.Wasm != nil:
		if parser, ok := p[WasmMsgParserRouteWasm]; ok {
			return parser.Parse(contractAddr, msg)
		} else {
			return nil, ErrNoRegisteredParser(WasmMsgParserRouteWasm)
		}
	}
	return nil, sdk.ErrInternal("failed to parse empty msg")
}
