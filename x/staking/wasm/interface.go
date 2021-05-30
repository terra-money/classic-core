package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = WasmQuerier{}
var _ wasm.WasmMsgParserInterface = WasmMsgParser{}

// WasmMsgParser - wasm msg parser for staking msgs
type WasmMsgParser struct{}

// NewWasmMsgParser returns staking wasm msg parser
func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

// Parse implements wasm staking msg parser
func (parser WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmTypes.CosmosMsg) (msgs []sdk.Msg, err error) {
	msg := wasmMsg.Staking

	if msg.Delegate != nil {
		validator, err := sdk.ValAddressFromBech32(msg.Delegate.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Delegate.Validator)
		}

		coin, err := wasm.ParseToCoin(msg.Delegate.Amount)
		if err != nil {
			return nil, err
		}

		sdkMsg := staking.MsgDelegate{
			DelegatorAddress: contractAddr,
			ValidatorAddress: validator,
			Amount:           coin,
		}

		if err := sdkMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, sdkMsg)
	}

	if msg.Redelegate != nil {
		src, err := sdk.ValAddressFromBech32(msg.Redelegate.SrcValidator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Redelegate.SrcValidator)
		}
		dst, err := sdk.ValAddressFromBech32(msg.Redelegate.DstValidator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Redelegate.DstValidator)
		}
		coin, err := wasm.ParseToCoin(msg.Redelegate.Amount)
		if err != nil {
			return nil, err
		}

		sdkMsg := staking.MsgBeginRedelegate{
			DelegatorAddress:    contractAddr,
			ValidatorSrcAddress: src,
			ValidatorDstAddress: dst,
			Amount:              coin,
		}

		if err := sdkMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, sdkMsg)
	}

	if msg.Undelegate != nil {
		validator, err := sdk.ValAddressFromBech32(msg.Undelegate.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Undelegate.Validator)
		}

		coin, err := wasm.ParseToCoin(msg.Undelegate.Amount)
		if err != nil {
			return nil, err
		}

		sdkMsg := staking.MsgUndelegate{
			DelegatorAddress: contractAddr,
			ValidatorAddress: validator,
			Amount:           coin,
		}

		if err := sdkMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, sdkMsg)
	}

	if msg.Withdraw != nil && len(msg.Withdraw.Recipient) != 0 {
		rcpt, err := sdk.AccAddressFromBech32(msg.Withdraw.Recipient)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Recipient)
		}

		sdkMsg := distribution.MsgSetWithdrawAddress{
			DelegatorAddress: contractAddr,
			WithdrawAddress:  rcpt,
		}

		if err := sdkMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, sdkMsg)
	}

	if msg.Withdraw != nil && len(msg.Withdraw.Validator) != 0 {
		var err error

		validator, err := sdk.ValAddressFromBech32(msg.Withdraw.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Validator)
		}

		sdkMsg := distribution.MsgWithdrawDelegatorReward{
			DelegatorAddress: contractAddr,
			ValidatorAddress: validator,
		}

		if err := sdkMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, sdkMsg)
	}

	if len(msgs) == 0 {
		return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Staking")
	}

	return
}

// ParseCustom implements custom parser
func (parser WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	return nil, nil
}

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	keeper staking.Keeper
}

// NewWasmQuerier returns staking wasm querier
func NewWasmQuerier(keeper staking.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error) {
	if request.Staking.BondedDenom != nil {
		res := wasmTypes.BondedDenomResponse{
			Denom: querier.keeper.BondDenom(ctx),
		}

		return json.Marshal(res)
	}

	if request.Staking.Validators != nil {
		validators := querier.keeper.GetBondedValidatorsByPower(ctx)
		wasmVals := make([]wasmTypes.Validator, len(validators))

		for i, v := range validators {
			wasmVals[i] = wasmTypes.Validator{
				Address:       v.OperatorAddress.String(),
				Commission:    v.Commission.Rate.String(),
				MaxCommission: v.Commission.MaxRate.String(),
				MaxChangeRate: v.Commission.MaxChangeRate.String(),
			}
		}

		res := wasmTypes.ValidatorsResponse{
			Validators: wasmVals,
		}

		return json.Marshal(res)
	}

	if request.Staking.AllDelegations != nil {
		delegator, err := sdk.AccAddressFromBech32(request.Staking.AllDelegations.Delegator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Staking.AllDelegations.Delegator)
		}

		delegations := querier.keeper.GetAllDelegatorDelegations(ctx, delegator)

		responseDelegations, err := querier.encodeDelegations(ctx, delegations)
		if err != nil {
			return nil, err
		}

		res := wasmTypes.AllDelegationsResponse{
			Delegations: responseDelegations,
		}

		return json.Marshal(res)
	}

	if request.Staking.Delegation != nil {
		delegator, err := sdk.AccAddressFromBech32(request.Staking.Delegation.Delegator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Staking.Delegation.Delegator)
		}
		validator, err := sdk.ValAddressFromBech32(request.Staking.Delegation.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Staking.Delegation.Validator)
		}

		var responseFullDelegation wasmTypes.FullDelegation
		delegation, found := querier.keeper.GetDelegation(ctx, delegator, validator)
		if found {
			responseFullDelegation, err = querier.encodeDelegation(ctx, delegation)
			if err != nil {
				return nil, err
			}
		}

		res := wasmTypes.DelegationResponse{
			Delegation: &responseFullDelegation,
		}

		return json.Marshal(res)
	}

	return nil, wasmTypes.UnsupportedRequest{Kind: "unknown Staking variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}

// encdoe cosmos delegations to wasm delegations
func (querier WasmQuerier) encodeDelegations(ctx sdk.Context, dels staking.Delegations) (wasmTypes.Delegations, error) {
	bondDenom := querier.keeper.BondDenom(ctx)

	var responseDelegations wasmTypes.Delegations
	for _, del := range dels {
		val, found := querier.keeper.GetValidator(ctx, del.ValidatorAddress)
		if !found {
			return nil, sdkerrors.Wrap(staking.ErrNoValidatorFound, "can't load validator for delegation")
		}

		amount := sdk.NewCoin(bondDenom, val.TokensFromShares(del.Shares).TruncateInt())

		responseDelegations = append(responseDelegations, wasmTypes.Delegation{
			Delegator: del.DelegatorAddress.String(),
			Validator: del.ValidatorAddress.String(),
			Amount:    wasm.EncodeSdkCoin(amount),
		})
	}
	return responseDelegations, nil
}

// encode cosmos staking to wasm delegation
func (querier WasmQuerier) encodeDelegation(ctx sdk.Context, del staking.Delegation) (wasmTypes.FullDelegation, error) {
	val, found := querier.keeper.GetValidator(ctx, del.ValidatorAddress)
	if !found {
		return wasmTypes.FullDelegation{}, sdkerrors.Wrap(staking.ErrNoValidatorFound, "can't load validator for delegation")
	}

	bondDenom := querier.keeper.BondDenom(ctx)
	amount := sdk.NewCoin(bondDenom, val.TokensFromShares(del.Shares).TruncateInt())

	return wasmTypes.FullDelegation{
		Delegator: del.DelegatorAddress.String(),
		Validator: del.ValidatorAddress.String(),
		Amount:    wasm.EncodeSdkCoin(amount),
		// TODO: AccumulatedRewards
		AccumulatedRewards: wasmTypes.NewCoin(0, bondDenom),
		// TODO: Determine redelegate
		CanRedelegate: wasmTypes.NewCoin(0, bondDenom),
	}, nil
}
