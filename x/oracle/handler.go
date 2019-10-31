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
		switch msg := msg.(type) {
		case MsgPricePrevote:
			return handleMsgPricePrevote(ctx, k, msg)
		case MsgPriceVote:
			return handleMsgPriceVote(ctx, k, msg)
		case MsgDelegateFeederPermission:
			return handleMsgDelegateFeederPermission(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized oracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgPricePrevote handles a MsgPricePrevote
func handleMsgPricePrevote(ctx sdk.Context, keeper Keeper, ppm MsgPricePrevote) sdk.Result {
	if !ppm.Feeder.Equals(ppm.Validator) {
		delegate := keeper.GetFeedDelegate(ctx, ppm.Validator)
		if !delegate.Equals(ppm.Feeder) {
			return ErrNoVotingPermission(keeper.Codespace(), ppm.Feeder, ppm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, ppm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	prevote := NewPricePrevote(ppm.Hash, ppm.Denom, ppm.Validator, ctx.BlockHeight())
	keeper.AddPrevote(ctx, prevote)

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

// handleMsgPriceVote handles a MsgPriceVote
func handleMsgPriceVote(ctx sdk.Context, keeper Keeper, pvm MsgPriceVote) sdk.Result {
	if !pvm.Feeder.Equals(pvm.Validator) {
		delegate := keeper.GetFeedDelegate(ctx, pvm.Validator)
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
	prevote, err := keeper.GetPrevote(ctx, pvm.Denom, pvm.Validator)
	if err != nil {
		return ErrNoPrevote(keeper.Codespace(), pvm.Validator, pvm.Denom).Result()
	}

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
		return ErrNotRevealPeriod(keeper.Codespace()).Result()
	}

	// If there is an prevote, we verify a price with prevote hash and move prevote to vote with given price
	bz, _ := hex.DecodeString(prevote.Hash) // prevote hash
	bz2, err2 := VoteHash(pvm.Salt, pvm.Price, prevote.Denom, prevote.Voter)
	if err2 != nil {
		return ErrVerificationFailed(keeper.Codespace(), bz, []byte{}).Result()
	}

	if !bytes.Equal(bz, bz2) {
		return ErrVerificationFailed(keeper.Codespace(), bz, bz2).Result()
	}

	// Add the vote to the store
	vote := NewPriceVote(pvm.Price, prevote.Denom, prevote.Voter)
	keeper.DeletePrevote(ctx, prevote)
	keeper.AddVote(ctx, vote)

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

// handleMsgDelegateFeederPermission handles a MsgDelegateFeederPermission
func handleMsgDelegateFeederPermission(ctx sdk.Context, keeper Keeper, dfpm MsgDelegateFeederPermission) sdk.Result {
	signer := dfpm.Operator

	// Check the delegator is a validator
	val := keeper.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return staking.ErrNoValidatorFound(keeper.Codespace()).Result()
	}

	// Set the delegation
	keeper.SetFeedDelegate(ctx, signer, dfpm.Delegatee)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeedDeleate,
			sdk.NewAttribute(types.AttributeKeyOperator, dfpm.Operator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, dfpm.Delegatee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
