package budget

import (
	"strconv"

	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/util"
	"github.com/terra-project/core/x/budget/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tally returns votePower = yesVotes minus NoVotes for program, as well as the total votes.
// Power is denominated in validator bonded tokens (Luna stake size)
func tally(ctx sdk.Context, k Keeper, targetProgramID uint64) (votePower sdk.Int, totalPower sdk.Int) {
	votePower = sdk.ZeroInt()
	totalPower = k.valset.TotalBondedTokens(ctx)

	voteCount := 0
	targetProgramIDPrefix := keyVote(targetProgramID, sdk.AccAddress{})
	k.IterateVotesWithPrefix(ctx, targetProgramIDPrefix, func(programID uint64, voter sdk.AccAddress, option bool) (stop bool) {
		voteCount++
		valAddr := sdk.ValAddress(voter)

		if validator := k.valset.Validator(ctx, valAddr); validator != nil {
			bondSize := validator.GetBondedTokens()
			if option {
				votePower = votePower.Add(bondSize)
			} else {
				votePower = votePower.Sub(bondSize)
			}
		} else {
			k.DeleteVote(ctx, targetProgramID, voter)
		}

		return false
	})

	return
}

// clearsThreshold returns true if totalPower * threshold < votePower
func clearsThreshold(votePower, totalPower sdk.Int, threshold sdk.Dec) bool {
	return votePower.GTE(threshold.MulInt(totalPower).RoundInt())
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (claims types.ClaimPool, resTags sdk.Tags) {
	params := k.GetParams(ctx)
	claims = types.ClaimPool{}
	resTags = sdk.EmptyTags()

	k.CandQueueIterateExpired(ctx, ctx.BlockHeight(), func(programID uint64) (stop bool) {
		program, err := k.GetProgram(ctx, programID)
		if err != nil {
			return false
		}

		// Did not pass the tally, delete program
		votePower, totalPower := tally(ctx, k, programID)

		if !clearsThreshold(votePower, totalPower, params.ActiveThreshold) {
			k.DeleteProgram(ctx, programID)
			k.DeleteVote(ctx, programID, sdk.AccAddress{})
			resTags.AppendTag(tags.Action, tags.ActionProgramRejected)
		} else {
			resTags.AppendTag(tags.Action, tags.ActionProgramPassed)
		}

		resTags.AppendTags(
			sdk.NewTags(
				tags.ProgramID, strconv.FormatUint(programID, 10),
				tags.Weight, votePower.String(),
			),
		)

		k.CandQueueRemove(ctx, program.getVotingEndBlock(ctx, k), programID)
		return false
	})

	// Time to re-weight programs
	if !util.IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	claims = types.ClaimPool{}

	// iterate programs and weight them
	k.IteratePrograms(ctx, true, func(programID uint64, program Program) (stop bool) {
		votePower, totalPower := tally(ctx, k, programID)

		//  Need to check if the program should be legacied
		if !clearsThreshold(votePower, totalPower, params.LegacyThreshold) {
			k.DeleteProgram(ctx, programID)
			k.DeleteVote(ctx, programID, sdk.AccAddress{})
			resTags.AppendTag(tags.Action, tags.ActionProgramLegacied)
		} else {
			claims = append(claims, types.NewClaim(types.BudgetClaimClass, votePower, program.Executor))
			resTags.AppendTag(tags.Action, tags.ActionProgramGranted)
		}

		resTags.AppendTags(
			sdk.NewTags(
				tags.ProgramID, strconv.FormatUint(programID, 10),
				tags.Weight, votePower.String(),
			),
		)

		return false
	})
	return
}
