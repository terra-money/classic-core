package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

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
func (parser WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg wasmvmtypes.CosmosMsg) (msgs sdk.Msg, err error) {
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

		cosmosMsg := stakingtypes.NewMsgDelegate(
			contractAddr,
			validator,
			coin,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
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

		cosmosMsg := stakingtypes.NewMsgBeginRedelegate(
			contractAddr,
			src,
			dst,
			coin,
		)

		return cosmosMsg, cosmosMsg.ValidateBasic()
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

		cosmosMsg := stakingtypes.NewMsgUndelegate(
			contractAddr,
			validator,
			coin,
		)

		if err := cosmosMsg.ValidateBasic(); err != nil {
			return nil, err
		}

		return cosmosMsg, nil
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Staking")
}

type CosmosQuery struct {
	Parameters *struct{} `json:"parameters,omitempty"`
}

// ParseCustom implements custom parser
func (parser WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) (sdk.Msg, error) {
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
func (querier WasmQuerier) Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error) {
	if request.Staking.BondedDenom != nil {
		res := wasmvmtypes.BondedDenomResponse{
			Denom: querier.stakingKeeper.BondDenom(ctx),
		}

		return json.Marshal(res)
	}

	if request.Staking.AllValidators != nil {
		validators := querier.stakingKeeper.GetBondedValidatorsByPower(ctx)
		wasmValidators := make([]wasmvmtypes.Validator, len(validators))

		for i, v := range validators {
			wasmValidators[i] = wasmvmtypes.Validator{
				Address:       v.OperatorAddress,
				Commission:    v.Commission.Rate.String(),
				MaxCommission: v.Commission.MaxRate.String(),
				MaxChangeRate: v.Commission.MaxChangeRate.String(),
			}
		}

		res := wasmvmtypes.AllValidatorsResponse{
			Validators: wasmValidators,
		}

		return json.Marshal(res)
	}

	if request.Staking.Validator != nil {
		validatorAddr, err := sdk.ValAddressFromBech32(request.Staking.Validator.Address)
		if err != nil {
			return nil, err
		}

		v, found := querier.stakingKeeper.GetValidator(ctx, validatorAddr)

		res := wasmvmtypes.ValidatorResponse{}
		if found {
			res.Validator = &wasmvmtypes.Validator{
				Address:       v.OperatorAddress,
				Commission:    v.Commission.Rate.String(),
				MaxCommission: v.Commission.MaxRate.String(),
				MaxChangeRate: v.Commission.MaxChangeRate.String(),
			}
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

		res := wasmvmtypes.AllDelegationsResponse{
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

		var responseFullDelegation *wasmvmtypes.FullDelegation
		delegation, found := querier.stakingKeeper.GetDelegation(ctx, delegator, validator)
		if found {
			responseFullDelegation, err = querier.encodeDelegation(ctx, delegation)
			if err != nil {
				return nil, err
			}
		}

		res := wasmvmtypes.DelegationResponse{
			Delegation: responseFullDelegation,
		}

		return json.Marshal(res)
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Staking variant"}
}

type StakingParametersResponse struct {
	// unbonding_time is the time duration of unbonding in a seconds.
	UnbondingTime uint64 `json:"unbonding_time"`
	// max_validators is the maximum number of validators.
	MaxValidators uint32 `json:"max_validators,omitempty"`
	// max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio).
	MaxEntries uint32 `json:"max_entries,omitempty"`
	// historical_entries is the number of historical entries to persist.
	HistoricalEntries uint32 `json:"historical_entries,omitempty"`
	// bond_denom defines the bondable coin denomination.
	BondDenom string `json:"bond_denom,omitempty"`
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	var bz []byte
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if query.Parameters != nil {
		parameters := querier.stakingKeeper.GetParams(ctx)
		bz, err = json.Marshal(StakingParametersResponse{
			UnbondingTime:     uint64(parameters.UnbondingTime.Seconds()),
			MaxValidators:     parameters.MaxValidators,
			MaxEntries:        parameters.MaxEntries,
			HistoricalEntries: parameters.HistoricalEntries,
			BondDenom:         parameters.BondDenom,
		})
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

// encdoe cosmos delegations to wasm delegations
func (querier WasmQuerier) encodeDelegations(ctx sdk.Context, delegations stakingtypes.Delegations) (wasmvmtypes.Delegations, error) {
	bondDenom := querier.stakingKeeper.BondDenom(ctx)

	var responseDelegations wasmvmtypes.Delegations
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

		responseDelegations = append(responseDelegations, wasmvmtypes.Delegation{
			Delegator: del.DelegatorAddress,
			Validator: del.ValidatorAddress,
			Amount:    wasm.EncodeSdkCoin(amount),
		})
	}
	return responseDelegations, nil
}

// encode cosmos staking to wasm delegation
func (querier WasmQuerier) encodeDelegation(ctx sdk.Context, del stakingtypes.Delegation) (*wasmvmtypes.FullDelegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

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
	delegationCoin := wasm.EncodeSdkCoin(amount)

	accumulatedRewards, err := querier.getAccumulatedRewards(ctx, del)
	if err != nil {
		return nil, err
	}

	// if this (val, delegate) pair is receiving a redelegation, it cannot redelegate more.
	// otherwise, it can redelegate the full amount
	redelegateCoin := wasmvmtypes.NewCoin(0, bondDenom)
	if !querier.stakingKeeper.HasReceivingRedelegation(ctx, delAddr, valAddr) {
		redelegateCoin = delegationCoin
	}

	return &wasmvmtypes.FullDelegation{
		Delegator:          del.DelegatorAddress,
		Validator:          del.ValidatorAddress,
		Amount:             delegationCoin,
		AccumulatedRewards: accumulatedRewards,
		CanRedelegate:      redelegateCoin,
	}, nil
}

func (querier WasmQuerier) getAccumulatedRewards(ctx sdk.Context, delegation stakingtypes.Delegation) (wasmvmtypes.Coins, error) {
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
	rewards := make(wasmvmtypes.Coins, len(res.Rewards))
	for i, r := range res.Rewards {
		rewards[i] = wasmvmtypes.Coin{
			Denom:  r.Denom,
			Amount: r.Amount.TruncateInt().String(),
		}
	}
	return rewards, nil
}
