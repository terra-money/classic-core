package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Market errors
var (
	ErrRecursiveSwap     = sdkerrors.Register(ModuleName, 2, "recursive swap")
	ErrNoEffectivePrice  = sdkerrors.Register(ModuleName, 3, "no price registered with oracle")
	ErrInvalidAddress    = sdkerrors.Register(ModuleName, 4, "route address is invalid")
	ErrDuplicateRoute    = sdkerrors.Register(ModuleName, 5, "routes have duplicated daddress")
	ErrInvalidWeight     = sdkerrors.Register(ModuleName, 6, "route weight is zero or negative value")
	ErrInvalidWeightsSum = sdkerrors.Register(ModuleName, 7, "route weights sum exceeds one")
)
