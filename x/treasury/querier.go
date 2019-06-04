package treasury

import (
	"github.com/terra-project/core/types"
	"github.com/terra-project/core/types/util"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the treasury Querier
const (
	QueryTaxRate             = "tax-rate"
	QueryTaxCap              = "tax-cap"
	QueryMiningRewardWeight  = "reward-weight"
	QuerySeigniorageProceeds = "seigniorage-proceeds"
	QueryActiveClaims        = "active-claims"
	QueryCurrentEpoch        = "current-epoch"
	QueryParams              = "params"
	QueryIssuance            = "issuance"
	QueryTaxProceeds         = "tax-proceeds"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTaxRate:
			return queryTaxRate(ctx, path[1:], req, keeper)
		case QueryTaxCap:
			return queryTaxCap(ctx, path[1:], req, keeper)
		case QueryMiningRewardWeight:
			return queryMiningRewardWeight(ctx, path[1:], req, keeper)
		case QueryTaxProceeds:
			return queryTaxProceeds(ctx, path[1:], req, keeper)
		case QuerySeigniorageProceeds:
			return querySeigniorageProceeds(ctx, path[1:], req, keeper)
		case QueryIssuance:
			return queryIssuance(ctx, path[1:], req, keeper)
		case QueryCurrentEpoch:
			return queryCurrentEpoch(ctx, req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown treasury query endpoint")
		}
	}
}

// nolint: unparam
func queryTaxRate(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	taxRate := keeper.GetTaxRate(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxRate)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryTaxCap(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]
	taxCap := keeper.GetTaxCap(ctx, denom)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, taxCap)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryIssuance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]

	curDay := ctx.BlockHeight() / util.BlocksPerDay
	issuance := keeper.mtk.GetIssuance(ctx, denom, sdk.NewInt(curDay))
	bz, err := codec.MarshalJSONIndent(keeper.cdc, issuance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryMiningRewardWeight(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	rewardWeight := keeper.GetRewardWeight(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, rewardWeight)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryTaxProceeds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	pool := keeper.PeekTaxProceeds(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func querySeigniorageProceeds(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	pool := keeper.mtk.PeekEpochSeigniorage(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryCurrentEpoch(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	curEpoch := util.GetEpoch(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, curEpoch)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
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
