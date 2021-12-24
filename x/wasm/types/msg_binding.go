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
	WasmMsgParserRouteGov          = "gov"
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

// IBCWasmMsgParserInterface - stargate msg parsers
type IBCWasmMsgParserInterface interface {
	Parse(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (sdk.Msg, error)
}

// WasmCustomMsg - wasm custom msg parser
type WasmCustomMsg struct {
	Route   string          `json:"route"`
	MsgData json.RawMessage `json:"msg_data"`
}

// MsgParser - holds multiple module msg parsers
type MsgParser struct {
	Parsers        map[string]WasmMsgParserInterface
	IBCParser      IBCWasmMsgParserInterface
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

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteDistribution)
	case msg.Gov != nil:
		if parser, ok := p.Parsers[WasmMsgParserRouteGov]; ok {
			return parser.Parse(contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, WasmMsgParserRouteGov)
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
		if p.IBCParser != nil {
			return p.IBCParser.Parse(ctx, contractAddr, msg)
		}

		return nil, sdkerrors.Wrap(ErrNoRegisteredParser, "IBC")
	}

	return nil, sdkerrors.Wrap(ErrInvalidMsg, "failed to parse empty msg")
}
