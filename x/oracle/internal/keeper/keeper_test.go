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

	prevote := types.NewPricePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote)

	KPrevote, err := input.OracleKeeper.GetPrevote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, prevote, KPrevote)

	input.OracleKeeper.DeletePrevote(input.Ctx, prevote)
	_, err = input.OracleKeeper.GetPrevote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestPrevoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	prevote1 := types.NewPricePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote1)

	prevote2 := types.NewPricePrevote("", core.MicroSDRDenom, sdk.ValAddress(Addrs[1]), 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IteratePrevotes(input.Ctx, func(p types.PricePrevote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, prevote1, p)
		} else {
			require.Equal(t, prevote2, p)
		}

		i++
		return false
	})

	prevote3 := types.NewPricePrevote("", core.MicroLunaDenom, sdk.ValAddress(Addrs[2]), 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote3)

	input.OracleKeeper.iteratePrevotesWithPrefix(input.Ctx, types.GetPrevoteKey(core.MicroLunaDenom, sdk.ValAddress{}), func(p types.PricePrevote) (stop bool) {
		require.Equal(t, prevote3, p)
		return false
	})
}

func TestVoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	price := sdk.NewDec(1700)
	vote := types.NewPriceVote(price, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddVote(input.Ctx, vote)

	KVote, err := input.OracleKeeper.getVote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, vote, KVote)

	input.OracleKeeper.DeleteVote(input.Ctx, vote)
	_, err = input.OracleKeeper.getVote(input.Ctx, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestVoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	price := sdk.NewDec(1700)
	vote1 := types.NewPriceVote(price, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddVote(input.Ctx, vote1)

	vote2 := types.NewPriceVote(price, core.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddVote(input.Ctx, vote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateVotes(input.Ctx, func(p types.PriceVote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote1, p)
		} else {
			require.Equal(t, vote2, p)
		}

		i++
		return false
	})

	vote3 := types.NewPriceVote(price, core.MicroLunaDenom, sdk.ValAddress(Addrs[2]))
	input.OracleKeeper.AddVote(input.Ctx, vote3)

	input.OracleKeeper.iterateVotesWithPrefix(input.Ctx, types.GetVoteKey(core.MicroLunaDenom, sdk.ValAddress{}), func(p types.PriceVote) (stop bool) {
		require.Equal(t, vote3, p)
		return false
	})
}

func TestVoteCollect(t *testing.T) {
	input := CreateTestInput(t)

	price := sdk.NewDec(1700)
	vote1 := types.NewPriceVote(price, core.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddVote(input.Ctx, vote1)

	vote2 := types.NewPriceVote(price, core.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddVote(input.Ctx, vote2)

	vote3 := types.NewPriceVote(price, core.MicroLunaDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddVote(input.Ctx, vote3)

	vote4 := types.NewPriceVote(price, core.MicroLunaDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddVote(input.Ctx, vote4)

	collectedVotes := input.OracleKeeper.CollectVotes(input.Ctx)

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

func TestPrice(t *testing.T) {
	input := CreateTestInput(t)

	cnyPrice := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	gbpPrice := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	krwPrice := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	lunaPrice := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)

	// Set & get prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroCNYDenom, cnyPrice)
	price, err := input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroCNYDenom)
	require.NoError(t, err)
	require.Equal(t, cnyPrice, price)

	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroGBPDenom, gbpPrice)
	price, err = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroGBPDenom)
	require.NoError(t, err)
	require.Equal(t, gbpPrice, price)

	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, krwPrice)
	price, err = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroKRWDenom)
	require.NoError(t, err)
	require.Equal(t, krwPrice, price)

	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroLunaDenom, lunaPrice)
	price, _ = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroLunaDenom)
	require.Equal(t, sdk.OneDec(), price)

	input.OracleKeeper.DeletePrice(input.Ctx, core.MicroKRWDenom)
	_, err = input.OracleKeeper.GetLunaPrice(input.Ctx, core.MicroKRWDenom)
	require.Error(t, err)
}

func TestIterateLunaPrices(t *testing.T) {
	input := CreateTestInput(t)

	cnyPrice := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	gbpPrice := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	krwPrice := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)
	lunaPrice := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(core.MicroUnit)

	// Set & get prices
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroCNYDenom, cnyPrice)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroGBPDenom, gbpPrice)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, krwPrice)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroLunaDenom, lunaPrice)

	input.OracleKeeper.IterateLunaPrices(input.Ctx, func(denom string, price sdk.Dec) (stop bool) {
		switch denom {
		case core.MicroCNYDenom:
			require.Equal(t, cnyPrice, price)
		case core.MicroGBPDenom:
			require.Equal(t, gbpPrice, price)
		case core.MicroKRWDenom:
			require.Equal(t, krwPrice, price)
		case core.MicroLunaDenom:
			require.Equal(t, lunaPrice, price)
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

	KFees := input.OracleKeeper.getRewardPool(input.Ctx)
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
	rewardDistributionPeriod := int64(10000000000000)

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:               votePeriod,
		VoteThreshold:            voteThreshold,
		RewardBand:               oracleRewardBand,
		RewardDistributionPeriod: rewardDistributionPeriod,
	}
	input.OracleKeeper.SetParams(input.Ctx, newParams)

	storedParams := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, newParams, storedParams)
}

func TestFeederDelegation(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetFeedDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, delegate, Addrs[0])

	input.OracleKeeper.SetFeedDelegate(input.Ctx, ValAddrs[0], Addrs[1])
	delegate = input.OracleKeeper.GetFeedDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, delegate, Addrs[1])
}

func TestIterateFeederDelegations(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetFeedDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, delegate, Addrs[0])

	input.OracleKeeper.SetFeedDelegate(input.Ctx, ValAddrs[0], Addrs[1])

	var delegators []sdk.ValAddress
	var delegatees []sdk.AccAddress
	input.OracleKeeper.IterateFeederDelegations(input.Ctx, func(delegator sdk.ValAddress, delegatee sdk.AccAddress) (stop bool) {
		delegators = append(delegators, delegator)
		delegatees = append(delegatees, delegatee)
		return false
	})

	require.Equal(t, len(delegators), 1)
	require.Equal(t, len(delegatees), 1)
	require.Equal(t, delegators[0], ValAddrs[0])
	require.Equal(t, delegatees[0], Addrs[1])
}
