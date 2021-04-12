package keeper

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-project/core/x/wasm/types"
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
func (WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
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

		cosmosMsg := types.NewMsgExecuteContract(contractAddr, destContractAddr, msg.Execute.Msg, coins)
		return cosmosMsg, nil
	}

	if msg.Instantiate != nil {
		coins, err := types.ParseToCoins(msg.Instantiate.Send)
		if err != nil {
			return nil, err
		}

		// The contract instantiated from the other contract, always migratable
		cosmosMsg := types.NewMsgInstantiateContract(contractAddr, msg.Instantiate.CodeID, msg.Instantiate.Msg, coins, true)
		return cosmosMsg, nil
	}

	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown variant of Wasm")
}

// ParseCustom implements custom parser
func (parser WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) (sdk.Msg, error) {
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
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error) {
	if request.Wasm != nil {
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

			return querier.keeper.queryToStore(ctx, addr, request.Wasm.Raw.Key), nil
		}
	}

	if request.Stargate != nil {
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

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown WasmQuery variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}
