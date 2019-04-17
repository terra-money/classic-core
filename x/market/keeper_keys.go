package market

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	paramStoreKeyParams = []byte("params")
)

func paramKeyTable() params.KeyTable {
	return params.NewKeyTable(
		paramStoreKeyParams, Params{},
	)
}
