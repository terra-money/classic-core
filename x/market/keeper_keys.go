package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
var (
	PrefixSeigniorage = []byte("seigniorage")
)

// KeySeigniorage is in format of PrefixTaxProceeds:denom
func KeySeigniorage(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixSeigniorage, epoch))
}
