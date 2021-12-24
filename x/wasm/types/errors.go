package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Wasm Errors
var (
	ErrAccountExists             = sdkerrors.Register(ModuleName, 2, "contract account already exists")
	ErrInstantiateFailed         = sdkerrors.Register(ModuleName, 3, "instantiate wasm contract failed")
	ErrExecuteFailed             = sdkerrors.Register(ModuleName, 4, "execute wasm contract failed")
	ErrGasLimit                  = sdkerrors.Register(ModuleName, 5, "insufficient gas")
	ErrInvalidGenesis            = sdkerrors.Register(ModuleName, 6, "invalid genesis")
	ErrNotFound                  = sdkerrors.Register(ModuleName, 7, "not found")
	ErrInvalidMsg                = sdkerrors.Register(ModuleName, 8, "invalid Msg from the contract")
	ErrNoRegisteredQuerier       = sdkerrors.Register(ModuleName, 9, "failed to find querier for route")
	ErrNoRegisteredParser        = sdkerrors.Register(ModuleName, 10, "failed to find parser for route")
	ErrMigrationFailed           = sdkerrors.Register(ModuleName, 11, "migrate wasm contract failed")
	ErrNotMigratable             = sdkerrors.Register(ModuleName, 12, "the contract is not migratable ")
	ErrStoreCodeFailed           = sdkerrors.Register(ModuleName, 13, "store wasm contract failed")
	ErrContractQueryFailed       = sdkerrors.Register(ModuleName, 14, "contract query failed")
	ErrExceedMaxContractSize     = sdkerrors.Register(ModuleName, 15, "exceeds max contract size limit")
	ErrExceedMaxContractMsgSize  = sdkerrors.Register(ModuleName, 16, "exceeds max contract msg size limit")
	ErrExceedMaxContractDataSize = sdkerrors.Register(ModuleName, 17, "exceeds max contract data size limit")
	ErrReplyFailed               = sdkerrors.Register(ModuleName, 18, "reply wasm contract failed")
	ErrExceedMaxQueryDepth       = sdkerrors.Register(ModuleName, 19, "exceed max query depth")
	ErrEmpty                     = sdkerrors.Register(ModuleName, 20, "empty")
	ErrUnsupportedForContract    = sdkerrors.Register(ModuleName, 21, "unsupported for this contract")
	ErrInvalid                   = sdkerrors.Register(ModuleName, 22, "invalid")
)
