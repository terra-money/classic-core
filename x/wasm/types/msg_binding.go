package types

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Routes of pre-determined wasm querier
const (
	WasmMsgParserRouteBank         = "bank"
	WasmMsgParserRouteStaking      = "staking"
	WasmMsgParserRouteDistribution = "distribution"
	WasmMsgParserRouteMarket       = "market"
	WasmMsgParserRouteWasm         = "wasm"
)

// WasmMsgParserInterface - msg parsers of each module
type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) (sdk.Msg, error)
}

// StargateWasmMsgParserInterface - stargate msg parsers
type StargateWasmMsgParserInterface interface {
	Parse(msg wasmvmtypes.CosmosMsg) (sdk.Msg, error)
}

// WasmCustomMsg - wasm custom msg parser
type WasmCustomMsg struct {
	Route   string          `json:"route"`
	MsgData json.RawMessage `json:"msg_data"`
}

// MsgParser - holds multiple module msg parsers
type MsgParser struct {
	Parsers        map[string]WasmMsgParserInterface
	StargateParser StargateWasmMsgParserInterface
}

// NewWasmMsgParser returns wasm msg parser
func NewWasmMsgParser() MsgParser {
	return MsgParser{
		Parsers: make(map[string]WasmMsgParserInterface),
	}
}

// Parse convert Wasm raw msg to chain msg
func (p MsgParser) Parse(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	switch {
	case msg.Bank != nil:
		if msg.Bank.Burn != nil {
			return nil, sdkerrors.Wrap(ErrNoRegisteredParser, "Burn not supported")
		}

		if parser, ok := p.Parsers[WasmMsgParserRouteBank]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteBank)
	case msg.Custom != nil:
		var customMsg WasmCustomMsg
		err := json.Unmarshal(msg.Custom, &customMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if parser, ok := p.Parsers[customMsg.Route]; ok {
			return parser.ParseCustom(contractAddr, customMsg.MsgData)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, customMsg.Route)
	case msg.Staking != nil:
		if parser, ok := p.Parsers[WasmMsgParserRouteStaking]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteStaking)
	case msg.Distribution != nil:
		if parser, ok := p.Parsers[WasmMsgParserRouteDistribution]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, "distribution")
	case msg.Wasm != nil:
		if parser, ok := p.Parsers[WasmMsgParserRouteWasm]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteWasm)
	case msg.Stargate != nil:
		if p.StargateParser != nil {
			return p.StargateParser.Parse(msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, "stargate")
	case msg.IBC != nil:
		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, "IBC not supported")
	}

	return nil, sdkerrors.Wrap(ErrInvalidMsg, "failed to parse empty msg")
}
