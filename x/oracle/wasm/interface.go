package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/terra-money/core/x/oracle/internal/keeper"
	wasm "github.com/terra-money/core/x/wasm/exported"
)

var _ wasm.WasmQuerierInterface = WasmQuerier{}

// WasmQuerier - staking query interface for wasm contract
type WasmQuerier struct {
	keeper keeper.Keeper
}

// NewWasmQuerier return bank wasm query interface
func NewWasmQuerier(keeper keeper.Keeper) WasmQuerier {
	return WasmQuerier{keeper}
}

// Query - implement query function
func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

// ExchangeRateQueryParams query request params for exchange rates
type ExchangeRateQueryParams struct {
	BaseDenom   string   `json:"base_denom"`
	QuoteDenoms []string `json:"quote_denoms"`
}

// CosmosQuery custom query interface for oracle querier
type CosmosQuery struct {
	ExchangeRates ExchangeRateQueryParams `json:"exchange_rates"`
}

// ExchangeRatesQueryResponseItem - exchange rates query response item
type exchangeRateItem struct {
	ExchangeRate string `json:"exchange_rate"`
	QuoteDenom   string `json:"quote_denom"`
}

// ExchangeRatesQueryResponse - exchange rates query response for wasm module
type ExchangeRatesQueryResponse struct {
	ExchangeRates []exchangeRateItem `json:"exchange_rates"`
	BaseDenom     string             `json:"base_denom"`
}

// QueryCustom implements custom query interface
func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var params CosmosQuery
	err := json.Unmarshal(data, &params)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// LUNA / BASE_DENOM
	baseDenomExchangeRate, err := querier.keeper.GetLunaExchangeRate(ctx, params.ExchangeRates.BaseDenom)
	if err != nil {
		return nil, err
	}

	var items []exchangeRateItem
	for _, quoteDenom := range params.ExchangeRates.QuoteDenoms {
		quoteDenomExchangeRate, err := querier.keeper.GetLunaExchangeRate(ctx, quoteDenom)
		if err != nil {
			return nil, err
		}

		// (BASE_DENOM / LUNA) / (DENOM / LUNA) = BASE_DENOM / QUOTE_DENOM
		items = append(items, exchangeRateItem{
			ExchangeRate: baseDenomExchangeRate.Quo(quoteDenomExchangeRate).String(),
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
