package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Denom is the object to hold configurations of each denom
type Denom struct {
	Name     string  `json:"name" yaml:"name"`
	TobinTax sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`
}

// String implements fmt.Stringer interface
func (d Denom) String() string {
	return fmt.Sprintf(`
Name:           %s
TobinTax:       %s
`, d.Name, d.TobinTax)
}

// DenomList is array of Denom
type DenomList []Denom

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	for _, d := range dl {
		out += d.String() + "\n"
	}
	return strings.TrimSpace(out)
}
