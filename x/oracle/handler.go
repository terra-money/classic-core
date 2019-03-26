package oracle

import (
	"terra/types"
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

// Get all active oracle asset denoms from the store
func getActiveDenoms(ctx sdk.Context, k Keeper) (denoms []string) {
	denoms = []string{}

	store := ctx.KVStore(k.key)
	iter := sdk.KVStorePrefixIterator(store, prefixPrice)
	for ; iter.Valid(); iter.Next() {
		n := len(prefixPrice) + 1
		denom := string(iter.Key()[n:])
		denoms = append(denoms, denom)
	}
	iter.Close()

	return
}

// Drop the ballot. If the ballot drops params.DropThreshold times sequentially, then blacklist
func dropBallot(ctx sdk.Context, k Keeper, denom string, params Params) sdk.Tags {
	actionTag := tags.ActionTallyDropped

	// Not enough votes received
	dropCounter := k.incrementDropCounter(ctx, denom)
	if dropCounter.GTE(params.DropThreshold) {

		// Too many drops, blacklist currency
		k.deletePrice(ctx, denom)
		k.resetDropCounter(ctx, denom)

		actionTag = tags.ActionBlacklist
	}

	return sdk.NewTags(
		tags.Action, actionTag,
		tags.Denom, []byte(denom),
	)
}

// ballot for the asset is passing the threshold amount of voting power
func ballotIsPassing(totalPower sdk.Int, ballot PriceBallot, params Params) bool {
	thresholdVotes := params.VoteThreshold.MulInt(totalPower).RoundInt()
	return ballot.TotalPower().GTE(thresholdVotes)
}

// at the block height for a tally
func isTimeForTally(ctx sdk.Context, params Params) bool {
	return sdk.NewInt(ctx.BlockHeight()).Mod(params.VotePeriod).Equal(sdk.ZeroInt())
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (rewardees types.ClaimPool, resTags sdk.Tags) {
	params := k.GetParams(ctx)

	if !isTimeForTally(ctx, params) {
		return
	}

	rewardees = types.ClaimPool{}
	actives := getActiveDenoms(ctx, k)
	votes := k.collectVotes(ctx)

	totalBondedTokens := k.valset.TotalBondedTokens(ctx)

	// Iterate through active oracle assets and drop assets that have no votes received.
	for _, activeDenom := range actives {
		if _, found := votes[activeDenom]; !found {
			dropBallot(ctx, k, activeDenom, params)
		}
	}

	// Iterate through votes and update prices; drop if not enough votes have been achieved.
	for denom, filteredVotes := range votes {
		if ballotIsPassing(totalBondedTokens, filteredVotes, params) {
			// Get weighted median prices, and faithful respondants
			mod, ballotWinners := filteredVotes.tally()

			// Append ballot winners for the denom
			rewardees = append(rewardees, ballotWinners...)

			actionTag := tags.ActionPriceUpdate
			if _, err := k.GetLunaSwapRate(ctx, denom); err != nil {
				actionTag = tags.ActionWhitelist
			}

			// Set price to the store
			k.SetLunaSwapRate(ctx, denom, mod)

			// Reset drop counter for the passed ballot
			k.resetDropCounter(ctx, denom)

			resTags = resTags.AppendTags(
				sdk.NewTags(
					tags.Action, actionTag,
					tags.Denom, []byte(denom),
					tags.Price, mod.Bytes(),
				),
			)
		} else {
			dropBallot(ctx, k, denom, params)
		}

		// Clear all votes
		k.iterateVotes(ctx, func(vote PriceVote) (stop bool) { k.deleteVote(ctx, vote); return false })
	}

	// Sort rewardees before we return
	rewardees.Sort()

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
