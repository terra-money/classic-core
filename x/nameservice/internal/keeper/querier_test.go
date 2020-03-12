package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/terra-project/core/x/nameservice/internal/types"
	"testing"
	"time"
)

func TestNewQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewQuerier(input.NameserviceKeeper)

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

	res, errRes := queryParameters(input.Ctx, input.NameserviceKeeper)
	require.NoError(t, errRes)

	err := cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.NameserviceKeeper.GetParams(input.Ctx), params)
}

func TestQueryAuctions(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	querier := NewQuerier(input.NameserviceKeeper)

	validName := types.Name("wallet.terra")
	validNameHash, _ := validName.NameHash()
	validName2 := types.Name("chai.terra")
	validNameHash2, _ := validName2.NameHash()
	endTime := time.Now().UTC()

	auction := types.NewAuction(validName, types.AuctionStatusBid, endTime)
	auction2 := types.NewAuction(validName2, types.AuctionStatusReveal, endTime)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash, auction)
	input.NameserviceKeeper.SetAuction(input.Ctx, validNameHash2, auction2)

	var auctions []types.Auction

	// empty data will give all auctions
	bz, err := cdc.MarshalJSON(types.QueryAuctionsParams{nil, types.AuctionStatusNil})
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryAuctions}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &auctions)
	require.NoError(t, err)
	require.Equal(t, 2, len(auctions))

	// hash query will give 1 result
	bz, err = cdc.MarshalJSON(types.QueryAuctionsParams{validNameHash, types.AuctionStatusNil})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryAuctions}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &auctions)
	require.NoError(t, err)
	require.Equal(t, 1, len(auctions))
	require.Equal(t, auction.Name, auctions[0].Name)

	// status query will give 1 result
	bz, err = cdc.MarshalJSON(types.QueryAuctionsParams{nil, types.AuctionStatusBid})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryAuctions}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &auctions)
	require.NoError(t, err)
	require.Equal(t, 1, len(auctions))
	require.Equal(t, auction.Name, auctions[0].Name)

	// status query will give 1 result
	bz, err = cdc.MarshalJSON(types.QueryAuctionsParams{nil, types.AuctionStatusReveal})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryAuctions}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &auctions)
	require.NoError(t, err)
	require.Equal(t, 1, len(auctions))
	require.Equal(t, auction2.Name, auctions[0].Name)
}

func TestQueryBids(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	querier := NewQuerier(input.NameserviceKeeper)

	validName := types.Name("wallet.terra")
	validNameHash, _ := validName.NameHash()

	amount := sdk.NewInt64Coin("foo", 123)
	bidHash := types.GetBidHash("salt", validName, amount, Addrs[0])
	bid := types.NewBid(bidHash, amount, Addrs[0])
	bid2 := types.NewBid(bidHash, amount, Addrs[1])

	var bids []types.Bid

	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, bid)
	input.NameserviceKeeper.SetBid(input.Ctx, validNameHash, bid2)

	// empty data will occur error
	bz, err := cdc.MarshalJSON(types.QueryBidsParams{nil, nil})
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	_, errRes := querier(input.Ctx, []string{types.QueryBids}, query)
	require.Error(t, errRes)

	// name hash query will give 2 bids
	bz, err = cdc.MarshalJSON(types.QueryBidsParams{validNameHash, nil})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryBids}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &bids)
	require.NoError(t, err)
	require.Equal(t, 2, len(bids))

	// name hash, bidder query will give 1 bid
	bz, err = cdc.MarshalJSON(types.QueryBidsParams{validNameHash, Addrs[0]})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryBids}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &bids)
	require.NoError(t, err)
	require.Equal(t, 1, len(bids))
}

func TestQueryResolve(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	querier := NewQuerier(input.NameserviceKeeper)

	name := "wallet.terra"
	name2 := "dokwon.wallet.terra"
	nameHash, childNameHash := types.Name(name).NameHash()
	_, childNameHash2 := types.Name(name2).NameHash()

	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childNameHash, Addrs[0])
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childNameHash2, Addrs[1])

	var addr sdk.AccAddress

	// empty data will occur error
	bz, err := cdc.MarshalJSON(types.QueryResolveParams{nil, nil})
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	_, errRes := querier(input.Ctx, []string{types.QueryResolve}, query)
	require.Error(t, errRes)

	// query for addr 0
	bz, err = cdc.MarshalJSON(types.QueryResolveParams{nameHash, childNameHash})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryResolve}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &addr)
	require.NoError(t, err)
	require.Equal(t, Addrs[0], addr)

	// query for addr 1
	bz, err = cdc.MarshalJSON(types.QueryResolveParams{nameHash, childNameHash2})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryResolve}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &addr)
	require.NoError(t, err)
	require.Equal(t, Addrs[1], addr)
}

func TestQueryReverse(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	querier := NewQuerier(input.NameserviceKeeper)

	name := types.Name("wallet.terra")
	nameHash, _ := name.NameHash()

	registry := types.NewRegistry(name, Addrs[0], time.Now().UTC())
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, Addrs[0], nameHash)
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, Addrs[1], nameHash)

	var resRegistry types.Registry

	// empty data will occur error
	bz, err := cdc.MarshalJSON(types.QueryReverseParams{nil})
	require.NoError(t, err)

	query := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	_, errRes := querier(input.Ctx, []string{types.QueryReverse}, query)
	require.Error(t, errRes)

	// query for addr 0
	bz, err = cdc.MarshalJSON(types.QueryReverseParams{Addrs[0]})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes := querier(input.Ctx, []string{types.QueryReverse}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &resRegistry)
	require.NoError(t, err)
	require.Equal(t, registry, resRegistry)

	// query for addr 1
	bz, err = cdc.MarshalJSON(types.QueryReverseParams{Addrs[1]})
	require.NoError(t, err)

	query = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, errRes = querier(input.Ctx, []string{types.QueryReverse}, query)
	require.NoError(t, errRes)

	err = cdc.UnmarshalJSON(res, &resRegistry)
	require.NoError(t, err)
	require.Equal(t, registry, resRegistry)
}
