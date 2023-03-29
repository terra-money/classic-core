package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/classic-terra/core/x/feeshare/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/FeeShare keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// FeeShares returns all FeeShares that have been registered for fee distribution
func (q Querier) FeeShares(
	c context.Context,
	req *types.QueryFeeSharesRequest,
) (*types.QueryFeeSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var feeshares []types.FeeShare
	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.KeyPrefixFeeShare)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var feeshare types.FeeShare
		if err := q.cdc.Unmarshal(value, &feeshare); err != nil {
			return err
		}
		feeshares = append(feeshares, feeshare)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryFeeSharesResponse{
		Feeshare:   feeshares,
		Pagination: pageRes,
	}, nil
}

// FeeShare returns the FeeShare that has been registered for fee distribution for a given
// contract
func (q Querier) FeeShare(
	c context.Context,
	req *types.QueryFeeShareRequest,
) (*types.QueryFeeShareResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// check if the contract is a non-zero hex address
	contract, err := sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid format for contract %s, should be bech32 ('terra...')", req.ContractAddress,
		)
	}

	feeshare, found := q.GetFeeShare(ctx, contract)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"fees registered contract '%s'",
			req.ContractAddress,
		)
	}

	return &types.QueryFeeShareResponse{Feeshare: feeshare}, nil
}

// Params returns the fees module params
func (q Querier) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// DeployerFeeShares returns all contracts that have been registered for fee
// distribution by a given deployer
func (q Querier) DeployerFeeShares( // nolint: dupl
	c context.Context,
	req *types.QueryDeployerFeeSharesRequest,
) (*types.QueryDeployerFeeSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	deployer, err := sdk.AccAddressFromBech32(req.DeployerAddress)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid format for deployer %s, should be bech32 ('terra...')", req.DeployerAddress,
		)
	}

	var contracts []string
	store := prefix.NewStore(
		ctx.KVStore(q.storeKey),
		types.GetKeyPrefixDeployer(deployer),
	)

	pageRes, err := query.Paginate(store, req.Pagination, func(key, _ []byte) error {
		contracts = append(contracts, sdk.AccAddress(key).String())
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDeployerFeeSharesResponse{
		ContractAddresses: contracts,
		Pagination:        pageRes,
	}, nil
}

// WithdrawerFeeShares returns all fees for a given withdraw address
func (q Querier) WithdrawerFeeShares( // nolint: dupl
	c context.Context,
	req *types.QueryWithdrawerFeeSharesRequest,
) (*types.QueryWithdrawerFeeSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	deployer, err := sdk.AccAddressFromBech32(req.WithdrawerAddress)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid format for withdraw addr %s, should be bech32 ('terra...')", req.WithdrawerAddress,
		)
	}

	var contracts []string
	store := prefix.NewStore(
		ctx.KVStore(q.storeKey),
		types.GetKeyPrefixWithdrawer(deployer),
	)

	pageRes, err := query.Paginate(store, req.Pagination, func(key, _ []byte) error {
		contracts = append(contracts, sdk.AccAddress(key).String())

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryWithdrawerFeeSharesResponse{
		ContractAddresses: contracts,
		Pagination:        pageRes,
	}, nil
}
