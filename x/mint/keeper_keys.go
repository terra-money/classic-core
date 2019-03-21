package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
var (
	prefixIssuance        = []byte("issuance")
	prefixSeignioragePool = []byte("seigniorage_pool")
)

func keyIssuance(denom string, epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixIssuance, denom, epoch))
}

func keySeignioragePool(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixSeignioragePool, epoch))
}
