package oracle

import (
	"terra/x/oracle/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (rewardablePower sdk.Int, rewardees map[string]sdk.Int, resTags sdk.Tags) {
	params := k.GetParams(ctx)
	votes := k.collectVotes(ctx)

	rewardablePower = sdk.ZeroInt()
	rewardees = map[string]sdk.Int{}

	// Tally vote for oracle prices
	if sdk.NewInt(ctx.BlockHeight()).Mod(params.VotePeriod).Equal(sdk.ZeroInt()) {
		for denom, filteredVotes := range votes {

			var actionTagForDenom string
			votePower := filteredVotes.TotalPower()
			rewardablePower = rewardablePower.Add(votePower)

			thresholdVotes := params.VoteThreshold.MulInt(k.valset.TotalBondedTokens(ctx)).TruncateInt()

			// Cleared the vote
			if votePower.GTE(thresholdVotes) {
				// Get weighted median prices, and faithful respondants
				mod, rewardableVotes := filteredVotes.tally()

				for _, rewardableVote := range rewardableVotes {
					voterStr := rewardableVote.Voter.String()
					if val, ok := rewardees[voterStr]; ok {
						rewardees[voterStr] = val.Add(rewardableVote.Power)
					} else {
						rewardees[voterStr] = rewardableVote.Power
					}
				}

				// Emit whitelist tag if the price is coming in for the first time
				_, err := k.GetPrice(ctx, denom)
				if err != nil {
					actionTagForDenom = tags.ActionWhitelist
				} else {
					actionTagForDenom = tags.ActionPriceUpdate
				}

				// Set the price for the asset
				k.SetPrice(ctx, denom, mod)

				// Emit price update tag
				resTags = resTags.AppendTags(
					sdk.NewTags(
						tags.Action, actionTagForDenom,
						tags.Denom, []byte(denom),
						tags.Price, mod.Bytes(),
					),
				)
			} else {
				// Not enough votes received
				dropCounter := k.incrementDropCounter(ctx, denom)
				if dropCounter.GT(params.DropThreshold) {

					// Too many drops, blacklist currency
					k.deletePrice(ctx, denom)
					k.resetDropCounter(ctx, denom)

					actionTagForDenom = tags.ActionBlacklist
				} else {
					actionTagForDenom = tags.ActionTallyDropped
				}

				resTags = resTags.AppendTags(
					sdk.NewTags(
						tags.Action, actionTagForDenom,
						tags.Denom, []byte(denom),
					),
				)
			}
		}

		// Clear all votes
		k.iterateVotes(ctx, func(vote PriceVote) (stop bool) { k.deleteVote(ctx, vote); return false })
	}

	return
}

// handleMsgPriceFeed handles a MsgPriceFeed
func handleMsgPriceFeed(ctx sdk.Context, keeper Keeper, pfm MsgPriceFeed) sdk.Result {
	valset := keeper.valset
	signer := pfm.Feeder

	// Check the feeder is a validator
	val := valset.Validator(ctx, sdk.ValAddress(signer.Bytes()))
	if val == nil {
		return ErrNotValidator(DefaultCodespace, pfm.Feeder).Result()
	}

	// Add the vote to the store
	vote := NewPriceVote(pfm.Price, pfm.Denom, val.GetBondedTokens(), signer)
	keeper.addVote(ctx, vote)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Denom, pfm.Denom,
			tags.Voter, pfm.Feeder.Bytes(),
			tags.Power, val.GetBondedTokens().String(),
			tags.Price, pfm.Price.Bytes(),
		),
	}
}
