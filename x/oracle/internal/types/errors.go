package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type codeType = sdk.CodeType

// Oracle error codes
const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownDenom        codeType = 1
	CodeInvalidExchangeRate codeType = 2
	CodeVoterNotValidator   codeType = 3
	CodeInvalidVote         codeType = 4
	CodeNoVotingPermission  codeType = 5
	CodeInvalidHashLength   codeType = 6
	CodeInvalidPrevote      codeType = 7
	CodeVerificationFailed  codeType = 8
	CodeNotRevealPeriod     codeType = 9
	CodeInvalidSaltLength   codeType = 10
	CodeInvalidMsgFormat    codeType = 11
	CodeNoAggregatePrevote  codeType = 12
	CodeNoAggregateVote     codeType = 13
	CodeNoTobinTax          codeType = 14
)

// ----------------------------------------
// Error constructors

// ErrInvalidHashLength called when the given hash has invalid length
func ErrInvalidHashLength(codespace sdk.CodespaceType, hashLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashLength, fmt.Sprintf("The hash length should equal %d but given %d", tmhash.TruncatedSize, hashLength))
}

// ErrUnknownDenomination called when the denom is not known
func ErrUnknownDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownDenom, fmt.Sprintf("The denom is not known: %s", denom))
}

// ErrInvalidExchangeRate called when the rate submitted is not valid
func ErrInvalidExchangeRate(codespace sdk.CodespaceType, rate sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidExchangeRate, fmt.Sprintf("ExchangeRate is invalid: %s", rate.String()))
}

// ErrVerificationFailed called when the given prevote has different hash from the retrieved one
func ErrVerificationFailed(codespace sdk.CodespaceType, hash []byte, retrivedHash []byte) sdk.Error {
	return sdk.NewError(codespace, CodeVerificationFailed, fmt.Sprintf("Retrieved hash [%s] differs from prevote hash [%s]", retrivedHash, hash))
}

// ErrNoPrevote called when no prevote exists
func ErrNoPrevote(codespace sdk.CodespaceType, voter sdk.ValAddress, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrevote, fmt.Sprintf("No prevote exists from %s with denom: %s", voter, denom))
}

// ErrNoVote called when no vote exists
func ErrNoVote(codespace sdk.CodespaceType, voter sdk.ValAddress, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("No vote exists from %s with denom: %s", voter, denom))
}

// ErrNoVotingPermission called when the feeder has no permission to submit a vote for the given operator
func ErrNoVotingPermission(codespace sdk.CodespaceType, feeder sdk.AccAddress, operator sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoVotingPermission, fmt.Sprintf("Feeder %s not permitted to vote on behalf of: %s", feeder.String(), operator.String()))
}

// ErrInvalidRevealPeriod called when the feeder submit rate reveal vote in wrong period.
func ErrInvalidRevealPeriod(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotRevealPeriod, fmt.Sprintf("invalid reveal period."))
}

// ErrInvalidSaltLength called when the salt length is not in 1~4
func ErrInvalidSaltLength(codespace sdk.CodespaceType, saltLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSaltLength, fmt.Sprintf("Salt legnth should be 1~4, but given %d", saltLength))
}

// ErrNoAggregatePrevote called when no prevote exists
func ErrNoAggregatePrevote(codespace sdk.CodespaceType, voter sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoAggregatePrevote, fmt.Sprintf("No aggregate prevote exists from %s", voter))
}

// ErrNoAggregateVote called when no prevote exists
func ErrNoAggregateVote(codespace sdk.CodespaceType, voter sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoAggregateVote, fmt.Sprintf("No aggregate vote exists from %s", voter))
}

// ErrNoTobinTax called when no tobin tax exists for the given denom
func ErrNoTobinTax(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeNoAggregateVote, fmt.Sprintf("No tobin tax exists for %s", denom))
}
