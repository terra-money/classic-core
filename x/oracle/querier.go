package oracle

import (
	"strings"

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

// JSON response format
type QueryPriceResponse struct {
	Price sdk.Dec `json:"price"`
}

func (r QueryPriceResponse) String() (out string) {
	out = r.Price.String()
	return strings.TrimSpace(out)
}

func queryPrice(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]

	price, err := keeper.GetLunaSwapRate(ctx, denom)
	if err != nil {
		return nil, ErrUnknownDenomination(DefaultCodespace, denom)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, QueryPriceResponse{Price: price})
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

// JSON response format
type QueryActiveResponse struct {
	Actives DenomList `json:"actives"`
}

func (r QueryActiveResponse) String() (out string) {
	out = r.Actives.String()
	return strings.TrimSpace(out)
}

func queryActive(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denoms := keeper.getActiveDenoms(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryActiveResponse{Actives: denoms})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// QueryVoteParams for query 'custom/oracle/votes'
type QueryVotesParams struct {
	Voter sdk.ValAddress
	Denom string
}

// NewQueryVotesParams creates a new instance of QueryVotesParams
func NewQueryVotesParams(voter sdk.ValAddress, denom string) QueryVotesParams {
	return QueryVotesParams{
		Voter: voter,
		Denom: denom,
	}
}

// JSON response format
type QueryVotesResponse struct {
	Votes PriceVotes `json:"votes"`
}

func (r QueryVotesResponse) String() (out string) {
	out = r.Votes.String()
	return strings.TrimSpace(out)
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryVotesParams
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

	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryVotesResponse{Votes: filteredVotes})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// QueryPrevotesParams for query 'custom/oracle/prevotes'
type QueryPrevotesParams QueryVotesParams

// NewQueryPrevotesParams creates a new instance of QueryPrevotesParams
func NewQueryPrevotesParams(voter sdk.ValAddress, denom string) QueryPrevotesParams {
	return QueryPrevotesParams{
		Voter: voter,
		Denom: denom,
	}
}

// JSON response format
type QueryPrevotesResponse struct {
	Prevotes PricePrevotes `json:"prevotes"`
}

func (r QueryPrevotesResponse) String() (out string) {
	out = r.Prevotes.String()
	return strings.TrimSpace(out)
}

func queryPrevotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPrevotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredPrevotes := PricePrevotes{}

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

	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryPrevotesResponse{Prevotes: filteredPrevotes})
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

// JSON response format
type QueryFeederDelegationResponse struct {
	Delegatee sdk.AccAddress `json:"delegatee"`
}

func (r QueryFeederDelegationResponse) String() (out string) {
	out = r.Delegatee.String()
	return strings.TrimSpace(out)
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryFeederDelegationParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	delegatee := keeper.GetFeedDelegate(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, QueryFeederDelegationResponse{Delegatee: delegatee})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
