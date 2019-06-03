package budget

import (
	"reflect"
	"strconv"

	"github.com/terra-project/core/x/budget/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// NewHandler creates a new handler for all budget type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitProgram:
			return handleMsgSubmitProgram(ctx, k, msg)
		case MsgWithdrawProgram:
			return handleMsgWithdrawProgram(ctx, k, msg)
		case MsgVoteProgram:
			return handleMsgVoteProgram(ctx, k, msg)

		default:
			errMsg := "Unrecognized budget Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgVoteProgram handles the logic of a MsgSubmitProgram
func handleMsgSubmitProgram(ctx sdk.Context, k Keeper, msg MsgSubmitProgram) sdk.Result {

	// Subtract coins from the submitter balance and updates it
	depositErr := k.PayDeposit(ctx, msg.Submitter)
	if depositErr != nil {
		return depositErr.Result()
	}

	// Create and add program
	programID := k.NewProgramID(ctx)
	program := NewProgram(
		programID,
		msg.Title,
		msg.Description,
		msg.Submitter,
		msg.Executor,
		ctx.BlockHeight(),
	)

	k.StoreProgram(ctx, program)
	k.CandQueueInsert(ctx, program.getVotingEndBlock(ctx, k), programID)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.ProgramID, strconv.FormatUint(programID, 10),
		),
	}
}

// handleMsgWithdrawProgram handles the logic of a MsgWithdrawProgram
func handleMsgWithdrawProgram(ctx sdk.Context, k Keeper, msg MsgWithdrawProgram) sdk.Result {
	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Only submitters can withdraw the program submission
	if !program.Submitter.Equals(msg.Submitter) {
		return ErrInvalidSubmitter(msg.Submitter).Result()
	}

	// Remove from candidate queue if not yet active
	prgmEndBlock := program.getVotingEndBlock(ctx, k)
	if k.CandQueueHas(ctx, prgmEndBlock, msg.ProgramID) {
		k.CandQueueRemove(ctx, prgmEndBlock, msg.ProgramID)
		err := k.RefundDeposit(ctx, program.Submitter)

		if err != nil {
			return ErrRefundFailed(msg.Submitter, msg.ProgramID).Result()
		}

		// Delete all votes on target program
		k.DeleteVotesForProgram(ctx, msg.ProgramID)
	}

	k.DeleteProgram(ctx, msg.ProgramID)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.ProgramID, strconv.FormatUint(msg.ProgramID, 10),
		),
	}
}

// handleMsgVoteProgram handles the logic of a MsgVoteProgram
func handleMsgVoteProgram(ctx sdk.Context, k Keeper, msg MsgVoteProgram) sdk.Result {
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

	k.AddVote(ctx, msg.ProgramID, msg.Voter, msg.Option)

	return sdk.Result{
		Tags: resTags.AppendTags(
			sdk.NewTags(
				tags.ProgramID, strconv.FormatUint(msg.ProgramID, 10),
				tags.Voter, msg.Voter.String(),
				tags.Option, strconv.FormatBool(msg.Option),
			),
		),
	}
}
