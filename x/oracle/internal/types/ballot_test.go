package types

import (
	"fmt"
	"math"
	"strconv"

	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestToMap(t *testing.T) {
	tests := struct {
		votes   []VoteForTally
		isValid []bool
	}{

		[]VoteForTally{
			{
				ExchangeRateVote{
					Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
					Denom:        core.MicroKRWDenom,
					ExchangeRate: sdk.NewDec(1600),
				},
				100,
			},
			{
				ExchangeRateVote{
					Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
					Denom:        core.MicroKRWDenom,
					ExchangeRate: sdk.ZeroDec(),
				},
				100,
			},
			{
				ExchangeRateVote{
					Voter:        sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()),
					Denom:        core.MicroKRWDenom,
					ExchangeRate: sdk.NewDec(1500),
				},
				100,
			},
		},
		[]bool{true, false, true},
	}

	pb := ExchangeRateBallot(tests.votes)
	mapData := pb.ToMap()
	for i, vote := range tests.votes {
		exchangeRate, ok := mapData[string(vote.Voter)]
		if tests.isValid[i] {
			require.True(t, ok)
			require.Equal(t, exchangeRate, vote.ExchangeRate)
		} else {
			require.False(t, ok)
		}
	}
}

func TestToCrossRate(t *testing.T) {
	data := []struct {
		base     sdk.Dec
		quote    sdk.Dec
		expected sdk.Dec
	}{
		{
			base:     sdk.NewDec(1600),
			quote:    sdk.NewDec(100),
			expected: sdk.NewDec(16),
		},
		{
			base:     sdk.NewDec(0),
			quote:    sdk.NewDec(100),
			expected: sdk.NewDec(16),
		},
		{
			base:     sdk.NewDec(1600),
			quote:    sdk.NewDec(0),
			expected: sdk.NewDec(16),
		},
	}

	pbBase := ExchangeRateBallot{}
	pbQuote := ExchangeRateBallot{}
	cb := ExchangeRateBallot{}
	for _, data := range data {
		valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
		if !data.base.IsZero() {
			pbBase = append(pbBase, NewVoteForTally(NewExchangeRateVote(data.base, core.MicroKRWDenom, valAddr), 100))
		}

		pbQuote = append(pbQuote, NewVoteForTally(NewExchangeRateVote(data.quote, core.MicroKRWDenom, valAddr), 100))

		if !data.base.IsZero() && !data.quote.IsZero() {
			cb = append(cb, NewVoteForTally(NewExchangeRateVote(data.base.Quo(data.quote), core.MicroKRWDenom, valAddr), 100))
		} else {
			cb = append(cb, NewVoteForTally(NewExchangeRateVote(sdk.ZeroDec(), core.MicroKRWDenom, valAddr), 0))
		}
	}

	baseMapBallot := pbBase.ToMap()
	require.Equal(t, cb, pbQuote.ToCrossRate(baseMapBallot))
}

func TestSqrt(t *testing.T) {
	num := sdk.NewDecWithPrec(144, 4)
	floatNum, err := strconv.ParseFloat(num.String(), 64)
	require.NoError(t, err)

	floatNum = math.Sqrt(floatNum)
	num, err = sdk.NewDecFromStr(fmt.Sprintf("%f", floatNum))
	require.NoError(t, err)

	require.Equal(t, sdk.NewDecWithPrec(12, 2), num)
}

func TestPBPower(t *testing.T) {

	ctx := sdk.NewContext(nil, abci.Header{}, false, nil)
	_, valAccAddrs, sk := GenerateRandomTestCase()
	pb := ExchangeRateBallot{}
	ballotPower := int64(0)

	for i := 0; i < len(sk.Validators()); i++ {
		power := sk.Validator(ctx, valAccAddrs[i]).GetConsensusPower()
		vote := NewVoteForTally(
			NewExchangeRateVote(
				sdk.ZeroDec(),
				core.MicroSDRDenom,
				valAccAddrs[i],
			),
			power,
		)

		pb = append(pb, vote)

		require.NotEqual(t, int64(0), vote.Power)

		ballotPower += vote.Power
	}

	require.Equal(t, ballotPower, pb.Power())

	// Mix in a fake validator, the total power should not have changed.
	pubKey := secp256k1.GenPrivKey().PubKey()
	faceValAddr := sdk.ValAddress(pubKey.Address())
	fakeVote := NewVoteForTally(
		NewExchangeRateVote(
			sdk.OneDec(),
			core.MicroSDRDenom,
			faceValAddr,
		),
		0,
	)

	pb = append(pb, fakeVote)
	require.Equal(t, ballotPower, pb.Power())
}

func TestPBWeightedMedian(t *testing.T) {
	tests := []struct {
		inputs      []int64
		weights     []int64
		isValidator []bool
		median      sdk.Dec
	}{
		{
			// Supermajority one number
			[]int64{2, 1, 10, 100000},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(10),
		},
		{
			// Adding fake validator doesn't change outcome
			[]int64{1, 2, 10, 100000, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDec(10),
		},
		{
			// Tie votes
			[]int64{1, 2, 3, 4},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDec(2),
		},
		{
			// No votes
			[]int64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDec(0),
		},
	}

	for _, tc := range tests {
		pb := ExchangeRateBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			if !tc.isValidator[i] {
				power = 0
			}

			vote := NewVoteForTally(
				NewExchangeRateVote(
					sdk.NewDec(int64(input)),
					core.MicroSDRDenom,
					valAddr,
				),
				power,
			)

			pb = append(pb, vote)
		}

		require.Equal(t, tc.median, pb.WeightedMedian())
	}
}

func TestPBStandardDeviation(t *testing.T) {
	tests := []struct {
		inputs            []float64
		weights           []int64
		isValidator       []bool
		standardDeviation sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(49995000363, oracleDecPrecision),
		},
		{
			// Adding fake validator doesn't change outcome
			[]float64{1.0, 2.0, 10.0, 100000.0, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDecWithPrec(4472135950751006, oracleDecPrecision),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(1224745, oracleDecPrecision),
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	base := math.Pow10(oracleDecPrecision)
	for _, tc := range tests {
		pb := ExchangeRateBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			if !tc.isValidator[i] {
				power = 0
			}

			vote := NewVoteForTally(
				NewExchangeRateVote(
					sdk.NewDecWithPrec(int64(input*base), int64(oracleDecPrecision)),
					core.MicroSDRDenom,
					valAddr,
				),
				power,
			)

			pb = append(pb, vote)
		}

		require.Equal(t, tc.standardDeviation, pb.StandardDeviation())
	}
}
