package budget

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the budgeternance Querier
const (
	QueryParams   = "params"
	QueryPrograms = "programs"
	QueryProgram  = "program"
	QueryVotes    = "votes"
	QueryVote     = "vote"
	QueryTally    = "tally"

	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, path[1:], req, keeper)
		case QueryPrograms:
			return queryPrograms(ctx, path[1:], req, keeper)
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
	switch path[0] {
	case ParamVoting:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetVotingParams(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case ParamTallying:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetTallyParams(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
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

	program := keeper.GetProgram(ctx, params.ProgramID)
	if program == nil {
		return nil, ErrUnknownProgram(DefaultCodespace, params.ProgramID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, program)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// Params for query 'custom/budget/deposit'
type QueryDepositParams struct {
	ProgramID uint64
	Depositor sdk.AccAddress
}

// creates a new instance of QueryDepositParams
func NewQueryDepositParams(ProgramID uint64, depositor sdk.AccAddress) QueryDepositParams {
	return QueryDepositParams{
		ProgramID: ProgramID,
		Depositor: depositor,
	}
}

// nolint: unparam
func queryDeposit(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryDepositParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	deposit, _ := keeper.GetDeposit(ctx, params.ProgramID, params.Depositor)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, deposit)
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

	program := keeper.GetProgram(ctx, ProgramID)
	if program == nil {
		return nil, ErrUnknownProgram(DefaultCodespace, ProgramID)
	}

	var tallyResult TallyResult

	if program.GetStatus() == StatusDepositPeriod {
		tallyResult = EmptyTallyResult()
	} else if program.GetStatus() == StatusPassed || program.GetStatus() == StatusRejected {
		tallyResult = program.GetTallyResult()
	} else {
		// program is in voting period
		_, tallyResult = tally(ctx, keeper, program)
	}

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

	var votes []Vote
	votesIterator := keeper.GetVotes(ctx, params.ProgramID)
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := Vote{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(votesIterator.Value(), &vote)
		votes = append(votes, vote)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// Params for query 'custom/budget/program'
type QueryProgramParams struct {
	Voter        sdk.AccAddress
	ProgramState ProgramState
	Limit        uint64
}

// creates a new instance of QueryProgramParams
func NewQueryProgramParams(state ProgramState, limit uint64, voter) QueryProgramParams {
	return QueryProgramParams{
		Voter:        voter,
		ProgramState: state,
		Limit:        limit,
	}
}

// nolint: unparam
func queryPrograms(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryProgramsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	Programs := keeper.GetProgramsFiltered(ctx, params.Voter, params.Depositor, params.ProgramStatus, params.Limit)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, Programs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
