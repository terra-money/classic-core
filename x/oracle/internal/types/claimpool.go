package types

import (
	"fmt"
)

// ClaimPool is a list of Claims
type ClaimPool []Claim

// Sort sorts the ClaimPool
func (cp ClaimPool) Sort() ClaimPool {
	sortBuf := map[string]Claim{}

	for _, claim := range cp {
		addrStr := claim.Recipient.String()
		if val, ok := sortBuf[addrStr]; ok {
			val.Weight = val.Weight + claim.Weight
			sortBuf[addrStr] = val
		} else {
			sortBuf[addrStr] = claim
		}
	}

	i := 0
	cp = make([]Claim, len(sortBuf))
	for _, claim := range sortBuf {
		cp[i] = claim
		i++
	}

	return cp
}

// String implements fmt.Stringer interface
func (cp ClaimPool) String() (out string) {
	out = "ClaimPool "
	for _, claim := range cp {
		out += fmt.Sprintf("\n  %s", claim.String())
	}
	return out
}
