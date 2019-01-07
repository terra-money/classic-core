package budget

import (
	"encoding/binary"
	"reflect"
	"terra/types/assets"
	"terra/x/treasury"

	"terra/x/budget/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
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
func EndBlocker(ctx sdk.Context, k Keeper) (resTags sdk.Tags) {
	resTags = sdk.NewTags()

	// Clean out expired inactive programs 
	inactiveIterator := k.InactiveProgramQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; inactiveIterator.Valid(); inactiveIterator.Next() {
		var programID uint64

		k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &programID)
		program := k.GetProgram(ctx, programID)

		k.RemoveFromInactiveProgramQueue(ctx, ctx.BlockHeader().Time, programID)
		k.DeleteProgram(ctx, programID)

		resTags = resTags.AppendTag(tags.Action, tags.ActionProgramDropped)
		resTags = resTags.AppendTag(tags.ProgramID, []byte(string(programID)))
	}
	inactiveIterator.Close()

	// Add claims to re-weight claims in accordance with voting results 
	if ctx.BlockHeight()%k.GetParams(ctx).VotePeriod == 0 {
		programIterator := sdk.KVStorePrefixIterator(k.key, PrefixProgram)
		for ; programIterator.Valid(); programIterator.Next() {
			
			var programID uint64
			var program Program
			k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Key(), &programID)
			k.cdc.MustUnmarshalBinaryLengthPrefixed(inactiveIterator.Value(), &program)

			k.tk.AddClaim(ctx, treasury.NewBaseClaim(
				treasury.BudgetShareID,
				program.weight(),
				program.Executor,
			))

			resTags = resTags.AppendTag(tags.ProgramID, []byte(string(programID)))
		}

		programIterator.Close()
	}

	return resTags
}

// handleVoteMsg handles the logic of a SubmitProgramMsg
func handleSubmitProgramMsg(ctx sdk.Context, k Keeper, msg SubmitProgramMsg) sdk.Result {

	// If deposit is sufficient
	if msg.Deposit.AmountOf(assets.TerraDenom).GT(sdk.NewInt(k.GetParams(ctx).MinDeposit)) {
		// Subtract coins from the submitter balance and updates it
		_, _, err := k.ck.SubtractCoins(ctx, msg.Submitter, msg.Deposit)
		if err != nil {
			return err.Result()
		}

		program := NewProgram(
			msg.Title,
			msg.Description,
			msg.Submitter,
			msg.Executor,
			ctx.BlockHeight(),
			msg.Deposit)

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
	if program.Submitter != msg.Submitter {
		return ErrInvalidSubmissiter().Result
	}

	// Refund the deposit
	k.RefundDeposit(ctx, msg.Submitter)

	// Only allow inactive programs to be withdrawn
	votingEndTime := program.getVotingEndTime(k.GetParams(ctx).VotingPeriod)
	if k.ProgramExistsInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID) {
		k.RemoveFromInactiveProgramQueue(ctx, votingEndTime, msg.ProgramID)
	}
	k.DeleteProgram(ctx, msg.ProgramID)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionProgramWithdrawn,
			tags.ProgramID, msg.ProgramID,
			tags.Submitter, msg.Submitter.Bytes(),
		),
	}
}

// handleVoteMsg handles the logic of a VoteMsg
func handleVoteMsg(ctx sdk.Context, k Keeper, msg VoteMsg) sdk.Result {
	program, err := k.GetProgram(ctx, msg.ProgramID)
	if err != nil {
		return ErrProgramNotFound(msg.ProgramID).Result()
	}

	// Voting period is closed
	if ctx.BlockHeight() > program.SubmitBlock+k.GetParams(ctx).VotePeriod {
		return ErrVotingPeriodClosed().Result()
	}

	// Check the voter is a validater
	val := k.valset.Validator(ctx, sdk.ValAddress(program.Submitter))
	if val == nil {
		return stake.ErrNoDelegatorForAddress(DefaultCodespace).Result()
	}

	// Override existing vote
	oldOption, err := k.GetVote(ctx, msg.ProgramID, msg.Voter)
	if err != nil {
		program.updateTally(oldOption, val.GetPower().Neg())
	}

	// update new vote
	err = program.updateTally(msg.Option, val.GetPower())

	// Needs to be activated 
	if ProgramExistsInactiveProgramQueue(ctx, program.getVotingEndTime(k.GetParams(ctx).VotePeriod), msg.ProgramID) {
		if program.weight().GT(k.GetParams(ctx).ActiveThreshold) {
			// Refund deposit
			k.RefundDeposit(ctx, msg.ProgramID)

			k.RemoveFromInactiveProgramQueue(ctx, endTime, msg.ProgramID)
		}
	} 
	
	// Needs to be legacied
	else if program.weight().LT(k.GetParams(ctx).LegacyThreshold) {
		k.DeleteProgram(ctx, msg.ProgramID)
		// Burn the deposit
	}
	
	// TODO: why does the vote need to be stored? 
	k.SetVote(ctx, msg.ProgramID, msg.Voter, msg.Option)

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Action, tags.ActionProgramVote,
			tags.ProgramID, int64ToBytes(msg.ProgramID),
			tags.Voter, msg.Voter.Bytes(),
			tags.Option, []byte(msg.Option),
		),
	}
}
