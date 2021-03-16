package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

		sdkMsg := stakingtypes.NewMsgDelegate(
			contractAddr,
			validator,
			coin,
		)

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

		sdkMsg := stakingtypes.NewMsgBeginRedelegate(
			contractAddr,
			src,
			dst,
			coin,
		)

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

		sdkMsg := stakingtypes.NewMsgUndelegate(
			contractAddr,
			validator,
			coin,
		)

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

		sdkMsg := distrtypes.NewMsgSetWithdrawAddress(
			contractAddr,
			rcpt,
		)

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

		sdkMsg := distrtypes.NewMsgWithdrawDelegatorReward(
			contractAddr,
			validator,
		)

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
	stakingKeeper stakingkeeper.Keeper
	distrKeeper   distrkeeper.Keeper
}

// NewWasmQuerier returns staking wasm querier
func NewWasmQuerier(stakingKeeper stakingkeeper.Keeper, distrKeeper distrkeeper.Keeper) WasmQuerier {
	return WasmQuerier{stakingKeeper, distrKeeper}
}

// Query - implement query function
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error) {
	if request.Staking.BondedDenom != nil {
		res := wasmTypes.BondedDenomResponse{
			Denom: querier.stakingKeeper.BondDenom(ctx),
		}

		return json.Marshal(res)
	}

	if request.Staking.Validators != nil {
		validators := querier.stakingKeeper.GetBondedValidatorsByPower(ctx)
		wasmVals := make([]wasmTypes.Validator, len(validators))

		for i, v := range validators {
			wasmVals[i] = wasmTypes.Validator{
				Address:       v.OperatorAddress,
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

		delegations := querier.stakingKeeper.GetAllDelegatorDelegations(ctx, delegator)

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

		var responseFullDelegation *wasmTypes.FullDelegation
		delegation, found := querier.stakingKeeper.GetDelegation(ctx, delegator, validator)
		if found {
			responseFullDelegation, err = querier.encodeDelegation(ctx, delegation)
			if err != nil {
				return nil, err
			}
		}

		res := wasmTypes.DelegationResponse{
			Delegation: responseFullDelegation,
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
func (querier WasmQuerier) encodeDelegations(ctx sdk.Context, delegations stakingtypes.Delegations) (wasmTypes.Delegations, error) {
	bondDenom := querier.stakingKeeper.BondDenom(ctx)

	var responseDelegations wasmTypes.Delegations
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
		}

		val, found := querier.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, "can't load validator for delegation")
		}

		amount := sdk.NewCoin(bondDenom, val.TokensFromShares(del.Shares).TruncateInt())

		responseDelegations = append(responseDelegations, wasmTypes.Delegation{
			Delegator: del.DelegatorAddress,
			Validator: del.ValidatorAddress,
			Amount:    wasm.EncodeSdkCoin(amount),
		})
	}
	return responseDelegations, nil
}

// encode cosmos staking to wasm delegation
func (querier WasmQuerier) encodeDelegation(ctx sdk.Context, del stakingtypes.Delegation) (*wasmTypes.FullDelegation, error) {
	valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	val, found := querier.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, "can't load validator for delegation")
	}

	bondDenom := querier.stakingKeeper.BondDenom(ctx)
	amount := sdk.NewCoin(bondDenom, val.TokensFromShares(del.Shares).TruncateInt())

	// TODO pass reward to
	_, err = querier.getAccumulatedRewards(ctx, del)
	if err != nil {
		return nil, err
	}

	return &wasmTypes.FullDelegation{
		Delegator: del.DelegatorAddress,
		Validator: del.ValidatorAddress,
		Amount:    wasm.EncodeSdkCoin(amount),
		// TODO: AccumulatedRewards
		AccumulatedRewards: wasmTypes.NewCoin(0, bondDenom),
		// TODO: Determine redelegate
		CanRedelegate: wasmTypes.NewCoin(0, bondDenom),
	}, nil
}

func (querier WasmQuerier) getAccumulatedRewards(ctx sdk.Context, delegation stakingtypes.Delegation) ([]wasmTypes.Coin, error) {
	// Try to get *delegator* reward info!
	params := distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delegation.DelegatorAddress,
		ValidatorAddress: delegation.ValidatorAddress,
	}
	cache, _ := ctx.CacheContext()
	res, err := querier.distrKeeper.DelegationRewards(sdk.WrapSDKContext(cache), &params)
	if err != nil {
		return nil, err
	}

	// now we have it, convert it into wasm types
	rewards := make([]wasmTypes.Coin, len(res.Rewards))
	for i, r := range res.Rewards {
		rewards[i] = wasmTypes.Coin{
			Denom:  r.Denom,
			Amount: r.Amount.TruncateInt().String(),
		}
	}
	return rewards, nil
}
