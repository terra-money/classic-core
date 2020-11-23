package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
	"strings"
)

// Denom is the object to hold configurations of each denom
type Denom struct {
	Name     string  `json:"name" yaml:"name"`
	TobinTax sdk.Dec `json:"tobin_tax" yaml:"tobin_tax"`
}

// String implements fmt.Stringer interface
func (d Denom) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
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
