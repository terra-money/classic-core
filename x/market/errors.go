package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// market error codes
const (
	DefaultCodespace sdk.CodespaceType = "market"

	CodeInsufficientSwap sdk.CodeType = 1
	CodeNoEffectivePrice sdk.CodeType = 2
	CodeRecursiveSwap    sdk.CodeType = 3
	CodeExceedsSwapLimit sdk.CodeType = 4
)

// ----------------------------------------
// Error constructors

// ErrNoEffectivePrice called when a price for the asset is not registered with the oracle
func ErrNoEffectivePrice(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeNoEffectivePrice, "No price registered with the oracle for asset: "+denom)
}

// ErrInsufficientSwapCoins called when not enough coins are being requested for a swap
func ErrInsufficientSwapCoins(codespace sdk.CodespaceType, rval sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientSwap, "Not enough coins for a swap: "+rval.String())
}

// ErrRecursiveSwap called when Ask and Offer coin denominatioins are equal
func ErrRecursiveSwap(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeRecursiveSwap, "Can't swap tokens with the same denomination: "+denom)
}

// ErrExceedsDailySwapLimit called when the coin swap exceeds the daily swap limit for Luna
func ErrExceedsDailySwapLimit(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeExceedsSwapLimit, "Exceeded the daily swap limit for Luna")
}
