package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type codeType = sdk.CodeType

// market error codes
const (
	DefaultCodespace sdk.CodespaceType = "market"

	CodeInvalidOfferCoin codeType = 1
	CodeNoEffectivePrice codeType = 2
	CodeRecursiveSwap    codeType = 3
)

// ----------------------------------------
// Error constructors

// ErrNoEffectivePrice called when a price for the asset is not registered with the oracle
func ErrNoEffectivePrice(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeNoEffectivePrice, "No price registered with the oracle for asset: "+denom)
}

// ErrInvalidOfferCoin called when not enough or too huge coins are being requested for a swap
func ErrInvalidOfferCoin(codespace sdk.CodespaceType, rval sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOfferCoin, "Invalid offer coin for a swap: "+rval.String())
}

// ErrRecursiveSwap called when Ask and Offer coin denominatioins are equal
func ErrRecursiveSwap(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeRecursiveSwap, "Can't swap tokens with the same denomination: "+denom)
}
