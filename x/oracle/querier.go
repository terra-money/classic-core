package oracle

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the oracle Querier
const (
	QueryPrice            = "price"
	QueryVotes            = "votes"
	QueryPrevotes         = "prevotes"
	QueryActive           = "active"
	QueryParams           = "params"
	QueryFeederDelegation = "feeder"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPrice:
			return queryPrice(ctx, path[1:], req, keeper)
		case QueryActive:
			return queryActive(ctx, req, keeper)
		case QueryVotes:
			return queryVotes(ctx, req, keeper)
		case QueryPrevotes:
			return queryPrevotes(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		case QueryFeederDelegation:
			return queryFeederDelegation(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryPrice(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]

	price, err := keeper.GetLunaSwapRate(ctx, denom)
	if err != nil {
		return nil, ErrUnknownDenomination(DefaultCodespace, denom)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, price)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryActive(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denoms := keeper.getActiveDenoms(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, denoms)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// QueryVoteParams for query 'custom/oracle/votes'
type QueryVoteParams struct {
	Voter sdk.ValAddress
	Denom string
}

// NewQueryVoteParams creates a new instance of QueryVoteParams
func NewQueryVoteParams(voter sdk.ValAddress, denom string) QueryVoteParams {
	return QueryVoteParams{
		Voter: voter,
		Denom: denom,
	}
}

// QueryPrevoteParams for query 'custom/oracle/prevotes'
type QueryPrevoteParams QueryVoteParams

// NewQueryPrevoteParams creates a new instance of QueryVoteParams
func NewQueryPrevoteParams(voter sdk.ValAddress, denom string) QueryPrevoteParams {
	return QueryPrevoteParams{
		Voter: voter,
		Denom: denom,
	}
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredVotes := PriceVotes{}

	// collects all votes without filter
	prefix := prefixVote
	handler := func(vote PriceVote) (stop bool) {
		filteredVotes = append(filteredVotes, vote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = keyVote(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = keyVote(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(vote PriceVote) (stop bool) {

			if vote.Voter.Equals(params.Voter) {
				filteredVotes = append(filteredVotes, vote)
			}

			return false
		}
	}

	keeper.iterateVotesWithPrefix(ctx, prefix, handler)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, filteredVotes)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryPrevotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredPrevotes := []PricePrevote{}

	// collects all votes without filter
	prefix := prefixPrevote
	handler := func(prevote PricePrevote) (stop bool) {
		filteredPrevotes = append(filteredPrevotes, prevote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = keyPrevote(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = keyPrevote(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(prevote PricePrevote) (stop bool) {

			if prevote.Voter.Equals(params.Voter) {
				filteredPrevotes = append(filteredPrevotes, prevote)
			}

			return false
		}
	}

	keeper.iteratePrevotesWithPrefix(ctx, prefix, handler)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, filteredPrevotes)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
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

// QueryFeederDelegationParams for query 'custom/oracle/feeder-delegation'
type QueryFeederDelegationParams struct {
	Validator sdk.ValAddress
}

// NewQueryFeederDelegationParams creates a new instance of QueryFeederDelegationParams
func NewQueryFeederDelegationParams(validator sdk.ValAddress) QueryFeederDelegationParams {
	return QueryFeederDelegationParams{
		Validator: validator,
	}
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryFeederDelegationParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetFeedDelegate(ctx, params.Validator))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
