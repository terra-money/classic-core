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
func EndBlocker(ctx sdk.Context, k Keeper) (prices map[string]sdk.Dec, rewardees map[string]sdk.Int, resTags sdk.Tags) {

	params := k.GetParams(ctx)
	votes := k.getVotes(ctx)

	// Tally vote for oracle prices
	if sdk.NewInt(ctx.BlockHeight()).Mod(params.VotePeriod).Equal(sdk.ZeroInt()) {
		for denom, filVotes := range votes {
			votePower := filVotes.totalPower()

			// Not enough validators have voted, skip
			thresholdVotes := params.VoteThreshold.MulInt(k.valset.TotalBondedTokens(ctx)).TruncateInt()
			if votePower.LT(thresholdVotes) {

				resTags = resTags.AppendTags(
					sdk.NewTags(
						tags.Action, tags.ActionTallyDropped,
						tags.Denom, []byte(denom),
					),
				)

				dropCounter := k.getDropCounter(ctx, denom)
				if dropCounter.GT(params.DropThreshold) {

					// Too many drops, blacklist currency
					k.deletePrice(ctx, denom)
					k.deleteDropCounter(ctx, denom)

					resTags = resTags.AppendTags(
						sdk.NewTags(
							tags.Action, tags.ActionBlacklist,
							tags.Denom, []byte(denom),
						),
					)

				} else {
					dropCounter = dropCounter.Add(sdk.OneInt())
					k.setDropCounter(ctx, denom, dropCounter)
				}

				continue
			}

			// Get weighted median prices, and faithful respondants
			mod, loyalVotes := filVotes.tally()

			prices[denom] = mod
			for _, lv := range loyalVotes {
				voterAddrStr := lv.Voter.String()
				rewardees[voterAddrStr] = rewardees[voterAddrStr].Add(lv.Power)
			}

			// Emit whitelist tag if the price is coming in for the first time
			_, err := k.GetPrice(ctx, denom)
			if err != nil {
				resTags = resTags.AppendTags(
					sdk.NewTags(
						tags.Action, tags.ActionWhitelist,
						tags.Denom, []byte(denom),
					),
				)
			}

			// Set the price for the asset
			k.setPrice(ctx, denom, mod)

			// Emit price update tag
			resTags = resTags.AppendTags(
				sdk.NewTags(
					tags.Action, tags.ActionPriceUpdate,
					tags.Denom, []byte(denom),
					tags.Price, mod.Bytes(),
				),
			)
		}

		// Clear all votes
		k.iterateVotes(ctx, func(vote PriceVote) (stop bool) { k.deleteVote(ctx, vote); return false })
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

	// Add the vote to the store
	vote := NewPriceVote(pfm.Price, pfm.Denom, val.GetBondedTokens(), signer)
	keeper.addVote(ctx, vote)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionVoteSubmitted,
			tags.Denom, []byte(pfm.Denom),
			tags.Voter, pfm.Feeder.Bytes(),
			tags.Power, val.GetBondedTokens(),
			tags.Price, pfm.Price.Bytes(),
		),
	}
}
