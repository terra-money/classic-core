package oracle

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceBallot is a convinience wrapper arounda a PriceVote slice
type PriceBallot []PriceVote

// // TotalPower gets the total amount of voting power in the ballot
// func (pb PriceBallot) TotalPower() sdk.Int {
// 	totalPower := sdk.ZeroInt()
// 	for _, vote := range pb {
// 		totalPower = totalPower.Add(vote.Power)
// 	}
// 	return totalPower
// }

// Returns the median weighted by the Power of the PriceVote.
func (pb PriceBallot) weightedMedian(totalPower sdk.Int) sdk.Dec {
	if pb.Len() > 0 {
		if !sort.IsSorted(pb) {
			sort.Sort(pb)
		}

		pivot := sdk.ZeroInt()
		for _, v := range pb {
			pivot = pivot.Add(v.Power)

			if pivot.GTE(totalPower.DivRaw(2)) {
				return v.Price
			}
		}
	}

	return sdk.ZeroDec()
}

// Computes the mean (in price) of the ballot
func (pb PriceBallot) mean() sdk.Dec {
	if pb.Len() > 0 {
		sumPrice := sdk.ZeroDec()
		for _, v := range pb {
			sumPrice = sumPrice.Add(v.Price)
		}

		return sumPrice.QuoInt64(int64(pb.Len()))
	}
	return sdk.ZeroDec()
}

const precision = 4

// Computes the stdDev (in price) of the ballot
func (pb PriceBallot) stdDev() sdk.Dec {
	if pb.Len() > 0 {
		x := []float64{}
		weights := []float64{}
		base := math.Pow10(precision)

		for _, v := range pb {
			x = append(x, float64(v.Price.MulInt64(int64(base)).TruncateInt64())/base)
			weights = append(weights, float64(v.Power.Int64()))
		}

		stdDevFlt := stat.StdDev(x, weights)

		return sdk.NewDecWithPrec(int64(stdDevFlt*base), precision)
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
