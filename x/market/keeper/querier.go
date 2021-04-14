package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/market/types"
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

// Params queries params of market module
func (q querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

// Swap queries for swap simulation
func (q querier) Swap(c context.Context, req *types.QuerySwapRequest) (*types.QuerySwapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err := sdk.ValidateDenom(req.AskDenom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ask denom")
	}

	if err := req.OfferCoin.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	retCoin, err := q.simulateSwap(ctx, req.OfferCoin, req.AskDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySwapResponse{ReturnCoin: retCoin}, nil
}

// MintPoolDelta queries mint pool delta
func (q querier) MintPoolDelta(c context.Context, req *types.QueryMintPoolDeltaRequest) (*types.QueryMintPoolDeltaResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	mintPoolDelta := q.GetMintPoolDelta(ctx)
	return &types.QueryMintPoolDeltaResponse{MintPoolDelta: mintPoolDelta}, nil
}

// BurnPoolDelta queries burn pool delta
func (q querier) BurnPoolDelta(c context.Context, req *types.QueryBurnPoolDeltaRequest) (*types.QueryBurnPoolDeltaResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	burnPoolDelta := q.GetBurnPoolDelta(ctx)
	return &types.QueryBurnPoolDeltaResponse{BurnPoolDelta: burnPoolDelta}, nil
}
