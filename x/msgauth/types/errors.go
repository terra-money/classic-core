package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/msgauth module sentinel errors
var (
	ErrInvalidPeriod        = sdkerrors.Register(ModuleName, 2, "period of authorization should be positive time duration")
	ErrInvalidMsgType       = sdkerrors.Register(ModuleName, 3, "given msg type is not grantable")
	ErrInvalidAuthorization = sdkerrors.Register(ModuleName, 4, "given authorization is not valid")
	ErrInvalidMsg           = sdkerrors.Register(ModuleName, 5, "given msg is not valid")
	ErrGrantExists          = sdkerrors.Register(ModuleName, 6, "grant already exists")
)
