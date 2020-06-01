package keeper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPrevoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	prevote := types.NewExchangeRatePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, prevote)

	KPrevote, err := input.OracleKeeper.GetExchangeRatePrevote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, prevote, KPrevote)

	input.OracleKeeper.DeleteExchangeRatePrevote(input.Ctx, prevote)
	_, err = input.OracleKeeper.GetExchangeRatePrevote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestPrevoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	prevote1 := types.NewExchangeRatePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, prevote1)

	prevote2 := types.NewExchangeRatePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[1]), 0)
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, prevote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateExchangeRatePrevotes(input.Ctx, func(p types.ExchangeRatePrevote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, prevote1, p)
		} else {
			require.Equal(t, prevote2, p)
		}

		i++
		return false
	})

	prevote3 := types.NewExchangeRatePrevote("", core.MicroLunaDenom, sdk.ValAddress(Addrs[2]), 0)
	input.OracleKeeper.AddExchangeRatePrevote(input.Ctx, prevote3)

	input.OracleKeeper.iterateExchangeRatePrevotesWithPrefix(input.Ctx, types.GetExchangeRatePrevoteKey(core.MicroLunaDenom, sdk.ValAddress{}), func(p types.ExchangeRatePrevote) (stop bool) {
		require.Equal(t, prevote3, p)
		return false
	})
}

func TestVoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote := types.NewExchangeRateVote(rate, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote)

	KVote, err := input.OracleKeeper.getExchangeRateVote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, vote, KVote)

	input.OracleKeeper.DeleteExchangeRateVote(input.Ctx, vote)
	_, err = input.OracleKeeper.getExchangeRateVote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestVoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote1 := types.NewExchangeRateVote(rate, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote1)

	vote2 := types.NewExchangeRateVote(rate, core.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateExchangeRateVotes(input.Ctx, func(p types.ExchangeRateVote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote1, p)
		} else {
			require.Equal(t, vote2, p)
		}

		i++
		return false
	})

	vote3 := types.NewExchangeRateVote(rate, core.MicroLunaDenom, sdk.ValAddress(Addrs[2]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote3)

	input.OracleKeeper.iterateExchangeRateVotesWithPrefix(input.Ctx, types.GetVoteKey(core.MicroLunaDenom, sdk.ValAddress{}), func(p types.ExchangeRateVote) (stop bool) {
		require.Equal(t, vote3, p)
		return false
	})
}

func TestVoteCollect(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote1 := types.NewExchangeRateVote(rate, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote1)

	vote2 := types.NewExchangeRateVote(rate, core.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote2)

	vote3 := types.NewExchangeRateVote(rate, core.MicroLunaDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote3)

	vote4 := types.NewExchangeRateVote(rate, core.MicroLunaDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote4)

	collectedVotes := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	pb1 := collectedVotes[core.MicroSDRDenom]
	pb2 := collectedVotes[core.MicroLunaDenom]

	bigger := bytes.Compare(Addrs[0], Addrs[1])
	for i, v := range pb1 {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote1, v)
		} else {
			require.Equal(t, vote2, v)
		}
	}

	for i, v := range pb2 {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote3, v)
		} else {
			require.Equal(t, vote4, v)
		}
	}
}

func TestExchangeRate(t *testing.T) {
	input := CreateTestInput(t)

	cnyExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	gbpExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	krwExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	lunaExchangeRate := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroCNYDenom, cnyExchangeRate)
	rate, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroCNYDenom)
	require.NoError(t, err)
	require.Equal(t, cnyExchangeRate, rate)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroGBPDenom, gbpExchangeRate)
	rate, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroGBPDenom)
	require.NoError(t, err)
	require.Equal(t, gbpExchangeRate, rate)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, krwExchangeRate)
	rate, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.NoError(t, err)
	require.Equal(t, krwExchangeRate, rate)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroLunaDenom, lunaExchangeRate)
	rate, _ = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroLunaDenom)
	require.Equal(t, sdk.OneDec(), rate)

	input.OracleKeeper.DeleteLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	_, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.Error(t, err)

	numExchangeRates := 0
	handler := func(denom string, exchangeRate sdk.Dec) (stop bool) {
		numExchangeRates = numExchangeRates + 1
		return false
	}
	input.OracleKeeper.IterateLunaExchangeRates(input.Ctx, handler)

	require.True(t, numExchangeRates == 3)
}

