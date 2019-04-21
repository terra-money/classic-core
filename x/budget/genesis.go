package budget

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Params            Params             `json:"params"` // budget params
	SubmitProgramMsgs []MsgSubmitProgram `json:"programs"`
}

func NewGenesisState(params Params, msgs []MsgSubmitProgram) GenesisState {
	return GenesisState{
		Params:            params,
		SubmitProgramMsgs: msgs,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	params := DefaultParams()
	return GenesisState{
		Params:            params,
		SubmitProgramMsgs: nil,
	}
}

// new oracle genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, msg := range data.SubmitProgramMsgs {

		if err := msg.ValidateBasic(); err != nil {
			panic(err)
		}

		// Create and add program
		programID := keeper.NewProgramID(ctx)
		program := NewProgram(
			programID,
			msg.Title,
			msg.Description,
			msg.Submitter,
			msg.Executor,
			ctx.BlockHeight(),
		)

		keeper.SetProgram(ctx, programID, program)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, and validator/delegator distribution info's
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)

	msgs := []MsgSubmitProgram{}
	keeper.IteratePrograms(ctx, true, func(programID uint64, program Program) bool {
		msgs = append(msgs, NewMsgSubmitProgram(
			program.Title,
			program.Description,
			program.Submitter,
			program.Executor,
		))
		return false
	})

	return NewGenesisState(params, msgs)
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	return validateParams(data.Params)
}
