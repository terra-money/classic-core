package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/terra-money/core/x/market/internal/keeper"
	"github.com/terra-money/core/x/market/internal/types"
	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = WasmQuerier{}
var _ wasm.WasmMsgParserInterface = WasmMsgParser{}

// WasmMsgParser - wasm msg parser for staking msgs
type WasmMsgParser struct{}

// NewWasmMsgParser returns bank wasm msg parser
func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

// Parse implements wasm staking msg parser
func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

// CosmosMsg only contains swap msg
type CosmosMsg struct {
	Swap     *types.MsgSwap     `json:"swap,omitempty"`
	SwapSend *types.MsgSwapSend `json:"swap_send,omitempty"`
}

// ParseCustom implements custom parser
func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse market custom msg")
	}

	if sdkMsg.Swap != nil {
		return []sdk.Msg{*sdkMsg.Swap}, sdkMsg.Swap.ValidateBasic()
	} else if sdkMsg.SwapSend != nil {
		return []sdk.Msg{*sdkMsg.SwapSend}, sdkMsg.SwapSend.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Market")
}

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

// CosmosQuery only contains swap simulation
type CosmosQuery struct {
	Swap types.QuerySwapParams `json:"swap"`
}

// SwapQueryResponse - swap simulation query response for wasm module
type SwapQueryResponse struct {
	Receive wasmTypes.Coin `json:"receive"`
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	retCoin, err := keeper.QuerySwap(ctx, params.Swap, querier.keeper)
	if err != nil {
		return nil, err
	}

	bz, err := json.Marshal(SwapQueryResponse{Receive: wasm.EncodeSdkCoin(retCoin)})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, err
}
