package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/wasm/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestLegacyContractState(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	anyAddr := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    anyAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, contractID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)

	contractModel := []types.Model{
		{Key: []byte("foo"), Value: []byte(`"bar"`)},
		{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)},
	}

	keeper.SetContractStore(ctx, addr, contractModel)

	querier := NewLegacyQuerier(keeper, input.Cdc)

	// query store []byte("foo")
	bz, err := input.Cdc.MarshalJSON(types.NewQueryRawStoreParams(addr, []byte("foo")))
	require.NoError(t, err)

	res, err := querier(ctx, []string{types.QueryRawStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, []byte(`"bar"`), res)

	// query store []byte{0x0, 0x1}
	bz, err = input.Cdc.MarshalJSON(types.NewQueryRawStoreParams(addr, []byte{0x0, 0x1}))
	require.NoError(t, err)

	res, err = querier(ctx, []string{types.QueryRawStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, []byte(`{"count":8}`), res)

	// query contract []byte(`{"verifier":{}}`)
	bz, err = input.Cdc.MarshalJSON(types.NewQueryContractParams(addr, []byte(`{"verifier":{}}`)))
	require.NoError(t, err)

	res, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()), string(res))

	// query contract []byte(`{"raw":{"key":"config"}}`
	bz, err = input.Cdc.MarshalJSON(types.NewQueryContractParams(addr, []byte(`{"raw":{"key":"config"}}`)))
	require.NoError(t, err)

	_, err = querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
	require.Error(t, err)

	bz, err = input.Cdc.MarshalJSON(types.NewQueryCodeIDParams(contractID))
	_, err = querier(ctx, []string{types.QueryGetByteCode}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)

	_, err = querier(ctx, []string{types.QueryGetCodeInfo}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)

	bz, err = input.Cdc.MarshalJSON(types.NewQueryContractAddressParams(addr))
	_, err = querier(ctx, []string{types.QueryGetContractInfo}, abci.RequestQuery{Data: []byte(bz)})
	require.NoError(t, err)

	_, err = querier(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)
}

func TestLegacyParams(t *testing.T) {
	input := CreateTestInput(t)

	var params types.Params

	res, errRes := queryParameters(input.Ctx, input.WasmKeeper, input.Cdc)
	require.NoError(t, errRes)

	err := input.Cdc.UnmarshalJSON(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.WasmKeeper.GetParams(input.Ctx), params)
}

func TestLegacyMultipleGoroutines(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit.Add(deposit...))
	anyAddr := createFakeFundedAccount(ctx, accKeeper, bankKeeper, topUp)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	contractID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initMsg := HackatomExampleInitMsg{
		Verifier:    anyAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, contractID, creator, sdk.AccAddress{}, initMsgBz, deposit)
	require.NoError(t, err)

	contractModel := []types.Model{
		{Key: []byte("foo"), Value: []byte(`"bar"`)},
		{Key: []byte{0x0, 0x1}, Value: []byte(`{"count":8}`)},
	}

	keeper.SetContractStore(ctx, addr, contractModel)

	querier := NewLegacyQuerier(keeper, input.Cdc)

	wg := &sync.WaitGroup{}
	testCases := 100
	wg.Add(testCases)
	for n := 0; n < testCases; n++ {
		go func() {
			// query contract []byte(`{"verifier":{}}`)
			bz, err := input.Cdc.MarshalJSON(types.NewQueryContractParams(addr, []byte(`{"verifier":{}}`)))
			require.NoError(t, err)

			res, err := querier(ctx, []string{types.QueryContractStore}, abci.RequestQuery{Data: []byte(bz)})
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()), string(res))

			wg.Done()
		}()
	}
	wg.Wait()
}
