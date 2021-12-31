package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-money/core/types"
)

type IbcReflectExampleInitMsg struct {
	ReflectCodeID uint64 `json:"reflect_code_id"`
}

type ibcReflectWhoAmIPayload struct {
	Account string `json:"account"`
}

type ibcReflectDispatchPayload struct {
	Ok interface{} `json:"ok"`
}

func TestDontBindPortNonIBCContract(t *testing.T) {
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

	contractInfo, err := keeper.GetContractInfo(ctx, addr)
	require.NoError(t, err)
	require.Empty(t, contractInfo.GetIBCPortID())
}

func TestBindPort(t *testing.T) {
	input := CreateTestInput(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	// bind test
	bindPort := "bindTest"
	err := keeper.bindIbcPort(ctx, bindPort)
	require.NoError(t, err)
	// bind another port
	bindPort2 := "bindTest2"
	err = keeper.bindIbcPort(ctx, bindPort2)
	require.NoError(t, err)
	// bind with same port will panic
	defer func() {
		if r := recover(); r != nil {
			require.Equal(t, "port bindTest is already bound", r)
		}
	}()
	err = keeper.bindIbcPort(ctx, bindPort)
}

func TestEnsureIBCPort(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload ibc reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// upload ibc reflect send code
	ibcReflectSendCode, err := ioutil.ReadFile("./testdata/ibc_reflect_send.wasm")
	require.NoError(t, err)
	ibcReflectSendID, err := keeper.StoreCode(ctx, creator, ibcReflectSendCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), ibcReflectSendID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// now we set contract as verifier of an ibc-reflect
	initMsg := IbcReflectExampleInitMsg{
		ReflectCodeID: reflectID,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	ibcReflectSendStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000))
	ibcReflectSendAddr, _, err := keeper.InstantiateContract(ctx, ibcReflectSendID, creator, sdk.AccAddress{}, initMsgBz, ibcReflectSendStart)
	require.NoError(t, err)
	require.NotEmpty(t, ibcReflectSendAddr)

	// check each ensured port always same for same contract
	ensuredPort1, err := keeper.ensureIbcPort(ctx, ibcReflectSendAddr)
	require.NoError(t, err)
	require.NotEmpty(t, ensuredPort1)
	require.Equal(t, "wasm."+ibcReflectSendAddr.String(), ensuredPort1)
	ensuredPort2, err := keeper.ensureIbcPort(ctx, ibcReflectSendAddr)
	require.NoError(t, err)
	require.Equal(t, ensuredPort1, ensuredPort2)

	// check ensured port is different for each contract
	reflectEnsuredPort, err := keeper.ensureIbcPort(ctx, reflectAddr)
	require.NoError(t, err)
	require.Equal(t, "wasm."+ibcReflectSendAddr.String(), ensuredPort2)
	require.NotEmpty(t, reflectEnsuredPort)
	require.NotEqual(t, reflectEnsuredPort, ensuredPort1)

	contractInfo, err := keeper.GetContractInfo(ctx, reflectAddr)
	require.NoError(t, err)
	require.Empty(t, contractInfo.GetIBCPortID())

	contractInfo, err = keeper.GetContractInfo(ctx, ibcReflectSendAddr)
	require.NoError(t, err)
	require.NotEmpty(t, contractInfo.GetIBCPortID())
}
