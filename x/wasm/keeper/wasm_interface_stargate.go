package keeper

import (
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/wasm/types"
)

var _ types.StargateWasmQuerierInterface = StargateWasmQuerier{}
var _ types.StargateWasmMsgParserInterface = StargateWasmMsgParser{}

// StargateWasmMsgParser - wasm msg parser for stargate msgs
type StargateWasmMsgParser struct {
	unpacker codectypes.AnyUnpacker
}

// NewStargateWasmMsgParser returns stargate wasm msg parser
func NewStargateWasmMsgParser(unpacker codectypes.AnyUnpacker) StargateWasmMsgParser {
	return StargateWasmMsgParser{unpacker}
}

// Parse implements wasm stargate msg parser
func (parser StargateWasmMsgParser) Parse(wasmMsg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	msg := wasmMsg.Stargate

	any := codectypes.Any{
		TypeUrl: msg.TypeURL,
		Value:   msg.Value,
	}

	var cosmosMsg sdk.Msg
	if err := parser.unpacker.UnpackAny(&any, &cosmosMsg); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("Cannot unpack proto message with type URL: %s", msg.TypeURL))
	}

	if err := codectypes.UnpackInterfaces(cosmosMsg, parser.unpacker); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("UnpackInterfaces inside msg: %s", err))
	}

	return cosmosMsg, nil
}

// StargateWasmQuerier - wasm query interface for wasm contract
type StargateWasmQuerier struct {
	keeper Keeper
}

// NewStargateWasmQuerier returns stargate wasm querier
func NewStargateWasmQuerier(keeper Keeper) StargateWasmQuerier {
	return StargateWasmQuerier{keeper}
}

var queryBlackList = []string{
	"/cosmos.tx",
	"/cosmos.base.tendermint",
}

// Query - implement query function
func (querier StargateWasmQuerier) Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error) {
	for _, b := range queryBlackList {
		if strings.Contains(request.Stargate.Path, b) {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", request.Stargate.Path)}
		}
	}

	route := querier.keeper.queryRouter.Route(request.Stargate.Path)
	if route == nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("No route to query '%s'", request.Stargate.Path)}
	}

	res, err := route(ctx, abci.RequestQuery{
		Data: request.Stargate.Data,
		Path: request.Stargate.Path,
	})

	if err != nil {
		return nil, err
	}

	return res.Value, nil
}
