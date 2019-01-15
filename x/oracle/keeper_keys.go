package oracle

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	prefixTargetVote    = []byte("target_vote")
	prefixObservedVote  = []byte("observed_vote")
	prefixTargetPrice   = []byte("target_price")
	prefixObservedPrice = []byte("observed_price")
	KeyWhitelist        = []byte("whitelist")
	KeyDelimiter        = []byte(":")

	ParamStoreKeyParams = []byte("params")
)

// PrefixTargetVote is in format of prefix||denom
func PrefixTargetVote(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTargetVote, denom))
}

// KeyTargetVote Key is in format of PrefixVote||denom||voter.AccAddress
func KeyTargetVote(denom string, voter sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		prefixTargetVote,
		[]byte(denom),
		voter.Bytes(),
	}, KeyDelimiter)
}

// prefixObservedVote is in format of prefix||denom
func PrefixObservedVote(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixObservedVote, denom))
}

// KeyObservedVote Key is in format of PrefixVote||denom||voter.AccAddress
func KeyObservedVote(denom string, voter sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		prefixObservedVote,
		[]byte(denom),
		voter.Bytes(),
	}, KeyDelimiter)
}

// KeyTargetPrice is in format of PrefixTargetPrice||denom
func KeyTargetPrice(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTargetPrice, denom))
}

// KeyObservedPrice is in format of PrefixObservedPrice||denom
func KeyObservedPrice(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixObservedPrice, denom))
}

// ParamTypeTable for oracle module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable(
		ParamStoreKeyParams, Params{},
	)
}

const (
	// default paramspace for params keeper
	DefaultParamspace = "oracle"
)
