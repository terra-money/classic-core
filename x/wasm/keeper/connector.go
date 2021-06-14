package keeper

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/terra-money/core/custom/auth/ante"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

func (k Keeper) dispatchAll(ctx sdk.Context, contractAddr sdk.AccAddress, subMsgs []wasmvmtypes.SubMsg, msgs []wasmvmtypes.CosmosMsg) error {
	// first dispatch all submessages (and the replies).
	err := k.dispatchSubmessages(ctx, contractAddr, subMsgs...)
	if err != nil {
		return err
	}

	// then dispatch all the normal messages
	return k.dispatchMessages(ctx, contractAddr, msgs...)
}

// dispatchSubmessages builds a sandbox to execute these messages and returns the execution result to the contract
// that dispatched them, both on success as well as failure
func (k Keeper) dispatchSubmessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs ...wasmvmtypes.SubMsg) error {
	for _, msg := range msgs {
		// first, we build a sub-context which we can use inside the submessages
		subCtx, commit := ctx.CacheContext()

		// check how much gas left locally, optionally wrap the gas meter
		gasRemaining := ctx.GasMeter().Limit() - ctx.GasMeter().GasConsumed()
		limitGas := msg.GasLimit != nil && (*msg.GasLimit < gasRemaining)

		var err error
		var events sdk.Events
		var data []byte
		if limitGas {
			events, data, err = k.dispatchMessageWithGasLimit(subCtx, contractAddr, msg.Msg, *msg.GasLimit)
		} else {
			events, data, err = k.dispatchMessage(subCtx, contractAddr, msg.Msg)
		}

		// if it succeeds, commit state changes from submessage, and pass on events to Event Manager
		// on failure, revert state from sandbox, and ignore events (just skip doing the above)
		if err == nil {
			commit()
			ctx.EventManager().EmitEvents(events)
		}

		// we only callback if requested. Short-circuit here the two cases we don't want to
		// 1. reply on success but error happened
		// 2. reply on failure but no error happened
		if (msg.ReplyOn == wasmvmtypes.ReplySuccess && err != nil) ||
			(msg.ReplyOn == wasmvmtypes.ReplyError && err == nil) {
			return err
		}

		// otherwise, we create a SubcallResult and pass it into the calling contract
		var result wasmvmtypes.SubcallResult
		if err == nil {
			// just take the first one for now if there are multiple sub-sdk messages
			// and safely return nothing if no data
			result = wasmvmtypes.SubcallResult{
				Ok: &wasmvmtypes.SubcallResponse{
					Events: types.EncodeSdkEvents(events),
					Data:   data,
				},
			}
		} else {
			result = wasmvmtypes.SubcallResult{
				Err: err.Error(),
			}
		}

		// now handle the reply, we use the parent context, and abort on error
		reply := wasmvmtypes.Reply{
			ID:     msg.ID,
			Result: result,
		}

		// we can ignore any result returned as the events are
		// already in the ctx.EventManager()
		err = k.reply(ctx, contractAddr, reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) dispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs ...wasmvmtypes.CosmosMsg) error {
	for _, msg := range msgs {
		events, _, err := k.dispatchMessage(ctx, contractAddr, msg)
		if err != nil {
			return err
		}

		// redispatch all events, (type sdk.EventTypeMessage will be filtered out in the handler)
		ctx.EventManager().EmitEvents(events)
	}

	return nil
}

// dispatchMessageWithGasLimit does not emit events to prevent duplicate emission
func (k Keeper) dispatchMessageWithGasLimit(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg, gasLimit uint64) (events sdk.Events, data []byte, err error) {
	subCtx := ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))

	// catch out of gas panic and just charge the entire gas limit
	defer func() {
		if r := recover(); r != nil {
			// if it's not an OutOfGas error, raise it again
			if _, ok := r.(sdk.ErrorOutOfGas); !ok {
				// log it to get the original stack trace somewhere (as panic(r) keeps message but stacktrace to here
				k.Logger(ctx).Info("SubMsg rethrow panic: %#v", r)
				panic(r)
			}

			ctx.GasMeter().ConsumeGas(gasLimit, "Sub-Message OutOfGas panic")
			err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "SubMsg hit gas limit")
		}
	}()

	events, data, err = k.dispatchMessage(subCtx, contractAddr, msg)

	// make sure we charge the parent what was spent
	spent := subCtx.GasMeter().GasConsumed()
	ctx.GasMeter().ConsumeGas(spent, "From limited Sub-Message")

	return events, data, err
}

// dispatchMessage does not emit events to prevent duplicate emission
func (k Keeper) dispatchMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (events sdk.Events, data []byte, err error) {
	sdkMsg, err := k.msgParser.Parse(ctx, contractAddr, msg)
	if err != nil {
		return nil, nil, err
	}

	if sdkMsg == nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInvalidMsg, "failed to parse msg")
	}

	// Charge tax on result msg
	taxes := ante.FilterMsgAndComputeTax(ctx, k.treasuryKeeper, sdkMsg)
	if !taxes.IsZero() {
		contractAcc := k.accountKeeper.GetAccount(ctx, contractAddr)
		if err := cosmosante.DeductFees(k.bankKeeper, ctx, contractAcc, taxes); err != nil {
			return nil, nil, err
		}
	}

	res, err := k.handleSdkMessage(ctx, contractAddr, sdkMsg)
	if err != nil {
		return nil, nil, err
	}

	// set data
	data = make([]byte, len(res.Data))
	copy(data, res.Data)

	// convert Tendermint.Events to sdk.Event
	sdkEvents := make(sdk.Events, len(res.Events))
	for i := range res.Events {
		sdkEvents[i] = sdk.Event(res.Events[i])
	}

	// append message action attribute
	events = append(events, sdkEvents...)

	return
}

func (k Keeper) handleSdkMessage(ctx sdk.Context, contractAddr sdk.AccAddress, msg sdk.Msg) (*sdk.Result, error) {
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract doesn't have permission")
		}
	}

	// find the handler and execute it
	h := k.serviceRouter.Handler(msg)
	if h == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, sdk.MsgTypeURL(msg))
	}

	res, err := h(ctx, msg)
	if err != nil {
		return nil, err
	}

	return res, nil
}
