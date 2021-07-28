package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/types"
)

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)

	querier := NewQuerier(input.MarketKeeper)
	res, err := querier.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)

	require.Equal(t, input.MarketKeeper.GetParams(input.Ctx), res.Params)
}

func TestQuerySwap(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)
	querier := NewQuerier(input.MarketKeeper)

	price := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, price)

	var err error

	// empty request cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{})
	require.Error(t, err)

	// empty ask denom cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: sdk.Coin{Denom: core.MicroSDRDenom, Amount: sdk.NewInt(100)}.String()})
	require.Error(t, err)

	// empty offer coin cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{AskDenom: core.MicroSDRDenom})
	require.Error(t, err)

	// recursive query
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10)).String()
	res, err := querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: offerCoin, AskDenom: core.MicroLunaDenom})
	require.Error(t, err)

	// overflow query
	overflowAmt, _ := sdk.NewIntFromString("1000000000000000000000000000000000")
	overflowOfferCoin := sdk.NewCoin(core.MicroLunaDenom, overflowAmt).String()
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: overflowOfferCoin, AskDenom: core.MicroSDRDenom})
	require.Error(t, err)

	// valid query
	res, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: offerCoin, AskDenom: core.MicroSDRDenom})
	require.NoError(t, err)

	require.Equal(t, core.MicroSDRDenom, res.ReturnCoin.Denom)
	require.True(t, sdk.NewInt(17).GTE(res.ReturnCoin.Amount))
	require.True(t, res.ReturnCoin.Amount.IsPositive())
}

func TestQueryMintPoolDelta(t *testing.T) {

	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)
	querier := NewQuerier(input.MarketKeeper)

	poolDelta := sdk.NewDecWithPrec(17, 1)
	input.MarketKeeper.SetTerraPoolDelta(input.Ctx, poolDelta)

	res, errRes := querier.TerraPoolDelta(ctx, &types.QueryTerraPoolDeltaRequest{})
	require.NoError(t, errRes)

	require.Equal(t, poolDelta, res.TerraPoolDelta)
}
