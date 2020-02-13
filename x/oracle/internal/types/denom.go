package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Denom is the object to hold configurations of each denom
type Denom struct {
	Name           string  `json:"name" yaml:"name"`
	IlliquidFactor sdk.Dec `json:"illiquid_factor" yaml:"illiquid_factor"`
}

// String implements fmt.Stringer interface
func (d Denom) String() string {
	return fmt.Sprintf(`
Name:           %s
IlliquidFactor: %s
`, d.Name, d.IlliquidFactor)
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
