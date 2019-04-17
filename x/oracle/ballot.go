package oracle

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceBallot is a convinience wrapper arounda a PriceVote slice
type PriceBallot []PriceVote

// Returns the total amount of voting power in the ballot
func (pb PriceBallot) power(ctx sdk.Context, valset sdk.ValidatorSet) sdk.Int {
	totalPower := sdk.ZeroInt()
	for _, vote := range pb {
		votePower, err := vote.getPower(ctx, valset)
		if err == nil {
			totalPower = totalPower.Add(votePower)
		}
	}
	return totalPower
}

// Returns the median weighted by the power of the PriceVote.
func (pb PriceBallot) weightedMedian(ctx sdk.Context, valset sdk.ValidatorSet) sdk.Dec {
	totalPower := pb.power(ctx, valset)
	if pb.Len() > 0 {
		if !sort.IsSorted(pb) {
			sort.Sort(pb)
		}

		pivot := sdk.ZeroInt()
		for _, v := range pb {
			votePower, err := v.getPower(ctx, valset)
			if err != nil {
				continue
			}

			pivot = pivot.Add(votePower)
			if pivot.GTE(totalPower.QuoRaw(2)) {
				return v.Price
			}
		}
	}
	return sdk.ZeroDec()
}

// Len implements sort.Interface
func (pb PriceBallot) Len() int {
	return len(pb)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pb PriceBallot) Less(i, j int) bool {
	return pb[i].Price.LTE(pb[j].Price)
}

// Swap implements sort.Interface.
func (pb PriceBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

// String implements fmt.Stringer interface
func (pb PriceBallot) String() (out string) {
	out = fmt.Sprintf("PriceBallot of %d votes\n", pb.Len())
	for _, pv := range pb {
		out += fmt.Sprintf("\n  %s", pv.String())
	}
	return
}
