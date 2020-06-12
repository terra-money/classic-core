package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Market errors
var (
	ErrNoEffectivePrice      = sdkerrors.Register(ModuleName, 1, "no price registered with oracle")
	ErrInvalidOfferCoin      = sdkerrors.Register(ModuleName, 2, "invalid offer coin")
	ErrRecursiveSwap         = sdkerrors.Register(ModuleName, 3, "recursive swap")
	ErrNoEffectiveCrossPrice = sdkerrors.Register(ModuleName, 4, "no price registered with the oracle for denom pair")
)
