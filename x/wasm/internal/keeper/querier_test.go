package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/internal/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryContractState(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	cdc := codec.New()
	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit.Add(deposit...))
	anyAddr := createFakeFundedAccount(ctx, accKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := InitMsg{
		Verifier:    anyAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, err := keeper.InstantiateContract(ctx, contractID, creator, initMsgBz, deposit, true)
	require.NoError(t, err)

	contractModel := []types.Model{
		{Key: []byte("foo"), Value: []byte(`"bar"`)},
		{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)},
	}

	keeper.SetContractStore(ctx, addr, contractModel)

	querier := NewQuerier(keeper)

	// query store []byte("foo")
	bz, err := cdc.MarshalJSON(types.NewQueryRawStoreParams(addr, []byte("foo")))
	require.NoError(t, err)

	res, err := querier(ctx, []string{types.QueryRawStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, []byte(`"bar"`), res)

	// query store []byte{0x0, 0x1}
	bz, err = cdc.MarshalJSON(types.NewQueryRawStoreParams(addr, []byte{0x0, 0x1}))
	require.NoError(t, err)

	res, err = querier(ctx, []string{types.QueryRawStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, []byte(`{"count":8}`), res)

	// query contract []byte(`{"verifier":{}}`)
	bz, err = cdc.MarshalJSON(types.NewQueryContractParams(addr, []byte(`{"verifier":{}}`)))
	require.NoError(t, err)

	res, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()), string(res))

	// query contract []byte(`{"raw":{"key":"config"}}`
	bz, err = cdc.MarshalJSON(types.NewQueryContractParams(addr, []byte(`{"raw":{"key":"config"}}`)))
	require.NoError(t, err)

	_, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
	require.Error(t, err)
}

func TestQueryParams(t *testing.T) {
	cdc := codec.New()
	input := CreateTestInput(t)

	var params types.Params

	res, errRes := queryParameters(input.Ctx, input.WasmKeeper)
	require.NoError(t, errRes)

	err := cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.WasmKeeper.GetParams(input.Ctx), params)
}
