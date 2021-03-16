package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	wasm "github.com/terra-project/core/x/wasm/exported"
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
func (WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	msg := wasmMsg.Bank

	if msg.Send == nil {
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Bank")
	}

	if len(msg.Send.Amount) == 0 {
		return nil, nil
	}

	_, stderr := sdk.AccAddressFromBech32(msg.Send.FromAddress)
	if stderr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Send.FromAddress)
	}

	_, stderr = sdk.AccAddressFromBech32(msg.Send.ToAddress)
	if stderr != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Send.ToAddress)
	}

	amount, err := wasm.ParseToCoins(msg.Send.Amount)
	if err != nil {
		return nil, err
	}

	sdkMsg := types.MsgSend{
		FromAddress: msg.Send.FromAddress,
		ToAddress:   msg.Send.ToAddress,
		Amount:      amount,
	}

	return []sdk.Msg{&sdkMsg}, sdkMsg.ValidateBasic()
}

// ParseCustom implements custom parser
func (WasmMsgParser) ParseCustom(_ sdk.AccAddress, _ json.RawMessage) ([]sdk.Msg, error) {
	return nil, nil
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
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error) {
	if request.Bank.AllBalances != nil {
		addr, err := sdk.AccAddressFromBech32(request.Bank.AllBalances.Address)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Bank.AllBalances.Address)
		}
		coins := querier.keeper.GetAllBalances(ctx, addr)
		res := wasmTypes.AllBalancesResponse{
			Amount: wasm.EncodeSdkCoins(coins),
		}
		return json.Marshal(res)
	}
	if request.Bank.Balance != nil {
		addr, err := sdk.AccAddressFromBech32(request.Bank.Balance.Address)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Bank.Balance.Address)
		}

		coin := querier.keeper.GetBalance(ctx, addr, request.Bank.Balance.Denom)
		res := wasmTypes.BalanceResponse{
			Amount: wasm.EncodeSdkCoin(coin),
		}
		return json.Marshal(res)
	}
	return nil, wasmTypes.UnsupportedRequest{Kind: "unknown BankQuery variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}
