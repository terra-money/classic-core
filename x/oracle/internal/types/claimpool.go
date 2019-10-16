package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Weight    int64          `json:"weight"`
	Recipient sdk.ValAddress `json:"recipient"`
}

// NewClaim generates a Claim instance.
func NewClaim(weight int64, recipient sdk.ValAddress) Claim {
	return Claim{
		Weight:    weight,
		Recipient: recipient,
	}
}

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
