package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	wasm "github.com/terra-project/core/x/wasm/exported"
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
func (parser WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
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

		return []sdk.Msg{sdkMsg}, nil
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

		return []sdk.Msg{sdkMsg}, nil
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

		return []sdk.Msg{sdkMsg}, nil
	}

	if msg.Withdraw != nil {
		var err error
		rcpt := contractAddr

		if len(msg.Withdraw.Recipient) != 0 {
			rcpt, err = sdk.AccAddressFromBech32(msg.Withdraw.Recipient)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Recipient)
			}
		}

		validator, err := sdk.ValAddressFromBech32(msg.Withdraw.Validator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Withdraw.Validator)
		}

		setMsg := distribution.MsgSetWithdrawAddress{
			DelegatorAddress: contractAddr,
			WithdrawAddress:  rcpt,
		}

		withdrawMsg := distribution.MsgWithdrawDelegatorReward{
			DelegatorAddress: contractAddr,
			ValidatorAddress: validator,
		}

		return []sdk.Msg{setMsg, withdrawMsg}, nil
	}
	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Staking")
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

	if request.Staking.Delegations != nil {
		delegator, err := sdk.AccAddressFromBech32(request.Staking.Delegations.Delegator)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Staking.Delegations.Delegator)
		}

		var delegations []staking.Delegation
		if len(request.Staking.Delegations.Validator) == 0 {
			delegations = querier.keeper.GetAllDelegatorDelegations(ctx, delegator)
		} else {
			var validator sdk.ValAddress
			validator, err = sdk.ValAddressFromBech32(request.Staking.Delegations.Validator)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Staking.Delegations.Validator)
			}

			delegation, found := querier.keeper.GetDelegation(ctx, delegator, validator)
			if found {
				delegations = append(delegations, delegation)
			}
		}

		var responseDelegations wasmTypes.Delegations
		bondDenom := querier.keeper.BondDenom(ctx)
		for _, del := range delegations {
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

		res := wasmTypes.DelegationsResponse{
			Delegations: responseDelegations,
		}

		return json.Marshal(res)
	}

	return nil, wasmTypes.UnsupportedRequest{Kind: "unknown Staking variant"}
}

// QueryCustom implements custom query interface
func (WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	return nil, nil
}
