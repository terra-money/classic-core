package oracle

import (
	"bytes"
	"encoding/hex"

	"github.com/terra-project/core/x/oracle/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
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
			errMsg := "Unrecognized oracle Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgPricePrevote handles a MsgPricePrevote
func handleMsgPricePrevote(ctx sdk.Context, keeper Keeper, ppm MsgPricePrevote) sdk.Result {
	valset := keeper.valset

	if !ppm.Feeder.Equals(ppm.Validator) {
		delegate := keeper.GetFeedDelegate(ctx, ppm.Validator)
		if !delegate.Equals(ppm.Feeder) {
			return ErrNoVotingPermission(DefaultCodespace, ppm.Feeder, ppm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := valset.Validator(ctx, ppm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	prevote := NewPricePrevote(ppm.Hash, ppm.Denom, ppm.Validator, ctx.BlockHeight())
	keeper.addPrevote(ctx, prevote)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, ppm.Denom,
			tags.Voter, ppm.Validator.String(),
			tags.FeedDelegate, ppm.Feeder.String(),
		),
	}
}

// handleMsgPriceVote handles a MsgPriceVote
func handleMsgPriceVote(ctx sdk.Context, keeper Keeper, pvm MsgPriceVote) sdk.Result {
	valset := keeper.valset

	if !pvm.Feeder.Equals(pvm.Validator) {
		delegate := keeper.GetFeedDelegate(ctx, pvm.Validator)
		if !delegate.Equals(pvm.Feeder) {
			return ErrNoVotingPermission(DefaultCodespace, pvm.Feeder, pvm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := valset.Validator(ctx, pvm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	params := keeper.GetParams(ctx)

	// Get prevote
	prevote, err := keeper.getPrevote(ctx, pvm.Denom, pvm.Validator)
	if err != nil {
		return ErrNoPrevote(DefaultCodespace, pvm.Validator, pvm.Denom).Result()
	}

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
		return ErrNotRevealPeriod(DefaultCodespace).Result()
	}

	// If there is an prevote, we verify a price with prevote hash and move prevote to vote with given price
	bz, _ := hex.DecodeString(prevote.Hash) // prevote hash
	bz2, err2 := VoteHash(pvm.Salt, pvm.Price, prevote.Denom, prevote.Voter)
	if err2 != nil {
		return ErrVerificationFailed(DefaultCodespace, bz, []byte{}).Result()
	}

	if !bytes.Equal(bz, bz2) {
		return ErrVerificationFailed(DefaultCodespace, bz, bz2).Result()
	}

	// Add the vote to the store
	vote := NewPriceVote(pvm.Price, prevote.Denom, prevote.Voter)
	keeper.deletePrevote(ctx, prevote)
	keeper.addVote(ctx, vote)

	log := NewLog()
	log = log.append(LogKeyPrice, pvm.Price.String())

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, pvm.Denom,
			tags.Voter, pvm.Validator.String(),
			tags.FeedDelegate, pvm.Feeder.String(),
		),
		Log: log.String(),
	}
}

// handleMsgDelegateFeederPermission handles a MsgDelegateFeederPermission
func handleMsgDelegateFeederPermission(ctx sdk.Context, keeper Keeper, dfpm MsgDelegateFeederPermission) sdk.Result {
	valset := keeper.valset
	signer := dfpm.Operator

	// Check the delegator is a validator
	val := valset.Validator(ctx, signer)
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	// Set the delegation
	keeper.SetFeedDelegate(ctx, signer, dfpm.FeedDelegate)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Operator, dfpm.Operator.String(),
			tags.FeedDelegate, dfpm.FeedDelegate.String(),
		),
	}
}
