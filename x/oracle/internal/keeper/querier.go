package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-project/core/x/oracle/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryPrice:
			return queryPrice(ctx, req, keeper)
		case types.QueryActives:
			return queryActives(ctx, keeper)
		case types.QueryVotes:
			return queryVotes(ctx, req, keeper)
		case types.QueryPrevotes:
			return queryPrevotes(ctx, req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		case types.QueryFeederDelegation:
			return queryFeederDelegation(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryPrice(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryPriceParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	price, err := keeper.GetLunaPrice(ctx, params.Denom)
	if err != nil {
		return nil, types.ErrUnknownDenomination(types.DefaultCodespace, params.Denom)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, price)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryActives(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	denoms := keeper.GetActiveDenoms(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, denoms)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryVotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredVotes := types.PriceVotes{}

	// collects all votes without filter
	prefix := types.VoteKey
	handler := func(vote types.PriceVote) (stop bool) {
		filteredVotes = append(filteredVotes, vote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = types.GetVoteKey(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = types.GetVoteKey(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(vote types.PriceVote) (stop bool) {

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
	var params types.QueryPrevotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	filteredPrevotes := types.PricePrevotes{}

	// collects all votes without filter
	prefix := types.PrevoteKey
	handler := func(prevote types.PricePrevote) (stop bool) {
		filteredPrevotes = append(filteredPrevotes, prevote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = types.GetPrevoteKey(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = types.GetPrevoteKey(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(prevote types.PricePrevote) (stop bool) {

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

func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeederDelegationParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	delegatee := keeper.GetFeedDelegate(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, delegatee)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
