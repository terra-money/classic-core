package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInstantiateExceedMaxGas(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	// must panic
	require.Panics(t, func() {
		params := keeper.GetParams(ctx)
		params.MaxContractGas = types.InstantiateContractCosts(0) + 1
		keeper.SetParams(ctx, params)
		NewMsgServerImpl(keeper).InstantiateContract(ctx.Context(), types.NewMsgInstantiateContract(creator, sdk.AccAddress{}, codeID, initMsgBz, nil))
	})
}

func TestExecuteExceedMaxGas(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, nil)

	// must panic
	require.Panics(t, func() {
		params := keeper.GetParams(ctx)
		params.MaxContractGas = types.InstantiateContractCosts(0) + 1
		keeper.SetParams(ctx, params)
		NewMsgServerImpl(keeper).ExecuteContract(ctx.Context(), types.NewMsgExecuteContract(creator, addr, []byte(`{"release":{}}`), nil))
	})
}

func TestMigrateExceedMaxGas(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, initMsgBz, nil)

	// must panic
	require.Panics(t, func() {
		params := keeper.GetParams(ctx)
		params.MaxContractGas = types.InstantiateContractCosts(0) + 1
		keeper.SetParams(ctx, params)
		NewMsgServerImpl(keeper).MigrateContract(ctx.Context(), types.NewMsgMigrateContract(creator, addr, codeID, []byte(`{"release":{}}`)))
	})
}

func TestStoreCodeMsgServer(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	// Create contract
	msgRes, err := NewMsgServerImpl(keeper).StoreCode(sdk.WrapSDKContext(ctx), &types.MsgStoreCode{
		Sender:       creator.String(),
		WASMByteCode: wasmCode,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), msgRes.CodeID)

	// Verify content
	storedCode, err := keeper.GetByteCode(ctx, msgRes.CodeID)
	require.NoError(t, err)
	require.Equal(t, wasmCode, storedCode)
}

func TestMigrateCodeMsgServer(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	fakeAccount := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	codeID := uint64(1)
	keeper.SetCodeInfo(ctx, codeID, types.CodeInfo{
		CodeID:   1,
		CodeHash: []byte{},
		Creator:  creator.String(),
	})

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	_, err = NewMsgServerImpl(keeper).MigrateCode(sdk.WrapSDKContext(ctx), &types.MsgMigrateCode{
		CodeID:       1,
		Sender:       fakeAccount.String(),
		WASMByteCode: wasmCode,
	})
	require.Error(t, err)

	_, err = NewMsgServerImpl(keeper).MigrateCode(sdk.WrapSDKContext(ctx), &types.MsgMigrateCode{
		CodeID:       1,
		Sender:       creator.String(),
		WASMByteCode: wasmCode,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

}

func TestMigrateContract(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	originalCodeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	newCodeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, originalCodeID, creator, creator, initMsgBz, nil)
	require.NoError(t, err)

	_, err = NewMsgServerImpl(keeper).MigrateContract(sdk.WrapSDKContext(ctx), types.NewMsgMigrateContract(creator, addr, newCodeID, initMsgBz))
	require.NoError(t, err)
}

func TestUpdateContractAdmin(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	newAdmin := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, creator, initMsgBz, nil)
	require.NoError(t, err)

	_, err = NewMsgServerImpl(keeper).UpdateContractAdmin(sdk.WrapSDKContext(ctx), types.NewMsgUpdateContractAdmin(creator, newAdmin, addr))
	require.NoError(t, err)
}

func TestClearContractAdmin(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	wasmCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)

	codeID, err := keeper.StoreCode(ctx, creator, wasmCode)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := HackatomExampleInitMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	addr, _, err := keeper.InstantiateContract(ctx, codeID, creator, creator, initMsgBz, nil)
	require.NoError(t, err)

	_, err = NewMsgServerImpl(keeper).ClearContractAdmin(sdk.WrapSDKContext(ctx), types.NewMsgClearContractAdmin(creator, addr))
	require.NoError(t, err)

}
