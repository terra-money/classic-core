package oracle

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/terra-project/core/types/assets"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	mcVal "github.com/terra-project/core/types/mock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func generateRandomTestCase() (prices []float64, valAccAddrs []sdk.AccAddress, mockValset mcVal.MockValset) {
	mockValset = mcVal.NewMockValSet()
	valAccAddrs = []sdk.AccAddress{}
	base := math.Pow10(oracleDecPrecision)

	rand.Seed(int64(time.Now().Nanosecond()))
	numInputs := 10 + (rand.Int() % 100)
	for i := 0; i < numInputs; i++ {
		price := float64(int64(rand.Float64()*base)) / base
		prices = append(prices, price)

		valAccAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		valAccAddrs = append(valAccAddrs, valAccAddr)

		power := sdk.NewInt(rand.Int63() % 1000)
		mockValAddr := sdk.ValAddress(valAccAddr.Bytes())
		mockVal := mcVal.NewMockValidator(mockValAddr, power)

		mockValset.Validators = append(mockValset.Validators, mockVal)
	}

	return
}

func checkFloatEquality(a sdk.Dec, b float64, precision int) bool {
	base := math.Pow10(precision)

	a2 := a.MulInt64(int64(base)).TruncateInt64()
	b2 := int64(b * base)

	return a2 == b2
}

func TestPBPower(t *testing.T) {
	input := createTestInput(t)

	_, valAccAddrs, mockValset := generateRandomTestCase()
	pb := PriceBallot{}
	ballotPower := sdk.ZeroInt()

	for i := 0; i < len(mockValset.Validators); i++ {
		vote := NewPriceVote(sdk.ZeroDec(), assets.MicroSDRDenom, sdk.ValAddress(valAccAddrs[i]))
		pb = append(pb, vote)

		valPower, err := vote.getPower(input.ctx, mockValset)
		require.Nil(t, err)

		ballotPower = ballotPower.Add(valPower)
	}

	require.Equal(t, ballotPower, pb.power(input.ctx, mockValset))

	// Mix in a fake validator, the total power should not have changed.
	fakeVote := NewPriceVote(sdk.OneDec(), assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	pb = append(pb, fakeVote)
	require.Equal(t, ballotPower, pb.power(input.ctx, mockValset))
}

func TestPBWeightedMedian(t *testing.T) {
	input := createTestInput(t)
	tests := []struct {
		inputs      []float64
		weights     []int64
		isValidator []bool
		median      sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Adding fake validator doesn't change outcome
			[]float64{1.0, 2.0, 10.0, 100000.0, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(2, 0),
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	mockValset := mcVal.NewMockValSet()
	base := math.Pow10(oracleDecPrecision)
	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			valAccAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := sdk.NewInt(tc.weights[i])
			mockValAddr := sdk.ValAddress(valAccAddr.Bytes())
			mockVal := mcVal.NewMockValidator(mockValAddr, power)

			if tc.isValidator[i] {
				mockValset.Validators = append(mockValset.Validators, mockVal)
			}
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*base), int64(oracleDecPrecision)), assets.MicroSDRDenom, sdk.ValAddress(valAccAddr))
			pb = append(pb, vote)
		}

		require.Equal(t, tc.median, pb.weightedMedian(input.ctx, mockValset))
	}
}

// func TestPBTally(t *testing.T) {
// 	_, addrs, _, _ := mock.CreateGenAccounts(4, sdk.Coins{})
// 	tests := []struct {
// 		inputs    []float64
// 		weights   []int64
// 		rewardees []sdk.AccAddress
// 	}{
// 		{
// 			// Supermajority one number
// 			[]float64{1.0, 2.0, 10.0, 100000.0},
// 			[]int64{1, 1, 100, 1},
// 			[]sdk.AccAddress{addrs[2]},
// 		},
// 		{
// 			// Tie votes
// 			[]float64{1.0, 2.0, 3.0, 4.0},
// 			[]int64{1, 100, 100, 1},
// 			[]sdk.AccAddress{addrs[1]},
// 		},
// 		{
// 			// No votes
// 			[]float64{},
// 			[]int64{},
// 			[]sdk.AccAddress{},
// 		},

// 		{
// 			// Lots of random votes
// 			[]float64{1.0, 78.48, 78.11, 79.0},
// 			[]int64{1, 51, 79, 33},
// 			[]sdk.AccAddress{addrs[1], addrs[2], addrs[3]},
// 		},
// 	}

// 	for _, tc := range tests {
// 		pb := PriceBallot{}
// 		for i, input := range tc.inputs {
// 			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*100), 2), "",
// 				sdk.NewInt(tc.weights[i]), addrs[i])
// 			pb = append(pb, vote)
// 		}

// 		_, rewardees := pb.tally()
// 		require.Equal(t, len(tc.rewardees), len(rewardees))
// 	}
// }
