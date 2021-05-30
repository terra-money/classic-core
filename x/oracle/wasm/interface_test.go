package wasm

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/oracle/internal/keeper"
)

func TestQueryExchangeRates(t *testing.T) {
	input := keeper.CreateTestInput(t)

	KRWExchangeRate := sdk.NewDec(1700)
	USDExchangeRate := sdk.NewDecWithPrec(17, 1)
	SDRExchangeRate := sdk.NewDecWithPrec(19, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, KRWExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroUSDDenom, USDExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, SDRExchangeRate)

	querier := NewWasmQuerier(input.OracleKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := ExchangeRateQueryParams{
		BaseDenom:   core.MicroLunaDenom,
		QuoteDenoms: []string{core.MicroMNTDenom},
	}
	bz, err := json.Marshal(CosmosQuery{
		ExchangeRates: queryParams,
	})
	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// not existing base denom query
	queryParams = ExchangeRateQueryParams{
		BaseDenom:   core.MicroCNYDenom,
		QuoteDenoms: []string{core.MicroKRWDenom, core.MicroUSDDenom, core.MicroSDRDenom},
	}
	bz, err = json.Marshal(CosmosQuery{
		ExchangeRates: queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query luna exchange rates
	queryParams = ExchangeRateQueryParams{
		BaseDenom:   core.MicroKRWDenom,
		QuoteDenoms: []string{core.MicroLunaDenom, core.MicroUSDDenom, core.MicroSDRDenom},
	}
	bz, err = json.Marshal(CosmosQuery{
		ExchangeRates: queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var exchangeRatesResponse ExchangeRatesQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, ExchangeRatesQueryResponse{
		BaseDenom: core.MicroKRWDenom,
		ExchangeRates: []exchangeRateItem{
			{
				ExchangeRate: KRWExchangeRate.String(),
				QuoteDenom:   core.MicroLunaDenom,
			},
			{
				ExchangeRate: KRWExchangeRate.Quo(USDExchangeRate).String(),
				QuoteDenom:   core.MicroUSDDenom,
			},
			{
				ExchangeRate: KRWExchangeRate.Quo(SDRExchangeRate).String(),
				QuoteDenom:   core.MicroSDRDenom,
			},
		},
	})
}
