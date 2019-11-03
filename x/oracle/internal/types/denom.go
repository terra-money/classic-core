package types

import (
	"strings"
)

// DenomList is array of denom
type DenomList []string

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	strings.Join(dl, "\n")
	return
}
