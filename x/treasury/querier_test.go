package treasury

import (
	"strings"
	"testing"

	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/types/util"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const custom = "custom"

func getQueriedTaxRate(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryTaxRate}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryTaxRate, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var reseponse QueryTaxRateResponse
	err2 := cdc.UnmarshalJSON(bz, &reseponse)
	require.Nil(t, err2)

	return reseponse.TaxRate
}

func getQueriedTaxCap(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryTaxCap}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryTaxCap, denom}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QueryTaxCapResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response.TaxCap
}

func getQueriedRewardWeight(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryMiningRewardWeight}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryMiningRewardWeight, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QueryMiningRewardWeightResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response.RewardWeight
}

func getQueriedTaxProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Coins {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryTaxProceeds}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryTaxProceeds, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QueryTaxProceedsResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response.TaxProceeds
}

func getQueriedSeigniorageProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Coin {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QuerySeigniorageProceeds}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QuerySeigniorageProceeds, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QuerySeigniorageProceedsResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return sdk.NewCoin(assets.MicroLunaDenom, response.SeigniorageProceeds)
}

func getQueriedIssuance(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryIssuance}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryIssuance, denom}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	bz, err = querier(ctx, []string{QueryIssuance, denom, "0"}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QueryIssuanceResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response.Issuance
}

func getQueriedCurrentEpoch(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryCurrentEpoch}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryCurrentEpoch}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var response QueryCurrentEpochResponse
	err2 := cdc.UnmarshalJSON(bz, &response)
	require.Nil(t, err2)

	return response.CurrentEpoch
}

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) Params {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryParams}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryParams}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var params Params
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func TestQueryParams(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	params := DefaultParams()
	input.treasuryKeeper.SetParams(input.ctx, params)

	queriedParams := getQueriedParams(t, input.ctx, input.cdc, querier)

	require.Equal(t, queriedParams, params)
}

func TestQueryRewardWeight(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	rewardWeight := sdk.NewDecWithPrec(77, 2)
	input.treasuryKeeper.SetRewardWeight(input.ctx, rewardWeight)

	queriedRewardWeight := getQueriedRewardWeight(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, queriedRewardWeight, rewardWeight)
}

func TestQueryTaxRate(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	taxRate := sdk.NewDecWithPrec(1, 3)
	input.treasuryKeeper.SetTaxRate(input.ctx, taxRate)

	queriedTaxRate := getQueriedTaxRate(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, queriedTaxRate, taxRate)
}

func TestQueryTaxCap(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	params := input.treasuryKeeper.GetParams(input.ctx)

	// Get a currency super random; should default to policy coin.
	queriedTaxCap := getQueriedTaxCap(t, input.ctx, input.cdc, querier, "hello")

	require.Equal(t, queriedTaxCap, params.TaxPolicy.Cap.Amount)
}

func TestQueryCurrentEpoch(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	curEpoch := util.GetEpoch(input.ctx)

	queriedCurEpoch := getQueriedCurrentEpoch(t, input.ctx, input.cdc, querier)

	require.Equal(t, queriedCurEpoch, curEpoch)
}

func TestQueryTaxProceeds(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	taxProceeds := sdk.Coins{
		sdk.NewCoin(assets.MicroSDRDenom, sdk.NewInt(1000).MulRaw(assets.MicroUnit)),
	}
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, taxProceeds)

	queriedTaxProceeds := getQueriedTaxProceeds(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, queriedTaxProceeds, taxProceeds)
}

func TestQuerySeigniorageProceeds(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	seigniorageProceeds := sdk.NewCoin(assets.MicroLunaDenom, sdk.NewInt(10).MulRaw(assets.MicroUnit))
	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, seigniorageProceeds.Amount))

	getQueriedSeigniorageProceeds(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	input.ctx = input.ctx.WithBlockHeight(util.BlocksPerEpoch)
	input.mintKeeper.Burn(input.ctx, addrs[0], sdk.NewCoin(assets.MicroLunaDenom, seigniorageProceeds.Amount))

	queriedSeigniorageProceeds := getQueriedSeigniorageProceeds(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, seigniorageProceeds, queriedSeigniorageProceeds)
}

func TestQueryIssuance(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	issuance := sdk.NewInt(1000).MulRaw(assets.MicroUnit)
	err := input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.MicroSDRDenom, issuance))
	require.Nil(t, err)

	queriedIssuance := getQueriedIssuance(t, input.ctx, input.cdc, querier, assets.MicroSDRDenom)

	require.Equal(t, issuance, queriedIssuance)
}
