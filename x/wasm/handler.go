package wasm

import (
	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"
)

// NewHandler returns a handler for "wasm" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		var (
			res proto.Message
			err error
		)

		switch msg := msg.(type) {
		case *types.MsgStoreCode:
			res, err = msgServer.StoreCode(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMigrateCode:
			res, err = msgServer.MigrateCode(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgInstantiateContract:
			res, err = msgServer.InstantiateContract(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgExecuteContract:
			res, err = msgServer.ExecuteContract(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMigrateContract:
			res, err = msgServer.MigrateContract(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgUpdateContractAdmin:
			res, err = msgServer.UpdateContractAdmin(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgClearContractAdmin:
			res, err = msgServer.ClearContractAdmin(sdk.WrapSDKContext(ctx), msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasm message type: %T", msg)
		}

		ctx = ctx.WithEventManager(filterMessageEvents(ctx))
		return sdk.WrapServiceResult(ctx, res, err)
	}
}

// filterMessageEvents returns the same events with all of type == EventTypeMessage removed except
// for wasm message types.
// this is so only our top-level message event comes through
func filterMessageEvents(ctx sdk.Context) *sdk.EventManager {
	m := sdk.NewEventManager()
	for _, e := range ctx.EventManager().Events() {
		if e.Type == sdk.EventTypeMessage &&
			!hasWasmModuleAttribute(e.Attributes) {
			continue
		}
		m.EmitEvent(e)
	}
	return m
}

func hasWasmModuleAttribute(attrs []abci.EventAttribute) bool {
	for _, a := range attrs {
		if sdk.AttributeKeyModule == string(a.Key) &&
			types.ModuleName == string(a.Value) {
			return true
		}
	}
	return false
}
