package keeper

import (
	"strings"
	"testing"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const custom = "custom"

func getQueriedTaxRate(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch int64) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTaxRate}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryTaxRate}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response sdk.Dec
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response
}

func getQueriedTaxCap(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Int {
	params := types.QueryTaxCapParams{
		Denom: denom,
	}

	bz, err := cdc.MarshalJSON(params)
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTaxCap}, "/"),
		Data: bz,
	}

	bz, err = querier(ctx, []string{types.QueryTaxCap}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response
}

func getQueriedRewardWeight(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch int64) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryRewardWeight}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryRewardWeight}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response sdk.Dec
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response
}

func getQueriedTaxProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch int64) sdk.Coins {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTaxProceeds}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryTaxProceeds}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response sdk.Coins
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response
}

func getQueriedSeigniorageProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch int64) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QuerySeigniorageProceeds}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QuerySeigniorageProceeds}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response
}

func getQueriedParameters(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.Params {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParameters}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{types.QueryParameters}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var params types.Params
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	params := types.DefaultParams()
	input.TreasuryKeeper.SetParams(input.Ctx, params)

	queriedParams := getQueriedParameters(t, input.Ctx, input.Cdc, querier)

	require.Equal(t, queriedParams, params)
}

func TestQueryRewardWeight(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	rewardWeight := sdk.NewDecWithPrec(77, 2)
	input.TreasuryKeeper.SetRewardWeight(input.Ctx, rewardWeight)

	queriedRewardWeight := getQueriedRewardWeight(t, input.Ctx, input.Cdc, querier, input.TreasuryKeeper.GetEpoch(input.Ctx))

	require.Equal(t, queriedRewardWeight, rewardWeight)
}

func TestQueryTaxRate(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	taxRate := sdk.NewDecWithPrec(1, 3)
	input.TreasuryKeeper.SetTaxRate(input.Ctx, taxRate)

	queriedTaxRate := getQueriedTaxRate(t, input.Ctx, input.Cdc, querier, input.TreasuryKeeper.GetEpoch(input.Ctx))

	require.Equal(t, queriedTaxRate, taxRate)
}

func TestQueryTaxCap(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	params := input.TreasuryKeeper.GetParams(input.Ctx)

	// Get a currency super random; should default to policy coin.
	queriedTaxCap := getQueriedTaxCap(t, input.Ctx, input.Cdc, querier, "hello")

	require.Equal(t, queriedTaxCap, params.TaxPolicy.Cap.Amount)
}

func TestQueryTaxProceeds(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	taxProceeds := sdk.Coins{
		sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(1000).MulRaw(core.MicroUnit)),
	}
	input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, taxProceeds)

	queriedTaxProceeds := getQueriedTaxProceeds(t, input.Ctx, input.Cdc, querier, input.TreasuryKeeper.GetEpoch(input.Ctx))

	require.Equal(t, queriedTaxProceeds, taxProceeds)
}

func TestQuerySeigniorageProceeds(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	targetIssuance := sdk.NewInt(1000)
	targetSeigniorage := sdk.NewInt(10)
	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance)))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)
	input.TreasuryKeeper.RecordEpochInitialIssuance(input.Ctx)

	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerWeek)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, targetIssuance.Sub(targetSeigniorage))))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	queriedSeigniorageProceeds := getQueriedSeigniorageProceeds(t, input.Ctx, input.Cdc, querier, input.TreasuryKeeper.GetEpoch(input.Ctx))

	require.Equal(t, targetSeigniorage, queriedSeigniorageProceeds)
}
