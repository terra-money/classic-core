package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// market errors
const (
	DefaultCodespace sdk.CodespaceType = "market"

	CodeInsufficientSwap sdk.CodeType = 1
	CodeUnknownDenom     sdk.CodeType = 2
	CodeNoEffectivePrice sdk.CodeType = 3
	CodeRecursiveSwap    sdk.CodeType = 4
)

// ----------------------------------------
// Error constructors

// ErrUnknownDenomination called when the denom is not whitelisted by the oracle
func ErrUnknownDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownDenom, "Unknown denom in SwapMsg: "+denom)
}

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
