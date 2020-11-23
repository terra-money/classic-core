package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters do salsh any operator who over criteria & clear all operators miss counter to zero
func SlashAndResetMissCounters(ctx sdk.Context, k Keeper) {
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	votePeriodsPerWindow := sdk.NewDec(k.SlashWindow(ctx)).QuoInt64(k.VotePeriod(ctx)).TruncateInt64()
	minValidPerWindow := k.MinValidPerWindow(ctx)
	slashFraction := k.SlashFraction(ctx)
	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter int64) bool {

		// Calculate valid vote rate; (SlashWindow - MissCounter)/SlashWindow
		validVoteRate := sdk.NewDecFromInt(
			sdk.NewInt(votePeriodsPerWindow - missCounter)).
			QuoInt64(votePeriodsPerWindow)

		// Penalize the validator whose the valid vote rate is smaller than min threshold
		if validVoteRate.LT(minValidPerWindow) {
			validator := k.StakingKeeper.Validator(ctx, operator)
			if validator.IsBonded() && !validator.IsJailed() {
				k.StakingKeeper.Slash(
					ctx, validator.GetConsAddr(),
					distributionHeight, validator.GetConsensusPower(), slashFraction,
				)
				k.StakingKeeper.Jail(ctx, validator.GetConsAddr())
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
