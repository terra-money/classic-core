package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/terra-money/core/x/oracle/types"
)

// SlashAndResetMissCounters do slash any operator who over criteria & clear all operators miss counter to zero
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	// slash_window / vote_period
	votePeriodsPerWindow := uint64(
		sdk.NewDec(int64(k.SlashWindow(ctx))).
			QuoInt64(int64(k.VotePeriod(ctx))).
			TruncateInt64(),
	)
	minValidPerWindow := k.MinValidPerWindow(ctx)
	slashFraction := k.SlashFraction(ctx)
	powerReduction := k.StakingKeeper.PowerReduction(ctx)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {

		// Calculate valid vote rate; (SlashWindow - MissCounter)/SlashWindow
		validVoteRate := sdk.NewDecFromInt(
			sdk.NewInt(int64(votePeriodsPerWindow - missCounter))).
			QuoInt64(int64(votePeriodsPerWindow))

		// Penalize the validator whose the valid vote rate is smaller than min threshold
		if validVoteRate.LT(minValidPerWindow) {
			validator := k.StakingKeeper.Validator(ctx, operator)
			if validator.IsBonded() && !validator.IsJailed() {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}

				power := validator.GetConsensusPower(powerReduction)
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						slashingtypes.EventTypeSlash,
						sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
						sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", power)),
						sdk.NewAttribute(slashingtypes.AttributeKeyReason, types.AttributeValueMissingOracleVote),
						sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consAddr.String()),
					),
				)

				k.StakingKeeper.Slash(
					ctx, consAddr,
					distributionHeight, validator.GetConsensusPower(powerReduction), slashFraction,
				)
				k.StakingKeeper.Jail(ctx, consAddr)

				logger := k.Logger(ctx)
				logger.Info(
					"slashing and jailing validator due to oracle vote liveness fault",
					"height", height,
					"validator", consAddr.String(),
					"slashed", slashFraction.String(),
				)
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
