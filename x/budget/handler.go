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

// NewHandler creates a new handler for all simple_gov type messages.
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

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) (claims map[string]sdk.Int, resTags sdk.Tags) {
	// Clean out expired inactive programs
	inactiveIterator := k.InactiveProgramQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var programID uint64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &programID)

		k.RemoveFromInactiveProgramQueue(ctx, ctx.BlockHeader().Time, programID)
		k.DeleteProgram(ctx, programID)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProgramDropped)
		resTags = resTags.AppendTag(tags.ProgramID, string(programID))
	}
	inactiveIterator.Close()

	// Add claims to re-weight claims in accordance with voting results
	if ctx.BlockHeight()%int64(k.GetParams(ctx).VotePeriod) == 0 {
		programIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.key), PrefixProgram)
		for ; programIterator.Valid(); programIterator.Next() {

			var programID uint64
			var program Program
			k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Key(), &programID)
			k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &program)

			claimantAddr := program.Executor.String()
			claims[claimantAddr] = claims[claimantAddr].Add(program.weight())

			resTags = resTags.AppendTags(
				sdk.NewTags(
					tags.Action, tags.ActionProgramGranted,
					tags.ProgramID, string(programID),
					tags.Submitter, program.Submitter.String(),
					tags.Executor, program.Executor.String(),
					tags.Weight, program.weight().String(),
				),
			)
		}

		programIterator.Close()
	}

	return
}

// handleVoteMsg handles the logic of a SubmitProgramMsg
func handleSubmitProgramMsg(ctx sdk.Context, k Keeper, msg SubmitProgramMsg) sdk.Result {

	// If deposit is sufficient
	if msg.Deposit.AmountOf(assets.SDRDenom).GT(sdk.NewInt(k.GetParams(ctx).MinDeposit)) {
		// Subtract coins from the submitter balance and updates it
		_, _, err := k.bk.SubtractCoins(ctx, msg.Submitter, msg.Deposit)
		if err != nil {
			return err.Result()
		}

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

		return sdk.Result{
			Tags: sdk.NewTags(
				tags.Action, tags.ActionProgramSubmitted,
				tags.ProgramID, uint64ToBytes(programID),
				tags.Submitter, msg.Submitter.Bytes(),
				tags.Executor, msg.Executor.Bytes(),
			),
		}
	}
	return ErrMinimumDeposit().Result()
}

// handleVoteMsg handles the logic of a SubmitProgramMsg
func handleWithdrawProgramMsg(ctx sdk.Context, k Keeper, msg WithdrawProgramMsg) sdk.Result {
	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Only submitters can withdraw the program submission
	if program.Submitter.Equals(msg.Submitter) {
		return ErrInvalidSubmissiter(msg.Submitter).Result()
	}

	// Refund the deposit
	k.RefundDeposit(ctx, msg.ProgramID)

	// Only allow inactive programs to be withdrawn
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotePeriod)
	if k.ProgramExistsInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID) {
		k.RemoveFromInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID)
	}
	program.State = LegacyProgramState
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
	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Check the voter is a validater
	val := k.valset.Validator(ctx, sdk.ValAddress(program.Submitter))
	if val == nil {
		return staking.ErrNoDelegatorForAddress(DefaultCodespace).Result()
	}

	// Override existing vote
	oldOption, err := k.GetVote(ctx, msg.ProgramID, msg.Voter)
	if err != nil {
		program.updateTally(oldOption, val.GetBondedTokens().Neg())
	}

	// update new vote
	err = program.updateTally(msg.Option, val.GetBondedTokens())

	// Needs to be activated
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotePeriod)
	if k.ProgramExistsInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID) {

		activationThreshold := k.GetParams(ctx).ActiveThreshold.MulInt(k.valset.TotalBondedTokens(ctx)).TruncateInt()
		if program.weight().GT(activationThreshold) {
			// Refund deposit
			k.RefundDeposit(ctx, msg.ProgramID)

			k.RemoveFromInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID)

			program.State = ActiveProgramState
			k.SetProgram(ctx, msg.ProgramID, program)
		}
	} else {
		legacyThreshold := k.GetParams(ctx).LegacyThreshold.MulInt(k.valset.TotalBondedTokens(ctx)).TruncateInt()
		if program.weight().LT(legacyThreshold) {
			program.State = LegacyProgramState

			k.DeleteProgram(ctx, msg.ProgramID)
			// Burn the deposit
		}
	}

	// TODO: why does the vote need to be stored?
	k.SetVote(ctx, msg.ProgramID, msg.Voter, msg.Option)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionProgramVote,
			tags.ProgramID, uint64ToBytes(msg.ProgramID),
			tags.Voter, msg.Voter.Bytes(),
			tags.Option, []byte(msg.Option),
		),
	}
}
