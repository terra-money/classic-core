package budget

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the budget module
const (
	RouterKey = "budget"

	// QuerierRoute is the querier route for budget
	QuerierRoute = "budget"
)

// Program defines the basic properties of a staking Program
type Program struct {
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Validator address of the proposer
	Executor    sdk.AccAddress `json:"executor"`    // Account address of the executor
	SubmitTime  time.Time      `json:"submit_time"` // Block height from which the Program is open for votations
	Deposit     sdk.Coin       `json:"deposit"`     // Coins deposited in escrow
	Tally       sdk.Int        `json:"tally_result"`
}

// NewProgram validates deposit and creates a new Program
func NewProgram(
	title string,
	description string,
	submitter sdk.AccAddress,
	executor sdk.AccAddress,
	submitTime time.Time,
	deposit sdk.Coin) Program {
	return Program{
		Title:       title,
		Description: description,
		Submitter:   submitter,
		Executor:    executor,
		SubmitTime:  submitTime,
		Deposit:     deposit,
		Tally:       sdk.ZeroInt(),
	}
}

func (p *Program) getVotingEndTime(votingPeriod time.Duration) time.Time {
	return p.SubmitTime.Add(votingPeriod)
}

// updateTally updates the counter for each of the available options
func (p *Program) updateTally(option bool, power sdk.Int) {
	if option {
		p.Tally = p.Tally.Add(power)
	} else {
		p.Tally = p.Tally.Sub(power)
	}
}

// String implements fmt.Stringer
func (p Program) String() string {
	return fmt.Sprintf("Program{ Title: %s, Description: %s, Submitter: %v, Executor: %v, SubmitTime: %v, Deposit: %v}",
		p.Title, p.Description, p.Submitter, p.Executor, p.SubmitTime, p.Deposit)
}

//--------------------------------------------------------
//--------------------------------------------------------

//SubmitProgramMsg defines a message to create a Program
type SubmitProgramMsg struct {
	Title       string         // Title of the Program
	Description string         // Description of the Program
	Deposit     sdk.Coin       // Deposit paid by submitter. Must be > MinDeposit to enter voting period
	Submitter   sdk.AccAddress // Address of the submitter
	Executor    sdk.AccAddress // Address of the executor
}

// NewSubmitProgramMsg submits a message with a new Program
func NewSubmitProgramMsg(title string, description string, deposit sdk.Coin,
	submitter sdk.AccAddress, executor sdk.AccAddress) SubmitProgramMsg {
	return SubmitProgramMsg{
		Title:       title,
		Description: description,
		Deposit:     deposit,
		Submitter:   submitter,
		Executor:    executor,
	}
}

// Route Implements Msg
func (msg SubmitProgramMsg) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg SubmitProgramMsg) Type() string { return "submitprogram" }

// Implements Msg
func (msg SubmitProgramMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg SubmitProgramMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

// Implements Msg
func (msg SubmitProgramMsg) ValidateBasic() sdk.Error {
	if len(msg.Submitter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Submitter.String())
	}
	if len(msg.Executor) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Executor.String())
	}
	if len(strings.TrimSpace(msg.Title)) <= 0 {
		return ErrInvalidTitle()
	}

	if len(strings.TrimSpace(msg.Description)) <= 0 {
		return ErrInvalidDescription()
	}

	if !msg.Deposit.IsPositive() {
		return sdk.ErrInvalidCoins("Deposit is not valid")
	}

	return nil
}

func (msg SubmitProgramMsg) String() string {
	return fmt.Sprintf("SubmitProgramMsg{%v, %v}", msg.Title, msg.Description)
}

//--------------------------------------------------------
//--------------------------------------------------------

// WithdrawProgramMsg defines the msg of a staker containing the vote option to an
// specific Program
type WithdrawProgramMsg struct {
	ProgramID uint64         // ID of the Program
	Submitter sdk.AccAddress // Address of the voter
}

// NewVoteMsg creates a VoteMsg instance
func NewWithdrawProgramMsg(programID uint64, submitter sdk.AccAddress) WithdrawProgramMsg {
	return WithdrawProgramMsg{
		ProgramID: programID,
		Submitter: submitter,
	}
}

// Route Implements Msg
func (msg WithdrawProgramMsg) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg WithdrawProgramMsg) Type() string { return "withdraw" }

// Implements Msg
func (msg WithdrawProgramMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg WithdrawProgramMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

// Implements Msg
func (msg WithdrawProgramMsg) ValidateBasic() sdk.Error {
	if len(msg.Submitter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Submitter.String())
	}
	return nil
}

// Implements Msg
func (msg WithdrawProgramMsg) String() string {
	return fmt.Sprintf("WithdrawProgramMsg{%v, %v}", msg.ProgramID, msg.Submitter)
}

//--------------------------------------------------------
//--------------------------------------------------------

// VoteMsg defines the msg of a staker containing the vote option to an
// specific Program
type VoteMsg struct {
	ProgramID uint64         // ID of the Program
	Option    bool           // Option chosen by voter
	Voter     sdk.AccAddress // Address of the voter
}

// NewVoteMsg creates a VoteMsg instance
func NewVoteMsg(programID uint64, option bool, voter sdk.AccAddress) VoteMsg {
	return VoteMsg{
		ProgramID: programID,
		Option:    option,
		Voter:     voter,
	}
}

// Route Implements Msg
func (msg VoteMsg) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg VoteMsg) Type() string { return "vote" }

// Implements Msg
func (msg VoteMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg VoteMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// Implements Msg
func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Voter.String())
	}
	if msg.ProgramID <= 0 {
		return ErrInvalidProgramID("ProgramID cannot be negative")
	}

	return nil
}

// Implements Msg
func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{%v, %v, %v}", msg.ProgramID, msg.Voter, msg.Option)
}
