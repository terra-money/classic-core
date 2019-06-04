package budget

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const custom = "custom"

func getQueriedProgram(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, programID uint64) Program {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryProgram}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryProgram, strconv.FormatUint(programID, 10)}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var program Program
	err2 := cdc.UnmarshalJSON(bz, &program)
	require.Nil(t, err2)

	return program
}

func getQueriedVotes(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, voter sdk.AccAddress, programID uint64) Votes {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryVotes}, "/"),
		Data: cdc.MustMarshalJSON(NewQueryVotesParams(voter, programID)),
	}

	bz, err := querier(ctx, []string{QueryVotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votes Votes
	err2 := cdc.UnmarshalJSON(bz, &votes)
	require.Nil(t, err2)

	return votes
}

func getQueriedActiveList(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) Programs {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryActiveList}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryActiveList}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var activeList Programs
	err2 := cdc.UnmarshalJSON(bz, &activeList)
	require.Nil(t, err2)

	return activeList
}

func getQueriedCandidateList(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) Programs {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryCandidateList}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryCandidateList}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var candidateList Programs
	err2 := cdc.UnmarshalJSON(bz, &candidateList)
	require.Nil(t, err2)

	return candidateList
}

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) Params {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryParams}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryParams}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var params Params
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func TestQueryParams(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.budgetKeeper)

	params := DefaultParams()
	input.budgetKeeper.SetParams(input.ctx, params)

	queriedParams := getQueriedParams(t, input.ctx, input.cdc, querier)

	require.Equal(t, queriedParams, params)
}

func TestQueryProgram(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.budgetKeeper)

	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

	queriedProgram := getQueriedProgram(t, input.ctx, input.cdc, querier, testProgram.ProgramID)

	require.Equal(t, queriedProgram, testProgram)
}

func TestQueryVotes(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.budgetKeeper)

	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

	var votes Votes
	for _, addr := range addrs {
		vote := NewVote(testProgram.ProgramID, true, addr)
		votes = append(votes, vote)

		input.budgetKeeper.AddVote(input.ctx, vote.ProgramID, vote.Voter, vote.Option)
	}

	// queriedVotes without filter
	queriedVotes := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.AccAddress{}, 0)
	require.Equal(t, len(queriedVotes), len(votes))

	// queriedVotes with programID filter
	queriedVotesWithProgramID := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.AccAddress{}, testProgram.ProgramID)
	require.Equal(t, len(queriedVotesWithProgramID), len(votes))

	// queriedVotes with voter filter
	queriedVotesWithVoter := getQueriedVotes(t, input.ctx, input.cdc, querier, addrs[0], 0)
	require.Equal(t, queriedVotesWithVoter, votes[:1])

	// queriedVotes with programID and voter filter
	queriedVotesWithBoth := getQueriedVotes(t, input.ctx, input.cdc, querier, addrs[1], testProgram.ProgramID)
	require.Equal(t, queriedVotesWithBoth, votes[1:2])
}

func TestQueryActiveList(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.budgetKeeper)

	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.StoreProgram(input.ctx, testProgram)

	queriedActiveList := getQueriedActiveList(t, input.ctx, input.cdc, querier)

	require.Equal(t, queriedActiveList, Programs{testProgram})
}

func TestQueryCandidateList(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.budgetKeeper)

	testProgram := generateTestProgram(input.ctx, input.budgetKeeper)
	input.budgetKeeper.StoreProgram(input.ctx, testProgram)
	input.budgetKeeper.CandQueueInsert(input.ctx, 0, testProgram.ProgramID)

	queriedCandidateList := getQueriedCandidateList(t, input.ctx, input.cdc, querier)

	require.Equal(t, queriedCandidateList, Programs{testProgram})
}
