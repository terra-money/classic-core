package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/types"
)

func TestNewQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewQuerier(input.MarketKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)

	_, err = querier(input.Ctx, []string{"INVALID_PATH"}, query)
	require.Error(t, err)
}

func TestQueryParams(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	var params types.Params

	res, errRes := queryParameters(input.Ctx, input.MarketKeeper)
	require.NoError(t, errRes)

	err := cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.MarketKeeper.GetParams(input.Ctx), params)
}

func TestQuerySwap(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	price := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, price)
	input.MarketKeeper.UpdatePools(input.Ctx)

	querier := NewQuerier(input.MarketKeeper)
	var err error

	// empty data will occur error
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	res, err := querier(input.Ctx, []string{types.QuerySwap}, query)
	require.Error(t, err)

	// recursive query
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	queryParams := types.NewQuerySwapParams(offerCoin, core.MicroLunaDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QuerySwap}, query)
	require.Error(t, err)

	// valid query
	queryParams = types.NewQuerySwapParams(offerCoin, core.MicroSDRDenom)
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QuerySwap}, query)
	require.NoError(t, err)

	var swapCoin sdk.Coin
	err = cdc.UnmarshalJSON(res, &swapCoin)
	require.NoError(t, err)
	require.Equal(t, core.MicroSDRDenom, swapCoin.Denom)
	require.True(t, sdk.NewInt(17).GTE(swapCoin.Amount))
	require.True(t, swapCoin.Amount.IsPositive())
}

func TestQueryTerraPool(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	pool := sdk.NewDecWithPrec(17, 1)
	input.MarketKeeper.SetTerraPool(input.Ctx, pool)

	querier := NewQuerier(input.MarketKeeper)
	query := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryTerraPool}, query)
	require.NoError(t, errRes)

	var retPool sdk.Dec
	err := cdc.UnmarshalJSON(res, &retPool)
	require.NoError(t, err)
	require.Equal(t, pool, retPool)
}

func TestQueryLastUpdateHeight(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	input.MarketKeeper.SetLastUpdateHeight(input.Ctx, 1)

	querier := NewQuerier(input.MarketKeeper)
	query := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryLastUpdateHeight}, query)
	require.NoError(t, errRes)

	var retPool int64
	err := cdc.UnmarshalJSON(res, &retPool)
	require.NoError(t, err)
	require.Equal(t, int64(1), retPool)
}

func TestQueryBasePool(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	querier := NewQuerier(input.MarketKeeper)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	res, errRes := querier(input.Ctx, []string{types.QueryBasePool}, query)

	require.NoError(t, errRes)
	var retBasePool sdk.Dec
	err := cdc.UnmarshalJSON(res, &retBasePool)
	require.NoError(t, err)

	require.Equal(t, input.MarketKeeper.GetBasePool(input.Ctx), retBasePool)
}
