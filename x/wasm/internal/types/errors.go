package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Codes for wasm contract errors
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeCreatedFailed       sdk.CodeType = 1
	CodeAccountExists       sdk.CodeType = 2
	CodeInstantiateFailed   sdk.CodeType = 3
	CodeExecuteFailed       sdk.CodeType = 4
	CodeGasLimit            sdk.CodeType = 5
	CodeInvalidGenesis      sdk.CodeType = 6
	CodeNotFound            sdk.CodeType = 7
	CodeInvalidMsg          sdk.CodeType = 8
	CodeNoRegisteredQuerier sdk.CodeType = 9
	CodeNoRegisteredParser  sdk.CodeType = 10
)

// ErrCreateFailed error for wasm code that has already been uploaded or failed
func ErrCreateFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCreatedFailed, fmt.Sprintf("create wasm contract failed: %s", err.Error()))
}

// ErrAccountExists error for a contract account that already exists
func ErrAccountExists(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAccountExists, fmt.Sprintf("contract account %s already exists", addr.String()))
}

// ErrInstantiateFailed error for rust instantiate contract failure
func ErrInstantiateFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInstantiateFailed, fmt.Sprintf("instantiate wasm contract failed: %s", err.Error()))
}

// ErrExecuteFailed error for rust execution contract failure
func ErrExecuteFailed(err error) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeExecuteFailed, fmt.Sprintf("execute wasm contract failed: %s", err.Error()))
}

// ErrGasLimit error for out of gas
func ErrGasLimit(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeGasLimit, fmt.Sprintf("insufficient gas: %s", msg))
}

// ErrInvalidGenesis error for out of gas
func ErrInvalidGenesis(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidGenesis, fmt.Sprintf("invalid genesis: %s", msg))
}

// ErrNotFound error for an entry not found in the store
func ErrNotFound(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNotFound, fmt.Sprintf("not found: %s", msg))
}

// ErrInvalidMsg error when we cannot process the error returned from the contract
func ErrInvalidMsg(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidMsg, fmt.Sprintf("invalid CosmosMsg from the contract: %s", msg))
}

// ErrNoRegisteredQuerier error when we cannot find querier
func ErrNoRegisteredQuerier(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNoRegisteredQuerier, fmt.Sprintf("failed to find querier for route %s", msg))
}

// ErrNoRegisteredParser error when we cannot find msg parser
func ErrNoRegisteredParser(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNoRegisteredParser, fmt.Sprintf("failed to find parser for route %s", msg))
}
