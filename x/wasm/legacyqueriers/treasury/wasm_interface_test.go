package treasury

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
)

func TestQueryTaxRate(t *testing.T) {
	rate := sdk.ZeroDec()

	querier := NewWasmQuerier()
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(sdk.Context{}, []byte{})
	require.Error(t, err)

	// tax rate query
	bz, err := json.Marshal(CosmosQuery{
		TaxRate: &struct{}{},
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(sdk.Context{}, bz)
	require.NoError(t, err)

	var taxRateResponse TaxRateQueryResponse
	require.NoError(t, json.Unmarshal(res, &taxRateResponse))
	require.Equal(t, rate.String(), taxRateResponse.Rate)
}

func TestQueryTaxCap(t *testing.T) {

	cap := sdk.ZeroInt()

	querier := NewWasmQuerier()
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(sdk.Context{}, []byte{})
	require.Error(t, err)

	// tax rate query
	bz, err := json.Marshal(CosmosQuery{
		TaxCap: &QueryTaxCapParams{
			Denom: core.MicroSDRDenom,
		},
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(sdk.Context{}, bz)
	require.NoError(t, err)

	var taxCapResponse TaxCapQueryResponse
	require.NoError(t, json.Unmarshal(res, &taxCapResponse))
	require.Equal(t, cap.String(), taxCapResponse.Cap)
}
