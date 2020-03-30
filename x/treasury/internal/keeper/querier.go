package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-project/core/x/treasury/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryCurrentEpoch(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	curEpoch := keeper.GetEpoch(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, curEpoch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTaxRate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	taxRate := keeper.GetTaxRate(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTaxCap(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTaxCapParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	taxCap := keeper.GetTaxCap(ctx, params.Denom)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxCap)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRewardWeight(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	taxRate := keeper.GetRewardWeight(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func querySeigniorageProceeds(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	seigniorage := keeper.PeekEpochSeigniorage(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, seigniorage)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTaxProceeds(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	proceeds := keeper.PeekEpochTaxProceeds(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, proceeds)
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
