package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle error codes
const (
	DefaultCodespace sdk.CodespaceType = "treasury"

	CodeInvalidEpoch sdk.CodeType = 1
)

// ----------------------------------------
// Error constructors

// ErrInvalidEpoch called when the denom is not known
func ErrInvalidEpoch(codespace sdk.CodespaceType, curEpoch, epoch int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidEpoch, fmt.Sprintf("The query epoch should be between [0, %d] but given %d", curEpoch, epoch))
}
