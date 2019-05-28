package oracle

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	mcVal "github.com/terra-project/core/types/mock"
)

var (
	randomPrice        = sdk.NewDecWithPrec(1049, 2) // swap rate
	anotherRandomPrice = sdk.NewDecWithPrec(4882, 2) // swap rate
)

func setup(t *testing.T) (testInput, sdk.Handler) {
	input := createTestInput(t)
	h := NewHandler(input.oracleKeeper)

	defaultOracleParams := DefaultParams()
	defaultOracleParams.VotePeriod = int64(1) // Set to one block for convinience
	input.oracleKeeper.SetParams(input.ctx, defaultOracleParams)

	return input, h
}

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	_, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1], sdk.ValAddress(addrs[1]))
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[2], sdk.ValAddress(addrs[2]))
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)

	// A new validator joins, we are now below threshold. Price update should now fail
	newValidator := mock.NewMockValidator(sdk.ValAddress(addrs[2].Bytes()), sdk.NewInt(30))
	input.valset.Validators = append(input.valset.Validators, newValidator)
	input.oracleKeeper.valset = input.valset

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	h(input.ctx, msg)

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[1], sdk.ValAddress(addrs[1]))
	h(input.ctx, msg)

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[1], sdk.ValAddress(addrs[1]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, randomPrice, addrs[2], sdk.ValAddress(addrs[2]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[1], sdk.ValAddress(addrs[1]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed(assets.MicroSDRDenom, anotherRandomPrice, addrs[2], sdk.ValAddress(addrs[2]))
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, price, anotherRandomPrice)
}

func TestOracleWhitelist(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(assets.MicroKRWDenom, randomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)
	EndBlocker(input.ctx, input.oracleKeeper)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	dropThreshold := input.oracleKeeper.GetParams(input.ctx).DropThreshold
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, randomPrice)

	msg := NewMsgPriceFeed(assets.MicroKRWDenom, randomPrice, addrs[0], sdk.ValAddress(addrs[0]))
	h(input.ctx, msg)

	input.ctx = input.ctx.WithBlockHeight(1)
	for i := 0; i < int(dropThreshold.Int64())-1; i++ {
		EndBlocker(input.ctx, input.oracleKeeper)
	}

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroKRWDenom)
	require.Nil(t, err)
	require.Equal(t, price, randomPrice)

	// Going over dropthreshold should blacklist the price
	EndBlocker(input.ctx, input.oracleKeeper)

	_, err = input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroKRWDenom)
	require.NotNil(t, err)
}

func generateValset(valWeights []int64) mock.MockValset {
	mockValset := mock.NewMockValSet()

	for i := 0; i < len(valWeights); i++ {
		valAccAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

		power := sdk.NewInt(valWeights[i])
		mockValAddr := sdk.ValAddress(valAccAddr.Bytes())
		mockVal := mcVal.NewMockValidator(mockValAddr, power)

		mockValset.Validators = append(mockValset.Validators, mockVal)
	}

	return mockValset
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := PriceBallot{}
	prices, valAccAddrs, mockValset := generateRandomTestCase()
	input.oracleKeeper.valset = mockValset
	h := NewHandler(input.oracleKeeper)
	for i, price := range prices {

		decPrice := sdk.NewDecWithPrec(int64(price*math.Pow10(oracleDecPrecision)), int64(oracleDecPrecision))
		pfm := NewMsgPriceFeed(
			assets.MicroSDRDenom,
			decPrice,
			valAccAddrs[i],
			sdk.ValAddress(valAccAddrs[i]),
		)

		vote := NewPriceVote(decPrice, assets.MicroSDRDenom, sdk.ValAddress(valAccAddrs[i]))
		ballot = append(ballot, vote)

		res := h(input.ctx, pfm)
		require.True(t, res.IsOK())

		// change power of every three validator
		if i%3 == 0 {
			mockValset.Validators[i].Power = sdk.NewInt(int64(i + 1))
		}
	}

	rewardees := []sdk.AccAddress{}
	weightedMedian := ballot.weightedMedian(input.ctx, mockValset)
	maxSpread := input.oracleKeeper.GetParams(input.ctx).OracleRewardBand.QuoInt64(2)

	for _, vote := range ballot {
		if vote.Price.GTE(weightedMedian.Sub(maxSpread)) && vote.Price.LTE(weightedMedian.Add(maxSpread)) {
			rewardees = append(rewardees, sdk.AccAddress(vote.Voter))
		}
	}

	tallyMedian, tallyClaims := tally(input.ctx, input.oracleKeeper, ballot)

	require.Equal(t, len(tallyClaims), len(rewardees))
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the addrs vote for the block ... not last period block yet, so tally fails
	for _, addr := range addrs {
		pfm := NewMsgPriceFeed(
			assets.MicroSDRDenom,
			sdk.OneDec(),
			addr,
			sdk.ValAddress(addr),
		)

		res := h(input.ctx, pfm)
		require.True(t, res.IsOK())
	}

	params := input.oracleKeeper.GetParams(input.ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convinience
	input.oracleKeeper.SetParams(input.ctx, params)

	require.Equal(t, 0, int(input.ctx.BlockHeight()))
	rewardees, _ := EndBlocker(input.ctx, input.oracleKeeper)
	require.Equal(t, 0, len(rewardees))

	input.ctx = input.ctx.WithBlockHeight(params.VotePeriod - 1)

	rewardees, _ = EndBlocker(input.ctx, input.oracleKeeper)
	require.Equal(t, len(addrs), len(rewardees))
}
