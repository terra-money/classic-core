package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-project/core/x/treasury/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryTaxRate:
			return queryTaxRate(ctx, keeper)
		case types.QueryTaxCap:
			return queryTaxCap(ctx, req, keeper)
		case types.QueryRewardWeight:
			return queryRewardWeight(ctx, keeper)
		case types.QuerySeigniorageProceeds:
			return querySeigniorageProceeds(ctx, keeper)
		case types.QueryTaxProceeds:
			return queryTaxProceeds(ctx, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown treasury query endpoint")
		}
	}
}

func queryCurrentEpoch(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	curEpoch := keeper.GetEpoch(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, curEpoch)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTaxRate(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	taxRate := keeper.GetTaxRate(ctx)
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

func queryRewardWeight(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	taxRate := keeper.GetRewardWeight(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func querySeigniorageProceeds(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	seigniorage := keeper.PeekEpochSeigniorage(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, seigniorage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTaxProceeds(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	proceeds := keeper.PeekEpochTaxProceeds(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, proceeds)
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
