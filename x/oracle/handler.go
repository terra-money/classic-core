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
		case MsgPriceFeed:
			return handleMsgPriceFeed(ctx, k, msg)
		case MsgDelegateFeederPermission:
			return handleMsgDelegateFeederPermission(ctx, k, msg)
		default:
			errMsg := "Unrecognized oracle Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgPriceFeed handles a MsgPriceFeed
func handleMsgPriceFeed(ctx sdk.Context, keeper Keeper, pfm MsgPriceFeed) sdk.Result {
	valset := keeper.valset

	if !pfm.Feeder.Equals(pfm.Validator) {
		delegate := keeper.GetFeedDelegate(ctx, pfm.Validator)
		if !delegate.Equals(pfm.Feeder) {
			return ErrNoVotingPermission(DefaultCodespace, pfm.Feeder, pfm.Validator).Result()
		}
	}

	// Check that the given validator exists
	val := valset.Validator(ctx, pfm.Validator)
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	params := keeper.GetParams(ctx)

	// For zero price, it will just replace prevote without checking the price and submitting vote. It is useful to change price before vote period
	if !pfm.Price.Equal(sdk.ZeroDec()) {

		// Get prevote
		if prevote, err := keeper.getPrevote(ctx, pfm.Denom, pfm.Validator); err == nil {

			// Check a msg is submitted porper period
			if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
				return ErrNotRevealPeriod(DefaultCodespace).Result()
			}

			// If there is an prevote, we verify a price with prevote hash and move prevote to vote with given price
			bz, _ := hex.DecodeString(prevote.Hash) // prevote hash
			bz2, err := VoteHash(pfm.Salt, pfm.Price, prevote.Denom, prevote.Voter)
			if err != nil {
				return ErrVerificationFailed(DefaultCodespace, bz, []byte{}).Result()
			}

			if !bytes.Equal(bz, bz2) {
				return ErrVerificationFailed(DefaultCodespace, bz, bz2).Result()
			}

			// Add the vote to the store
			vote := NewPriceVote(pfm.Price, prevote.Denom, prevote.Voter)
			keeper.deletePrevote(ctx, prevote)
			keeper.addVote(ctx, vote)
		}

	}

	// Add the prevote to the store
	if len(pfm.Hash) != 0 {
		prevote := NewPricePrevote(pfm.Hash, pfm.Denom, pfm.Validator, ctx.BlockHeight())
		keeper.addPrevote(ctx, prevote)
	}

	log := NewLog()
	log = log.append(LogKeyPrice, pfm.Price.String())

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, pfm.Denom,
			tags.Voter, pfm.Validator.String(),
			tags.FeedDelegate, pfm.Feeder.String(),
		),
		Log: log.String(),
	}
}

// handleMsgPriceFeed handles a MsgPriceFeed
func handleMsgDelegateFeederPermission(ctx sdk.Context, keeper Keeper, pfm MsgDelegateFeederPermission) sdk.Result {
	valset := keeper.valset
	signer := pfm.Operator

	// Check the delegator is a validator
	val := valset.Validator(ctx, signer)
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	// Set the delegation
	keeper.SetFeedDelegate(ctx, signer, pfm.FeedDelegate)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Operator, pfm.Operator.String(),
			tags.FeedDelegate, pfm.FeedDelegate.String(),
		),
	}
}
