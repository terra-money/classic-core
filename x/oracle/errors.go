package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle errors reserve 1101-1199
const (
	DefaultCodespace sdk.CodespaceType = "ORA"

	CodeNotValidator   sdk.CodeType = 1
	CodeNotEnoughVotes sdk.CodeType = 2
	CodeUnknownRequest sdk.CodeType = sdk.CodeUnknownRequest
)

// ----------------------------------------
// Error constructors

// ErrWrongDenomination called when the signer of a Msg is not a validator
func ErrWrongDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeNotValidator, denom)
}

// ErrNotValidator called when the signer of a Msg is not a validator
func ErrNotValidator(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNotValidator, address.String())
}

// ErrAlreadyProcessed called when a payload is already processed
func ErrNotEnoughVotes(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotEnoughVotes, "")
}
