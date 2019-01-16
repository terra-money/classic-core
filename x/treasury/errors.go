package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle errors reserve 1101-1199
const (
	DefaultCodespace sdk.CodespaceType = "treasury"

	CodeWrongTaxDenom   sdk.CodeType = 1
	CodeExcessiveWeight sdk.CodeType = 2
	CodeNoShareFound    sdk.CodeType = 3
)

// ----------------------------------------
// Error constructors

// ErrWrongTaxDenomination called when wrong denom coin used to pay taxes to the treasury
func ErrWrongTaxDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeWrongTaxDenom, denom)
}

// ErrExcessiveWeight called when a claim is added with an excessive weight (total > 1)
func ErrExcessiveWeight(codespace sdk.CodespaceType, weight sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeExcessiveWeight, weight.String())
}

func ErrNoShareFound(codespace sdk.CodespaceType, shareID string) sdk.Error {
	return sdk.NewError(codespace, CodeNoShareFound, shareID)
}
