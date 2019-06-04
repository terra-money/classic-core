package oracle

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/terra-project/core/types/assets"

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
	err2 := cdc.UnmarshalJSON(bz, &price)
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

func getQueriedVotes(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, voter sdk.ValAddress, denom string) PriceVotes {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryVotes}, "/"),
		Data: cdc.MustMarshalJSON(NewQueryVoteParams(voter, denom)),
	}

	bz, err := querier(ctx, []string{QueryVotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votes PriceVotes
	err2 := cdc.UnmarshalJSON(bz, &votes)
	require.Nil(t, err2)
	return votes
}

func getQueriedPrevotes(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, voter sdk.ValAddress, denom string) PricePrevotes {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryPrevotes}, "/"),
		Data: cdc.MustMarshalJSON(NewQueryPrevoteParams(voter, denom)),
	}

	bz, err := querier(ctx, []string{QueryPrevotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var prevotes PricePrevotes
	err2 := cdc.UnmarshalJSON(bz, &prevotes)
	require.Nil(t, err2)
	return prevotes
}

func getQueriedFeederDelegation(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, validator sdk.ValAddress) sdk.AccAddress {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QuerierRoute, QueryFeederDelegation}, "/"),
		Data: cdc.MustMarshalJSON(NewQueryFeederDelegationParams(validator)),
	}

	bz, err := querier(ctx, []string{QueryFeederDelegation}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var delegate sdk.AccAddress
	err2 := cdc.UnmarshalJSON(bz, &delegate)
	require.Nil(t, err2)
	return delegate
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

	testPrice := sdk.NewDecWithPrec(48842, 4)
	input.oracleKeeper.SetLunaSwapRate(input.ctx, assets.MicroKRWDenom, testPrice)

	price := getQueriedPrice(t, input.ctx, input.cdc, querier, assets.MicroKRWDenom)

	require.Equal(t, testPrice, price)
}

func TestQueryActives(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	empty := getQueriedActive(t, input.ctx, input.cdc, querier)
	require.Equal(t, 0, len(empty))

	testPrice := sdk.NewDecWithPrec(48842, 4)
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

	votes := PriceVotes{
		// first voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0])),
		NewPriceVote(testPrice, assets.MicroKRWDenom, sdk.ValAddress(addrs[0])),
		NewPriceVote(testPrice, assets.MicroUSDDenom, sdk.ValAddress(addrs[0])),

		// Second voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[1])),
		NewPriceVote(testPrice, assets.MicroKRWDenom, sdk.ValAddress(addrs[1])),
		NewPriceVote(testPrice, assets.MicroGBPDenom, sdk.ValAddress(addrs[1])),

		// Third voter votes
		NewPriceVote(testPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[2])),
		NewPriceVote(testPrice, assets.MicroCNYDenom, sdk.ValAddress(addrs[2])),
		NewPriceVote(testPrice, assets.MicroGBPDenom, sdk.ValAddress(addrs[2])),
	}

	for _, vote := range votes {
		input.oracleKeeper.addVote(input.ctx, vote)
	}

	voterOneSDR := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.ValAddress(addrs[0]), assets.MicroSDRDenom)
	require.Equal(t, 1, len(voterOneSDR))

	voterOne := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.ValAddress(addrs[0]), "")
	require.Equal(t, 3, len(voterOne))

	assetKRW := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.ValAddress{}, assets.MicroKRWDenom)
	require.Equal(t, 2, len(assetKRW))

	noFilters := getQueriedVotes(t, input.ctx, input.cdc, querier, sdk.ValAddress{}, "")
	require.Equal(t, 9, len(noFilters))
}

func TestQueryPrevotes(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	testPrice := sdk.NewDecWithPrec(48842, 4)

	hash, _ := VoteHash("abcd", testPrice, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	hexHash := hex.EncodeToString(hash)
	prevotes := PricePrevotes{
		// first voter votes
		NewPricePrevote(hexHash, assets.MicroSDRDenom, sdk.ValAddress(addrs[0]), 1),
		NewPricePrevote(hexHash, assets.MicroKRWDenom, sdk.ValAddress(addrs[0]), 1),
		NewPricePrevote(hexHash, assets.MicroUSDDenom, sdk.ValAddress(addrs[0]), 1),

		// Second voter votes
		NewPricePrevote(hexHash, assets.MicroSDRDenom, sdk.ValAddress(addrs[1]), 1),
		NewPricePrevote(hexHash, assets.MicroKRWDenom, sdk.ValAddress(addrs[1]), 1),
		NewPricePrevote(hexHash, assets.MicroGBPDenom, sdk.ValAddress(addrs[1]), 1),

		// Third voter votes
		NewPricePrevote(hexHash, assets.MicroSDRDenom, sdk.ValAddress(addrs[2]), 1),
		NewPricePrevote(hexHash, assets.MicroCNYDenom, sdk.ValAddress(addrs[2]), 1),
		NewPricePrevote(hexHash, assets.MicroGBPDenom, sdk.ValAddress(addrs[2]), 1),
	}

	for _, prevote := range prevotes {
		input.oracleKeeper.addPrevote(input.ctx, prevote)
	}

	voterOneSDR := getQueriedPrevotes(t, input.ctx, input.cdc, querier, sdk.ValAddress(addrs[0]), assets.MicroSDRDenom)
	require.Equal(t, 1, len(voterOneSDR))

	voterOne := getQueriedPrevotes(t, input.ctx, input.cdc, querier, sdk.ValAddress(addrs[0]), "")
	require.Equal(t, 3, len(voterOne))

	assetKRW := getQueriedPrevotes(t, input.ctx, input.cdc, querier, sdk.ValAddress{}, assets.MicroKRWDenom)
	require.Equal(t, 2, len(assetKRW))

	noFilters := getQueriedPrevotes(t, input.ctx, input.cdc, querier, sdk.ValAddress{}, "")
	require.Equal(t, 9, len(noFilters))
}

func TestQueryFeederDelegations(t *testing.T) {
	input := createTestInput(t)
	querier := NewQuerier(input.oracleKeeper)

	input.oracleKeeper.SetFeedDelegate(input.ctx, sdk.ValAddress(addrs[0]), addrs[1])

	delegate := getQueriedFeederDelegation(t, input.ctx, input.cdc, querier, sdk.ValAddress(addrs[0]))

	require.Equal(t, sdk.AccAddress(sdk.ValAddress(addrs[1])), delegate)
	require.Equal(t, sdk.AccAddress(sdk.ValAddress(addrs[2])), addrs[2])
	require.NotEqual(t, sdk.AccAddress(sdk.ValAddress(addrs[2])), addrs[1])
}
