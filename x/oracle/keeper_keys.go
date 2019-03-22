package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

//nolint
const (
	DefaultParamspace = "oracle"
)

var (
	prefixVote          = []byte("vote")
	prefixPrice         = []byte("price")
	prefixDropCounter   = []byte("drop")
	paramStoreKeyParams = []byte("params")
)

func keyVote(denom string, voter sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixVote, denom, voter))
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
