package budget

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // budget params

	ActivePrograms    Programs `json:"active_programs"`
	CandidatePrograms Programs `json:"candidate_programs"`

	Votes Votes `json:"votes"`
}

func NewGenesisState(params Params, activePrograms,
	candidatePrograms Programs, votes Votes) GenesisState {
	return GenesisState{
		Params: params,

		ActivePrograms:    activePrograms,
		CandidatePrograms: candidatePrograms,
		Votes:             votes,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),

		ActivePrograms:    Programs{},
		CandidatePrograms: Programs{},
		Votes:             Votes{},
	}
}

// new oracle genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, program := range data.ActivePrograms {
		keeper.StoreProgram(ctx, program)
	}

	for _, program := range data.CandidatePrograms {
		keeper.StoreProgram(ctx, program)
		keeper.CandQueueInsert(ctx, data.Params.VotePeriod, program.ProgramID)
	}

	for _, vote := range data.Votes {
		keeper.AddVote(ctx, vote.ProgramID, vote.Voter, vote.Option)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	var activePrograms Programs
	keeper.IteratePrograms(ctx, true, func(program Program) (stop bool) {
		activePrograms = append(activePrograms, program)
		return false
	})

	var candidatePrograms Programs
	keeper.CandQueueIterate(ctx, func(programID uint64) (stop bool) {
		program, err := keeper.GetProgram(ctx, programID)
		if err != nil {
			return false
		}

		candidatePrograms = append(candidatePrograms, program)
		return false
	})

	var votes Votes
	keeper.IterateVotes(ctx, func(programID uint64, voterAddr sdk.AccAddress, option bool) (stop bool) {
		votes = append(votes, NewVote(programID, option, voterAddr))
		return false
	})

	return NewGenesisState(params, activePrograms, candidatePrograms, votes)
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	err := validateParams(data.Params)

	if err != nil {
		return err
	}

	var programMap map[uint64]bool
	for _, program := range data.ActivePrograms {
		if program.SubmitBlock != 0 {
			return ErrInvalidSubmitBlockHeight(program.SubmitBlock)
		}

		// duplicate program ID check
		if _, ok := programMap[program.ProgramID]; ok {
			return ErrDuplicateProgramID(program.ProgramID)
		}

		programMap[program.ProgramID] = true
	}

	for _, program := range data.CandidatePrograms {
		if program.SubmitBlock != 0 {
			return ErrInvalidSubmitBlockHeight(program.SubmitBlock)
		}

		// duplicate program ID check
		if _, ok := programMap[program.ProgramID]; ok {
			return ErrDuplicateProgramID(program.ProgramID)
		}

		programMap[program.ProgramID] = true
	}

	for _, vote := range data.Votes {
		if _, ok := programMap[vote.ProgramID]; !ok {
			return ErrProgramNotFound(vote.ProgramID)
		}
	}

	return nil
}