func TestIterateLunaExchangeRates(t *testing.T) {
	input := CreateTestInput(t)

	cnyExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	gbpExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	krwExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	lunaExchangeRate := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroCNYDenom, cnyExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroGBPDenom, gbpExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, krwExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroLunaDenom, lunaExchangeRate)

	input.OracleKeeper.IterateLunaExchangeRates(input.Ctx, func(denom string, rate sdk.Dec) (stop bool) {
		switch denom {
		case core.MicroCNYDenom:
			require.Equal(t, cnyExchangeRate, rate)
		case core.MicroGBPDenom:
			require.Equal(t, gbpExchangeRate, rate)
		case core.MicroKRWDenom:
			require.Equal(t, krwExchangeRate, rate)
		case core.MicroLunaDenom:
			require.Equal(t, lunaExchangeRate, rate)
		}
		return false
	})

}

func TestRewardPool(t *testing.T) {
	input := CreateTestInput(t)

	fees := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1000)))
	acc := input.SupplyKeeper.GetModuleAccount(input.Ctx, types.ModuleName)
	err := acc.SetCoins(fees)
	if err != nil {
		panic(err) // nerver occurs
	}

	input.SupplyKeeper.SetModuleAccount(input.Ctx, acc)

	KFees := input.OracleKeeper.GetRewardPool(input.Ctx)
	require.Equal(t, fees, KFees)
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	// Test default params setting
	input.OracleKeeper.SetParams(input.Ctx, types.DefaultParams())
	params := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := int64(10)
	voteThreshold := sdk.NewDecWithPrec(1, 10)
	oracleRewardBand := sdk.NewDecWithPrec(1, 2)
	rewardDistributionWindow := int64(10000000000000)
	slashFraction := sdk.NewDecWithPrec(1, 2)
	slashWindow := int64(1000)
	minValidPerWindow := sdk.NewDecWithPrec(1, 4)
	whilelist := types.DenomList{
		core.MicroSDRDenom,
		core.MicroKRWDenom,
	}

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:               votePeriod,
		VoteThreshold:            voteThreshold,
		RewardBand:               oracleRewardBand,
		RewardDistributionWindow: rewardDistributionWindow,
		Whitelist:                whilelist,
		SlashFraction:            slashFraction,
		SlashWindow:              slashWindow,
		MinValidPerWindow:        minValidPerWindow,
	}
	input.OracleKeeper.SetParams(input.Ctx, newParams)

	storedParams := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, storedParams, newParams)
}

func TestFeederDelegation(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, ValAddrs[0], Addrs[1])
	delegate = input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[1], delegate)
}

func TestIterateFeederDelegations(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, ValAddrs[0], Addrs[1])

	var delegators []sdk.ValAddress
	var delegates []sdk.AccAddress
	input.OracleKeeper.IterateOracleDelegates(input.Ctx, func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool) {
		delegators = append(delegators, delegator)
		delegates = append(delegates, delegate)
		return false
	})

	require.Equal(t, 1, len(delegators))
	require.Equal(t, 1, len(delegates))
	require.Equal(t, ValAddrs[0], delegators[0])
	require.Equal(t, Addrs[1], delegates[0])
}

func TestMissCounter(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)

	missCounter := int64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], missCounter)
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, missCounter, counter)

	input.OracleKeeper.DeleteMissCounter(input.Ctx, ValAddrs[0])
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)
}

func TestIterateMissCounters(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)

	missCounter := int64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[1], missCounter)

	var operators []sdk.ValAddress
	var missCounters []int64
	input.OracleKeeper.IterateMissCounters(input.Ctx, func(delegator sdk.ValAddress, missCounter int64) (stop bool) {
		operators = append(operators, delegator)
		missCounters = append(missCounters, missCounter)
		return false
	})

	require.Equal(t, 1, len(operators))
	require.Equal(t, 1, len(missCounters))
	require.Equal(t, ValAddrs[1], operators[0])
	require.Equal(t, missCounter, missCounters[0])
}
