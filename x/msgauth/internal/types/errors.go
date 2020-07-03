package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrInvalidPeriod  = sdkerrors.Register(ModuleName, 3, "period of authorization should be positive time duration")
	ErrInvalidMsgType = sdkerrors.Register(ModuleName, 4, "given msg type is not grantable")
)
