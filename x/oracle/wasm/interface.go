package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/terra-money/core/x/oracle/keeper"
	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = Querier{}

// Querier - staking query interface for wasm contract
type Querier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) Querier {
	return Querier{keeper}
}

// Query - implement query function
func (Querier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

// ExchangeRateQueryParams query request params for exchange rates
type ExchangeRateQueryParams struct {
	BaseDenom   string   `json:"base_denom"`
	QuoteDenoms []string `json:"quote_denoms"`
}

// CosmosQuery custom query interface for oracle querier
type CosmosQuery struct {
	ExchangeRates *ExchangeRateQueryParams `json:"exchange_rates,omitempty"`
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

// QueryCustom implements custom query interface
func (querier Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.ExchangeRates != nil {
		// LUNA / BASE_DENOM
		baseDenomExchangeRate, err := querier.keeper.GetLunaExchangeRate(ctx, params.ExchangeRates.BaseDenom)
		if err != nil {
			return nil, err
		}

		var items []ExchangeRateItem
		for _, quoteDenom := range params.ExchangeRates.QuoteDenoms {
			// LUNA / QUOTE_DENOM
			quoteDenomExchangeRate, err := querier.keeper.GetLunaExchangeRate(ctx, quoteDenom)
			if err != nil {
				continue
			}

			// (LUNA / QUOTE_DENOM) / (BASE_DENOM / LUNA) = BASE_DENOM / QUOTE_DENOM
			items = append(items, ExchangeRateItem{
				ExchangeRate: quoteDenomExchangeRate.Quo(baseDenomExchangeRate).String(),
				QuoteDenom:   quoteDenom,
			})
		}

		bz, err := json.Marshal(ExchangeRatesQueryResponse{
			BaseDenom:     params.ExchangeRates.BaseDenom,
			ExchangeRates: items,
		})

		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}

		return bz, nil
	}

	return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Oracle variant"}
}
