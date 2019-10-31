package types

// DenomList is array of denom
type DenomList map[string]bool

// String implements fmt.Stringer interface
func (dl DenomList) String() (out string) {
	out = dl.String()
	return
}
