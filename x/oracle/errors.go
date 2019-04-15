package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle error codes
const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownDenom      sdk.CodeType = 1
	CodeInvalidPrice      sdk.CodeType = 2
	CodeVoterNotValidator sdk.CodeType = 3
	CodeInvalidVote       sdk.CodeType = 4
)

// ----------------------------------------
// Error constructors

// ErrUnknownDenomination called when the denom is not known
func ErrUnknownDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownDenom, fmt.Sprintf("The denom is not known: %s", denom))
}

// ErrInvalidPrice called when the price submitted is not valid
func ErrInvalidPrice(codespace sdk.CodespaceType, price sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrice, fmt.Sprintf("Price is invalid: %s", price.String()))
}

// ErrVoterNotValidator called when the voter is not a validator
func ErrVoterNotValidator(codespace sdk.CodespaceType, voter sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeVoterNotValidator, fmt.Sprintf("Voter is not a validator: %s", voter.String()))
}

// ErrNoVote called when no vote exists
func ErrNoVote(codespace sdk.CodespaceType, voter sdk.AccAddress, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("No vote exists from %s with denom: %s", voter, denom))
}
