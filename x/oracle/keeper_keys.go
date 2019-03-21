package oracle

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	prefixVote          = []byte("vote")
	prefixPrice         = []byte("price")
	prefixDropCounter   = []byte("drop")
	keyDelimiter        = []byte(":")
	paramStoreKeyParams = []byte("params")
)

func keyVote(denom string, voter sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		prefixVote,
		[]byte(denom),
		voter.Bytes(),
	}, keyDelimiter)
}

func keyPrice(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixPrice, denom))
}

func keyDropCounter(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixDropCounter, denom))
}

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
