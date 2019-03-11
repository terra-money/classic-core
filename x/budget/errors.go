package budget

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "budget"

	// Budget errors
	CodeInvalidOption        CodeType = 1
	CodeInvalidProgramID     CodeType = 2
	CodeVotingPeriodClosed   CodeType = 3
	CodeEmptyProgramQueue    CodeType = 4
	CodeInvalidTitle         CodeType = 5
	CodeInvalidDescription   CodeType = 6
	CodeInvalidVotingWindow  CodeType = 7
	CodeProgramNotFound      CodeType = 8
	CodeVoteNotFound         CodeType = 9
	CodeProgramQueueNotFound CodeType = 10
	CodeInvalidDeposit       CodeType = 11
	CodeInvalidSubmitter     CodeType = 12
)

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInvalidOption:
		return "Invalid option"
	case CodeInvalidProgramID:
		return "Invalid ProgramID"
	case CodeVotingPeriodClosed:
		return "Voting Period Closed"
	case CodeEmptyProgramQueue:
		return "ProgramQueue is empty"
	case CodeInvalidTitle:
		return "Invalid program title"
	case CodeInvalidDescription:
		return "Invalid program description"
	case CodeInvalidVotingWindow:
		return "Invalid voting window"
	case CodeProgramNotFound:
		return "program not found"
	case CodeVoteNotFound:
		return "Option not found"
	case CodeProgramQueueNotFound:
		return "program Queue not found"
	case CodeInvalidDeposit:
		return "Invalid deposit"
	case CodeInvalidSubmitter:
		return "Invalid submitter"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

// nolint
func ErrInvalidOption(msg string) sdk.Error {
	if msg != "" {
		return newError(DefaultCodespace, CodeInvalidOption, msg)
	}
	return newError(DefaultCodespace, CodeInvalidOption, "The chosen option is invalid")
}

// nolint
func ErrInvalidProgramID(msg string) sdk.Error {
	if msg != "" {
		return newError(DefaultCodespace, CodeInvalidProgramID, msg)
	}
	return newError(DefaultCodespace, CodeInvalidProgramID, "ProgramID is not valid")
}

// nolint
func ErrInvalidTitle() sdk.Error {
	return newError(DefaultCodespace, CodeInvalidTitle, "Cannot submit a program with empty title")
}

// nolint
func ErrInvalidDescription() sdk.Error {
	return newError(DefaultCodespace, CodeInvalidDescription, "Cannot submit a program with empty description")
}

// nolint
func ErrVotingPeriodClosed() sdk.Error {
	return newError(DefaultCodespace, CodeVotingPeriodClosed, "Voting period is closed for this program")
}

// nolint
func ErrEmptyProgramQueue() sdk.Error {
	return newError(DefaultCodespace, CodeEmptyProgramQueue, "Can't get element from an empty program queue")
}

// nolint
func ErrProgramNotFound(ProgramID uint64) sdk.Error {
	return newError(DefaultCodespace, CodeProgramNotFound, "program with id "+
		strconv.Itoa(int(ProgramID))+" not found")
}

// nolint
func ErrVoteNotFound() sdk.Error {
	return newError(DefaultCodespace, CodeVoteNotFound, "Vote not found")
}

// nolint
func ErrProgramQueueNotFound() sdk.Error {
	return newError(DefaultCodespace, CodeProgramQueueNotFound, "program Queue not found")
}

// nolint
func ErrInvalidVotingWindow(msg string) sdk.Error {
	if msg != "" {
		return newError(DefaultCodespace, CodeInvalidVotingWindow, msg)
	}
	return newError(DefaultCodespace, CodeInvalidVotingWindow, "Voting window is not positive")
}

// nolint
func ErrMinimumDeposit() sdk.Error {
	return newError(DefaultCodespace, CodeInvalidDeposit, "Deposit is lower than the minimum")
}

// nolint
func ErrDepositDenom() sdk.Error {
	return newError(DefaultCodespace, CodeInvalidDeposit, "Deposit should be paid in TerraSDR")
}

// nolint
func ErrInvalidSubmitter(submitter sdk.AccAddress) sdk.Error {
	return newError(DefaultCodespace, CodeInvalidSubmitter, fmt.Sprintf("Submitter does not match %s", submitter))
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
