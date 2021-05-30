package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/terra-money/core/x/wasm/internal/types"
)

var _ types.WasmQuerierInterface = WasmQuerier{}
var _ types.WasmMsgParserInterface = WasmMsgParser{}

// WasmMsgParser - wasm msg parser for staking msgs
type WasmMsgParser struct{}

// NewWasmMsgParser returns wasm msg parser
func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

// Parse implements wasm staking msg parser
func (WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	msg := wasmMsg.Wasm

	if msg.Execute != nil {
		destContractAddr, err := sdk.AccAddressFromBech32(msg.Execute.ContractAddr)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Execute.ContractAddr)
		}
		coins, err := types.ParseToCoins(msg.Execute.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgExecuteContract{
			Sender:     contractAddr,
			Contract:   destContractAddr,
			ExecuteMsg: msg.Execute.Msg,
			Coins:      coins,
		}
		return []sdk.Msg{sdkMsg}, nil
	}

	if msg.Instantiate != nil {
		coins, err := types.ParseToCoins(msg.Instantiate.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgInstantiateContract{
			Owner:     contractAddr,
			CodeID:    msg.Instantiate.CodeID,
			InitMsg:   msg.Instantiate.Msg,
			InitCoins: coins,
		}
		return []sdk.Msg{sdkMsg}, nil
	}

	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Wasm")
}

// ParseCustom implements custom parser
func (parser WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	return nil, nil
}

// WasmQuerier - wasm query interface for wasm contract
type WasmQuerier struct {
	keeper Keeper
}

// NewWasmQuerier returns wasm querier
func NewWasmQuerier(keeper Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error) {
	if request.Wasm.Smart != nil {
		addr, err := sdk.AccAddressFromBech32(request.Wasm.Smart.ContractAddr)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Wasm.Smart.ContractAddr)
		}

		return querier.keeper.queryToContract(ctx, addr, request.Wasm.Smart.Msg)
	}
	if request.Wasm.Raw != nil {
		addr, err := sdk.AccAddressFromBech32(request.Wasm.Raw.ContractAddr)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Wasm.Raw.ContractAddr)
		}

		models := querier.keeper.queryToStore(ctx, addr, request.Wasm.Raw.Key)
		return json.Marshal(models)
	}

	return nil, wasmTypes.UnsupportedRequest{Kind: "unknown WasmQuery variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}
