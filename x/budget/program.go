package budget

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Program defines the basic properties of a staking Program
type Program struct {
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Validator address of the proposer
	Executor    sdk.AccAddress `json:"executor"`    // Account address of the executor
	SubmitBlock int64          `json:"submit_time"` // Block height from which the Program is open for votations
}

// NewProgram validates deposit and creates a new Program
func NewProgram(
	title string,
	description string,
	submitter sdk.AccAddress,
	executor sdk.AccAddress,
	submitBlock int64) Program {
	return Program{
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
	Title: %s
	Description: %s
	Submitter: %v
	Executor: %v
	SubmitBlock: %d`,
		p.Title, p.Description, p.Submitter, p.Executor, p.SubmitBlock)
}
