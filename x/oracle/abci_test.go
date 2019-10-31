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
	// Prevote without exchangeRate
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.NoError(t, err)

	prevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	require.True(t, res.IsOK())

	// Vote and new Prevote
	voteMsg := NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	salt = "1"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "3"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	exchangeRate, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomPrice, exchangeRate)

	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), sdk.Unbonded, val, false)

	salt = "1"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgVote(randomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	exchangeRate, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	salt := "1"
	bz, err := VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.NoError(t, err)

	prevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, randomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[0])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[1])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	bz, err = VoteHash(salt, anotherRandomPrice, core.MicroSDRDenom, keeper.ValAddrs[2])
	require.NoError(t, err)

	prevoteMsg = NewMsgPrevote(hex.EncodeToString(bz), core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Reveal Price
	input.Ctx = input.Ctx.WithBlockHeight(1)
	voteMsg := NewMsgVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgVote(anotherRandomPrice, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx, input.OracleKeeper)

	exchangeRate, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, exchangeRate, anotherRandomPrice)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, randomPrice)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomPrice, 0)

	// Immediately swap halt after an illiquid oracle vote
	EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.NotNil(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := ExchangeRateBallot{}
	exchangeRates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	h := NewHandler(input.OracleKeeper)
	for i, exchangeRate := range exchangeRates {

		decPrice := sdk.NewDecWithPrec(int64(exchangeRate*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))

		salt := string(i)
		bz, err := VoteHash(salt, decPrice, core.MicroSDRDenom, valAddrs[i])
		require.NoError(t, err)

		prevoteMsg := NewMsgPrevote(
			hex.EncodeToString(bz),
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
		require.True(t, res.IsOK())

		voteMsg := NewMsgVote(
			decPrice,
			salt,
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res = h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.True(t, res.IsOK())

		vote := NewVote(decPrice, core.MicroSDRDenom, valAddrs[i])
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
		if vote.ExchangeRate.GTE(weightedMedian.Sub(maxSpread)) && vote.ExchangeRate.LTE(weightedMedian.Add(maxSpread)) {
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
	for i := range keeper.Addrs {
		makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomPrice, i)
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convinience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(params.VotePeriod - 1)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.NoError(t, err)
}

func TestOracleRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomPrice, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomPrice, 1)

	rewardsAmt := sdk.NewInt(100000000)
	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, rewardsAmt)))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriod := input.OracleKeeper.VotePeriod(input.Ctx)
	rewardDistributionPeriod := input.OracleKeeper.RewardDistributionPeriod(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardsAmt.QuoRaw(2)).MulInt64(votePeriod).QuoInt64(rewardDistributionPeriod).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomPrice, 0)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomPrice, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomPrice, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomPrice, 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, rewardAmt)))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriod := input.OracleKeeper.VotePeriod(input.Ctx)
	rewardDistributedPeriod := input.OracleKeeper.RewardDistributionPeriod(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(2)).MulInt64(votePeriod).QuoInt64(rewardDistributedPeriod).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(4)).MulInt64(votePeriod).QuoInt64(rewardDistributedPeriod).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroSDRDenom).TruncateInt())
}

func makePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, denom string, exchangeRate sdk.Dec, idx int) {
	// Account 1, SDR
	salt := "1"
	bz, err := VoteHash(salt, exchangeRate, denom, keeper.ValAddrs[idx])
	require.NoError(t, err)

	prevoteMsg := NewMsgPrevote(hex.EncodeToString(bz), denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
	res := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.True(t, res.IsOK())

	voteMsg := NewMsgVote(exchangeRate, salt, denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
	res = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.True(t, res.IsOK())
}
