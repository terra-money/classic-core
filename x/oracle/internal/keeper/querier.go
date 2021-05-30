package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
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
		case types.QueryAggregatePrevote:
			return queryAggregatePrevote(ctx, req, keeper)
		case types.QueryAggregateVote:
			return queryAggregateVote(ctx, req, keeper)
		case types.QueryVoteTargets:
			return queryVoteTargets(ctx, keeper)
		case types.QueryTobinTax:
			return queryTobinTax(ctx, req, keeper)
		case types.QueryTobinTaxes:
			return queryTobinTaxes(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryExchangeRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryExchangeRateParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	rate, err := keeper.GetLunaExchangeRate(ctx, params.Denom)
	if err != nil {
		if core.IsWaitingForSoftfork(ctx, 1) {
			return nil, sdkerrors.Wrap(types.ErrInternal, "unknown denom")
		}

		return nil, sdkerrors.Wrap(types.ErrUnknownDenom, params.Denom)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, rate)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryExchangeRates(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var rates sdk.DecCoins

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rates)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryActives(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	denoms := []string{}

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryVotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryPrevotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryPrevotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryFeederDelegationParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegate := keeper.GetOracleDelegate(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, delegate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryMissCounter(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryMissCounterParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	missCounter := keeper.GetMissCounter(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, missCounter)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregatePrevote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryAggregatePrevoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	aggregateExchangeRatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, aggregateExchangeRatePrevote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregateVote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryAggregateVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	aggregateExchangeRateVote, err := keeper.GetAggregateExchangeRateVote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, aggregateExchangeRateVote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryVoteTargets(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	voteTargets := keeper.GetVoteTargets(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, voteTargets)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTobinTax(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTobinTaxParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	tobinTax, err := keeper.GetTobinTax(ctx, params.Denom)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tobinTax)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTobinTaxes(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var denoms types.DenomList

	keeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		denoms = append(denoms, types.Denom{Name: denom, TobinTax: tobinTax})
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
