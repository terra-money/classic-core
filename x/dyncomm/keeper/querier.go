package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/dyncomm/types"
)

// querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over q
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the market QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

var _ types.QueryServer = querier{}

// Params queries params of dyncomm module
func (q querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

// Rates queries Validator Rate of dyncomm module
func (q querier) Rate(c context.Context, req *types.QueryRateRequest) (*types.QueryRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	rate := q.GetDynCommissionRate(ctx, req.ValidatorAddr)
	target := q.GetTargetCommissionRate(ctx, req.ValidatorAddr)
	return &types.QueryRateResponse{Rate: &rate, Target: &target}, nil
}
