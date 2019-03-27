package budget

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitProgramMsg defines a message to create a Program
type MsgSubmitProgram struct {
	Title       string         `json:"title"`       // Title of the Program
	Description string         `json:"description"` // Description of the Program
	Submitter   sdk.AccAddress `json:"submitter"`   // Address of the submitter
	Executor    sdk.AccAddress `json:"executor"`    // Address of the executor
}

// NewMsgSubmitProgram submits a message with a new Program
func NewMsgSubmitProgram(title string, description string,
	submitter sdk.AccAddress, executor sdk.AccAddress) MsgSubmitProgram {
	return MsgSubmitProgram{
		Title:       title,
		Description: description,
		Submitter:   submitter,
		Executor:    executor,
	}
}

// Route Implements Msg
func (msg MsgSubmitProgram) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg MsgSubmitProgram) Type() string { return "submitprogram" }

// Implements Msg
func (msg MsgSubmitProgram) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg MsgSubmitProgram) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

// Implements Msg
func (msg MsgSubmitProgram) ValidateBasic() sdk.Error {
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

	return nil
}

// Implements Msg
func (msg MsgSubmitProgram) String() string {
	return fmt.Sprintf(`MsgSubmitProgram
	Title: %v
	Submitter: %v
	Executor: %v`, msg.Title, msg.Submitter, msg.Executor)
}

//--------------------------------------------------------
//--------------------------------------------------------

// WithdrawProgramMsg defines the msg of a staker containing the vote option to an
// specific Program
type MsgWithdrawProgram struct {
	ProgramID uint64         `json:"program_id"` // ID of the Program
	Submitter sdk.AccAddress `json:"submitter"`  // Address of the voter
}

// NewVoteMsg creates a VoteMsg instance
func NewMsgWithdrawProgram(programID uint64, submitter sdk.AccAddress) MsgWithdrawProgram {
	return MsgWithdrawProgram{
		ProgramID: programID,
		Submitter: submitter,
	}
}

// Route Implements Msg
func (msg MsgWithdrawProgram) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg MsgWithdrawProgram) Type() string { return "withdraw" }

// Implements Msg
func (msg MsgWithdrawProgram) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg MsgWithdrawProgram) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

// Implements Msg
func (msg MsgWithdrawProgram) ValidateBasic() sdk.Error {
	if len(msg.Submitter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Submitter.String())
	}
	return nil
}

// Implements Msg
func (msg MsgWithdrawProgram) String() string {
	return fmt.Sprintf(`MsgWithdrawProgram
	ProgramID: %v
	Submitter: %v`, msg.ProgramID, msg.Submitter)
}

//--------------------------------------------------------
//--------------------------------------------------------

// MsgVoteProgram defines the msg of a staker containing the vote option to an
// specific Program
type MsgVoteProgram struct {
	ProgramID uint64         `json:"program_id"` // ID of the Program
	Option    bool           `json:"option"`     // Option chosen by voter
	Voter     sdk.AccAddress `json:"voter"`      // Address of the voter
}

// NewMsgVoteProgram creates a MsgVoteProgram instance
func NewMsgVoteProgram(programID uint64, option bool, voter sdk.AccAddress) MsgVoteProgram {
	return MsgVoteProgram{
		ProgramID: programID,
		Option:    option,
		Voter:     voter,
	}
}

// Route Implements Msg
func (msg MsgVoteProgram) Route() string { return "budget" }

// Type implements sdk.Msg
func (msg MsgVoteProgram) Type() string { return "vote" }

// Implements Msg
func (msg MsgVoteProgram) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg
func (msg MsgVoteProgram) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// Implements Msg
func (msg MsgVoteProgram) ValidateBasic() sdk.Error {
	if len(msg.Voter) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Voter.String())
	}
	if msg.ProgramID <= 0 {
		return ErrInvalidProgramID(msg.ProgramID)
	}

	return nil
}

// Implements Msg
func (msg MsgVoteProgram) String() string {
	return fmt.Sprintf(`MsgVoteProgram
	ProgramID: %v
	Voter: %v
	Option: %v`, msg.ProgramID, msg.Voter, msg.Option)
}
