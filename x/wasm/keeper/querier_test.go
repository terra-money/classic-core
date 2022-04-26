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
)

func TestQueryContractState(t *testing.T) {
	input := CreateTestInput(t)
	goCtx := sdk.WrapSDKContext(input.Ctx)
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

	querier := NewQuerier(keeper)

	// query store []byte("foo")
	res, err := querier.RawStore(goCtx, &types.QueryRawStoreRequest{ContractAddress: addr.String(), Key: []byte("foo")})
	require.NoError(t, err)
	require.Equal(t, []byte(`"bar"`), res.Data)

	// query store []byte{0x0, 0x1}
	res, err = querier.RawStore(goCtx, &types.QueryRawStoreRequest{ContractAddress: addr.String(), Key: []byte{0x0, 0x1}})
	require.NoError(t, err)
	require.Equal(t, []byte(`{"count":8}`), res.Data)

	// invalid argurment
	_, err = querier.RawStore(goCtx, &types.QueryRawStoreRequest{ContractAddress: ``, Key: []byte{0x0, 0x1}})
	require.Error(t, err)

	// query contract []byte(`{"verifier":{}}`)
	res2, err := querier.ContractStore(goCtx, &types.QueryContractStoreRequest{ContractAddress: addr.String(), QueryMsg: []byte(`{"verifier":{}}`)})
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()), string(res2.QueryResult))

	// query contract []byte(`{"raw":{"key":"config"}}`
	_, err = querier.ContractStore(goCtx, &types.QueryContractStoreRequest{ContractAddress: addr.String(), QueryMsg: []byte(`{"raw":{"key":"config"}}`)})
	require.Error(t, err)

	// invalid arguemnt
	_, err = querier.ContractStore(goCtx, &types.QueryContractStoreRequest{ContractAddress: ``, QueryMsg: []byte(`{"raw":{"key":"config"}}`)})
	require.Error(t, err)

	// out-of-gas panic
	keeper.wasmConfig.ContractQueryGasLimit = 0
	querier = NewQuerier(keeper)

	_, err = querier.ContractStore(goCtx, &types.QueryContractStoreRequest{ContractAddress: addr.String(), QueryMsg: []byte(`{"verifier":{}}`)})
	require.Error(t, err)

}

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	goCtx := sdk.WrapSDKContext(input.Ctx)
	querier := NewQuerier(input.WasmKeeper)
	res, err := querier.Params(goCtx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, input.WasmKeeper.GetParams(input.Ctx), res.Params)
}

func TestQueryMultipleGoroutines(t *testing.T) {
	input := CreateTestInput(t)
	goCtx := sdk.WrapSDKContext(input.Ctx)
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

	querier := NewQuerier(keeper)

	wg := &sync.WaitGroup{}
	testCases := 100
	wg.Add(testCases)
	for n := 0; n < testCases; n++ {
		go func() {
			// query contract []byte(`{"verifier":{}}`)
			res2, err := querier.ContractStore(goCtx, &types.QueryContractStoreRequest{ContractAddress: addr.String(), QueryMsg: []byte(`{"verifier":{}}`)})
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf(`{"verifier":"%s"}`, anyAddr.String()), string(res2.QueryResult))
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestQueryCodeAndContractInfo(t *testing.T) {
	input := CreateTestInput(t)
	goCtx := sdk.WrapSDKContext(input.Ctx)
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

	querier := NewQuerier(keeper)

	res, err := querier.CodeInfo(goCtx, &types.QueryCodeInfoRequest{CodeId: contractID})
	require.NoError(t, err)
	require.Equal(t, creator.String(), res.GetCodeInfo().Creator)
	require.Equal(t, contractID, res.GetCodeInfo().CodeID)

	res2, err := querier.ByteCode(goCtx, &types.QueryByteCodeRequest{CodeId: contractID})
	require.NoError(t, err)
	require.Equal(t, res2.GetByteCode(), wasmCode)

	res3, err := querier.ContractInfo(goCtx, &types.QueryContractInfoRequest{ContractAddress: addr.String()})
	require.NoError(t, err)
	require.Equal(t, addr.String(), res3.GetContractInfo().Address)
	require.Equal(t, "", res3.GetContractInfo().Admin)
	require.Equal(t, creator.String(), res3.GetContractInfo().Creator)
	require.Equal(t, contractID, res3.GetContractInfo().CodeID)
	queriedInitMsg, err := res3.GetContractInfo().InitMsg.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, initMsgBz, queriedInitMsg)
}
