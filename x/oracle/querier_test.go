package oracle

import (
	"strings"
	"terra/types/assets"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const custom = "custom"

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

func getQueriedPrice(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, denom string) sdk.Dec {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryPrice}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryPrice, denom}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var price sdk.Dec
	err2 := cdc.UnmarshalBinaryLengthPrefixed(bz, &price)
	require.Nil(t, err2)
	return price
}

func getQueriedActive(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) []string {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryActive}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryActive}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var actives []string
	err2 := cdc.UnmarshalJSON(bz, &actives)
	require.Nil(t, err2)
	return actives
}

func getQueriedVotes(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, voter sdk.AccAddress, denom string) PriceBallot {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryVotes}, "/"),
		Data: cdc.MustMarshalJSON(NewQueryVoteParams(voter, denom)),
	}

	bz, err := querier(ctx, []string{QueryVotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votes PriceBallot
	err2 := cdc.UnmarshalJSON(bz, &votes)
	require.Nil(t, err2)
	return votes
}

func TestQueryParams(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	defaultParams := DefaultParams()
	input.oracleKeeper.SetParams(input.ctx, defaultParams)

	params := getQueriedParams(t, input.ctx, input.cdc, querier)

	require.Equal(t, defaultParams, params)
}

func TestQueryPrice(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	testPrice := sdk.NewDecWithPrec(48842, 4).MulInt64(assets.MicroUnit)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, testPrice)

	price := getQueriedPrice(t, input.ctx, input.cdc, querier, assets.MicroKRWDenom)

	require.Equal(t, testPrice, price)
}

func TestQueryActives(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	empty := getQueriedActive(t, input.ctx, input.cdc, querier)
	require.Equal(t, 0, len(empty))

	testPrice := sdk.NewDecWithPrec(48842, 4).MulInt64(assets.MicroUnit)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, testPrice)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroUSDDenom, testPrice)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroSDRDenom, testPrice)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroGBPDenom, testPrice)

	actives := getQueriedActive(t, input.ctx, input.cdc, querier)

	require.Equal(t, 4, len(actives))
}

func TestQueryVotes(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	testPrice := sdk.NewDecWithPrec(48842, 4)

	votes := []PriceVote{
		// first voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[0]),
		NewPriceVote(testPrice, assets.MicroKRWDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[0]),
		NewPriceVote(testPrice, assets.MicroUSDDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[0]),

		// Second voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[1]),
		NewPriceVote(testPrice, assets.MicroKRWDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[1]),
		NewPriceVote(testPrice, assets.MicroGBPDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[1]),

		// Third voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[2]),
		NewPriceVote(testPrice, assets.MicroCNYDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[2]),
		NewPriceVote(testPrice, assets.MicroGBPDenom, sdk.OneInt().MulRaw(assets.MicroUnit), addrs[2]),
	}

	for _, vote := range votes {
		input.oracleKeeper.addVote(input.ctx, vote)
	}

	voterOneSDR := getQueriedVotes(t, input.ctx, input.cdc, querier, addrs[0], assets.MicroSDRDenom)
	require.Equal(t, 1, len(voterOneSDR))

	voterOne := getQueriedVotes(t, input.ctx, input.cdc, querier, addrs[0], "")
	require.Equal(t, 3, len(voterOne))

	assetKRW := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.AccAddress{}, assets.MicroKRWDenom)
	require.Equal(t, 2, len(assetKRW))

	noFilters := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.AccAddress{}, "")
	require.Equal(t, 9, len(noFilters))
}
