package oracle

import (
	"fmt"
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//-------------------------------------------------
//-------------------------------------------------

// PriceVote - struct to store a validator's vote on the price
type PriceVote struct {
	Price sdk.Dec        `json:"price"` // Price of Luna in target fiat currency
	Denom string         `json:"denom"` // Ticker name of target fiat currency
	Power sdk.Int        `json:"power"` // Total bonded tokens of validator
	Voter sdk.AccAddress `json:"voter"` // account address of validator
}

// NewPriceVote creates a PriceVote instance
func NewPriceVote(price sdk.Dec, denom string, power sdk.Int, voter sdk.AccAddress) PriceVote {
	return PriceVote{
		Price: price,
		Denom: denom,
		Power: power,
		Voter: voter,
	}
}

func (pv PriceVote) String() string {
	return fmt.Sprintf("PriceVote{ Denom %s, Voter %s, Power %v, Price %v }", pv.Denom, pv.Voter, pv.Power, pv.Price)
}

type PriceBallot []PriceVote

// String implements fmt.Stringer interface
func (pb PriceBallot) String() (out string) {
	for _, pv := range pb {
		out += fmt.Sprintf("\n  %s", pv.String())
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

func (pb PriceBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

func (pb PriceBallot) totalPower() sdk.Int {
	sumWeight := sdk.ZeroInt()
	for _, v := range pb {
		sumWeight = sumWeight.Add(v.Power)
	}
	return sumWeight
}

func (pb PriceBallot) allVoters() []sdk.AccAddress {
	allVoters := []sdk.AccAddress{}
	for _, v := range pb {
		allVoters = append(allVoters, v.Voter)
	}
	return allVoters
}

func (pb PriceBallot) weightedMedian() (i int64, mod PriceVote) {
	if len(pb) == 0 {
		mod = PriceVote{}
		return
	}

	if !sort.IsSorted(pb) {
		sort.Sort(pb)
	}

	voterTotalPower := pb.totalPower()
	sumWeight := sdk.ZeroInt()
	for _, v := range pb {
		if sumWeight.GTE(voterTotalPower.Div(sdk.NewInt(2))) {
			break
		}

		i++
		sumWeight = sumWeight.Add(v.Power)
		mod = v
	}

	return
}

func (pb PriceBallot) mean() (mean sdk.Dec) {
	if len(pb) == 0 {
		return sdk.ZeroDec()
	}

	sumPrice := sdk.ZeroDec()
	i := 0
	for _, v := range pb {
		i++
		sumPrice = sumPrice.Add(v.Price)
	}

	mean = sumPrice.QuoInt(sdk.NewInt(int64(i)))
	return
}

func (pb PriceBallot) stdDev() (stdDev sdk.Dec) {
	mean := pb.mean()
	sumCalc := sdk.ZeroDec()
	for _, v := range pb {
		spread := v.Price.Sub(mean)
		sumCalc = sumCalc.Add(spread.Mul(spread))
	}

	temp := sumCalc.Sqrt(big.NewInt(2))
	stdDev = sdk.NewDecFromBigInt(temp)
	return
}

func (pb PriceBallot) tally() (median sdk.Dec, rewardees []PriceVote) {
	sort.Sort(pb)

	_, modStub := pb.weightedMedian()
	median = modStub.Price

	stdDev := pb.stdDev()

	for _, stub := range pb {
		if stub.Price.GTE(median.Sub(stdDev)) || stub.Price.LTE(median.Add(stdDev)) {
			rewardees = append(rewardees, stub)
		}
	}

	return
}
