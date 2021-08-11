package wasm_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/wasm"
	"github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
)

func TestHandleStore(t *testing.T) {
	loadContracts()

	cases := map[string]struct {
		msg     sdk.Msg
		isValid bool
	}{
		"empty": {
			msg:     &types.MsgStoreCode{},
			isValid: false,
		},
		"invalid wasm": {
			msg: &types.MsgStoreCode{
				Sender:       addr1.String(),
				WASMByteCode: []byte("foobar"),
			},
			isValid: false,
		},
		"old wasm": {
			msg: &types.MsgStoreCode{
				Sender:       addr1.String(),
				WASMByteCode: oldEscrowContract,
			},
			isValid: false,
		},
		"valid wasm": {
			msg: &types.MsgStoreCode{
				Sender:       addr1.String(),
				WASMByteCode: testContract,
			},
			isValid: true,
		},
		"other valid wasm": {
			msg: &types.MsgStoreCode{
				Sender:       addr1.String(),
				WASMByteCode: reflectContract,
			},
			isValid: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			input := keeper.CreateTestInput(t)

			h := wasm.NewHandler(input.WasmKeeper)

			res, err := h(input.Ctx, tc.msg)
			if !tc.isValid {
				require.Error(t, err, "%#v", res)
				_, err := input.WasmKeeper.GetCodeInfo(input.Ctx, 1)
				require.Error(t, err)
				return
			}
			require.NoError(t, err, "%#v", res)
			_, err = input.WasmKeeper.GetCodeInfo(input.Ctx, 1)
			require.NoError(t, err)
		})
	}
}
func TestHandleInstantiate(t *testing.T) {
	loadContracts()

	input := keeper.CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, deposit)

	h := wasm.NewHandler(input.WasmKeeper)

	msg := types.NewMsgStoreCode(creator, testContract)
	_, err := h(input.Ctx, msg)
	require.NoError(t, err)

	bytecode, sdkErr := input.WasmKeeper.GetByteCode(input.Ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initMsg := initMsg{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
	}

	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	// create with no balance is also legal
	initCmd := types.NewMsgInstantiateContract(creator, sdk.AccAddress{}, 1, initMsgBz, nil)
	res, err := h(input.Ctx, initCmd)
	require.NoError(t, err)

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

	contractInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, sdk.AccAddress{}, initMsgBz)
	require.Equal(t, expectedContractInfo, contractInfo)

	iter := input.WasmKeeper.GetContractStoreIterator(input.Ctx, contractAddr)
	var models []types.Model
	for ; iter.Valid(); iter.Next() {
		models = append(models, types.Model{Key: iter.Key(), Value: iter.Value()})
	}

	expectedConfigStore := state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	}

	assertContractStore(t, models, expectedConfigStore)
}

func TestHandleExecute(t *testing.T) {
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

	initCmd := types.NewMsgInstantiateContract(creator, sdk.AccAddress{}, 1, initMsgBz, deposit)
	res, err := h(input.Ctx, initCmd)
	require.NoError(t, err)

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

	contractInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, sdk.AccAddress{}, initMsgBz)
	require.Equal(t, expectedContractInfo, contractInfo)

	// ensure bob doesn't exist
	bobAcct := input.AccKeeper.GetAccount(input.Ctx, bob)
	require.Nil(t, bobAcct)

	// ensure funder has reduced balance
	creatorAcct := input.AccKeeper.GetAccount(input.Ctx, creator)
	require.NotNil(t, creatorAcct)
	// we started at 2*deposit, should have spent one above
	assert.Equal(t, deposit, input.BankKeeper.GetAllBalances(input.Ctx, creator))

	// ensure contract has updated balance
	contractAcct := input.AccKeeper.GetAccount(input.Ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, deposit, input.BankKeeper.GetAllBalances(input.Ctx, contractAddr))

	execCmd := types.NewMsgExecuteContract(fred, contractAddr, []byte(`{"release":{}}`), topUp)
	_, err = h(input.Ctx, execCmd)
	require.NoError(t, err)

	// ensure bob now exists and got both payments released
	bobAcct = input.AccKeeper.GetAccount(input.Ctx, bob)
	require.NotNil(t, bobAcct)
	assert.Equal(t, deposit.Add(topUp...), input.BankKeeper.GetAllBalances(input.Ctx, bob))

	// ensure contract has updated balance
	contractAcct = input.AccKeeper.GetAccount(input.Ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins{}, input.BankKeeper.GetAllBalances(input.Ctx, contractAddr))

	iter := input.WasmKeeper.GetContractStoreIterator(input.Ctx, contractAddr)
	var models []types.Model
	for ; iter.Valid(); iter.Next() {
		models = append(models, types.Model{Key: iter.Key(), Value: iter.Value()})
	}

	expectedConfigStore := state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	}

	assertContractStore(t, models, expectedConfigStore)
}

