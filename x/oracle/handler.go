package oracle

import (
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

	// Add the vote to the store
	vote := NewPriceVote(pfm.Price, pfm.Denom, pfm.Validator)
	keeper.addVote(ctx, vote)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, pfm.Denom,
			tags.Voter, pfm.Validator.String(),
			tags.FeedDelegate, pfm.Feeder.String(),
			tags.Price, pfm.Price.String(),
		),
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
