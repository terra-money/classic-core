package wasm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm"
	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInitGenesis(t *testing.T) {
	loadContracts()

	input := keeper.CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, topUp)

	h := wasm.NewHandler(input.WasmKeeper)

	msg := types.NewMsgStoreCode(creator, testContract)
	_, err := h(input.Ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgStoreCode(creator, reflectContract)
	_, err = h(input.Ctx, msg)
	require.NoError(t, err)

	bytecode, sdkErr := input.WasmKeeper.GetByteCode(input.Ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initMsg := initMsg{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := types.NewMsgInstantiateContract(creator, creator, 1, initMsgBz, deposit)
	res, err := h(input.Ctx, initCmd)
	require.NoError(t, err)

	// Check contract address
	var contractAddr sdk.AccAddress
	for _, event := range res.Events {
		if event.Type == types.EventTypeInstantiateContract {
			for _, attr := range event.Attributes {
				if string(attr.GetKey()) == types.AttributeKeyContractAddress {
					contractAddr, err = sdk.AccAddressFromBech32(string(attr.GetValue()))
					require.NoError(t, err)
					break
				}
			}
		}
	}

	require.False(t, contractAddr.Empty())
	_, sdkErr = input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	require.NoError(t, sdkErr)

	execCmd := types.NewMsgExecuteContract(fred, contractAddr, []byte(`{"release":{}}`), topUp)
	_, err = h(input.Ctx, execCmd)
	require.NoError(t, err)

	// ensure all contract state is as after init
	bytecode, sdkErr = input.WasmKeeper.GetByteCode(input.Ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, creator, initMsgBz)
	contractInfo, sdkErr := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	require.NoError(t, sdkErr)
	require.Equal(t, expectedContractInfo, contractInfo)

	iter := input.WasmKeeper.GetContractStoreIterator(input.Ctx, contractAddr)
	var models []types.Model
	for ; iter.Valid(); iter.Next() {
		models = append(models, types.Model{Key: iter.Key(), Value: iter.Value()})
	}

	expectedConfigState := state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	}

	assertContractStore(t, models, expectedConfigState)

	// export into genstate
	genState := wasm.ExportGenesis(input.Ctx, input.WasmKeeper)

	// create new app to import genstate into
	newInput := keeper.CreateTestInput(t)

	// initialize new app with genstate
	wasm.InitGenesis(newInput.Ctx, newInput.WasmKeeper, genState)

	// run same checks again on newdata, to make sure it was reinitialized correctly
	bytecode, err = newInput.WasmKeeper.GetByteCode(newInput.Ctx, 1)
	require.NoError(t, err)
	require.Equal(t, testContract, bytecode)

	contractInfo, err = newInput.WasmKeeper.GetContractInfo(newInput.Ctx, contractAddr)
	require.NoError(t, err)
	require.Equal(t, expectedContractInfo, contractInfo)

	iter = newInput.WasmKeeper.GetContractStoreIterator(newInput.Ctx, contractAddr)
	models = []types.Model{}
	for ; iter.Valid(); iter.Next() {
		models = append(models, types.Model{Key: iter.Key(), Value: iter.Value()})
	}

	assertContractStore(t, models, expectedConfigState)
}
