package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle errors reserve 1101-1199
const (
	DefaultCodespace sdk.CodespaceType = "ORACLE"

	CodeNotValidator sdk.CodeType = 1
	CodeUnknownDenom sdk.CodeType = 2
)

// ----------------------------------------
// Error constructors

// ErrUnknownDenomination called when the signer of a Msg is not a validator
func ErrUnknownDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownDenom, denom)
}

// ErrNotValidator called when the signer of a Msg is not a validator
func ErrNotValidator(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotValidator, address.String())
}
