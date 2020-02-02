package types

import (
	"gopkg.in/yaml.v2"
)

// DenomList is array of denom
type DenomList []string

// String implements fmt.Stringer interface
func (dl DenomList) String() string {
	out, _ := yaml.Marshal(dl)
	return string(out)
}
