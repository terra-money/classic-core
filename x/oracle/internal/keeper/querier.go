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
		case types.QueryExchangeRate:
			return queryExchangeRate(ctx, req, keeper)
		case types.QueryExchangeRates:
			return queryExchangeRates(ctx, keeper)
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
		case types.QueryMissCounter:
			return queryMissCounter(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

func queryExchangeRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryExchangeRateParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	rate, err := keeper.GetLunaExchangeRate(ctx, params.Denom)
	if err != nil {
		return nil, types.ErrUnknownDenomination(types.DefaultCodespace, params.Denom)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, rate)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryExchangeRates(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	var rates sdk.DecCoins

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rates)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryActives(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	denoms := []string{}

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

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

	filteredVotes := types.ExchangeRateVotes{}

	// collects all votes without filter
	prefix := types.VoteKey
	handler := func(vote types.ExchangeRateVote) (stop bool) {
		filteredVotes = append(filteredVotes, vote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = types.GetVoteKey(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = types.GetVoteKey(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(vote types.ExchangeRateVote) (stop bool) {

			if vote.Voter.Equals(params.Voter) {
				filteredVotes = append(filteredVotes, vote)
			}

			return false
		}
	}

	keeper.iterateExchangeRateVotesWithPrefix(ctx, prefix, handler)

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

	filteredPrevotes := types.ExchangeRatePrevotes{}

	// collects all prevotes without filter
	prefix := types.PrevoteKey
	handler := func(prevote types.ExchangeRatePrevote) (stop bool) {
		filteredPrevotes = append(filteredPrevotes, prevote)
		return false
	}

	// applies filter
	if len(params.Denom) != 0 && !params.Voter.Empty() {
		prefix = types.GetExchangeRatePrevoteKey(params.Denom, params.Voter)
	} else if len(params.Denom) != 0 {
		prefix = types.GetExchangeRatePrevoteKey(params.Denom, sdk.ValAddress{})
	} else if !params.Voter.Empty() {
		handler = func(prevote types.ExchangeRatePrevote) (stop bool) {

			if prevote.Voter.Equals(params.Voter) {
				filteredPrevotes = append(filteredPrevotes, prevote)
			}

			return false
		}
	}

	keeper.iterateExchangeRatePrevotesWithPrefix(ctx, prefix, handler)

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

	delegatee := keeper.GetOracleDelegate(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, delegatee)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryMissCounter(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMissCounterParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	missCounter := keeper.GetMissCounter(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, missCounter)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
