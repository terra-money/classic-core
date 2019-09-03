package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-project/core/x/market/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QuerySwap:
			return querySwap(ctx, req, keeper)
		case types.QueryPrevDayIssuance:
			return queryPrevDayIssuance(ctx, req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown market query endpoint")
		}
	}
}

func querySwap(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySwapParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	if params.AskDenom == params.OfferCoin.Denom {
		return nil, types.ErrRecursiveSwap(types.DefaultCodespace, params.AskDenom)
	}

	swapCoin, spread, err := keeper.GetSwapCoin(ctx, params.OfferCoin, params.AskDenom, false)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Failed to get swapped coin amount", err.Error()))
	}

	if spread.IsPositive() {
		swapFeeAmt := spread.MulInt(swapCoin.Amount).TruncateInt()
		if swapFeeAmt.IsPositive() {
			swapFee := sdk.NewCoin(swapCoin.Denom, swapFeeAmt)
			swapCoin = swapCoin.Sub(swapFee)
		}
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, swapCoin)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err2.Error()))
	}

	return bz, nil
}

func queryPrevDayIssuance(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetPrevDayIssuance(ctx))
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
