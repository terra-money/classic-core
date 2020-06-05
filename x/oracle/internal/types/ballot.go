package types

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VoteForTally is a convinience wrapper to reduct redundant lookup cost
type VoteForTally struct {
	ExchangeRateVote
	Power int64
}

// NewVoteForTally returns a new VoteForTally instance
func NewVoteForTally(vote ExchangeRateVote, power int64) VoteForTally {
	return VoteForTally{
		vote,
		power,
	}
}

type CrossExchangeRate struct {
	Denom1            string  `json:"denom1"`              // Ticker name of target first terra currency
	Denom2            string  `json:"denom2"`              // Ticker name of target second terra currency
	CrossExchangeRate sdk.Dec `json:"cross_exchange_rate"` // CrossExchangeRate of Luna in target fiat currency
}

func GetDenomOrderAsc(denom1, denom2 string) (string, string) {
	if denom1 > denom2 {
		return denom2, denom1
	}
	return denom1, denom2
}

func NewCrossExchangeRate(denom1, denom2 string, exchangeRate sdk.Dec) CrossExchangeRate {
	// swap ascending order for deterministic kv indexing
	denom1, denom2 = GetDenomOrderAsc(denom1, denom2)
	return CrossExchangeRate{
		denom1,
		denom2,
		exchangeRate,
	}
}

func (cer CrossExchangeRate) DenomPair() string {
	return cer.Denom1 + "_" + cer.Denom2
}

// String implements fmt.Stringer interface
func (cer CrossExchangeRate) String() string {
	return fmt.Sprintf(`CrossExchangeRate
	Denom1:             %s, 
	Denom2:             %s, 
	CrossExchangeRate:  %s`,
		cer.Denom1, cer.Denom2, cer.CrossExchangeRate)
}

type CrossExchangeRates []CrossExchangeRate

// String implements fmt.Stringer interface
func (v CrossExchangeRates) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ExchangeRateBallot is a convinience wrapper around a ExchangeRateVote slice
type ExchangeRateBallot []VoteForTally

// Power returns the total amount of voting power in the ballot
func (pb ExchangeRateBallot) Power() int64 {
	totalPower := int64(0)
	for _, vote := range pb {
		totalPower += vote.Power
	}

	return totalPower
}

// WeightedMedian returns the median weighted by the power of the ExchangeRateVote.
func (pb ExchangeRateBallot) WeightedMedian() sdk.Dec {
	totalPower := pb.Power()
	if pb.Len() > 0 {
		if !sort.IsSorted(pb) {
			sort.Sort(pb)
		}

		pivot := int64(0)
		for _, v := range pb {
			votePower := v.Power

			pivot += votePower
			if pivot >= (totalPower / 2) {
				return v.ExchangeRate
			}
		}
	}
	return sdk.ZeroDec()
}

// StandardDeviation returns the standard deviation by the power of the ExchangeRateVote.
func (pb ExchangeRateBallot) StandardDeviation() (standardDeviation sdk.Dec) {
	if len(pb) == 0 {
		return sdk.ZeroDec()
	}

	median := pb.WeightedMedian()

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