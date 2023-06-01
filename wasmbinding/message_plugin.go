package wasmbinding

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	//	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/classic-terra/core/v2/wasmbinding/bindings"
	marketkeeper "github.com/classic-terra/core/v2/x/market/keeper"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(market *marketkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:      old,
			marketKeeper: market,
		}
	}
}

type CustomMessenger struct {
	wrapped      wasmkeeper.Messenger
	marketKeeper *marketkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg bindings.TerraMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "terra msg")
		}

		switch {
		case contractMsg.Swap != nil:
			_, bz, err := m.swap(ctx, contractAddr, contractMsg.Swap)
			if err != nil {
				return nil, nil, sdkerrors.Wrap(err, "swap msg failed")
			}
			return nil, bz, nil

		case contractMsg.SwapSend != nil:
			_, bz, err := m.swapSend(ctx, contractAddr, contractMsg.SwapSend)
			if err != nil {
				return nil, nil, sdkerrors.Wrap(err, "swap msg failed")
			}
			return nil, bz, nil

		default:
			return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown terra msg variant"}
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

// swap wraps around performing market swap
func (m *CustomMessenger) swap(ctx sdk.Context, contractAddr sdk.AccAddress, contractMsg *bindings.Swap) ([]sdk.Event, [][]byte, error) {
	res, err := PerformSwap(m.marketKeeper, ctx, contractAddr, contractMsg)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform swap")
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "error marshal swap response")
	}

	return nil, [][]byte{bz}, nil
}

// PerformSwap performs market swap
func PerformSwap(f *marketkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, contractMsg *bindings.Swap) (*markettypes.MsgSwapResponse, error) {
	if contractMsg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "market swap msg was null"}
	}

	msgServer := marketkeeper.NewMsgServerImpl(*f)

	msgSwap := markettypes.NewMsgSwap(contractAddr, contractMsg.OfferCoin, contractMsg.AskDenom)

	if err := msgSwap.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed validating MsgSwap")
	}

	// swap
	res, err := msgServer.Swap(
		sdk.WrapSDKContext(ctx),
		msgSwap,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "swapping")
	}
	return res, nil
}

// swap wraps around performing market swap
func (m *CustomMessenger) swapSend(ctx sdk.Context, contractAddr sdk.AccAddress, contractMsg *bindings.SwapSend) ([]sdk.Event, [][]byte, error) {
	res, err := PerformSwapSend(m.marketKeeper, ctx, contractAddr, contractMsg)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform swap send")
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "error marshal swap send response")
	}

	return nil, [][]byte{bz}, nil
}

// PerformSwapSend performs market swap
func PerformSwapSend(f *marketkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, contractMsg *bindings.SwapSend) (*markettypes.MsgSwapSendResponse, error) {
	if contractMsg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "market swap send msg was null"}
	}

	msgServer := marketkeeper.NewMsgServerImpl(*f)

	toAddr, err := sdk.AccAddressFromBech32(contractMsg.ToAddress)
	if err != nil {
		return nil, err
	}

	msgSwapSend := markettypes.NewMsgSwapSend(contractAddr, toAddr, contractMsg.OfferCoin, contractMsg.AskDenom)

	if err := msgSwapSend.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed validating MsgSwapSend")
	}

	// swap
	res, err := msgServer.SwapSend(
		sdk.WrapSDKContext(ctx),
		msgSwapSend,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "swapping and sending")
	}
	return res, nil
}
