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
	QueryMRL                = "mrl"
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
			return queryMiningRewardWeight(ctx, path[1:], req, keeper)
		case QueryBalance:
			return queryTreasuryBalance(ctx, path[1:], req, keeper)
		case QueryActiveClaims:
			return queryActiveClaims(ctx, req, keeper)
		case QueryIssuance:
			return queryIssunace(ctx, path[1:], req, keeper)
		case QueryMRL:
			return queryMRL(ctx, path[1:], req, keeper)
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
func queryMRL(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	mrl := MRL(ctx, keeper, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, mrl)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// nolint: unparam
func queryTreasuryBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	epoch, ok := sdk.NewIntFromString(path[0])
	if !ok {
		return nil, sdk.ErrInternal("epoch parameter is not correctly formatted")
	}

	pool := keeper.mtk.PeekSeignioragePool(ctx, epoch)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryActiveClaims(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	claims := types.ClaimPool{}
	keeper.IterateClaims(ctx, func(claim types.Claim) (stop bool) {
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
