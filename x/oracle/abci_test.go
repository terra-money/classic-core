package oracle

import (
	"encoding/hex"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/keeper"
	"github.com/terra-project/core/x/oracle/internal/types"
)

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	// Prevote without price
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	require.True(t, res.IsOK())

	// Vote and new Prevote
	voteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err = input.OracleKeeper.GetLunaPrice(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	salt = "1"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "3"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	price, err := input.OracleKeeper.GetLunaPrice(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, price)

	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), sdk.Unbonded, val, false)

	salt = "1"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	price, err = input.OracleKeeper.GetLunaPrice(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Reveal Price
	input.Ctx = input.Ctx.WithBlockHeight(1)
	voteMsg := NewMsgPriceVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgPriceVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgPriceVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx, input.OracleKeeper)

	price, err := input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, price, anotherRandomPrice)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, randomPrice)

	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroKRWDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroKRWDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx, prevoteMsg)

	input.Ctx = input.Ctx.WithBlockHeight(1)
	voteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroKRWDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx, voteMsg)

	// Immediately swap halt after an illiquid oracle vote
	EndBlocker(input.Ctx, input.OracleKeeper)

	_, err = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroKRWDenom)
	require.NotNil(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := PriceBallot{}
	prices, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	h := NewHandler(input.OracleKeeper)
	for i, price := range prices {

		decPrice := sdk.NewDecWithPrec(int64(price*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))

		salt := string(i)
		bz, err := VoteHash(salt, decPrice, core.MicroSDRDenom, valAddrs[i])
		require.Nil(t, err)

		prevoteMsg := NewMsgPricePrevote(
			hex.EncodeToString(bz),
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
		require.True(t, res.IsOK())

		voteMsg := NewMsgPriceVote(
			decPrice,
			salt,
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res = h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.True(t, res.IsOK())

		vote := NewPriceVote(decPrice, core.MicroSDRDenom, valAddrs[i])
		ballot = append(ballot, vote)

		// change power of every three validator
		if i%3 == 0 {
			stakingKeeper.Validators()[i].SetPower(int64(i + 1))
		}
	}

	rewardees := []sdk.AccAddress{}
	weightedMedian := ballot.WeightedMedian(input.Ctx, stakingKeeper)
	standardDeviation := ballot.StandardDeviation(input.Ctx, stakingKeeper)
	maxSpread := input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2)

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	for _, vote := range ballot {
		if vote.Price.GTE(weightedMedian.Sub(maxSpread)) && vote.Price.LTE(weightedMedian.Add(maxSpread)) {
			rewardees = append(rewardees, sdk.AccAddress(vote.Voter))
		}
	}

	tallyMedian, ballotWinner := tally(input.Ctx, ballot, input.OracleKeeper)

	require.Equal(t, len(rewardees), len(ballotWinner))
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for _, addr := range keeper.Addrs {
		salt := "1"
		bz, err := VoteHash(salt, sdk.OneDec(), core.MicroSDRDenom, sdk.ValAddress(addr))
		require.Nil(t, err)

		prevoteMsg := NewMsgPricePrevote(
			hex.EncodeToString(bz),
			core.MicroSDRDenom,
			addr,
			sdk.ValAddress(addr),
		)

		res := h(input.Ctx, prevoteMsg)
		require.True(t, res.IsOK())

		voteMsg := NewMsgPriceVote(
			sdk.OneDec(),
			salt,
			core.MicroSDRDenom,
			addr,
			sdk.ValAddress(addr),
		)

		res = h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.True(t, res.IsOK())
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convinience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroSDRDenom)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(params.VotePeriod - 1)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroSDRDenom)
	require.NoError(t, err)
}

func TestOracleRewardDistribution(t *testing.T) {
	input, h := setup(t)

	salt := "1"
	bz, _ := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	bz, _ = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, stakingAmt.MulRaw(100))))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	expectedRewardAmt := sdk.NewDecFromInt(stakingAmt.MulRaw(50)).MulInt64(input.OracleKeeper.VotePeriod(input.Ctx)).QuoInt64(input.OracleKeeper.RewardDistributionPeriod(input.Ctx)).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	salt := "1"
	bz, _ := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	prevoteMsg := NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg := NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	// Account 1, KRW
	salt = "2"
	bz, _ = VoteHash(salt, randomPrice, core.MicroKRWDenom, keeper.ValAddrs[0])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroKRWDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroKRWDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	// Account 2, SDR
	salt = "3"
	bz, _ = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	// Account 3, KRW
	salt = "3"
	bz, _ = VoteHash(salt, randomPrice, core.MicroKRWDenom, keeper.ValAddrs[2])
	prevoteMsg = NewMsgPricePrevote(hex.EncodeToString(bz), core.MicroKRWDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgPriceVote(randomPrice, salt, core.MicroKRWDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, stakingAmt.MulRaw(100))))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriod := input.OracleKeeper.VotePeriod(input.Ctx)
	rewardDistributedPeriod := input.OracleKeeper.RewardDistributionPeriod(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(stakingAmt.MulRaw(50)).MulInt64(votePeriod).QuoInt64(rewardDistributedPeriod).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(stakingAmt.MulRaw(25)).MulInt64(votePeriod).QuoInt64(rewardDistributedPeriod).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
}
