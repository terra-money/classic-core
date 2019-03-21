package treasury

import (
	"terra/types"
	"terra/types/util"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryTaxRate            = "tax-rate"
	QueryTaxCap             = "tax-cap"
	QueryMiningRewardWeight = "mining-reward-weight"
	QueryBalance            = "balance"
	QueryActiveClaims       = "active-claims"
	QueryRewards            = "rewards"
	QueryParams             = "params"
	QueryIssuance           = "issuance"

	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTaxRate:
			return queryTaxRate(ctx, req, keeper)
		case QueryTaxCap:
			return queryTaxCap(ctx, path[1:], req, keeper)
		case QueryMiningRewardWeight:
			return queryMiningRewardWeight(ctx, req, keeper)
		case QueryBalance:
			return queryTreasuryBalance(ctx, req, keeper)
		case QueryActiveClaims:
			return queryActiveClaims(ctx, req, keeper)
		case QueryIssuance:
			return queryIssunace(ctx, path[1:], req, keeper)
		case QueryParams:
			return queryParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown oracle query endpoint")
		}
	}
}

// nolint: unparam
func queryTaxRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	taxRate := keeper.GetTaxRate(ctx)
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
func queryIssunace(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	denom := path[0]
	issuance := keeper.mtk.GetIssuance(ctx, denom, util.GetEpoch(ctx))
	bz, err := codec.MarshalJSONIndent(keeper.cdc, issuance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryMiningRewardWeight(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	rewardWeight := keeper.GetRewardWeight(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, rewardWeight)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryTreasuryBalance(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	pool := keeper.mtk.PeekSeigniorage(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryActiveClaims(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	claims := []types.Claim{}
	keeper.iterateClaims(ctx, func(claim types.Claim) (stop bool) {
		claims = append(claims, claim)
		return false
	})

	bz, err := codec.MarshalJSONIndent(keeper.cdc, claims)
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
