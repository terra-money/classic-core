package types

// ClaimPool is a list of Claims
type ClaimPool []Claim

// Sort sorts the ClaimPool
func (cp ClaimPool) Sort() ClaimPool {
	sortBuf := map[string]Claim{}

	for _, claim := range cp {
		addrStr := claim.Recipient.String()
		if val, ok := sortBuf[addrStr]; ok {
			val.Weight = val.Weight.Add(claim.Weight)
			sortBuf[addrStr] = val
		} else {
			sortBuf[addrStr] = claim
		}
	}

	cp = make([]Claim, len(sortBuf))

	for _, claim := range sortBuf {
		cp = append(cp, claim)
	}

	return cp
}
