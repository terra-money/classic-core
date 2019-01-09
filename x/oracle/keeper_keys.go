package oracle

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	prefixVote          = []byte("vote")
	prefixTargetPrice   = []byte("target_price")
	prefixObservedPrice = []byte("observed_price")
	KeyWhitelist        = []byte("whitelist")
	KeyDelimiter        = []byte(":")

	ParamStoreKeyParams = []byte("params")
)

// PrefixVote is in format of prefix||denom
func PrefixVote(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixVote, denom))
}

// KeyVote Key is in format of PrefixVote||denom||voter.AccAddress
func KeyVote(denom string, voter sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		prefixVote,
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
