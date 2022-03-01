package v05

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over q
type querier struct {
}

// NewQuerier returns an implementation of the market QueryServer interface
// for the provided Keeper.
func NewQuerier() QueryServer {
	return &querier{}
}

var _ QueryServer = querier{}

// Params queries params of distribution module
func (q querier) Params(c context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return &QueryParamsResponse{Params: DefaultParams()}, nil
}

// TaxRate return the current tax rate
func (q querier) TaxRate(c context.Context, req *QueryTaxRateRequest) (*QueryTaxRateResponse, error) {
	return &QueryTaxRateResponse{TaxRate: sdk.ZeroDec()}, nil
}

// TaxCap returns the tax cap of a denom
func (q querier) TaxCap(c context.Context, req *QueryTaxCapRequest) (*QueryTaxCapResponse, error) {
	return &QueryTaxCapResponse{TaxCap: sdk.ZeroInt()}, nil
}

// TaxCaps returns the all tax caps
func (q querier) TaxCaps(c context.Context, req *QueryTaxCapsRequest) (*QueryTaxCapsResponse, error) {
	var taxCaps []QueryTaxCapsResponseItem
	return &QueryTaxCapsResponse{TaxCaps: taxCaps}, nil
}

// RewardWeight return the current reward weight
func (q querier) RewardWeight(c context.Context, req *QueryRewardWeightRequest) (*QueryRewardWeightResponse, error) {
	return &QueryRewardWeightResponse{RewardWeight: sdk.ZeroDec()}, nil
}

// SeigniorageProceeds return the current seigniorage proceeds
func (q querier) SeigniorageProceeds(c context.Context, req *QuerySeigniorageProceedsRequest) (*QuerySeigniorageProceedsResponse, error) {
	return &QuerySeigniorageProceedsResponse{SeigniorageProceeds: sdk.ZeroInt()}, nil
}

// TaxProceeds return the current tax proceeds
func (q querier) TaxProceeds(c context.Context, req *QueryTaxProceedsRequest) (*QueryTaxProceedsResponse, error) {
	return &QueryTaxProceedsResponse{TaxProceeds: sdk.Coins{}}, nil
}

// Indicators return the current trl information
func (q querier) Indicators(c context.Context, req *QueryIndicatorsRequest) (*QueryIndicatorsResponse, error) {
	return &QueryIndicatorsResponse{
		TRLYear:  sdk.ZeroDec(),
		TRLMonth: sdk.ZeroDec(),
	}, nil
}
