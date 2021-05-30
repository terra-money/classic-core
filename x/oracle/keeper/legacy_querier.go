package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-money/core/x/oracle/types"
)

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryExchangeRate:
			return queryExchangeRate(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryExchangeRates:
			return queryExchangeRates(ctx, keeper, legacyQuerierCdc)
		case types.QueryActives:
			return queryActives(ctx, keeper, legacyQuerierCdc)
		case types.QueryParameters:
			return queryParameters(ctx, keeper, legacyQuerierCdc)
		case types.QueryFeederDelegation:
			return queryFeederDelegation(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryMissCounter:
			return queryMissCounter(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryAggregatePrevote:
			return queryAggregatePrevote(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryAggregatePrevotes:
			return queryAggregatePrevotes(ctx, keeper, legacyQuerierCdc)
		case types.QueryAggregateVote:
			return queryAggregateVote(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryAggregateVotes:
			return queryAggregateVotes(ctx, keeper, legacyQuerierCdc)
		case types.QueryVoteTargets:
			return queryVoteTargets(ctx, keeper, legacyQuerierCdc)
		case types.QueryTobinTax:
			return queryTobinTax(ctx, req, keeper, legacyQuerierCdc)
		case types.QueryTobinTaxes:
			return queryTobinTaxes(ctx, keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryExchangeRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryExchangeRateParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	rate, err := keeper.GetLunaExchangeRate(ctx, params.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnknownDenom, params.Denom)
	}

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, rate)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryExchangeRates(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var rates sdk.DecCoins

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, rates)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryActives(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	denoms := []string{}

	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryFeederDelegationParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegate := keeper.GetFeederDelegation(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryMissCounter(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryMissCounterParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	missCounter := keeper.GetMissCounter(ctx, params.Validator)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, missCounter)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregatePrevote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAggregatePrevoteParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	aggregateExchangeRatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, aggregateExchangeRatePrevote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregatePrevotes(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var aggregatePrevotes []types.AggregateExchangeRatePrevote
	keeper.IterateAggregateExchangeRatePrevotes(ctx, func(_ sdk.ValAddress, aggregatePrevote types.AggregateExchangeRatePrevote) bool {
		aggregatePrevotes = append(aggregatePrevotes, aggregatePrevote)
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, aggregatePrevotes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregateVote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAggregateVoteParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	aggregateExchangeRateVote, err := keeper.GetAggregateExchangeRateVote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, aggregateExchangeRateVote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAggregateVotes(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var aggregateVotes []types.AggregateExchangeRateVote
	keeper.IterateAggregateExchangeRateVotes(ctx, func(_ sdk.ValAddress, aggregateVote types.AggregateExchangeRateVote) bool {
		aggregateVotes = append(aggregateVotes, aggregateVote)
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, aggregateVotes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryVoteTargets(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	voteTargets := keeper.GetVoteTargets(ctx)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, voteTargets)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTobinTax(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTobinTaxParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	tobinTax, err := keeper.GetTobinTax(ctx, params.Denom)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, tobinTax)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTobinTaxes(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var denoms types.DenomList

	keeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		denoms = append(denoms, types.Denom{Name: denom, TobinTax: tobinTax})
		return false
	})

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, denoms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
