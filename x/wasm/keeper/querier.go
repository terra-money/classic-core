package keeper

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/terra-money/core/x/wasm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// Params queries params of wasm module
func (q querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

// CodeInfo returns the stored code info
func (q querier) CodeInfo(c context.Context, req *types.QueryCodeInfoRequest) (*types.QueryCodeInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	codeInfo, err := q.GetCodeInfo(ctx, req.CodeId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryCodeInfoResponse{CodeInfo: codeInfo}, nil
}

// ByteCode returns the stored byte code
func (q querier) ByteCode(c context.Context, req *types.QueryByteCodeRequest) (*types.QueryByteCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	byteCode, err := q.GetByteCode(ctx, req.CodeId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if len(byteCode) == 0 {
		return nil, status.Error(codes.NotFound, "Code not found")
	}

	return &types.QueryByteCodeResponse{ByteCode: byteCode}, nil
}

// ContractInfo returns the stored contract info
func (q querier) ContractInfo(c context.Context, req *types.QueryContractInfoRequest) (*types.QueryContractInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	contractAddr, err := sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	contractInfo, err := q.GetContractInfo(ctx, contractAddr)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryContractInfoResponse{ContractInfo: contractInfo}, nil
}

// ContractStore return smart query result from the contract
func (q querier) ContractStore(c context.Context, req *types.QueryContractStoreRequest) (res *types.QueryContractStoreResponse, err error) {
	ctx := sdk.UnwrapSDKContext(c)

	// external query gas limit must be specified here
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(q.wasmConfig.ContractQueryGasLimit))

	var contractAddr sdk.AccAddress
	contractAddr, err = sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// recover from out-of-gas panic
	defer func() {

		if r := recover(); r != nil {
			switch rType := r.(type) {
			// TODO: Use ErrOutOfGas instead of ErrorOutOfGas which would allow us
			// to keep the stracktrace.
			case sdk.ErrorOutOfGas:
				err = sdkerrors.Wrap(
					sdkerrors.ErrOutOfGas, fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, ctx.GasMeter().Limit(), ctx.GasMeter().GasConsumed(),
					),
				)

			default:
				err = sdkerrors.Wrap(
					sdkerrors.ErrPanic, fmt.Sprintf(
						"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
					),
				)
			}

			res = nil
		}
	}()

	bz, err := q.queryToContract(ctx, contractAddr, req.QueryMsg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res = &types.QueryContractStoreResponse{
		QueryResult: bz,
	}

	return res, nil
}

// RawStore return single key from the raw store data of a contract
func (q querier) RawStore(c context.Context, req *types.QueryRawStoreRequest) (*types.QueryRawStoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	contractAddr, err := sdk.AccAddressFromBech32(req.ContractAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res := q.queryToStore(ctx, contractAddr, req.Key)
	return &types.QueryRawStoreResponse{
		Data: res,
	}, nil
}
