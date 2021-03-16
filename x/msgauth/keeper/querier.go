package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/msgauth/types"
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

// Grants queries for grants between a granter-grantee pair
func (q querier) Grants(c context.Context, req *types.QueryGrantsRequest) (*types.QueryGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	granterAddr, err := sdk.AccAddressFromBech32(req.Granter)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	granteeAddr, err := sdk.AccAddressFromBech32(req.Grantee)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	grants := q.GetGrants(ctx, granterAddr, granteeAddr)

	return &types.QueryGrantsResponse{Grants: grants}, nil
}

// AllGrants queries for all grants of a granter
func (q querier) AllGrants(c context.Context, req *types.QueryAllGrantsRequest) (*types.QueryAllGrantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	granterAddr, err := sdk.AccAddressFromBech32(req.Granter)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	grants := q.GetAllGrants(ctx, granterAddr)

	return &types.QueryAllGrantsResponse{Grants: grants}, nil
}
