package oracle

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
var (
	PrefixVote        = []byte("vote")
	PrefixPrice       = []byte("price")
	PrefixDropCounter = []byte("drop")
	KeyDelimiter      = []byte(":")

	ParamStoreKeyParams = []byte("params")
)

// KeyVote Key is in format of PrefixVote||denom||voter.AccAddress
func KeyVote(denom string, voter sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		PrefixVote,
		[]byte(denom),
		voter.Bytes(),
	}, KeyDelimiter)
}

// KeyPrice is in format of PrefixPrice||denom
func KeyPrice(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixPrice, denom))
}

// KeyDropCounter is in format of PrefixDropCounter||denom
func KeyDropCounter(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixDropCounter, denom))
}

// ParamKeyTable for oracle module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyParams, Params{},
	)
}

const (
	// default paramspace for params keeper
	DefaultParamspace = "oracle"
)
