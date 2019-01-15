package oracle

import (
	"terra/x/oracle/tags"

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
func EndBlocker(ctx sdk.Context, k Keeper) (targetPrices map[string]sdk.Dec,
	observedPrices map[string]sdk.Dec,
	rewardees []PriceVote,
	resTags sdk.Tags) {

	targetPrices = map[string]sdk.Dec{}
	observedPrices = map[string]sdk.Dec{}
	votePeriod := k.GetParams(ctx).VotePeriod
	voteThreshold := k.GetParams(ctx).VoteThreshold
	whitelist := k.GetParams(ctx).Whitelist

	// Tally vote for oracle prices
	if sdk.NewInt(ctx.BlockHeight()).Mod(votePeriod).Equal(sdk.ZeroInt()) {
		resTags = resTags.AppendTag(sdk.TagAction, tags.ActionPriceUpdate)
		for _, denom := range whitelist {

			targetVotes := k.getTargetVotes(ctx, denom)
			observedVotes := k.getObservedVotes(ctx, denom)
			votePower := targetVotes.totalPower() // should be same for observed

			// Not enough validators have voted, skip
			if votePower.LT(k.valset.TotalPower(ctx).Mul(voteThreshold)) {
				resTags = resTags.AppendTag(denom, tags.ActionTallyDropped)
				continue
			}

			// Get weighted median prices, and faithful respondants
			targetMod, tRewardees := targetVotes.tally()
			observedMod, oRewardees := observedVotes.tally()

			targetPrices[denom] = targetMod
			observedPrices[denom] = observedMod

			rewardees = append(rewardees, tRewardees...)
			rewardees = append(rewardees, oRewardees...)

			// Clear all votes
			k.iterateTargetVotes(ctx, denom, func(vote PriceVote) (stop bool) { k.deleteTargetVote(ctx, vote); return false })
			k.iterateObservedVotes(ctx, denom, func(vote PriceVote) (stop bool) { k.deleteObservedVote(ctx, vote); return false })

			// Set the Target and Observed prices for the asset
			k.setPriceTarget(ctx, denom, targetMod)
			k.setPriceObserved(ctx, denom, observedMod)

			resTags = resTags.AppendTags(
				sdk.NewTags(
					sdk.TagAction, tags.ActionPriceUpdate,
					tags.Denom, []byte(denom),
					tags.TargetPrice, targetMod.Bytes(),
					tags.ObservedPrice, observedMod.Bytes(),
				),
			)
		}
	}

	return
}

// handlePriceFeedMsg is used by other modules to handle Msg
func handlePriceFeedMsg(ctx sdk.Context, keeper Keeper, pfm PriceFeedMsg) sdk.Result {
	valset := keeper.valset
	signer := pfm.Feeder

	// Check the feeder is a validater
	val := valset.Validator(ctx, sdk.ValAddress(signer.Bytes()))
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
	targetVote := NewPriceVote(pfm.TargetPrice, pfm.Denom, val.GetPower(), signer)
	keeper.addTargetVote(ctx, targetVote)

	observedVote := NewPriceVote(pfm.ObservedPrice, pfm.Denom, val.GetPower(), signer)
	keeper.addObservedVote(ctx, observedVote)

	return sdk.Result{
		Tags: sdk.NewTags(
			sdk.TagAction, tags.ActionVoteSubmitted,
			tags.Denom, []byte(pfm.Denom),
			tags.Voter, pfm.Feeder.Bytes(),
			tags.Power, val.GetPower().Bytes(),
			tags.TargetPrice, pfm.TargetPrice.Bytes(),
			tags.ObservedPrice, pfm.ObservedPrice.Bytes(),
		),
	}
}
