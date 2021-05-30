package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/wasm/types"
)

var _ types.WasmQuerierInterface = WasmQuerier{}
var _ types.WasmMsgParserInterface = WasmMsgParser{}

// WasmMsgParser - wasm msg parser for wasm msgs
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

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.Instantiate != nil {
		coins, err := types.ParseToCoins(msg.Instantiate.Send)
		if err != nil {
			return nil, err
		}

		adminAddr := sdk.AccAddress{}
		if msg.Instantiate.Admin != "" {
			adminAddr, err = sdk.AccAddressFromBech32(msg.Instantiate.Admin)
			if err != nil {
				return nil, err
			}
		}

		// The contract instantiated from the other contract, always migratable
		cosmosMsg := types.NewMsgInstantiateContract(
			contractAddr,
			adminAddr,
			msg.Instantiate.CodeID,
			msg.Instantiate.Msg,
			coins,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.Migrate != nil {
		targetContractAddr, err := sdk.AccAddressFromBech32(msg.Migrate.ContractAddr)
		if err != nil {
			return nil, err
		}

		cosmosMsg := types.NewMsgMigrateContract(contractAddr, targetContractAddr, msg.Migrate.NewCodeID, msg.Migrate.Msg)
		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.UpdateAdmin != nil {
		targetContractAddr, err := sdk.AccAddressFromBech32(msg.UpdateAdmin.ContractAddr)
		if err != nil {
			return nil, err
		}

		newAdminAddr, err := sdk.AccAddressFromBech32(msg.UpdateAdmin.Admin)
		if err != nil {
			return nil, err
		}

		// current admin must be contractAddr
		cosmosMsg := types.NewMsgUpdateContractAdmin(contractAddr, newAdminAddr, targetContractAddr)
		return cosmosMsg, cosmosMsg.ValidateBasic()
	}

	if msg.ClearAdmin != nil {
		targetContractAddr, err := sdk.AccAddressFromBech32(msg.ClearAdmin.ContractAddr)
		if err != nil {
			return nil, err
		}

		cosmosMsg := types.NewMsgClearContractAdmin(contractAddr, targetContractAddr)
		return cosmosMsg, cosmosMsg.ValidateBasic()
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

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown WasmQuery variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}
