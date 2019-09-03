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

func TestClaimPool(t *testing.T) {
	input := CreateTestInput(t)

	// Test addClaimPool
	claim := types.NewClaim(10, ValAddrs[0])
	claim2 := types.NewClaim(20, ValAddrs[1])
	claimPool := types.ClaimPool{claim, claim2}
	input.OracleKeeper.AddClaimPool(input.Ctx, claimPool)

	claim = types.NewClaim(15, ValAddrs[0])
	claim2 = types.NewClaim(30, ValAddrs[2])
	claimPool = types.ClaimPool{claim, claim2}
	input.OracleKeeper.AddClaimPool(input.Ctx, claimPool)

	// Test IterateClaimPool
	input.OracleKeeper.IterateClaimPool(input.Ctx, func(recipient sdk.ValAddress, weight int64) (stop bool) {
		if recipient.Equals(ValAddrs[0]) {
			require.Equal(t, int64(25), weight)
		} else if recipient.Equals(ValAddrs[1]) {
			require.Equal(t, int64(20), weight)
		} else if recipient.Equals(ValAddrs[2]) {
			require.Equal(t, int64(30), weight)
		}
		return false
	})
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
	votesWindow := int64(2000)
	minValidVotesPerWindow := sdk.NewDecWithPrec(1, 2)
	slashFraction := sdk.NewDecWithPrec(5, 2)
	rewardFraction := sdk.NewDecWithPrec(1, 2)

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:             votePeriod,
		VoteThreshold:          voteThreshold,
		RewardBand:             oracleRewardBand,
		VotesWindow:            votesWindow,
		MinValidVotesPerWindow: minValidVotesPerWindow,
		SlashFraction:          slashFraction,
		RewardFraction:         rewardFraction,
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

func TestVotingInfo(t *testing.T) {
	input := CreateTestInput(t)

	// voting info not found
	_, found := input.OracleKeeper.getVotingInfo(input.Ctx, ValAddrs[0])
	require.False(t, found)

	// register voting info
	votingInfo := types.NewVotingInfo(ValAddrs[0], 7, 1, 32)
	input.OracleKeeper.SetVotingInfo(input.Ctx, ValAddrs[0], votingInfo)

	KVotingInfo, found := input.OracleKeeper.getVotingInfo(input.Ctx, ValAddrs[0])
	require.True(t, found)
	require.Equal(t, votingInfo, KVotingInfo)

	votingInfo2 := types.NewVotingInfo(ValAddrs[1], 1, 2, 3)
	input.OracleKeeper.SetVotingInfo(input.Ctx, ValAddrs[1], votingInfo2)

	i := 0
	bigger := bytes.Compare(Addrs[0].Bytes(), Addrs[1].Bytes())
	input.OracleKeeper.IterateVotingInfos(input.Ctx, func(info types.VotingInfo) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, votingInfo, info)
		} else {
			require.Equal(t, votingInfo2, info)
		}
		i++
		return false
	})

}

func TestGetSetValidatorMissedBlockBitArray(t *testing.T) {
	input := CreateTestInput(t)
	missed := input.OracleKeeper.GetMissedVoteBitArray(input.Ctx, ValAddrs[0], 0)
	require.False(t, missed) // treat empty key as not missed
	input.OracleKeeper.SetMissedVoteBitArray(input.Ctx, ValAddrs[0], 0, true)
	missed = input.OracleKeeper.GetMissedVoteBitArray(input.Ctx, ValAddrs[0], 0)
	require.True(t, missed) // now should be missed

	// iterate 1 bit array
	input.OracleKeeper.IterateMissedVoteBitArray(input.Ctx, ValAddrs[0], func(index int64, missed bool) (stop bool) {
		require.Equal(t, int64(0), index)
		require.True(t, missed)
		return false
	})

	// clear vote bit array
	input.OracleKeeper.clearMissedVoteBitArray(input.Ctx, ValAddrs[0])
	missed = input.OracleKeeper.GetMissedVoteBitArray(input.Ctx, ValAddrs[0], 0)
	require.False(t, missed) // treat empty key as not missed
}
