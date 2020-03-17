package keeper

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestOrganize(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	got := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.True(t, got.IsOK())
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), core.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), core.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), core.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), core.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), core.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), core.MicroKRWDenom, ValAddrs[2]), power),
	}

	for _, vote := range sdrBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range krwBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[core.MicroSDRDenom])
	sort.Sort(ballotMap[core.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[core.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[core.MicroKRWDenom])
}

func TestOrganizeAggregate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	got := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.True(t, got.IsOK())
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), core.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), core.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), core.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), core.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), core.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), core.MicroKRWDenom, ValAddrs[2]), power),
	}

	for i := range sdrBallot {
		input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
			{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
			{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
		}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[core.MicroSDRDenom])
	sort.Sort(ballotMap[core.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[core.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[core.MicroKRWDenom])
}

func TestDuplicateVote(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	got := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.True(t, got.IsOK())
	got = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.True(t, got.IsOK())
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), core.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), core.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), core.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), core.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), core.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), core.MicroKRWDenom, ValAddrs[2]), power),
	}

	for i := range sdrBallot {

		// this will be ignored
		for _, vote := range sdrBallot {
			input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
		}
		for _, vote := range krwBallot {
			input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
		}

		input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
			{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
			{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
		}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[core.MicroSDRDenom])
	sort.Sort(ballotMap[core.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[core.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[core.MicroKRWDenom])
}
