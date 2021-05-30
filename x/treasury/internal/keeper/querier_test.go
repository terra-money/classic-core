package keeper

import (
	"strings"
	"testing"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/treasury/internal/types"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
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

func getQueriedTaxCaps(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.TaxCapsQueryResponse {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTaxCaps}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryTaxCaps}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response types.TaxCapsQueryResponse
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

func getQueriedIndicators(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.IndicatorQueryResonse {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryIndicators}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{types.QueryIndicators}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var indicators types.IndicatorQueryResonse
	err2 := cdc.UnmarshalJSON(bz, &indicators)
	require.Nil(t, err2)

	return indicators
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

func TestQueryTaxCaps(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)

	input.TreasuryKeeper.SetTaxCap(input.Ctx, "ukrw", sdk.NewInt(1000000000))
	input.TreasuryKeeper.SetTaxCap(input.Ctx, "usdr", sdk.NewInt(1000000))
	input.TreasuryKeeper.SetTaxCap(input.Ctx, "uusd", sdk.NewInt(1200000))

	// Get a currency super random; should default to policy coin.
	queriedTaxCaps := getQueriedTaxCaps(t, input.Ctx, input.Cdc, querier)

	require.Equal(t, queriedTaxCaps,
		types.TaxCapsQueryResponse{
			{
				Denom:  "ukrw",
				TaxCap: sdk.NewInt(1000000000),
			},
			{
				Denom:  "usdr",
				TaxCap: sdk.NewInt(1000000),
			},

			{
				Denom:  "uusd",
				TaxCap: sdk.NewInt(1200000),
			},
		},
	)
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

func TestQueryIndicators(t *testing.T) {
	input := CreateTestInput(t)
	querier := NewQuerier(input.TreasuryKeeper)
	sh := staking.NewHandler(input.StakingKeeper)

	stakingAmt := sdk.TokensFromConsensusPower(1)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	_, err := sh(input.Ctx, NewTestMsgCreateValidator(addr, val, stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, NewTestMsgCreateValidator(addr1, val1, stakingAmt))
	require.NoError(t, err)

	staking.EndBlocker(input.Ctx.WithBlockHeight(core.BlocksPerWeek-1), input.StakingKeeper)

	proceedsAmt := sdk.NewInt(1000000000000)
	taxProceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, proceedsAmt))
	input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, taxProceeds)

	targetIndicators := types.IndicatorQueryResonse{
		TRLYear:  proceedsAmt.ToDec().QuoInt(stakingAmt.MulRaw(2)),
		TRLMonth: proceedsAmt.ToDec().QuoInt(stakingAmt.MulRaw(2)),
	}

	queriedIndicators := getQueriedIndicators(t, input.Ctx, input.Cdc, querier)
	require.Equal(t, targetIndicators, queriedIndicators)

	// Update indicators
	input.TreasuryKeeper.UpdateIndicators(input.Ctx)

	// Record same tax proceeds to get same trl
	input.TreasuryKeeper.RecordEpochTaxProceeds(input.Ctx, taxProceeds)

	// Change context to next epoch
	input.Ctx = input.Ctx.WithBlockHeight(core.BlocksPerWeek)
	queriedIndicators = getQueriedIndicators(t, input.Ctx, input.Cdc, querier)
	require.Equal(t, targetIndicators, queriedIndicators)
}
