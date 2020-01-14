package budget

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Program defines the basic properties of a staking Program
type Program struct {
	ProgramID   uint64         `json:"program_id"`  // ID of the Program
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Validator address of the proposer
	Executor    sdk.AccAddress `json:"executor"`    // Account address of the executor
	SubmitBlock int64          `json:"submit_time"` // Block height from which the Program is open for votations
}

// NewProgram validates deposit and creates a new Program
func NewProgram(
	programID uint64,
	title string,
	description string,
	submitter sdk.AccAddress,
	executor sdk.AccAddress,
	submitBlock int64) Program {
	return Program{
		ProgramID:   programID,
		Title:       title,
		Description: description,
		Submitter:   submitter,
		Executor:    executor,
		SubmitBlock: submitBlock,
	}
}

func (p *Program) getVotingEndBlock(ctx sdk.Context, k Keeper) int64 {
	return p.SubmitBlock + k.GetParams(ctx).VotePeriod
}

// String implements fmt.Stringer
func (p Program) String() string {
	return fmt.Sprintf(`Program
	ProgramID: %d
	Title: %s
	Description: %s
	Submitter: %v
	Executor: %v
	SubmitBlock: %d`,
		p.ProgramID, p.Title, p.Description, p.Submitter, p.Executor, p.SubmitBlock)
}

// Programs is a collection of Program
type Programs []Program

func (p Programs) String() (out string) {
	for _, val := range p {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
