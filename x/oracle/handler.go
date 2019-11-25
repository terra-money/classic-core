package oracle

import (
	"bytes"
	"encoding/hex"
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
		default:
			errMsg := fmt.Sprintf("Unrecognized oracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgExchangeRatePrevote handles a MsgExchangeRatePrevote
func handleMsgExchangeRatePrevote(ctx sdk.Context, keeper Keeper, ppm MsgExchangeRatePrevote) sdk.Result {
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

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
		return ErrNotRevealPeriod(keeper.Codespace()).Result()
	}

	// If there is an prevote, we verify a exchange rate with prevote hash and move prevote to vote with given exchange rate
	bz, _ := hex.DecodeString(prevote.Hash) // prevote hash
	bz2, err2 := VoteHash(pvm.Salt, pvm.ExchangeRate, prevote.Denom, prevote.Voter)
	if err2 != nil {
		return ErrVerificationFailed(keeper.Codespace(), bz, []byte{}).Result()
	}

	if !bytes.Equal(bz, bz2) {
		return ErrVerificationFailed(keeper.Codespace(), bz, bz2).Result()
	}

	// Add the vote to the store
	vote := NewExchangeRateVote(pvm.ExchangeRate, prevote.Denom, prevote.Voter)
	keeper.DeleteExchangeRatePrevote(ctx, prevote)
	keeper.AddExchangeRateVote(ctx, vote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVote,
			sdk.NewAttribute(types.AttributeKeyDenom, pvm.Denom),
			sdk.NewAttribute(types.AttributeKeyVoter, pvm.Validator.String()),
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
			types.EventTypeFeedDeleate,
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
