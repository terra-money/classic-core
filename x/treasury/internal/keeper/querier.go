package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryCurrentEpoch:
			return queryCurrentEpoch(ctx, keeper)
		case types.QueryTaxRate:
			return queryTaxRate(ctx, req, keeper)
		case types.QueryTaxCap:
			return queryTaxCap(ctx, req, keeper)
		case types.QueryRewardWeight:
			return queryRewardWeight(ctx, req, keeper)
		case types.QuerySeigniorageProceeds:
			return querySeigniorageProceeds(ctx, req, keeper)
		case types.QueryTaxProceeds:
			return queryTaxProceeds(ctx, req, keeper)
		case types.QueryHistoricalIssuance:
			return queryHistoricalIssuance(ctx, req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown market query endpoint")
		}
	}
}

func queryCurrentEpoch(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	curEpoch := core.GetEpoch(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, curEpoch)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTaxRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTaxRateParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	curEpoch := core.GetEpoch(ctx)
	if 0 > params.Epoch || curEpoch < params.Epoch {
		return nil, types.ErrInvalidEpoch(types.DefaultCodespace, curEpoch, params.Epoch)
	}

	taxRate := keeper.GetTaxRate(ctx, params.Epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryTaxCap(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTaxCapParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	taxCap := keeper.GetTaxCap(ctx, params.Denom)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxCap)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryRewardWeight(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryRewardWeightParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	curEpoch := core.GetEpoch(ctx)
	if 0 > params.Epoch || curEpoch < params.Epoch {
		return nil, types.ErrInvalidEpoch(types.DefaultCodespace, curEpoch, params.Epoch)
	}

	taxRate := keeper.GetRewardWeight(ctx, params.Epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func querySeigniorageProceeds(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySeigniorageProceedsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	curEpoch := core.GetEpoch(ctx)
	if 0 > params.Epoch || curEpoch < params.Epoch {
		return nil, types.ErrInvalidEpoch(types.DefaultCodespace, curEpoch, params.Epoch)
	}

	seigniorage := keeper.PeekEpochSeigniorage(ctx, params.Epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, seigniorage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTaxProceeds(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTaxProceedsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	curEpoch := core.GetEpoch(ctx)
	if 0 > params.Epoch || curEpoch < params.Epoch {
		return nil, types.ErrInvalidEpoch(types.DefaultCodespace, curEpoch, params.Epoch)
	}

	proceeds := keeper.PeekTaxProceeds(ctx, params.Epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, proceeds)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryHistoricalIssuance(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryHistoricalIssuanceParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	curEpoch := core.GetEpoch(ctx)
	if 0 > params.Epoch || curEpoch < params.Epoch {
		return nil, types.ErrInvalidEpoch(types.DefaultCodespace, curEpoch, params.Epoch)
	}

	issuance := keeper.GetHistoricalIssuance(ctx, params.Epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, issuance)
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
