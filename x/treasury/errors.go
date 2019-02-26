package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "treasury"

	CodeWrongTaxDenom        sdk.CodeType = 1
	CodeExcessiveWeight      sdk.CodeType = 2
	CodeClaimIDConflict      sdk.CodeType = 3
	CodeInsufficientIssuance sdk.CodeType = 4
)

// ----------------------------------------
// Error constructors

// // ErrWrongTaxDenomination called when wrong denom coin used to pay taxes to the treasury
// func ErrWrongTaxDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
// 	return sdk.NewError(codespace, CodeWrongTaxDenom, denom)
// }

// // ErrExcessiveWeight called when a claim is added with an excessive weight (total > 1)
// func ErrExcessiveWeight(codespace sdk.CodespaceType, weight sdk.Dec) sdk.Error {
// 	return sdk.NewError(codespace, CodeExcessiveWeight, weight.String())
// }

// // ErrClaimIDConflict called when a claim being added has an id conflict with an existing claim
// func ErrClaimIDConflict(codespace sdk.CodespaceType, claimID string) sdk.Error {
// 	return sdk.NewError(codespace, CodeClaimIDConflict, claimID)
// }

// // ErrInsufficientIssuance called when the issuance of an asset is insufficient to meet a contraction
// func ErrInsufficientIssuance(codespace sdk.CodespaceType, denom string) sdk.Error {
// 	return sdk.NewError(codespace, CodeInsufficientIssuance, denom)
// }
