package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/internal/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgExchangeRatePrevote:
			return handleMsgExchangeRatePrevote(ctx, k, msg)
		case MsgExchangeRateVote:
			return handleMsgExchangeRateVote(ctx, k, msg)
		case MsgDelegateFeedConsent:
			return handleMsgDelegateFeedConsent(ctx, k, msg)
		case MsgAggregateExchangeRatePrevote:
			return handleMsgAggregateExchangeRatePrevote(ctx, k, msg)
		case MsgAggregateExchangeRateVote:
			return handleMsgAggregateExchangeRateVote(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// handleMsgExchangeRatePrevote handles a MsgExchangeRatePrevote

func handleMsgExchangeRatePrevote(ctx sdk.Context, keeper Keeper, msg MsgExchangeRatePrevote) (*sdk.Result, error) {

	// check the denom is in the vote target
	if !keeper.IsVoteTarget(ctx, msg.Denom) {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(ErrInternal, "unknown denom")
		}

		return nil, sdkerrors.Wrap(ErrUnknownDenom, msg.Denom)
	}

	err := keeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, !core.IsWaitingForSoftfork(ctx, 3))
	if err != nil {
		return nil, err
	}

	prevote := NewExchangeRatePrevote(msg.Hash, msg.Denom, msg.Validator, ctx.BlockHeight())
	keeper.AddExchangeRatePrevote(ctx, prevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePrevote,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgExchangeRateVote handles a MsgExchangeRateVote
func handleMsgExchangeRateVote(ctx sdk.Context, keeper Keeper, msg MsgExchangeRateVote) (*sdk.Result, error) {
	err := keeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, !core.IsWaitingForSoftfork(ctx, 3))
	if err != nil {
		return nil, err
	}

	params := keeper.GetParams(ctx)

	// Get prevote
	prevote, err := keeper.GetExchangeRatePrevote(ctx, msg.Denom, msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrNoPrevote, fmt.Sprintf("(%s, %s)", msg.Validator, msg.Denom))
	}

	// Check a msg is submitted proper period
	if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, ErrRevealPeriodMissMatch
	}

	// If there is an prevote, we verify a exchange rate with prevote hash and move prevote to vote with given exchange rate
	hash := GetVoteHash(msg.Salt, msg.ExchangeRate, msg.Denom, msg.Validator)
	if !prevote.Hash.Equal(hash) {
		return nil, sdkerrors.Wrap(ErrVerificationFailed, fmt.Sprintf("must be given %s not %s", prevote.Hash, hash))
	}

	// Add the vote to the store
	vote := NewExchangeRateVote(msg.ExchangeRate, msg.Denom, msg.Validator)
	keeper.DeleteExchangeRatePrevote(ctx, prevote)
	keeper.AddExchangeRateVote(ctx, vote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVote,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyExchangeRate, msg.ExchangeRate.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgDelegateFeedConsent handles a MsgDelegateFeedConsent
func handleMsgDelegateFeedConsent(ctx sdk.Context, keeper Keeper, msg MsgDelegateFeedConsent) (*sdk.Result, error) {
	signer := msg.Operator

	// Check the delegator is a validator
	val := keeper.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return nil, sdkerrors.Wrap(staking.ErrNoValidatorFound, signer.String())
	}

	// Set the delegation
	keeper.SetOracleDelegate(ctx, signer, msg.Delegate)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeedDelegate,
			sdk.NewAttribute(types.AttributeKeyOperator, msg.Operator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Delegate.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAggregateExchangeRatePrevote handles a MsgAggregateExchangeRatePrevote
func handleMsgAggregateExchangeRatePrevote(ctx sdk.Context, keeper Keeper, msg MsgAggregateExchangeRatePrevote) (*sdk.Result, error) {
	err := keeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, !core.IsWaitingForSoftfork(ctx, 3))
	if err != nil {
		return nil, err
	}

	aggregatePrevote := NewAggregateExchangeRatePrevote(msg.Hash, msg.Validator, ctx.BlockHeight())
	keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregatePrevote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAggregateExchangeRateVote handles a MsgAggregateExchangeRateVote
func handleMsgAggregateExchangeRateVote(ctx sdk.Context, keeper Keeper, msg MsgAggregateExchangeRateVote) (*sdk.Result, error) {
	err := keeper.ValidateFeeder(ctx, msg.Feeder, msg.Validator, !core.IsWaitingForSoftfork(ctx, 3))
	if err != nil {
		return nil, err
	}

	params := keeper.GetParams(ctx)

	aggregatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrNoAggregatePrevote, msg.Validator.String())
	}

	// Check a msg is submitted proper period
	if (ctx.BlockHeight()/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, ErrRevealPeriodMissMatch
	}

	exchangeRateTuples, err := types.ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}

	// check all denoms are in the vote target
	for _, tuple := range exchangeRateTuples {
		if !keeper.IsVoteTarget(ctx, tuple.Denom) {
			if core.IsWaitingForSoftfork(ctx, 1) {
				return nil, sdkerrors.Wrap(ErrInternal, "unknown denom")
			}

			return nil, sdkerrors.Wrap(ErrUnknownDenom, tuple.Denom)
		}
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := GetAggregateVoteHash(msg.Salt, msg.ExchangeRates, aggregatePrevote.Voter)
	if !aggregatePrevote.Hash.Equal(hash) {
		return nil, sdkerrors.Wrap(ErrVerificationFailed, fmt.Sprintf("must be given %s not %s", aggregatePrevote.Hash, hash))
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	keeper.AddAggregateExchangeRateVote(ctx, NewAggregateExchangeRateVote(exchangeRateTuples, aggregatePrevote.Voter))
	keeper.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregateVote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyExchangeRates, msg.ExchangeRates),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