func TestHandleExecuteEscrow(t *testing.T) {
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

	bytecode, sdkErr := input.WasmKeeper.GetByteCode(input.Ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initMsg := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := types.NewMsgInstantiateContract(creator, sdk.AccAddress{}, 1, initMsgBz, deposit)
	res, err := h(input.Ctx, initCmd)
	require.NoError(t, err)

	// Retrieve contract address from events
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

	contractInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, sdk.AccAddress{}, initMsgBz)
	require.Equal(t, expectedContractInfo, contractInfo)

	handleMsg := map[string]interface{}{
		"release": map[string]interface{}{},
	}

	handleMsgBz, err := json.Marshal(handleMsg)
	require.NoError(t, err)

	execCmd := types.NewMsgExecuteContract(fred, contractAddr, handleMsgBz, topUp)

	res, err = h(input.Ctx, execCmd)
	require.NoError(t, err)

	// ensure bob now exists and got both payments released
	bobAcct := input.AccKeeper.GetAccount(input.Ctx, bob)
	require.NotNil(t, bobAcct)
	assert.Equal(t, deposit.Add(topUp...), input.BankKeeper.GetAllBalances(input.Ctx, bob))

	// ensure contract has updated balance
	contractAcct := input.AccKeeper.GetAccount(input.Ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins{}, input.BankKeeper.GetAllBalances(input.Ctx, contractAddr))
}

func TestHandleMigrate(t *testing.T) {
	loadContracts()

	input := keeper.CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, topUp)

	h := wasm.NewHandler(input.WasmKeeper)

	storeMsg := types.NewMsgStoreCode(creator, testContract)

	// store two same code
	_, err := h(input.Ctx, storeMsg)
	require.NoError(t, err)
	res, err := h(input.Ctx, storeMsg)
	require.NoError(t, err)

	var newCodeID uint64
	for _, event := range res.Events {
		if event.Type == types.EventTypeStoreCode {
			for _, attr := range event.Attributes {
				if string(attr.GetKey()) == types.AttributeKeyCodeID {
					newCodeID, err = strconv.ParseUint(string(attr.GetValue()), 10, 64)
					require.NoError(t, err)
					break
				}
			}
		}
	}

	bytecode, sdkErr := input.WasmKeeper.GetByteCode(input.Ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initData := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initDataBz, err := json.Marshal(initData)
	require.NoError(t, err)

	initMsg := types.NewMsgInstantiateContract(creator, creator, 1, initDataBz, deposit)
	res, err = h(input.Ctx, initMsg)
	require.NoError(t, err)

	// Retrieve contract address from events
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

	migData := map[string]interface{}{
		"verifier": creator.String(),
	}
	migDataBz, err := json.Marshal(migData)
	require.NoError(t, err)

	migrateMsg := types.NewMsgMigrateContract(creator, contractAddr, newCodeID, migDataBz)
	_, err = h(input.Ctx, migrateMsg)
	require.NoError(t, err)

	cInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	require.NoError(t, err)
	assert.Equal(t, newCodeID, cInfo.CodeID)
}

func TestHandleUpdateAdmin(t *testing.T) {
	loadContracts()

	input := keeper.CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, topUp)

	h := wasm.NewHandler(input.WasmKeeper)

	storeMsg := types.NewMsgStoreCode(creator, testContract)

	// store two same code
	_, err := h(input.Ctx, storeMsg)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initData := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initDataBz, err := json.Marshal(initData)
	require.NoError(t, err)

	initMsg := types.NewMsgInstantiateContract(creator, creator, 1, initDataBz, deposit)
	res, err := h(input.Ctx, initMsg)
	require.NoError(t, err)

	// Retrieve contract address from events
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

	updateAdminMsg := types.NewMsgUpdateContractAdmin(creator, fred, contractAddr)
	_, err = h(input.Ctx, updateAdminMsg)
	require.NoError(t, err)

	cInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	require.NoError(t, err)
	require.Equal(t, fred.String(), cInfo.Admin)
}

func TestHandleClearAdmin(t *testing.T) {
	loadContracts()

	input := keeper.CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(input.Ctx, input.AccKeeper, input.BankKeeper, topUp)

	h := wasm.NewHandler(input.WasmKeeper)

	storeMsg := types.NewMsgStoreCode(creator, testContract)

	// store two same code
	_, err := h(input.Ctx, storeMsg)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initData := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initDataBz, err := json.Marshal(initData)
	require.NoError(t, err)

	initMsg := types.NewMsgInstantiateContract(creator, creator, 1, initDataBz, deposit)
	res, err := h(input.Ctx, initMsg)
	require.NoError(t, err)

	// Retrieve contract address from events
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

	clearAdminMsg := types.NewMsgClearContractAdmin(creator, contractAddr)
	_, err = h(input.Ctx, clearAdminMsg)
	require.NoError(t, err)

	cInfo, err := input.WasmKeeper.GetContractInfo(input.Ctx, contractAddr)
	require.NoError(t, err)
	require.True(t, len(cInfo.Admin) == 0)
}
