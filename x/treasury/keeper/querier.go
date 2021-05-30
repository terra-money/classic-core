package keeper

import (
	"context"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/types"
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

// Params queries params of distribution module
func (q querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

// TaxRate return the current tax rate
func (q querier) TaxRate(c context.Context, req *types.QueryTaxRateRequest) (*types.QueryTaxRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryTaxRateResponse{TaxRate: q.GetTaxRate(ctx)}, nil
}

// TaxCap returns the tax cap of a denom
func (q querier) TaxCap(c context.Context, req *types.QueryTaxCapRequest) (*types.QueryTaxCapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryTaxCapResponse{TaxCap: q.GetTaxCap(ctx, req.Denom)}, nil
}

// TaxCaps returns the all tax caps
func (q querier) TaxCaps(c context.Context, req *types.QueryTaxCapsRequest) (*types.QueryTaxCapsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var taxCaps []types.QueryTaxCapsResponseItem
	q.IterateTaxCap(ctx, func(denom string, taxCap sdk.Int) bool {
		taxCaps = append(taxCaps, types.QueryTaxCapsResponseItem{
			Denom:  denom,
			TaxCap: taxCap,
		})
		return false
	})

	return &types.QueryTaxCapsResponse{TaxCaps: taxCaps}, nil
}

// RewardWeight return the current reward weight
func (q querier) RewardWeight(c context.Context, req *types.QueryRewardWeightRequest) (*types.QueryRewardWeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryRewardWeightResponse{RewardWeight: q.GetRewardWeight(ctx)}, nil
}

// SeigniorageProceeds return the current seigniorage proceeds
func (q querier) SeigniorageProceeds(c context.Context, req *types.QuerySeigniorageProceedsRequest) (*types.QuerySeigniorageProceedsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QuerySeigniorageProceedsResponse{SeigniorageProceeds: q.PeekEpochSeigniorage(ctx)}, nil
}

// TaxProceeds return the current tax proceeds
func (q querier) TaxProceeds(c context.Context, req *types.QueryTaxProceedsRequest) (*types.QueryTaxProceedsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryTaxProceedsResponse{TaxProceeds: q.PeekEpochTaxProceeds(ctx)}, nil
}

// Indicators return the current trl informations
func (q querier) Indicators(c context.Context, req *types.QueryIndicatorsRequest) (*types.QueryIndicatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Compute Total Staked Luna (TSL)
	TSL := q.stakingKeeper.TotalBondedTokens(ctx)

	// Compute Tax Rewards (TR)
	taxRewards := sdk.NewDecCoinsFromCoins(q.PeekEpochTaxProceeds(ctx)...)
	TR := q.alignCoins(ctx, taxRewards, core.MicroSDRDenom)

	epoch := q.GetEpoch(ctx)
	var res types.QueryIndicatorsResponse
	if epoch == 0 {
		res = types.QueryIndicatorsResponse{
			TRLYear:  TR.QuoInt(TSL),
			TRLMonth: TR.QuoInt(TSL),
		}
	} else {
		params := q.GetParams(ctx)
		previousEpochCtx := ctx.WithBlockHeight(ctx.BlockHeight() - int64(core.BlocksPerWeek))
		trlYear := q.rollingAverageIndicator(previousEpochCtx, int64(params.WindowLong-1), TRL)
		trlMonth := q.rollingAverageIndicator(previousEpochCtx, int64(params.WindowShort-1), TRL)

		computedEpochForYear := int64(math.Min(float64(params.WindowLong-1), float64(epoch)))
		computedEpochForMonty := int64(math.Min(float64(params.WindowShort-1), float64(epoch)))

		trlYear = trlYear.MulInt64(computedEpochForYear).Add(TR.QuoInt(TSL)).QuoInt64(computedEpochForYear + 1)
		trlMonth = trlMonth.MulInt64(computedEpochForMonty).Add(TR.QuoInt(TSL)).QuoInt64(computedEpochForMonty + 1)

		res = types.QueryIndicatorsResponse{
			TRLYear:  trlYear,
			TRLMonth: trlMonth,
		}
	}

	return &res, nil
}
