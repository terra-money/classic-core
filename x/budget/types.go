package budget

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProgramVote string
type ProgramState string

const (
	YesVote     ProgramVote = "yes"
	NoVote      ProgramVote = "no"
	AbstainVote ProgramVote = "abstain"

	ActiveProgramState   ProgramState = "Inactive"
	RejectedProgramState ProgramState = "Rejected"
	LegacyProgramState   ProgramState = "Legacy"
	ActiveProgramState   ProgramState = "Active"
)

// Program defines the basic properties of a staking Program
type Program struct {
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Validator address of the proposer
	Executor    sdk.AccAddress `json:"executor"`    // Account address of the executor
	SubmitTime  time.Time      `json:"submit_time"` // Block height from which the Program is open for votations
	Deposit     sdk.Coins      `json:"deposit"`     // Coins deposited in escrow
	TallyResult TallyResult    `json:"tally_result"`
}

// NewProgram validates deposit and creates a new Program
func NewProgram(
	title string,
	description string,
	submitter sdk.AccAddress,
	executor sdk.AccAddress,
	submitTime time.Time,
	deposit sdk.Coins) Program {
	return Program{
		Title:       title,
		Description: description,
		Submitter:   submitter,
		SubmitTime:  submitTime,
		Deposit:     deposit,
		TallyResult: EmptyTallyResult(),
	}
}

func (p *Program) getVotingEndTime(votingPeriod time.Time) {
	return p.SubmitTime.Add(votingPeriod)
}

// updateTally updates the counter for each of the available options
func (p *Program) updateTally(option ProgramVote, power sdk.Dec) sdk.Error {
	switch option {
	case YesVote:
		p.TallyResult.Yes = p.TallyResult.Yes.Add(power)
		return nil
	case NoVote:
		p.TallyResult.No = p.TallyResult.No.Add(power)
		return nil
	case AbstainVote:
		p.TallyResult.Abstain = p.TallyResult.Abstain.Add(power)
		return nil
	default:
		return ErrInvalidOption("Invalid option: " + option)
	}
}

func (p *Program) weight() sdk.Dec {
	return p.TallyResult.Yes.Sub(p.TallyResult.No)
}

//--------------------------------------------------------
//--------------------------------------------------------

//SubmitProgramMsg defines a message to create a Program
type SubmitProgramMsg struct {
	Title       string         // Title of the Program
	Description string         // Description of the Program
	Deposit     sdk.Coins      // Deposit paid by submitter. Must be > MinDeposit to enter voting period
	Submitter   sdk.AccAddress // Address of the submitter
	Executor    sdk.AccAddress // Address of the executor
}

// NewSubmitProgramMsg submits a message with a new Program
func NewSubmitProgramMsg(title string, description string, votingWindow int64, deposit sdk.Coins,
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
	return []sdk.Address{msg.Submitter}
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

	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins("Deposit is not valid")
	}

	if !msg.Deposit.IsPositive() {
		return sdk.ErrInvalidCoins("Deposit cannot be negative")
	}

	return nil
}

func (msg SubmitProgramMsg) String() string {
	return fmt.Sprintf("SubmitProgramMsg{%v, %v}", msg.Title, msg.Description)
}

//--------------------------------------------------------
//--------------------------------------------------------

// TallyResult Tally Results
type TallyResult struct {
	Yes     sdk.Dec `json:"yes"`
	Abstain sdk.Dec `json:"abstain"`
	No      sdk.Dec `json:"no"`
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:     sdk.ZeroDec(),
		Abstain: sdk.ZeroDec(),
		No:      sdk.ZeroDec(),
	}
}

// checks if two proposals are equal
func (resultA TallyResult) Equals(resultB TallyResult) bool {
	return (resultA.Yes.Equal(resultB.Yes) &&
		resultA.Abstain.Equal(resultB.Abstain) &&
		resultA.No.Equal(resultB.No))
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
func NewWithdrawProgramMsg(programID int64, submitter sdk.AccAddress) WithdrawProgramMsg {
	return VoteMsg{
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
	return []sdk.Address{msg.Submitter}
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
	Option    ProgramVote    // Option chosen by voter
	Voter     sdk.AccAddress // Address of the voter
}

// NewVoteMsg creates a VoteMsg instance
func NewVoteMsg(ProgramID int64, option ProgramVote, voter sdk.AccAddress) VoteMsg {
	// by default a nil option is an abstention
	switch option {
	case YesVote:
	case NoVote:
		break
	default:
		option = AbstainVote
	}
	return VoteMsg{
		ProgramID: ProgramID,
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
	return []sdk.Address{msg.Voter}
}

func isValidOption(option ProgramVote) bool {
	switch option {
	case YesVote:
	case NoVote:
	case AbstainVote:
		return true
	}

	return false
}

// Implements Msg
func (msg VoteMsg) ValidateBasic() sdk.Error {
	if len(msg.Voter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Voter.String())
	}
	if msg.ProgramID <= 0 {
		return ErrInvalidProgramID("ProgramID cannot be negative")
	}
	if !isValidOption(msg.Option) {
		return ErrInvalidOption("Invalid voting option: " + msg.Option)
	}
	if len(strings.TrimSpace(msg.Option)) <= 0 {
		return ErrInvalidOption("Option can't be blank")
	}

	return nil
}

// Implements Msg
func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{%v, %v, %v}", msg.ProgramID, msg.Voter, msg.Option)
}
