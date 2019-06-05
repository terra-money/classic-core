package market

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the oracle Querier
const (
	QuerySwap   = "swap"
	QueryParams = "params"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QuerySwap:
			return querySwap(ctx, path[1:], req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown market query endpoint")
		}
	}
}

// QuerySwapParams for query 'custom/market/swap'
type QuerySwapParams struct {
	OfferCoin sdk.Coin
}

func NewQuerySwapParams(offerCoin sdk.Coin) QuerySwapParams {
	return QuerySwapParams{
		OfferCoin: offerCoin,
	}
}

func querySwap(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	askDenom := path[0]

	var params QuerySwapParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	swapCoin, spread, err := keeper.GetSwapCoin(ctx, params.OfferCoin, askDenom, false)
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

func queryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
