package types

import (
	"strings"
)

// DenomList is array of denom
type DenomList []string

func (dl DenomList) String() (out string) {
	out = strings.Join(dl, "\n")
	return
}
