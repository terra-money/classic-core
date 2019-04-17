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
		default:
			errMsg := "Unrecognized oracle Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgPriceFeed handles a MsgPriceFeed
func handleMsgPriceFeed(ctx sdk.Context, keeper Keeper, pfm MsgPriceFeed) sdk.Result {
	valset := keeper.valset
	signer := pfm.Feeder

	// Check the feeder is a validator
	val := valset.Validator(ctx, sdk.ValAddress(signer.Bytes()))
	if val == nil {
		return staking.ErrNoValidatorFound(DefaultCodespace).Result()
	}

	// Add the vote to the store
	vote := NewPriceVote(pfm.Price, pfm.Denom, signer)
	keeper.addVote(ctx, vote)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, pfm.Denom,
			tags.Voter, pfm.Feeder.String(),
			tags.Price, pfm.Price.String(),
		),
	}
}
