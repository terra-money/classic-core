package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters do salsh any operator who over criteria & clear all operators miss counter to zero
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	slashWindow := k.SlashWindow(ctx)
	minValidPerWindow := k.MinValidPerWindow(ctx)
	slashFraction := k.SlashFraction(ctx)
	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter int64) bool {

		// Calculate valid vote rate; (SlashWindow - MissCounter)/SlashWindow
		validVoteRate := sdk.NewDecFromInt(sdk.NewInt(slashWindow - missCounter)).QuoInt64(k.SlashWindow(ctx))
		fmt.Println(validVoteRate, minValidPerWindow)
		if validVoteRate.LT(minValidPerWindow) {
			validator := k.StakingKeeper.Validator(ctx, operator)
			k.StakingKeeper.Slash(ctx, validator.GetConsAddr(), distributionHeight, validator.GetConsensusPower(), slashFraction)
		}

		k.SetMissCounter(ctx, operator, 0)
		return false
	})
}
