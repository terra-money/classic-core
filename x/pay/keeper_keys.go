package pay

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//nolint
var (
	KeyTaxRate        = []byte("tax_rate")
	PrefixTaxProceeds = []byte("tax_proceeds")
	PrefixTaxCap      = []byte("tax_cap")
	PrefixIssuance    = []byte("issuance")
)

// KeyTaxProceeds is in format of PrefixTaxProceeds:denom
func KeyTaxProceeds(epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixTaxProceeds, epoch))
}

// KeyTaxCap is in format of PrefixTaxCap:denom
func KeyTaxCap(denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s", PrefixTaxCap, denom))
}

// KeyIssuance is in format of PrefixIssuance:denom
func KeyIssuance(denom string, epoch sdk.Int) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", PrefixIssuance, denom, epoch))
}
