package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Market errors
var (
	ErrRecursiveSwap    = sdkerrors.Register(ModuleName, 2, "recursive swap")
	ErrNoEffectivePrice = sdkerrors.Register(ModuleName, 3, "no price registered with oracle")
	ErrEmptyChanges     = sdkerrors.Register(ModuleName, 4, "submitted route changes are empty")
	ErrEmptyAddress     = sdkerrors.Register(ModuleName, 5, "route address is empty")
	ErrDuplicateRoute   = sdkerrors.Register(ModuleName, 6, "routes have duplicated daddress")
	ErrZeroWeight       = sdkerrors.Register(ModuleName, 7, "route weight is zero")
	ErrInvalidWeightSum = sdkerrors.Register(ModuleName, 8, "route weight sum exceeds one")
)
