package types

// query endpoints supported by the auth Querier
const (
	QueryCurrentEpoch        = "currentEpoch"
	QueryTaxRate             = "taxRate"
	QueryTaxCap              = "taxCap"
	QueryRewardWeight        = "rewardWeight"
	QuerySeigniorageProceeds = "seigniorageProceeds"
	QueryTaxProceeds         = "taxProceeds"
	QueryParameters          = "parameters"
	QueryHistoricalIssuance  = "historicalIssuance"
)

// QueryTaxCapParams for query
// - 'custom/treasury/taxRate
type QueryTaxCapParams struct {
	Denom string
}

func NewQueryTaxCapParams(denom string) QueryTaxCapParams {
	return QueryTaxCapParams{
		Denom: denom,
	}
}

// QueryTaxRateParams for query
// - 'custom/treasury/taxRate
type QueryTaxRateParams struct {
	Epoch int64
}

func NewQueryTaxRateParams(epoch int64) QueryTaxRateParams {
	return QueryTaxRateParams{
		Epoch: epoch,
	}
}

// QueryRewardWeightParams for query
// - 'custom/treasury/rewardWeight
type QueryRewardWeightParams struct {
	Epoch int64
}

func NewQueryRewardWeightParams(epoch int64) QueryRewardWeightParams {
	return QueryRewardWeightParams{
		Epoch: epoch,
	}
}

// QuerySeigniorageProceedsParams for query
// - 'custom/treasury/seigniorageProceeds
type QuerySeigniorageProceedsParams struct {
	Epoch int64
}

func NewQuerySeigniorageParams(epoch int64) QuerySeigniorageProceedsParams {
	return QuerySeigniorageProceedsParams{
		Epoch: epoch,
	}
}

// QueryTaxProceedsParams for query
// - 'custom/treasury/taxProceeds
type QueryTaxProceedsParams struct {
	Epoch int64
}

func NewQueryTaxProceedsParams(epoch int64) QueryTaxProceedsParams {
	return QueryTaxProceedsParams{
		Epoch: epoch,
	}
}

// QueryHistoricalIssuanceParams for query
// - 'custom/treasury/microLunaIssuance
type QueryHistoricalIssuanceParams struct {
	Epoch int64
}

func NewQueryHistoricalIssuanceParams(epoch int64) QueryHistoricalIssuanceParams {
	return QueryHistoricalIssuanceParams{
		Epoch: epoch,
	}
}
