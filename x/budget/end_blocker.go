package budget

import (
	"terra/types"
	"terra/x/budget/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func tally(ctx sdk.Context, k Keeper, program Program) (votePower sdk.Int, totalPower sdk.Int) {
	currValidators := make(map[string]sdk.Int)

	// fetch all the bonded validators, insert them into currValidators
	k.valset.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		currValidators[validator.GetOperator().String()] = validator.GetBondedTokens()
		return false
	})

	k.IterateVotes(ctx, func(programID uint64, voter sdk.AccAddress, option bool) (stop bool) {
		valAddrStr := sdk.ValAddress(voter).String()
		if bondSize, ok := currValidators[valAddrStr]; ok {
			if option == true {
				votePower.Add(bondSize)
			} else {
				votePower.Sub(bondSize)
			}
			totalPower.Add(bondSize)
		}
		return false
	})

	return
}

func clearsThreshold(votePower, totalPower sdk.Int, threshold sdk.Dec) bool {
	return votePower.GTE(threshold.MulInt(totalPower).RoundInt())
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (claims types.ClaimPool, resTags sdk.Tags) {
	params := k.GetParams(ctx)

	k.CandQueueIterateMature(ctx, ctx.BlockHeight(), func(programID uint64) (stop bool) {
		program, err := k.GetProgram(ctx, programID)
		if err != nil {
			return false
		}

		// Did not pass the tally, delete program
		votePower, totalPower := tally(ctx, k, program)
		if !clearsThreshold(votePower, totalPower, params.ActiveThreshold) {
			k.DeleteProgram(ctx, programID)
			resTags.AppendTag(tags.Action, tags.ActionProgramRejected)
		} else {
			resTags.AppendTag(tags.Action, tags.ActionProgramPassed)
		}

		resTags.AppendTags(
			sdk.NewTags(
				tags.ProgramID, sdk.Uint64ToBigEndian(programID),
				tags.Weight, votePower,
			),
		)

		k.CandQueueRemove(ctx, program, programID)
		return false
	})

	// Not time to review programs yet
	if ctx.BlockHeight()%k.GetParams(ctx).VotePeriod != 0 {
		return
	}

	claims = types.ClaimPool{}

	// iterate programs and weight them
	k.IteratePrograms(ctx, func(programID uint64, program Program) (stop bool) {
		votePower, totalPower := tally(ctx, k, program)

		// Need to legacy program
		if !clearsThreshold(votePower, totalPower, params.LegacyThreshold) {
			k.DeleteProgram(ctx, programID)
			resTags.AppendTag(tags.Action, tags.ActionProgramLegacied)

		} else {
			claims = append(claims, types.NewClaim(types.BudgetClaimClass, votePower, program.Executor))
			resTags.AppendTag(tags.Action, tags.ActionProgramGranted)
		}

		resTags.AppendTags(
			sdk.NewTags(
				tags.ProgramID, sdk.Uint64ToBigEndian(programID),
				tags.Weight, votePower,
			),
		)

		return false
	})

	return
}
