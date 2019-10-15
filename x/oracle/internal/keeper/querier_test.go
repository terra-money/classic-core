package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

func TestNewQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewQuerier(input.OracleKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)
}

func TestQueryParams(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	var params types.Params

	res, errRes := queryParameters(input.Ctx, input.OracleKeeper)
	require.NoError(t, errRes)

	err := cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.OracleKeeper.GetParams(input.Ctx), params)
}

func TestQueryPrevotes(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	prevote1 := types.NewPricePrevote("", core.MicroSDRDenom, ValAddrs[0], 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote1)
	prevote2 := types.NewPricePrevote("", core.MicroSDRDenom, ValAddrs[1], 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote2)
	prevote3 := types.NewPricePrevote("", core.MicroLunaDenom, ValAddrs[2], 0)
	input.OracleKeeper.AddPrevote(input.Ctx, prevote3)

	// voter denom both query params
	queryParams := types.NewQueryPrevotesParams(ValAddrs[0], core.MicroSDRDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryPrevotes}, req)
	require.NoError(t, err)

	var filteredPrevotes types.PricePrevotes
	err = cdc.UnmarshalJSON(res, &filteredPrevotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredPrevotes))
	require.Equal(t, prevote1, filteredPrevotes[0])

	// voter query params
	queryParams = types.NewQueryPrevotesParams(ValAddrs[0], "")
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryPrevotes}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &filteredPrevotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredPrevotes))
	require.Equal(t, prevote1, filteredPrevotes[0])

	// denom query params
	queryParams = types.NewQueryPrevotesParams(sdk.ValAddress{}, core.MicroLunaDenom)
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryPrevotes}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &filteredPrevotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredPrevotes))
	require.Equal(t, prevote3, filteredPrevotes[0])
}

func TestQueryVotes(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	vote1 := types.NewPriceVote(sdk.NewDec(1700), core.MicroSDRDenom, ValAddrs[0])
	input.OracleKeeper.AddVote(input.Ctx, vote1)
	vote2 := types.NewPriceVote(sdk.NewDec(1700), core.MicroSDRDenom, ValAddrs[1])
	input.OracleKeeper.AddVote(input.Ctx, vote2)
	vote3 := types.NewPriceVote(sdk.NewDec(1700), core.MicroLunaDenom, ValAddrs[2])
	input.OracleKeeper.AddVote(input.Ctx, vote3)

	// voter denom both query params
	queryParams := types.NewQueryVotesParams(ValAddrs[0], core.MicroSDRDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryVotes}, req)
	require.NoError(t, err)

	var filteredVotes types.PriceVotes
	err = cdc.UnmarshalJSON(res, &filteredVotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredVotes))
	require.Equal(t, vote1, filteredVotes[0])

	// voter query params
	queryParams = types.NewQueryVotesParams(ValAddrs[0], "")
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryVotes}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &filteredVotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredVotes))
	require.Equal(t, vote1, filteredVotes[0])

	// denom query params
	queryParams = types.NewQueryVotesParams(sdk.ValAddress{}, core.MicroLunaDenom)
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryVotes}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &filteredVotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(filteredVotes))
	require.Equal(t, vote3, filteredVotes[0])
}

func TestQueryPrice(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	price := sdk.NewDec(1700)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, price)

	// denom query params
	queryParams := types.NewQueryPriceParams(core.MicroSDRDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryPrice}, req)
	require.NoError(t, err)

	var rprice sdk.Dec
	err = cdc.UnmarshalJSON(res, &rprice)
	require.NoError(t, err)
	require.Equal(t, price, rprice)
}

func TestQueryActives(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	price := sdk.NewDec(1700)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroSDRDenom, price)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroKRWDenom, price)
	input.OracleKeeper.SetLunaPrice(input.Ctx, core.MicroUSDDenom, price)

	res, err := querier(input.Ctx, []string{types.QueryActives}, abci.RequestQuery{})
	require.NoError(t, err)

	targetDenoms := types.DenomList{
		core.MicroKRWDenom,
		core.MicroSDRDenom,
		core.MicroUSDDenom,
	}
	var denoms types.DenomList
	err2 := cdc.UnmarshalJSON(res, &denoms)
	require.NoError(t, err2)
	require.Equal(t, targetDenoms, denoms)
}

func TestQueryFeederDelegation(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	input.OracleKeeper.SetFeedDelegate(input.Ctx, ValAddrs[0], Addrs[1])

	queryParams := types.NewQueryFeederDelegationParams(ValAddrs[0])
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryFeederDelegation}, req)
	require.NoError(t, err)

	var delegatee sdk.AccAddress
	cdc.UnmarshalJSON(res, &delegatee)
	require.Equal(t, Addrs[1], delegatee)
}
