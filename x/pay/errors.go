// Pay TODO - mandatory update

package pay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCodespace no-lint
	DefaultCodespace sdk.CodespaceType = "pay"

	// CodePayDisabled no-lint
	CodePayDisabled sdk.CodeType = 101
)

// ErrPayDisabled is an error
func ErrPayDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodePayDisabled, "pay transactions are currently disabled")
}
