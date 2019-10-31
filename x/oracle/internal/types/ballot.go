package types

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExchangeRateBallot is a convinience wrapper around a ExchangeRateVote slice
type ExchangeRateBallot []ExchangeRateVote

// Power returns the total amount of voting power in the ballot
func (pb ExchangeRateBallot) Power(ctx sdk.Context, sk StakingKeeper) int64 {
	totalPower := int64(0)
	for _, vote := range pb {
		totalPower += vote.getPower(ctx, sk)
	}

	return totalPower
}

// WeightedMedian returns the median weighted by the power of the ExchangeRateVote.
func (pb ExchangeRateBallot) WeightedMedian(ctx sdk.Context, sk StakingKeeper) sdk.Dec {
	totalPower := pb.Power(ctx, sk)
	if pb.Len() > 0 {
		if !sort.IsSorted(pb) {
			sort.Sort(pb)
		}

		pivot := int64(0)
		for _, v := range pb {
			votePower := v.getPower(ctx, sk)

			pivot += votePower
			if pivot >= (totalPower / 2) {
				return v.ExchangeRate
			}
		}
	}
	return sdk.ZeroDec()
}

// StandardDeviation returns the standard deviation by the power of the ExchangeRateVote.
func (pb ExchangeRateBallot) StandardDeviation(ctx sdk.Context, sk StakingKeeper) (standardDeviation sdk.Dec) {
	if len(pb) == 0 {
		return sdk.ZeroDec()
	}

	median := pb.WeightedMedian(ctx, sk)

	sum := sdk.ZeroDec()
	for _, v := range pb {
		deviation := v.ExchangeRate.Sub(median)
		sum = sum.Add(deviation.Mul(deviation))
	}

	variance := sum.QuoInt64(int64(len(pb)))

	floatNum, _ := strconv.ParseFloat(variance.String(), 64)
	floatNum = math.Sqrt(floatNum)
	standardDeviation, _ = sdk.NewDecFromStr(fmt.Sprintf("%f", floatNum))

	return
}

// Len implements sort.Interface
func (pb ExchangeRateBallot) Len() int {
	return len(pb)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pb ExchangeRateBallot) Less(i, j int) bool {
	return pb[i].ExchangeRate.LTE(pb[j].ExchangeRate)
}

// Swap implements sort.Interface.
func (pb ExchangeRateBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

// String implements fmt.Stringer interface
func (pb ExchangeRateBallot) String() (out string) {
	out = fmt.Sprintf("ExchangeRateBallot of %d votes\n", pb.Len())
	for _, pv := range pb {
		out += fmt.Sprintf("\n  %s", pv.String())
	}
	return
}
