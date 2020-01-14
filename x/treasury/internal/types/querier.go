package types

// query endpoints supported by the auth Querier
const (
	QueryTaxRate             = "taxRate"
	QueryTaxCap              = "taxCap"
	QueryRewardWeight        = "rewardWeight"
	QuerySeigniorageProceeds = "seigniorageProceeds"
	QueryTaxProceeds         = "taxProceeds"
	QueryParameters          = "parameters"
)

// QueryTaxCapParams for query
// - 'custom/treasury/taxRate
type QueryTaxCapParams struct {
	Denom string
}

// NewQueryTaxCapParams returns new QueryTaxCapParams instance
func NewQueryTaxCapParams(denom string) QueryTaxCapParams {
	return QueryTaxCapParams{
		Denom: denom,
	}
}
