package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the auth Querier
const (
	QueryTaxRate             = "taxRate"
	QueryTaxCap              = "taxCap"
	QueryTaxCaps             = "taxCaps"
	QueryRewardWeight        = "rewardWeight"
	QuerySeigniorageProceeds = "seigniorageProceeds"
	QueryTaxProceeds         = "taxProceeds"
	QueryParameters          = "parameters"
	QueryIndicators          = "indicators"
)

// QueryTaxCapParams for query
// - 'custom/treasury/taxRate
type QueryTaxCapParams struct {
	Denom string `json:"denom"`
}

// NewQueryTaxCapParams returns new QueryTaxCapParams instance
func NewQueryTaxCapParams(denom string) QueryTaxCapParams {
	return QueryTaxCapParams{
		Denom: denom,
	}
}

// TaxCapsResponseItem query response item of tax caps querier
type TaxCapsResponseItem struct {
	Denom  string  `json:"denom"`
	TaxCap sdk.Int `json:"tax_cap"`
}

// TaxCapsQueryResponse query response body of tax caps querier
type TaxCapsQueryResponse []TaxCapsResponseItem

// IndicatorQueryResonse query response body
type IndicatorQueryResonse struct {
	TRLYear  sdk.Dec `json:"trl_year"`
	TRLMonth sdk.Dec `json:"trl_month"`
}

// String implements fmt.Stringer interface
func (res IndicatorQueryResonse) String() string {
	return fmt.Sprintf(`Treasury Params:
  TRL Year      : %s 
  TRL Month     : %s

  `, res.TRLYear, res.TRLMonth)
}
