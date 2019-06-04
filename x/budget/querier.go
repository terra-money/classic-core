package budget

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryProgram       = "program"
	QueryVotes         = "votes"
	QueryActiveList    = "active-list"
	QueryCandidateList = "candidate-list"
	QueryParams        = "params"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProgram:
			return queryProgram(ctx, path[1:], req, keeper)
		case QueryVotes:
			return queryVotes(ctx, req, keeper)
		case QueryActiveList:
			return queryActiveList(ctx, req, keeper)
		case QueryCandidateList:
			return queryCandidateList(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

// nolint: unparam
func queryProgram(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	programIDStr := path[0]
	programIDInt, strConvertError := strconv.Atoi(programIDStr)
	if strConvertError != nil {
		return nil, sdk.ErrInternal("ProgramID must be a valid int")
	}

	programID := uint64(programIDInt)
	program, pErr := keeper.GetProgram(ctx, programID)
	if pErr != nil {
		return nil, pErr
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, program)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result to JSON")
	}

	return bz, nil
}

// Params for query 'custom/oracle/votes'
type QueryVotesParams struct {
	Voter     sdk.AccAddress
	ProgramID uint64
}

// creates a new instance of QueryVoteParams
func NewQueryVotesParams(voter sdk.AccAddress, programID uint64) QueryVotesParams {
	return QueryVotesParams{
		Voter:     voter,
		ProgramID: programID,
	}
}

// nolint: unparam
func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredVotes := Votes{}
	prefix := prefixVote
	handler := func(programID uint64, voter sdk.AccAddress, option bool) (stop bool) {
		vote := NewVote(programID, option, voter)
		filteredVotes = append(filteredVotes, vote)

		return false
	}

	if params.ProgramID != 0 && !params.Voter.Empty() {
		prefix = keyVote(params.ProgramID, params.Voter)
	} else if params.ProgramID != 0 {
		prefix = keyVote(params.ProgramID, sdk.AccAddress{})
	} else if !params.Voter.Empty() {
		handler = func(programID uint64, voter sdk.AccAddress, option bool) (stop bool) {
			if params.Voter.Equals(voter) {
				vote := NewVote(programID, option, voter)
				filteredVotes = append(filteredVotes, vote)
			}

			return false
		}
	}

	keeper.IterateVotesWithPrefix(ctx, prefix, handler)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, filteredVotes)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

// nolint: unparam
func queryActiveList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	programs := Programs{}
	keeper.IteratePrograms(ctx, true, func(program Program) (stop bool) {
		programs = append(programs, program)
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, programs)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

// nolint: unparam
func queryCandidateList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	programs := Programs{}

	keeper.CandQueueIterate(ctx, func(programID uint64) (stop bool) {
		program, err := keeper.GetProgram(ctx, programID)
		if err != nil {
			return false
		}

		programs = append(programs, program)
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, programs)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
