package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/treasury module sentinel errors
//
var (
	ErrInvalidEpoch = sdkerrors.Register(ModuleName, 1, "The query epoch should be between [0, current epoch]")
)
