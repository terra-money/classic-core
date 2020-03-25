package oracle

import (
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
	// Prevote without exchange rate
	salt := "1"
	hash := GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	prevoteMsg := NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	require.True(t, res.IsOK())

	// Vote and new Prevote
	voteMsg := NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	salt = "1"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[1])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "3"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[2])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomExchangeRate, rate)

	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), sdk.Unbonded, val, false)

	salt = "1"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[1])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = NewMsgExchangeRateVote(randomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx.WithBlockHeight(1), core.MicroSDRDenom)
	require.NotNil(t, err)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	salt := "1"
	hash := GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	prevoteMsg := NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res := h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[1])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	hash = GetVoteHash(salt, randomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[2])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	hash = GetVoteHash(salt, anotherRandomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[0])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	hash = GetVoteHash(salt, anotherRandomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[1])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	hash = GetVoteHash(salt, anotherRandomExchangeRate, core.MicroSDRDenom, keeper.ValAddrs[2])

	prevoteMsg = NewMsgExchangeRatePrevote(hash, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, prevoteMsg)
	require.True(t, res.IsOK())

	// Reveal ExchangeRate
	input.Ctx = input.Ctx.WithBlockHeight(1)
	voteMsg := NewMsgExchangeRateVote(anotherRandomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[0], keeper.ValAddrs[0])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgExchangeRateVote(anotherRandomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[1], keeper.ValAddrs[1])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	voteMsg = NewMsgExchangeRateVote(anotherRandomExchangeRate, salt, core.MicroSDRDenom, keeper.Addrs[2], keeper.ValAddrs[2])
	res = h(input.Ctx, voteMsg)
	require.True(t, res.IsOK())

	EndBlocker(input.Ctx, input.OracleKeeper)

	rate, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, rate, anotherRandomExchangeRate)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, randomExchangeRate)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

	// Immediately swap halt after an illiquid oracle vote
	EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.NotNil(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := ExchangeRateBallot{}
	rates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	h := NewHandler(input.OracleKeeper)
	for i, rate := range rates {

		decExchangeRate := sdk.NewDecWithPrec(int64(rate*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))

		salt := string(i)
		hash := GetVoteHash(salt, decExchangeRate, core.MicroSDRDenom, valAddrs[i])

		prevoteMsg := NewMsgExchangeRatePrevote(
			hash,
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
		require.True(t, res.IsOK())

		voteMsg := NewMsgExchangeRateVote(
			decExchangeRate,
			salt,
			core.MicroSDRDenom,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		res = h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.True(t, res.IsOK())

		vote := NewVoteForTally(NewExchangeRateVote(decExchangeRate, core.MicroSDRDenom, valAddrs[i]), stakingAmt.QuoRaw(core.MicroUnit).Int64())
		ballot = append(ballot, vote)

		// change power of every three validator
		if i%3 == 0 {
			stakingKeeper.Validators()[i].SetPower(int64(i + 1))
		}
	}

	rewardees := []sdk.AccAddress{}
	weightedMedian := ballot.WeightedMedian()
	standardDeviation := ballot.StandardDeviation()
	maxSpread := weightedMedian.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	for _, vote := range ballot {
		if vote.ExchangeRate.GTE(weightedMedian.Sub(maxSpread)) && vote.ExchangeRate.LTE(weightedMedian.Add(maxSpread)) {
			rewardees = append(rewardees, sdk.AccAddress(vote.Voter))
		}
	}

	tallyMedian, ballotWinner := tally(input.Ctx, ballot, input.OracleKeeper.RewardBand(input.Ctx))

	require.Equal(t, len(rewardees), len(ballotWinner))
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for i := range keeper.Addrs {
		makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomExchangeRate, i)
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
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomExchangeRate, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomExchangeRate, 1)

	rewardsAmt := sdk.NewInt(100000000)
	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, rewardsAmt)))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.RewardDistributionWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	expectedRewardAmt := sdk.NewDecFromInt(rewardsAmt.QuoRaw(2)).QuoInt64(votePeriodsPerWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestOracleRewardBand(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, DefaultTobinTax)

	rewardSpread := randomExchangeRate.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	// no one will miss the vote
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread), 0)

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Add(rewardSpread), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// Account 1 will miss the vote due to raward band condition
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread.Add(sdk.OneDec())), 0)

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Add(rewardSpread), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomExchangeRate, 0)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, core.MicroSDRDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.SupplyKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), ModuleName)
	err := moduleAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, rewardAmt)))
	require.NoError(t, err)

	input.SupplyKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(4)).QuoInt64(rewardDistributedWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestInvalidVotesSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i < sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

		// Account 2, KRW, miss vote
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)), 1)

		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur ValAddrs[1] slashing
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

	// Account 2, KRW, miss vote
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)), 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestWhitelistSlashing(t *testing.T) {
	input, h := setup(t)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i < sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 2, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)
		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur Account 1 slashing

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)
	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestNotPassedBallotSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, DefaultTobinTax)

	input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

	EndBlocker(input.Ctx, input.OracleKeeper)
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i <= sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)

		// Account 2, KRW, abstain vote
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, sdk.ZeroDec(), 1)

		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())
}

func TestVoteTargets(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}, {Name: core.MicroSDRDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, DefaultTobinTax)

	// KRW
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	// no missing current
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// vote targets are {KRW, SDR}
	require.Equal(t, []string{core.MicroKRWDenom, core.MicroSDRDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	// tobin tax must be exists for SDR
	sdrTobinTax, err := input.OracleKeeper.GetTobinTax(input.Ctx, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, DefaultTobinTax, sdrTobinTax)

	// delete SDR
	params.Whitelist = types.DenomList{{Name: core.MicroKRWDenom, TobinTax: DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// KRW, missing
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 0, core.MicroKRWDenom, randomExchangeRate, 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// SDR must be deleted
	require.Equal(t, []string{core.MicroKRWDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	_, err = input.OracleKeeper.GetTobinTax(input.Ctx, core.MicroSDRDenom)
	require.Error(t, err)
}

func makePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, denom string, rate sdk.Dec, idx int) {
	// Account 1, SDR
	salt := "1"
	hash := GetVoteHash(salt, rate, denom, keeper.ValAddrs[idx])

	prevoteMsg := NewMsgExchangeRatePrevote(hash, denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
	res := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.True(t, res.IsOK())

	voteMsg := NewMsgExchangeRateVote(rate, salt, denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
	res = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.True(t, res.IsOK())
}

func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, rates sdk.DecCoins, idx int) {
	// Account 1, SDR
	salt := "1"
	hash := GetAggregateVoteHash(salt, rates.String(), keeper.ValAddrs[idx])

	prevoteMsg := NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	res := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.True(t, res.IsOK())

	voteMsg := NewMsgAggregateExchangeRateVote(salt, rates.String(), keeper.Addrs[idx], keeper.ValAddrs[idx])
	res = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.True(t, res.IsOK())
}
