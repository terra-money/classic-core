package oracle

import (
	"testing"

	"gonum.org/v1/gonum/stat"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestPBStdDev(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		inputs []float64
		mean   sdk.Dec
	}{
		{[]float64{1.98, 1.44, 47727.19, 34.39, 4375.22, 59.11, 3.44, 955.14, 0.29}, sdk.NewDecWithPrec(100, 2)},
	}

	for _, tc := range tests {
		pb := PriceBallot{}
		weights := []float64{}
		for _, i := range tc.inputs {
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(i*100), 2), "", sdk.OneDec(), addrs[0])
			weights = append(weights, 1.0)
			pb = append(pb, vote)
		}

		require.Equal(t, stat.StdDev(tc.inputs, weights), pb.stdDev())
	}
}

func TestPBMean(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		inputs []float64
		mean   sdk.Dec
	}{
		{[]float64{1.0, 1.0, 1.0, 1.0}, sdk.NewDecWithPrec(100, 2)},
		{[]float64{1.1, 1.9, 3.0}, sdk.NewDecWithPrec(200, 2)},
		{[]float64{0.0}, sdk.ZeroDec()},
		{[]float64{}, sdk.ZeroDec()},
	}

	for _, tc := range tests {
		pb := PriceBallot{}
		for _, i := range tc.inputs {
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(i*100), 2), "", sdk.ZeroDec(), addrs[0])
			pb = append(pb, vote)
		}

		require.Equal(t, tc.mean, pb.mean())
	}
}

func TestPBWeightedMedian(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	tests := []struct {
		inputs  []float64
		weights []float64
		median  sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]float64{1.0, 1.0, 100.0, 1.0},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]float64{1.0, 100.0, 100.0, 1.0},
			sdk.NewDecWithPrec(2, 0),
		},
		{
			// No votes
			[]float64{},
			[]float64{},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*100), 2), "",
				sdk.NewDecWithPrec(int64(tc.weights[i]*100), 2), addrs[0])
			pb = append(pb, vote)
		}

		_, median := pb.weightedMedian()
		require.Equal(t, tc.median, median.Price)
	}
}
