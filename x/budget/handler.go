package budget

import (
	"encoding/binary"
	"reflect"
	"terra/types/assets"
	"time"

	"terra/x/budget/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func uint64ToBytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

// NewHandler creates a new handler for all budget type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitProgramMsg:
			return handleSubmitProgramMsg(ctx, k, msg)
		case WithdrawProgramMsg:
			return handleWithdrawProgramMsg(ctx, k, msg)
		case VoteMsg:
			return handleVoteMsg(ctx, k, msg)
		default:
			errMsg := "Unrecognized budget Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func calculateThreshold(ctx sdk.Context, k Keeper, threshold sdk.Dec) sdk.Int {
	return threshold.MulInt(k.valset.TotalBondedTokens(ctx)).TruncateInt()
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (claims map[string]sdk.Int, resTags sdk.Tags) {
	params := k.GetParams(ctx)

	// Clean out expired programs
	k.IterateMatureCandidates(ctx, ctx.BlockHeader().Time, func(programID uint64, program Program) (stop bool) {

		k.CandidateQueueRemove(ctx, program, programID)

		// Program now activated.
		if program.Tally.GTE(calculateThreshold(ctx, k, params.ActiveThreshold)) {
			resTags = resTags.AppendTag(tags.Action, tags.ActionProgramPassed)
		} else {
			// Delete program
			k.DeleteProgram(ctx, programID)
			resTags = resTags.AppendTag(tags.Action, tags.ActionProgramRejected)
		}

		resTags = resTags.AppendTag(tags.ProgramID, string(programID))
		return false
	})

	// Add claims to re-weight claims in accordance with voting results
	if ctx.BlockHeight()%int64(k.GetParams(ctx).VotePeriod) == 0 {

		k.IterateActivePrograms(ctx, func(programID uint64, program Program) (stop bool) {
			claimantAddr := program.Executor.String()
			claims[claimantAddr] = claims[claimantAddr].Add(program.Tally)

			resTags = resTags.AppendTags(
				sdk.NewTags(
					tags.Action, tags.ActionProgramGranted,
					tags.ProgramID, string(programID),
					tags.Submitter, program.Submitter.String(),
					tags.Executor, program.Executor.String(),
					tags.Weight, program.Tally.String(),
				),
			)
			return false
		})
	}

	return
}

// handleVoteMsg handles the logic of a SubmitProgramMsg
func handleSubmitProgramMsg(ctx sdk.Context, k Keeper, msg SubmitProgramMsg) sdk.Result {

	// Deposit should be paid in TerraSDR
	if msg.Deposit.Denom != assets.SDRDenom {
		return ErrDepositDenom().Result()
	}

	// If deposit is sufficient
	if msg.Deposit.IsLT(k.GetParams(ctx).MinDeposit) {
		return ErrMinimumDeposit().Result()
	}

	// Subtract coins from the submitter balance and updates it
	_, _, err := k.bk.SubtractCoins(ctx, msg.Submitter, sdk.Coins{msg.Deposit})
	if err != nil {
		return err.Result()
	}

	// Create and add program
	program := NewProgram(
		msg.Title,
		msg.Description,
		msg.Submitter,
		msg.Executor,
		time.Now(),
		msg.Deposit,
	)
	programID := k.NewProgramID(ctx)
	k.SetProgram(ctx, programID, program)

	// Add to candidate program queue
	k.CandidateQueueInsert(ctx, program, programID)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionProgramSubmitted,
			tags.ProgramID, uint64ToBytes(programID),
			tags.Submitter, msg.Submitter.Bytes(),
			tags.Executor, msg.Executor.Bytes(),
		),
	}
}

// handleWithdrawProgramMsg handles the logic of a WithdrawProgramMsg
func handleWithdrawProgramMsg(ctx sdk.Context, k Keeper, msg WithdrawProgramMsg) sdk.Result {
	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Only submitters can withdraw the program submission
	if program.Submitter.Equals(msg.Submitter) {
		return ErrInvalidSubmitter(msg.Submitter).Result()
	}

	// Remove from candidate queue if not yet active
	if k.CandidateQueueHas(ctx, program, msg.ProgramID) {
		k.CandidateQueueRemove(ctx, program, msg.ProgramID)
	} else {
		// Only refund the deposit if the program is already inactive
		k.RefundDeposit(ctx, msg.ProgramID)
	}

	k.DeleteProgram(ctx, msg.ProgramID)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionProgramWithdrawn,
			tags.ProgramID, msg.ProgramID,
			tags.Submitter, msg.Submitter.Bytes(),
			tags.Executor, program.Executor.Bytes(),
		),
	}
}

// handleVoteMsg handles the logic of a VoteMsg
func handleVoteMsg(ctx sdk.Context, k Keeper, msg VoteMsg) sdk.Result {
	resTags := sdk.NewTags()

	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Check the voter is a validator
	val := k.valset.Validator(ctx, sdk.ValAddress(program.Submitter))
	if val == nil {
		return staking.ErrNoDelegatorForAddress(DefaultCodespace).Result()
	}

	// Override existing vote
	oldOption, err := k.GetVote(ctx, msg.ProgramID, msg.Voter)
	if err == nil {
		program.updateTally(oldOption, val.GetBondedTokens().Neg())
	}

	// update new vote
	program.updateTally(msg.Option, val.GetBondedTokens())

	// TODO: why does the vote need to be stored?
	k.SetVote(ctx, msg.ProgramID, msg.Voter, msg.Option)

	// The support level has now fallen below the legacy threshold; drop
	params := k.GetParams(ctx)
	if !k.CandidateQueueHas(ctx, program, msg.ProgramID) &&
		program.Tally.LT(calculateThreshold(ctx, k, params.LegacyThreshold)) {
		k.ClearVotesForProgram(ctx, msg.ProgramID)
		k.DeleteProgram(ctx, msg.ProgramID)

		resTags = resTags.AppendTags(
			sdk.NewTags(
				tags.Action, tags.ActionProgramLegacied,
				tags.ProgramID, uint64ToBytes(msg.ProgramID),
			),
		)
	}

	resTags = resTags.AppendTags(
		sdk.NewTags(
			tags.Action, tags.ActionProgramVote,
			tags.ProgramID, uint64ToBytes(msg.ProgramID),
			tags.Voter, msg.Voter.Bytes(),
			tags.Option, msg.Option,
		),
	)

	return sdk.Result{
		Tags: resTags,
	}
}
