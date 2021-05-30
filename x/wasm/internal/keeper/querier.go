package keeper

import (
	"fmt"
	"runtime/debug"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-money/core/x/wasm/internal/types"
)

// NewQuerier creates a new querier
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetByteCode:
			return queryByteCode(ctx, req, keeper)
		case types.QueryGetCodeInfo:
			return queryCodeInfo(ctx, req, keeper)
		case types.QueryGetContractInfo:
			return queryContractInfo(ctx, req, keeper)
		case types.QueryRawStore:
			return queryRawStore(ctx, req, keeper)
		case types.QueryContractStore:
			return queryContractStore(ctx, req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryByteCode(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryCodeIDParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	byteCode, err := keeper.GetByteCode(ctx, params.CodeID)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, byteCode)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryCodeInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryCodeIDParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	codeInfo, err := keeper.GetCodeInfo(ctx, params.CodeID)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, codeInfo)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryContractInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryContractAddressParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	contractInfo, err := keeper.GetContractInfo(ctx, params.ContractAddress)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, contractInfo)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryRawStore(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryRawStoreParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res := keeper.queryToStore(ctx, params.ContractAddress, params.Key)
	return res, nil
}

func queryContractStore(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (bz []byte, err error) {
	// external query gas limit must be specified here
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(keeper.wasmConfig.ContractQueryGasLimit))

	var params types.QueryContractParams
	err = types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
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

			bz = nil
		}
	}()

	bz, err = keeper.queryToContract(ctx, params.ContractAddress, params.Msg)

	return
}

func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
