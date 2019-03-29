package budget

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace nolint
	DefaultCodespace sdk.CodespaceType = "budget"

	// Budget errors
	CodeInvalidProgramID   sdk.CodeType = 1
	CodeInvalidTitle       sdk.CodeType = 2
	CodeInvalidDescription sdk.CodeType = 3
	CodeProgramNotFound    sdk.CodeType = 4
	CodeVoteNotFound       sdk.CodeType = 5
	CodeInvalidSubmitter   sdk.CodeType = 6
)

// nolint
func ErrInvalidTitle() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidTitle, "Cannot submit a program with empty title")
}

// nolint
func ErrInvalidDescription() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidDescription, "Cannot submit a program with empty description")
}

// nolint
func ErrProgramNotFound(ProgramID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeProgramNotFound, "program with id "+
		strconv.Itoa(int(ProgramID))+" not found")
}

// nolint
func ErrInvalidProgramID(ProgramID uint64) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidProgramID, "program id "+
		strconv.Itoa(int(ProgramID))+" invalid. Must be an uint")
}

// nolint
func ErrVoteNotFound() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeVoteNotFound, "Vote not found")
}

// nolint
func ErrInvalidSubmitter(submitter sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidSubmitter, fmt.Sprintf("Submitter does not match %s", submitter))
}
