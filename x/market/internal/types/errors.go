package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Market errors
var (
	ErrInvalidOfferCoin = sdkerrors.Register(ModuleName, 2, "invalid offer coin")
	ErrRecursiveSwap    = sdkerrors.Register(ModuleName, 3, "recursive swap")
	ErrNoEffectivePrice = sdkerrors.Register(ModuleName, 4, "no price registered with oracle")
)
