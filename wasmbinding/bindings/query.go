package bindings

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	markettypes "github.com/classic-terra/core/x/market/types"
	treasurytypes "github.com/classic-terra/core/x/treasury/types"
)

// ExchangeRateQueryParams query request params for exchange rates
type ExchangeRateQueryParams struct {
	BaseDenom   string   `json:"base_denom"`
	QuoteDenoms []string `json:"quote_denoms"`
}

// TerraQuery contains terra custom queries.
type TerraQuery struct {
	Swap          *markettypes.QuerySwapParams     `json:"swap,omitempty"`
	ExchangeRates *ExchangeRateQueryParams         `json:"exchange_rates,omitempty"`
	TaxRate       *struct{}                        `json:"tax_rate,omitempty"`
	TaxCap        *treasurytypes.QueryTaxCapParams `json:"tax_cap,omitempty"`
}

// SwapQueryResponse - swap simulation query response for wasm module
type SwapQueryResponse struct {
	Receive wasmvmtypes.Coin `json:"receive"`
}

// ExchangeRatesQueryResponseItem - exchange rates query response item
type ExchangeRateItem struct {
	ExchangeRate string `json:"exchange_rate"`
	QuoteDenom   string `json:"quote_denom"`
}

// ExchangeRatesQueryResponse - exchange rates query response for wasm module
type ExchangeRatesQueryResponse struct {
	ExchangeRates []ExchangeRateItem `json:"exchange_rates"`
	BaseDenom     string             `json:"base_denom"`
}

// TaxRateQueryResponse - tax rate query response for wasm module
type TaxRateQueryResponse struct {
	// decimal string, eg "0.02"
	Rate string `json:"rate"`
}

// TaxCapQueryResponse - tax cap query response for wasm module
type TaxCapQueryResponse struct {
	// uint64 string, eg "1000000"
	Cap string `json:"cap"`
}
