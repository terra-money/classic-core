package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/terra-project/core/x/oracle/internal/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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
			errMsg := fmt.Sprintf("Unrecognized oracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgExchangeRatePrevote handles a MsgExchangeRatePrevote
func handleMsgExchangeRatePrevote(ctx sdk.Context, keeper Keeper, ppm MsgExchangeRatePrevote) sdk.Result {

	// check the denom is in the vote target
	if !keeper.IsVoteTarget(ctx, ppm.Denom) {
		return ErrUnknownDenomination(keeper.Codespace(), ppm.Denom).Result()
	}

	if !ppm.Feeder.Equals(ppm.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, ppm.Validator)
		if !delegate.Equals(ppm.Feeder) {
			return ErrNoVotingPermission(keeper.Codespace(), ppm.Feeder, ppm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, ppm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	prevote := NewExchangeRatePrevote(ppm.Hash, ppm.Denom, ppm.Validator, ctx.BlockHeight())
	keeper.AddExchangeRatePrevote(ctx, prevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypePrevote,
			sdk.NewAttribute(types.AttributeKeyDenom, ppm.Denom),
			sdk.NewAttribute(types.AttributeKeyVoter, ppm.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, ppm.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgExchangeRateVote handles a MsgExchangeRateVote
func handleMsgExchangeRateVote(ctx sdk.Context, keeper Keeper, pvm MsgExchangeRateVote) sdk.Result {
	if !pvm.Feeder.Equals(pvm.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, pvm.Validator)
		if !delegate.Equals(pvm.Feeder) {
			return ErrNoVotingPermission(keeper.Codespace(), pvm.Feeder, pvm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, pvm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	params := keeper.GetParams(ctx)

	// Get prevote
	prevote, err := keeper.GetExchangeRatePrevote(ctx, pvm.Denom, pvm.Validator)
	if err != nil {
		return ErrNoPrevote(keeper.Codespace(), pvm.Validator, pvm.Denom).Result()
	}

	// Check a msg is submitted proper period
	if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
		return ErrInvalidRevealPeriod(keeper.Codespace()).Result()
	}

	// If there is an prevote, we verify a exchange rate with prevote hash and move prevote to vote with given exchange rate
	hash := GetVoteHash(pvm.Salt, pvm.ExchangeRate, pvm.Denom, pvm.Validator)
	if !prevote.Hash.Equal(hash) {
		return ErrVerificationFailed(keeper.Codespace(), prevote.Hash, hash).Result()
	}

	// Add the vote to the store
	vote := NewExchangeRateVote(pvm.ExchangeRate, pvm.Denom, pvm.Validator)
	keeper.DeleteExchangeRatePrevote(ctx, prevote)
	keeper.AddExchangeRateVote(ctx, vote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVote,
			sdk.NewAttribute(types.AttributeKeyDenom, pvm.Denom),
			sdk.NewAttribute(types.AttributeKeyVoter, pvm.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyExchangeRate, pvm.ExchangeRate.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, pvm.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgDelegateFeedConsent handles a MsgDelegateFeedConsent
func handleMsgDelegateFeedConsent(ctx sdk.Context, keeper Keeper, dfpm MsgDelegateFeedConsent) sdk.Result {
	signer := dfpm.Operator

	// Check the delegator is a validator
	val := keeper.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	// Set the delegation
	keeper.SetOracleDelegate(ctx, signer, dfpm.Delegate)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeedDelegate,
			sdk.NewAttribute(types.AttributeKeyOperator, dfpm.Operator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, dfpm.Delegate.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgAggregateExchangeRatePrevote handles a MsgAggregateExchangeRatePrevote
func handleMsgAggregateExchangeRatePrevote(ctx sdk.Context, keeper Keeper, ppm MsgAggregateExchangeRatePrevote) sdk.Result {
	if !ppm.Feeder.Equals(ppm.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, ppm.Validator)
		if !delegate.Equals(ppm.Feeder) {
			return ErrNoVotingPermission(keeper.Codespace(), ppm.Feeder, ppm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, ppm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	aggregatePrevote := NewAggregateExchangeRatePrevote(ppm.Hash, ppm.Validator, ctx.BlockHeight())
	keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregatePrevote,
			sdk.NewAttribute(types.AttributeKeyVoter, ppm.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, ppm.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgAggregateExchangeRateVote handles a MsgAggregateExchangeRateVote
func handleMsgAggregateExchangeRateVote(ctx sdk.Context, keeper Keeper, pvm MsgAggregateExchangeRateVote) sdk.Result {
	if !pvm.Feeder.Equals(pvm.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, pvm.Validator)
		if !delegate.Equals(pvm.Feeder) {
			return ErrNoVotingPermission(keeper.Codespace(), pvm.Feeder, pvm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, pvm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	params := keeper.GetParams(ctx)

	aggregatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, pvm.Validator)
	if err != nil {
		return ErrNoAggregatePrevote(keeper.Codespace(), pvm.Validator).Result()
	}

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return ErrInvalidRevealPeriod(keeper.Codespace()).Result()
	}

	exchangeRateTuples, err2 := types.ParseExchangeRateTuples(pvm.ExchangeRates)
	if err2 != nil {
		return sdk.ErrInvalidCoins(err2.Error()).Result()
	}

	// check all denoms are in the vote target
	for _, tuple := range exchangeRateTuples {
		if !keeper.IsVoteTarget(ctx, tuple.Denom) {
			return ErrUnknownDenomination(keeper.Codespace(), tuple.Denom).Result()
		}
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := GetAggregateVoteHash(pvm.Salt, pvm.ExchangeRates, aggregatePrevote.Voter)

	if !aggregatePrevote.Hash.Equal(hash) {
		return ErrVerificationFailed(keeper.Codespace(), aggregatePrevote.Hash, hash).Result()
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	keeper.AddAggregateExchangeRateVote(ctx, NewAggregateExchangeRateVote(exchangeRateTuples, aggregatePrevote.Voter))
	keeper.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregateVote,
			sdk.NewAttribute(types.AttributeKeyVoter, pvm.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyExchangeRates, pvm.ExchangeRates),
			sdk.NewAttribute(types.AttributeKeyFeeder, pvm.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
