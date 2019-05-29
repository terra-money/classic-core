package oracle

import (
	"encoding/hex"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/mock"
	mcVal "github.com/terra-project/core/types/mock"
)

func TestOracleWhitelist(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroKRWDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	// Less than the threshold signs, msg fails
	msg := NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroKRWDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	input.ctx = input.ctx.WithBlockHeight(1)

	msg = NewMsgPriceFeed("", salt, assets.MicroKRWDenom, addrs[0], sdk.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	EndBlocker(input.ctx, input.oracleKeeper)
}

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	// Prevote without price
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	msg := NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	res := h(input.ctx.WithBlockHeight(0), msg)
	require.True(t, res.IsOK())

	// Vote and new Prevote
	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), randomPrice)
	res = h(input.ctx.WithBlockHeight(1), msg)
	require.True(t, res.IsOK())

	EndBlocker(input.ctx.WithBlockHeight(1), input.oracleKeeper)

	_, err = input.oracleKeeper.GetLunaSwapRate(input.ctx.WithBlockHeight(1), assets.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	salt = "1"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	h(input.ctx.WithBlockHeight(0), msg)

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), randomPrice)
	h(input.ctx.WithBlockHeight(1), msg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[1]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), sdk.ZeroDec())
	h(input.ctx.WithBlockHeight(0), msg)

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), randomPrice)
	h(input.ctx.WithBlockHeight(1), msg)

	salt = "3"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[2]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[2], sdk.ValAddress(addrs[2]), sdk.ZeroDec())
	h(input.ctx.WithBlockHeight(0), msg)

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[2], sdk.ValAddress(addrs[2]), randomPrice)
	h(input.ctx.WithBlockHeight(1), msg)

	EndBlocker(input.ctx.WithBlockHeight(1), input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx.WithBlockHeight(1), assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)

	// A new validator joins, we are now below threshold. Price update should now fail
	newValidator := mock.NewMockValidator(sdk.ValAddress(addrs[2].Bytes()), sdk.NewInt(30))
	input.valset.Validators = append(input.valset.Validators, newValidator)
	input.oracleKeeper.valset = input.valset

	salt = "1"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	h(input.ctx.WithBlockHeight(0), msg)

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), randomPrice)
	h(input.ctx.WithBlockHeight(1), msg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[1]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), sdk.ZeroDec())
	h(input.ctx.WithBlockHeight(0), msg)

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), randomPrice)
	h(input.ctx.WithBlockHeight(1), msg)

	EndBlocker(input.ctx.WithBlockHeight(1), input.oracleKeeper)

	price, err = input.oracleKeeper.GetLunaSwapRate(input.ctx.WithBlockHeight(1), assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	msg := NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	res := h(input.ctx, msg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[1]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), sdk.ZeroDec())
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[2]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[2], sdk.ValAddress(addrs[2]), sdk.ZeroDec())
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[1]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), sdk.ZeroDec())
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[2]))
	msg = NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroSDRDenom, addrs[2], sdk.ValAddress(addrs[2]), sdk.ZeroDec())
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	// Reveal Price
	input.ctx = input.ctx.WithBlockHeight(1)
	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[0], sdk.ValAddress(addrs[0]), anotherRandomPrice)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[1], sdk.ValAddress(addrs[1]), anotherRandomPrice)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	msg = NewMsgPriceFeed("", salt, assets.MicroSDRDenom, addrs[2], sdk.ValAddress(addrs[2]), anotherRandomPrice)
	res = h(input.ctx, msg)
	require.True(t, res.IsOK())

	EndBlocker(input.ctx, input.oracleKeeper)

	price, err := input.oracleKeeper.GetLunaSwapRate(input.ctx, assets.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, price, anotherRandomPrice)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	dropThreshold := input.oracleKeeper.GetParams(input.ctx).DropThreshold
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, randomPrice)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, assets.MicroKRWDenom, sdk.ValAddress(addrs[0]))
	msg := NewMsgPriceFeed(hex.EncodeToString(bz), "", assets.MicroKRWDenom, addrs[0], sdk.ValAddress(addrs[0]), sdk.ZeroDec())
	h(input.ctx, msg)

	input.ctx = input.ctx.WithBlockHeight(1)
	msg = NewMsgPriceFeed("", salt, assets.MicroKRWDenom, addrs[0], sdk.ValAddress(addrs[0]), randomPrice)
	h(input.ctx, msg)

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

		salt := string(i)
		bz, err := VoteHash(salt, decPrice, assets.MicroSDRDenom, sdk.ValAddress(valAccAddrs[i]))
		require.Nil(t, err)

		pfm := NewMsgPriceFeed(
			hex.EncodeToString(bz),
			"",
			assets.MicroSDRDenom,
			valAccAddrs[i],
			sdk.ValAddress(valAccAddrs[i]),
			sdk.ZeroDec(),
		)

		res := h(input.ctx.WithBlockHeight(0), pfm)
		require.True(t, res.IsOK())

		pfm = NewMsgPriceFeed(
			"",
			salt,
			assets.MicroSDRDenom,
			valAccAddrs[i],
			sdk.ValAddress(valAccAddrs[i]),
			decPrice,
		)

		res = h(input.ctx.WithBlockHeight(1), pfm)
		require.True(t, res.IsOK())

		vote := NewPriceVote(decPrice, assets.MicroSDRDenom, sdk.ValAddress(valAccAddrs[i]))
		ballot = append(ballot, vote)

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
		salt := "1"
		bz, err := VoteHash(salt, sdk.OneDec(), assets.MicroSDRDenom, sdk.ValAddress(addr))
		require.Nil(t, err)

		pfm := NewMsgPriceFeed(
			hex.EncodeToString(bz),
			"",
			assets.MicroSDRDenom,
			addr,
			sdk.ValAddress(addr),
			sdk.ZeroDec(),
		)

		res := h(input.ctx, pfm)
		require.True(t, res.IsOK())

		pfm = NewMsgPriceFeed(
			"",
			salt,
			assets.MicroSDRDenom,
			addr,
			sdk.ValAddress(addr),
			sdk.OneDec(),
		)

		res = h(input.ctx.WithBlockHeight(1), pfm)
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
