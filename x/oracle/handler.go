package oracle

import (
	"fmt"
	"terra/x/treasury"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gonum.org/v1/gonum/stat"
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

			votes := getVotesForDenom(k, ctx, denom)
			votePower := getTotalVotePower(votes)

			// Not enough validators have voted, skip
			if votePower.LT(k.valset.TotalPower(ctx).Mul(voteThreshold)) {
				newTags.AppendTag(denom, []byte("no confidence"))
				continue
			}

			// Get weighted median prices, and faithful respondants
			targetMode, observedMode, rewardees := tallyVotes(votes)

			// Clear stale votes
			clearVotesForDenom(k, ctx, denom)

			// Set the Target and Observed prices for the asset
			k.setPriceTarget(ctx, denom, targetMode)
			k.setPriceObserved(ctx, denom, observedMode)

			// Pay out rewardees
			// TODO: handle cases where the reward is too small
			rewardeePower := getTotalVotePower(rewardees)
			for _, recipient := range rewardees {

				k.tk.AddClaim(ctx, treasury.NewBaseClaim(
					treasury.OracleShareID,
					recipient.Power.Quo(rewardeePower),
					recipient.FeedMsg.Feeder,
				),
				)
			}

			newTags.AppendTag(denom, []byte(fmt.Sprintf("target %v observed %v rewardees %v",
				targetMode, observedMode, rewardees)))
		}
	}

	return resTags
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
	priceVote := NewPriceVote(pfm, val.GetPower())
	keeper.addVote(ctx, priceVote)

	return sdk.Result{}
}

//------------------------------------------
// Util functions

func getVotesForDenom(k Keeper, ctx sdk.Context, denom string) (votes []PriceVote) {
	voteIterator := k.getVoteIterator(ctx, denom)
	for ; voteIterator.Valid(); voteIterator.Next() {
		var vote PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(voteIterator.Value(), &vote)

		votes = append(votes, vote)
	}
	voteIterator.Close()
	return votes
}

func clearVotesForDenom(k Keeper, ctx sdk.Context, denom string) {
	voteIterator := k.getVoteIterator(ctx, denom)
	for ; voteIterator.Valid(); voteIterator.Next() {
		var vote PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(voteIterator.Value(), &vote)

		k.deleteVote(ctx, vote)
	}
	voteIterator.Close()
}

func getTotalVotePower(votes []PriceVote) sdk.Dec {
	votePower := sdk.ZeroDec()
	for _, vote := range votes {
		votePower.Add(vote.Power)
	}

	return votePower
}

func decToFloat64(a sdk.Dec) float64 {
	// roundup
	b := a.MulInt(sdk.NewInt(10 ^ OracleDecPrec))
	c := b.TruncateInt64()

	return float64(c) / (10 ^ OracleDecPrec)
}

func float64ToDec(a float64) sdk.Dec {
	b := int64(a * (10 ^ OracleDecPrec))
	return sdk.NewDecWithPrec(b, 2)
}

func tallyVotes(votes []PriceVote) (targetMode sdk.Dec, observedMode sdk.Dec, rewardees []PriceVote) {
	vTarget := make([]float64, len(votes))
	vPower := make([]float64, len(votes))
	vObserved := make([]float64, len(votes))

	for _, vote := range votes {
		vPower = append(vPower, decToFloat64(vote.Power))
		vTarget = append(vTarget, decToFloat64(vote.FeedMsg.TargetPrice))
		vObserved = append(vObserved, decToFloat64(vote.FeedMsg.ObservedPrice))
	}

	tmode, _ := stat.Mode(vTarget, vPower)
	omode, _ := stat.Mode(vObserved, vPower)

	tsd := stat.StdDev(vTarget, vPower)
	osd := stat.StdDev(vTarget, vPower)

	for i, vote := range votes {
		if vTarget[i] >= tmode-tsd && vTarget[i] <= tmode+tsd &&
			vObserved[i] >= omode-osd && vObserved[i] <= omode+osd {
			rewardees = append(rewardees, vote)
		}
	}

	targetMode = float64ToDec(tmode)
	observedMode = float64ToDec(omode)
	return
}
