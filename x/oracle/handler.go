package oracle

import (
	"fmt"
	"terra/x/treasury"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case PriceFeedMsg:
			return handlePriceFeedMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized oracle Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	newTags := sdk.NewTags()

	votePeriod := k.GetParams(ctx).VotePeriod
	voteThreshold := k.GetParams(ctx).VoteThreshold
	whitelist := k.GetParams(ctx).Whitelist

	// Tally vote for oracle prices
	if sdk.NewInt(ctx.BlockHeight()).Mod(votePeriod).Equal(sdk.ZeroInt()) {
		newTags.AppendTag("action", []byte("price_update"))
		for _, denom := range whitelist {

			votes := k.getVotes(ctx, denom)
			votePower := getTotalVotePower(votes)

			// Not enough validators have voted, skip
			if votePower.LT(k.valset.TotalPower(ctx).Mul(voteThreshold)) {
				newTags.AppendTag(denom, []byte("no confidence"))
				continue
			}

			// Get weighted median prices, and faithful respondants
			targetMode, observedMode, rewardees := tallyVotes(votes)

			// Clear stale votes
			k.clearVotes(ctx, denom)

			// Set the Target and Observed prices for the asset
			k.setPriceTarget(ctx, denom, targetMode)
			k.setPriceObserved(ctx, denom, observedMode)

			// Pay out rewardees
			// TODO: handle cases where the reward is too small
			rewardeePower := getTotalVotePower(rewardees)
			for _, recipient := range rewardees {
				k.tk.AddClaim(ctx, treasury.Claim{
					Account: recipient.FeedMsg.Feeder,
					Weight:  recipient.Power.Quo(rewardeePower),
				})
			}

			newTags.AppendTag(denom, []byte(fmt.Sprintf("target %v observed %v rewardees %v",
				targetMode, observedMode, rewardees)))
		}
	}

	return newTags
}

// handlePriceFeedMsg is used by other modules to handle Msg
func handlePriceFeedMsg(ctx sdk.Context, keeper Keeper, pfm PriceFeedMsg) sdk.Result {
	valset := keeper.valset

	// Check the feeder is a validater
	val := valset.Validator(ctx, sdk.ValAddress(pfm.GetSigners()[0].Bytes()))
	if val == nil {
		return ErrNotValidator(DefaultCodespace, pfm.Feeder).Result()
	}

	// Check the vote is for a whitelisted asset
	whitelist := keeper.GetParams(ctx).Whitelist
	contains := false
	for _, denom := range whitelist {
		if denom == pfm.Denom {
			contains = true
			break
		}
	}
	if !contains {
		return ErrUnknownDenomination(DefaultCodespace, pfm.Denom).Result()
	}

	// Add the vote to the store
	priceVote := NewPriceVote(pfm, val.GetPower())
	keeper.addVote(ctx, priceVote)

	return sdk.Result{}
}
