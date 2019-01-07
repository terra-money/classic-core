package budget

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the budgeternance Querier
const (
	QueryParams  = "params"
	QueryProgram = "program"
	QueryVotes   = "votes"
	QueryVote    = "vote"
	QueryTally   = "tally"

	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, path[1:], req, keeper)
		case QueryProgram:
			return queryProgram(ctx, path[1:], req, keeper)
		case QueryVotes:
			return queryVotes(ctx, path[1:], req, keeper)
		case QueryVote:
			return queryVote(ctx, path[1:], req, keeper)
		case QueryTally:
			return queryTally(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown budget query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// Params for queries:
// - 'custom/budget/program'
// - 'custom/budget/deposits'
// - 'custom/budget/tally'
// - 'custom/budget/votes'
type QueryProgramParams struct {
	ProgramID uint64
}

// creates a new instance of QueryProgramParams
func NewQueryProgramParams(ProgramID uint64) QueryProgramParams {
	return QueryProgramParams{
		ProgramID: ProgramID,
	}
}

// nolint: unparam
func queryProgram(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryProgramParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	program, err := keeper.GetProgram(ctx, params.ProgramID)
	if err != nil {
		return nil, ErrProgramNotFound(params.ProgramID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, program)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// Params for query 'custom/budget/vote'
type QueryVoteParams struct {
	ProgramID uint64
	Voter     sdk.AccAddress
}

// creates a new instance of QueryVoteParams
func NewQueryVoteParams(ProgramID uint64, voter sdk.AccAddress) QueryVoteParams {
	return QueryVoteParams{
		ProgramID: ProgramID,
		Voter:     voter,
	}
}

// nolint: unparam
func queryVote(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	vote, _ := keeper.GetVote(ctx, params.ProgramID, params.Voter)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, vote)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryTally(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryProgramParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	ProgramID := params.ProgramID

	program, err := keeper.GetProgram(ctx, ProgramID)
	if err != nil {
		return nil, ErrProgramNotFound(ProgramID)
	}

	tallyResult := program.TallyResult

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tallyResult)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryVotes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryProgramParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	var votes []string
	votesIterator := sdk.KVStorePrefixIterator(ctx.KVStore(keeper.key), PrefixVote)
	for ; votesIterator.Valid(); votesIterator.Next() {
		var key string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Key(), &key)

		var option string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &option)
		votes = append(votes, fmt.Sprintf("%s:%s:%s", key, option))
	}
	votesIterator.Close()

	bz, err := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
