package wasm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/keeper"
	"github.com/terra-project/core/x/oracle/wasm"
)

func TestQueryExchangeRates(t *testing.T) {
	input := keeper.CreateTestInput(t)

	KRWExchangeRate := sdk.NewDec(1700)
	USDExchangeRate := sdk.NewDecWithPrec(17, 1)
	SDRExchangeRate := sdk.NewDecWithPrec(19, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, KRWExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroUSDDenom, USDExchangeRate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, SDRExchangeRate)

	querier := wasm.NewWasmQuerier(input.OracleKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// not existing quote denom query
	queryParams := wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroLunaDenom,
		QuoteDenoms: []string{core.MicroMNTDenom},
	}
	bz, err := json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// not existing base denom query
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroCNYDenom,
		QuoteDenoms: []string{core.MicroKRWDenom, core.MicroUSDDenom, core.MicroSDRDenom},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query luna exchange rates
	queryParams = wasm.ExchangeRateQueryParams{
		BaseDenom:   core.MicroKRWDenom,
		QuoteDenoms: []string{core.MicroLunaDenom, core.MicroUSDDenom, core.MicroSDRDenom},
	}
	bz, err = json.Marshal(wasm.CosmosQuery{
		ExchangeRates: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var exchangeRatesResponse wasm.ExchangeRatesQueryResponse
	err = json.Unmarshal(res, &exchangeRatesResponse)
	require.NoError(t, err)
	require.Equal(t, exchangeRatesResponse, wasm.ExchangeRatesQueryResponse{
		BaseDenom: core.MicroKRWDenom,
		ExchangeRates: []wasm.ExchangeRateItem{
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
