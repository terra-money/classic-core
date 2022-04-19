package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"
	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = Querier{}
var _ wasm.WasmMsgParserInterface = MsgParser{}

// MsgParser - wasm msg parser for staking msgs
type MsgParser struct{}

// NewWasmMsgParser returns bank wasm msg parser
func NewWasmMsgParser() MsgParser {
	return MsgParser{}
}

// Parse implements wasm staking msg parser
func (MsgParser) Parse(_ sdk.AccAddress, _ wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	return nil, nil
}

// CosmosMsg only contains swap msg
type CosmosMsg struct {
	Swap     *types.MsgSwap     `json:"swap,omitempty"`
	SwapSend *types.MsgSwapSend `json:"swap_send,omitempty"`
}

// ParseCustom implements custom parser
func (MsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) (sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse market custom msg")
	}

	if sdkMsg.Swap != nil {
		sdkMsg.Swap.Trader = contractAddr.String()
		return sdkMsg.Swap, sdkMsg.Swap.ValidateBasic()
	} else if sdkMsg.SwapSend != nil {
		sdkMsg.SwapSend.FromAddress = contractAddr.String()
		return sdkMsg.SwapSend, sdkMsg.SwapSend.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Market")
}

// WasmQuerier - staking query interface for wasm contract
type Querier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

// Query - implement query function
func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

// CosmosQuery only contains swap simulation
type CosmosQuery struct {
	Swap *types.QuerySwapParams `json:"swap,omitempty"`
}

// SwapQueryResponse - swap simulation query response for wasm module
type SwapQueryResponse struct {
	Receive wasmvmtypes.Coin `json:"receive"`
}

// QueryCustom implements custom query interface
func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	q := keeper.NewQuerier(querier.keeper)
	if params.Swap != nil {
		res, err := q.Swap(sdk.WrapSDKContext(ctx), &types.QuerySwapRequest{
			OfferCoin: params.Swap.OfferCoin.String(),
			AskDenom:  params.Swap.AskDenom,
		})
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(SwapQueryResponse{Receive: wasm.EncodeSdkCoin(res.ReturnCoin)})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}

		return bz, err
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Market variant"}
}
