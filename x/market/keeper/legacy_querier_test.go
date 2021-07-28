package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"
)

func TestNewLegacyQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewLegacyQuerier(input.MarketKeeper, input.Cdc)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)

	_, err = querier(input.Ctx, []string{"INVALID_PATH"}, query)
	require.Error(t, err)
}

func TestLegacyQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewLegacyQuerier(input.MarketKeeper, input.Cdc)

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryParameters}, req)
	require.NoError(t, err)

	var params types.Params
	err = input.Cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.MarketKeeper.GetParams(input.Ctx), params)
}

func TestLegacyQuerySwap(t *testing.T) {
	input := CreateTestInput(t)

	price := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, price)

	querier := NewLegacyQuerier(input.MarketKeeper, input.Cdc)
	var err error

	// empty data will occur error
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err = querier(input.Ctx, []string{types.QuerySwap}, query)
	require.Error(t, err)

	// recursive query
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	queryParams := types.NewQuerySwapParams(offerCoin, core.MicroLunaDenom)
	bz, err := input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QuerySwap}, query)
	require.Error(t, err)

	// overflow query
	overflowAmt, _ := sdk.NewIntFromString("1000000000000000000000000000000000")
	overflowOfferCoin := sdk.NewCoin(core.MicroLunaDenom, overflowAmt)
	queryParams = types.NewQuerySwapParams(overflowOfferCoin, core.MicroSDRDenom)
	bz, err = input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	_, err = querier(input.Ctx, []string{types.QuerySwap}, query)
	require.Error(t, err)

	// valid query
	queryParams = types.NewQuerySwapParams(offerCoin, core.MicroSDRDenom)
	bz, err = input.Cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QuerySwap}, query)
	require.NoError(t, err)

	var swapCoin sdk.Coin
	err = input.Cdc.UnmarshalJSON(res, &swapCoin)
	require.NoError(t, err)
	require.Equal(t, core.MicroSDRDenom, swapCoin.Denom)
	require.True(t, sdk.NewInt(17).GTE(swapCoin.Amount))
	require.True(t, swapCoin.Amount.IsPositive())
}

func TestLegacyQueryMintPool(t *testing.T) {

	input := CreateTestInput(t)

	poolDelta := sdk.NewDecWithPrec(17, 1)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, poolDelta)

	querier := NewLegacyQuerier(input.MarketKeeper, input.Cdc)
	query := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryTerraPoolDelta}, query)
	require.NoError(t, errRes)

	var retPool sdk.Dec
	err := input.Cdc.UnmarshalJSON(res, &retPool)
	require.NoError(t, err)
	require.Equal(t, poolDelta, retPool)
}
