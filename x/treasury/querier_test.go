package treasury

import (
	"strings"
	"terra/types"
	"terra/types/assets"
	"terra/types/util"
	"testing"

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

	var taxRate sdk.Dec
	err2 := cdc.UnmarshalJSON(bz, &taxRate)
	require.Nil(t, err2)

	return taxRate
}

func getQueriedTaxCap(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryTaxCap}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryTaxCap, denom}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var taxCap sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &taxCap)
	require.Nil(t, err2)

	return taxCap
}

func getQueriedRewardWeight(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryMiningRewardWeight}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryMiningRewardWeight, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var rewardWeight sdk.Dec
	err2 := cdc.UnmarshalJSON(bz, &rewardWeight)
	require.Nil(t, err2)

	return rewardWeight
}

func getQueriedTaxProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Coins {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryTaxProceeds}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryTaxProceeds, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var taxProceeds sdk.Coins
	err2 := cdc.UnmarshalJSON(bz, &taxProceeds)
	require.Nil(t, err2)

	return taxProceeds
}

func getQueriedSeigniorageProceeds(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, epoch sdk.Int) sdk.Coin {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QuerySeigniorageProceeds}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QuerySeigniorageProceeds, epoch.String()}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var seigniorageProceeds sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &seigniorageProceeds)
	require.Nil(t, err2)

	return sdk.NewCoin(assets.LunaDenom, seigniorageProceeds)
}

func getQueriedActiveClaims(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.ClaimPool {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryActiveClaims}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryActiveClaims}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var activeClaims types.ClaimPool
	err2 := cdc.UnmarshalJSON(bz, &activeClaims)
	require.Nil(t, err2)

	return activeClaims
}

func getQueriedIssuance(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryIssuance}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryIssuance, denom}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var issuance sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &issuance)
	require.Nil(t, err2)

	return issuance
}

func getQueriedCurrentEpoch(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) sdk.Int {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryCurrentEpoch}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryCurrentEpoch}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var curEpoch sdk.Int
	err2 := cdc.UnmarshalJSON(bz, &curEpoch)
	require.Nil(t, err2)

	return curEpoch
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
		sdk.NewCoin(assets.SDRDenom, sdk.NewInt(1000)),
	}
	input.treasuryKeeper.RecordTaxProceeds(input.ctx, taxProceeds)

	queriedTaxProceeds := getQueriedTaxProceeds(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, queriedTaxProceeds, taxProceeds)
}

func TestQuerySeigniorageProceeds(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	seigniorageProceeds := sdk.NewCoin(assets.LunaDenom, sdk.NewInt(10))
	input.mintKeeper.AddSeigniorage(input.ctx, seigniorageProceeds.Amount)

	queriedSeigniorageProceeds := getQueriedSeigniorageProceeds(t, input.ctx, input.cdc, querier, util.GetEpoch(input.ctx))

	require.Equal(t, queriedSeigniorageProceeds, seigniorageProceeds)
}

func TestQueryIssuance(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	issuance := sdk.NewInt(1000)
	input.mintKeeper.Mint(input.ctx, addrs[0], sdk.NewCoin(assets.SDRDenom, issuance))

	queriedIssuance := getQueriedIssuance(t, input.ctx, input.cdc, querier, assets.SDRDenom)

	require.Equal(t, queriedIssuance, issuance)
}

func TestQueryActiveClaims(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.treasuryKeeper)

	input.treasuryKeeper.AddClaim(input.ctx, types.NewClaim(
		types.OracleClaimClass, sdk.NewInt(10), addrs[0],
	))
	input.treasuryKeeper.AddClaim(input.ctx, types.NewClaim(
		types.BudgetClaimClass, sdk.NewInt(10), addrs[0],
	))
	input.treasuryKeeper.AddClaim(input.ctx, types.NewClaim(
		types.OracleClaimClass, sdk.NewInt(10), addrs[1],
	))
	input.treasuryKeeper.AddClaim(input.ctx, types.NewClaim(
		types.BudgetClaimClass, sdk.NewInt(10), addrs[1],
	))
	input.treasuryKeeper.AddClaim(input.ctx, types.NewClaim(
		types.OracleClaimClass, sdk.NewInt(10), addrs[2],
	))

	queriedActiveClaims := getQueriedActiveClaims(t, input.ctx, input.cdc, querier)

	require.Equal(t, 5, len(queriedActiveClaims))
}
