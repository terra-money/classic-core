package pay

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	keyTaxRate        = []byte("tax_rate")
	prefixTaxProceeds = []byte("tax_proceeds")
	prefixTaxCap      = []byte("tax_cap")
	prefixIssuance    = []byte("issuance")
)

func keyTaxProceeds(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxProceeds, epoch))
}

func keyTaxCap(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", prefixTaxCap, denom))
}

func keyIssuance(denom string, epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", prefixIssuance, denom, epoch))
}
