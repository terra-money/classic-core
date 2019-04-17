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

func keyIssuance(denom string, day sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixIssuance, denom, day))
}

func keySeignioragePool(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixSeignioragePool, epoch))
}
