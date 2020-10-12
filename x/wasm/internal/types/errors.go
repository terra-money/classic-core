package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Wasm Errors
var (
	ErrAccountExists       = sdkerrors.Register(ModuleName, 2, "contract account already exists")
	ErrInstantiateFailed   = sdkerrors.Register(ModuleName, 3, "instantiate wasm contract failed")
	ErrExecuteFailed       = sdkerrors.Register(ModuleName, 4, "execute wasm contract failed")
	ErrGasLimit            = sdkerrors.Register(ModuleName, 5, "insufficient gas")
	ErrInvalidGenesis      = sdkerrors.Register(ModuleName, 6, "invalid genesis")
	ErrNotFound            = sdkerrors.Register(ModuleName, 7, "not found")
	ErrInvalidMsg          = sdkerrors.Register(ModuleName, 8, "invalid Msg from the contract")
	ErrNoRegisteredQuerier = sdkerrors.Register(ModuleName, 9, "failed to find querier for route")
	ErrNoRegisteredParser  = sdkerrors.Register(ModuleName, 10, "failed to find parser for route")
	ErrMigrationFailed     = sdkerrors.Register(ModuleName, 11, "migrate wasm contract failed")
	ErrNotMigratable       = sdkerrors.Register(ModuleName, 12, "the contract is not migratable ")
	ErrStoreCodeFailed     = sdkerrors.Register(ModuleName, 13, "store wasm contract failed")
)
