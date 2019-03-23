package oracle

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"gonum.org/v1/gonum/stat"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func generateRandomTestCase() (prices []float64, weights []float64) {
	rand.Seed(int64(time.Now().Nanosecond()))
	numInputs := 10 + (rand.Int() % 100)
	for i := 0; i < numInputs; i++ {

		// Cut off at precision 4
		price := float64(int(rand.Float64()*10000)) / 10000

		// Wrap int in float
		weight := float64(rand.Int63())

		prices = append(prices, price)
		weights = append(weights, weight)
	}

	return
}

func TestPBStdDev(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	prices, weights := generateRandomTestCase()
	pb := PriceBallot{}
	for i, price := range prices {
		weight := sdk.NewDec(int64(weights[i])).TruncateInt()
		vote := NewPriceVote(sdk.NewDecWithPrec(int64(price*10000), 4), "", weight, addrs[0])
		pb = append(pb, vote)
	}

	statAnswerRaw := stat.StdDev(prices, weights)
	statAnswerRounded := float64(int64(statAnswerRaw*math.Pow10(precision))) / math.Pow10(precision)
	ballotAnswerRounded := float64(pb.stdDev().MulInt64(int64(math.Pow10(precision))).TruncateInt64()) / math.Pow10(precision)

	require.Equal(t, statAnswerRounded, ballotAnswerRounded)
}

func TestPBMean(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	prices, _ := generateRandomTestCase()
	weights := []float64{}
	pb := PriceBallot{}
	for _, price := range prices {
		vote := NewPriceVote(sdk.NewDecWithPrec(int64(price*10000), 4), "", sdk.OneInt(), addrs[0])
		weights = append(weights, 1.0)
		pb = append(pb, vote)
	}

	statAnswerRounded := float64(int64(stat.Mean(prices, weights)*10000)) / 10000
	ballotAnswerRounded := float64(pb.mean().MulTruncate(sdk.NewDec(10000)).TruncateInt64()) / 10000
	require.Equal(t, statAnswerRounded, ballotAnswerRounded)
}

func TestPBWeightedMedian(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		inputs  []float64
		weights []int64
		median  sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			sdk.NewDecWithPrec(2, 0),
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*100), 2), "",
				sdk.NewInt(tc.weights[i]), addrs[0])
			pb = append(pb, vote)
		}

		require.Equal(t, tc.median, pb.weightedMedian())
	}
}

func TestPBTally(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(4, sdk.Coins{})
	tests := []struct {
		inputs    []float64
		weights   []int64
		rewardees []sdk.AccAddress
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			[]sdk.AccAddress{addrs[2]},
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			[]sdk.AccAddress{addrs[1]},
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			[]sdk.AccAddress{},
		},

		{
			// Lots of random votes
			[]float64{1.0, 78.48, 78.11, 79.0},
			[]int64{1, 51, 79, 33},
			[]sdk.AccAddress{addrs[1], addrs[2], addrs[3]},
		},
	}

	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*100), 2), "",
				sdk.NewInt(tc.weights[i]), addrs[i])
			pb = append(pb, vote)
		}

		_, rewardees := pb.tally()
		require.Equal(t, len(tc.rewardees), len(rewardees))
	}
}
