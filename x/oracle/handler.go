package oracle

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case PriceFeedMsg:
			return handlePriceFeedMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	newTags := sdk.NewTags()

	// Update price elects
	if sdk.NewInt(ctx.BlockHeight()).Mod(k.GetVotePeriod(ctx)).Equal(sdk.ZeroInt()) {
		whitelist := k.GetWhitelist(ctx)
		newTags.AppendTag("action", []byte("price_update"))
		for _, denom := range whitelist {

			votes := k.GetAllVotes(ctx, denom)
			votePower := sdk.ZeroDec()
			for _, vote := range votes {
				votePower.Add(vote.Power)
			}

			// Not enough validators have voted, skip
			if votePower.LT(k.valset.TotalPower(ctx).Mul(k.GetThreshold(ctx))) {
				newTags.AppendTag(denom, []byte("no confidence"))
				continue
			}

			// Sort votes by price
			sort.Sort(votes)

			medPower := sdk.ZeroDec()
			median := PriceVote{}
			for i := 0; i < len(votes); i++ {
				medPower.Add(votes[i].Power)

				// Get the weighted median of the votes
				if medPower.GTE(votePower.Mul(sdk.NewDecWithPrec(5, 1))) {
					median = votes[i]
				}
			}

			k.SetElect(ctx, median)
			k.ClearVotes(ctx)

			newTags.AppendTag(denom, []byte(median.FeedMsg.CurrentPrice.String()))
		}
	}

	return newTags
}

// handlePriceFeedMsg is used by other modules to handle Msg
func handlePriceFeedMsg(ctx sdk.Context, keeper Keeper, pfm PriceFeedMsg) sdk.Result {
	valset := keeper.valset

	// Check the feeder is a validater
	val := valset.Validator(ctx, sdk.ValAddress(pfm.Feeder.Bytes()))
	if val == nil {
		return ErrNotValidator(DefaultCodespace, pfm.Feeder).Result()
	}

	// Check the vote is for a whitelisted asset
	whitelist := keeper.GetWhitelist(ctx)
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

	priceVote := NewPriceVote(pfm, val.GetPower())
	keeper.AddVote(ctx, priceVote)

	return sdk.Result{}
}
