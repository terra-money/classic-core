package oracle

import (
	"fmt"
	"math"
	"sort"
	"terra/types"

	"gonum.org/v1/gonum/stat"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceBallot is a convinience wrapper arounda a PriceVote slice
type PriceBallot []PriceVote

// TotalPower gets the total amount of voting power in the ballot
func (pb PriceBallot) TotalPower() sdk.Int {
	totalPower := sdk.ZeroInt()
	for _, vote := range pb {
		totalPower = totalPower.Add(vote.Power)
	}
	return totalPower
}

// Returns the median weighted by the Power of the PriceVote.
func (pb PriceBallot) weightedMedian() sdk.Dec {
	totalPower := pb.TotalPower()
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

// Calculates the median and returns the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median.
func (pb PriceBallot) tally() (weightedMedian sdk.Dec, ballotWinners types.ClaimPool) {
	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	ballotWinners = types.ClaimPool{}
	weightedMedian = pb.weightedMedian()

	maxSpread := weightedMedian.Mul(sdk.NewDecWithPrec(1, 2)) // 1%
	stdDev := pb.stdDev()

	if stdDev.LT(maxSpread) {
		maxSpread = stdDev
	}

	for _, vote := range pb {
		if vote.Price.GTE(weightedMedian.Sub(maxSpread)) && vote.Price.LTE(weightedMedian.Add(maxSpread)) {
			ballotWinners = append(ballotWinners, types.Claim{
				Recipient: vote.Voter,
				Weight:    vote.Power,
				Class:     types.OracleClaimClass,
			})
		}
	}

	return
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
	out = fmt.Sprintf("PriceBallot of %d votes with %s total power\n", pb.Len(), pb.TotalPower())
	for _, pv := range pb {
		out += fmt.Sprintf("\n  %s", pv.String())
	}
	return
}
