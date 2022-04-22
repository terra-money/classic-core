package stargatelayer_test

import (
	"encoding/hex"
	"testing"

	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/terra-money/core/x/wasm/stargatelayer"
	stargatebank "github.com/terra-money/core/x/wasm/stargatelayer/bank"
)

/**
 * Origin Response
 * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" >
 * "0a090a036261721202333012050a03666f6f"
 *
 * New Version Response
 * The binary built from the proto response with additional field address
 * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" > address:"cosmos1j6j5tsquq2jlw2af7l3xekyaq7zg4l8jsufu78"
 * "0a090a036261721202333012050a03666f6f1a2d636f736d6f73316a366a357473717571326a6c77326166376c3378656b796171377a67346c386a737566753738"

// Updated proto
message QueryAllBalancesResponse {
	// balances is the balances of all the coins.
	repeated cosmos.base.v1beta1.Coin balances = 1
	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];

	// pagination defines the pagination in the response.
	cosmos.base.query.v1beta1.PageResponse pagination = 2;

	// address is the address to query all balances for.
	string address = 3;
}

// Origin proto
message QueryAllBalancesResponse {
	// balances is the balances of all the coins.
	repeated cosmos.base.v1beta1.Coin balances = 1
	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];

	// pagination defines the pagination in the response.
	cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
*/

func TestDeterministic_AllBalances(t *testing.T) {
	originVersionBz, err := hex.DecodeString("0a090a036261721202333012050a03666f6f")
	require.NoError(t, err)

	newVersionBz, err := hex.DecodeString("0a090a036261721202333012050a03666f6f1a2d636f736d6f73316a366a357473717571326a6c77326166376c3378656b796171377a67346c386a737566753738")
	require.NoError(t, err)

	binding, ok := stargatelayer.StargateLayerBindings.Load("/cosmos.bank.v1beta1.Query/AllBalances")
	require.True(t, ok)

	// new version response should be changed into origin version response
	normalizedBz, err := stargatelayer.NormalizeResponse(binding, newVersionBz)
	require.NoError(t, err)

	require.Equal(t, originVersionBz, normalizedBz)
	require.NotEqual(t, newVersionBz, normalizedBz)

	// raw build also make same result
	expectedResponse := stargatebank.QueryAllBalancesResponse{
		Balances: sdk.NewCoins(sdk.NewCoin("bar", sdk.NewInt(30))),
		Pagination: &query.PageResponse{
			NextKey: []byte("foo"),
		},
	}
	expectedResponseBz, err := proto.Marshal(&expectedResponse)
	require.NoError(t, err)
	require.Equal(t, expectedResponseBz, normalizedBz)

	// should be clear
	data := binding.(*stargatebank.QueryAllBalancesResponse)
	require.Empty(t, data.Balances)
	require.Empty(t, data.Pagination)
}
