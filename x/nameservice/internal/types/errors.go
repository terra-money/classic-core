package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type codeType = sdk.CodeType

// Oracle error codes
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidName                 codeType = 1
	CodeInvalidSaltLength           codeType = 2
	CodeInvalidHashLength           codeType = 3
	CodeInvalidNameLength           codeType = 4
	CodeVerificationFailed          codeType = 5
	CodeAuctionNotExists            codeType = 6
	CodeBidNotExists                codeType = 7
	CodeRevealNotExists             codeType = 8
	CodeRegistryNotExists           codeType = 9
	CodeResolveNotExists            codeType = 10
	CodeReverseResolveNotExists     codeType = 11
	CodeInvalidRootName             codeType = 12
	CodeNameAlreadyTaken            codeType = 13
	CodeAuctionExists               codeType = 14
	CodeAuctionNotBidStatus         codeType = 15
	CodeBidAlreadyExists            codeType = 16
	CodeAuctionNotRevealStatus      codeType = 17
	CodeRevealAlreadyExists         codeType = 18
	CodeAddressAlreadyRegistered    codeType = 19
	CodeDepositSmallerThanBidAmount codeType = 20
	CodeLowDeposit                  codeType = 21
)

// ----------------------------------------
// Error constructors

// ErrInvalidHashLength called when the given hash has invalid length
func ErrInvalidName(codespace sdk.CodespaceType, name Name, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidName, fmt.Sprintf("Failed to validate name %s: %s", name, msg))
}

// ErrInvalidHashLength called when the given hash has invalid length
func ErrInvalidHashLength(codespace sdk.CodespaceType, hashLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashLength, fmt.Sprintf("The hash bytes length should equal %d but given %d", tmhash.TruncatedSize, hashLength))
}

// ErrVerificationFailed called when the given prevote has different hash from the retrieved one
func ErrVerificationFailed(codespace sdk.CodespaceType, hash, retrievedHash []byte) sdk.Error {
	return sdk.NewError(codespace, CodeVerificationFailed, fmt.Sprintf("Retrieved hash [%s] differs from prevote hash [%s]", retrievedHash, hash))
}

// ErrInvalidSaltLength called when the salt length is not in 1~4
func ErrInvalidSaltLength(codespace sdk.CodespaceType, saltLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSaltLength, fmt.Sprintf("Salt legnth should be 1~4, but given %d", saltLength))
}

// ErrInvalidNameLength called when the name length is smaller than params.MinNameLength
func ErrInvalidNameLength(codespace sdk.CodespaceType, minNameLength, nameLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidNameLength, fmt.Sprintf("Name legnth should be bigger than or equal to %d, but given %d", minNameLength, nameLength))
}

// ErrAuctionNotExists called when the auction is not exists
func ErrAuctionNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAuctionNotExists, "There is no auction for the given name")
}

// ErrBidNotExists called when the bid is not exists
func ErrBidNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBidNotExists, "There is no bid for then given name and address")
}

// ErrRevealNotExists called when the reveal is not exists
func ErrRevealNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRevealNotExists, "There is no reveal for then given name and address")
}

// ErrRegistryNotExists called when the reveal is not exists
func ErrRegistryNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRegistryNotExists, "There is no registry for then given name")
}

// ErrResolveNotExists called when the resolve entry is not exists
func ErrResolveNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeResolveNotExists, "There is no resolve entry for then given name")
}

// ErrReverseResolveNotExists called when the reverse resolve entry is not exists
func ErrReverseResolveNotExists(codespace sdk.CodespaceType, address sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeReverseResolveNotExists, fmt.Sprintf("There is no reverse resolve entry for %s", address))
}

// ErrInvalidRootName called when the given root name is different form the root name of param
func ErrInvalidRootName(codespace sdk.CodespaceType, paramRootName, rootName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRootName, fmt.Sprintf("The root name should be %s not %s", paramRootName, rootName))
}

// ErrNameAlreadyTaken called when the given parent name is already taken by registry
func ErrNameAlreadyTaken(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNameAlreadyTaken, fmt.Sprintf("The name already taken"))
}

// ErrAuctionExists called when the auction is exists for the given parent name
func ErrAuctionExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAuctionExists, fmt.Sprintf("The name auction exists"))
}

// ErrAuctionNotBidStatus called when the auction is not in bid status
func ErrAuctionNotBidStatus(codespace sdk.CodespaceType, status AuctionStatus) sdk.Error {
	return sdk.NewError(codespace, CodeAuctionNotBidStatus, fmt.Sprintf("The auction is not bid status, but %s status", status))
}

// ErrBidAlreadyExists called when the auction is not in bid status
func ErrBidAlreadyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBidAlreadyExists, fmt.Sprintf("Bid is already exists"))
}

// ErrAuctionNotRevealStatus called when the auction is not in bid status
func ErrAuctionNotRevealStatus(codespace sdk.CodespaceType, status AuctionStatus) sdk.Error {
	return sdk.NewError(codespace, CodeAuctionNotRevealStatus, fmt.Sprintf("The auction is not reveal status, but %s status", status))
}

// ErrRevealAlreadyExists called when the auction is not in bid status
func ErrRevealAlreadyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeRevealAlreadyExists, fmt.Sprintf("Reveal is already exists"))
}

// ErrAddressAlreadyRegistered called when the address is already registered in nameservice module
func ErrAddressAlreadyRegistered(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAddressAlreadyRegistered, fmt.Sprintf("Address is already registered"))
}

// ErrDepositSmallerThanBidAmount called when the deposit is smaller than bid amount
func ErrDepositSmallerThanBidAmount(codespace sdk.CodespaceType, deposit sdk.Coin) sdk.Error {
	return sdk.NewError(codespace, CodeDepositSmallerThanBidAmount, fmt.Sprintf("Bid amount must be greater than or equal to deposit %s", deposit))
}

// ErrLowDeposit called when the deposit is smaller than min deposit params
func ErrLowDeposit(codespace sdk.CodespaceType, minDeposit sdk.Coin) sdk.Error {
	return sdk.NewError(codespace, CodeLowDeposit, fmt.Sprintf("Deposit can not be smaller than %s", minDeposit))
}
