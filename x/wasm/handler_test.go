package wasm

import (
	"encoding/json"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/x/wasm/internal/types"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	core "github.com/terra-money/core/types"
)

func TestHandleStore(t *testing.T) {
	loadContracts()

	cases := map[string]struct {
		msg     sdk.Msg
		isValid bool
	}{
		"empty": {
			msg:     MsgStoreCode{},
			isValid: false,
		},
		"invalid wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: []byte("foobar"),
			},
			isValid: false,
		},
		"old wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: oldEscrowContract,
			},
			isValid: false,
		},
		"valid wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: testContract,
			},
			isValid: true,
		},
		"other valid wasm": {
			msg: MsgStoreCode{
				Sender:       addr1,
				WASMByteCode: maskContract,
			},
			isValid: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			data, cleanup := setupTest(t)
			defer cleanup()

			h := data.module.NewHandler()

			res, err := h(data.ctx, tc.msg)
			if !tc.isValid {
				require.Error(t, err, "%#v", res)
				_, err := data.keeper.GetCodeInfo(data.ctx, 1)
				require.Error(t, err)
				return
			}
			require.NoError(t, err, "%#v", res)
			_, err = data.keeper.GetCodeInfo(data.ctx, 1)
			require.NoError(t, err)
		})
	}
}
func TestHandleInstantiate(t *testing.T) {
	loadContracts()

	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}
	_, err := h(data.ctx, msg)
	require.NoError(t, err)

	bytecode, sdkErr := data.keeper.GetByteCode(data.ctx, 1)
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
	initCmd := MsgInstantiateContract{
		Owner:      creator,
		CodeID:     1,
		InitMsg:    initMsgBz,
		InitCoins:  nil,
		Migratable: true,
	}
	res, err := h(data.ctx, initCmd)
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

	contractInfo, err := data.keeper.GetContractInfo(data.ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz, true)
	require.Equal(t, expectedContractInfo, contractInfo)

	iter := data.keeper.GetContractStoreIterator(data.ctx, contractAddr)
	var models []Model
	for ; iter.Valid(); iter.Next() {
		models = append(models, Model{Key: iter.Key(), Value: iter.Value()})
	}

	expectedConfigStore := state{
		Verifier:    wasmTypes.CanonicalAddress(fred),
		Beneficiary: wasmTypes.CanonicalAddress(bob),
		Funder:      wasmTypes.CanonicalAddress(creator),
	}

	assertContractStore(t, models, expectedConfigStore)
}

func TestHandleExecute(t *testing.T) {
	loadContracts()

	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}
	_, err := h(data.ctx, msg)
	require.NoError(t, err)

	bytecode, sdkErr := data.keeper.GetByteCode(data.ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initMsg := initMsg{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := MsgInstantiateContract{
		Owner:      creator,
		CodeID:     1,
		InitMsg:    initMsgBz,
		InitCoins:  deposit,
		Migratable: true,
	}
	res, err := h(data.ctx, initCmd)
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

	contractInfo, err := data.keeper.GetContractInfo(data.ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz, true)
	require.Equal(t, expectedContractInfo, contractInfo)

	// ensure bob doesn't exist
	bobAcct := data.acctKeeper.GetAccount(data.ctx, bob)
	require.Nil(t, bobAcct)

	// ensure funder has reduced balance
	creatorAcct := data.acctKeeper.GetAccount(data.ctx, creator)
	require.NotNil(t, creatorAcct)
	// we started at 2*deposit, should have spent one above
	assert.Equal(t, deposit, creatorAcct.GetCoins())

	// ensure contract has updated balance
	contractAcct := data.acctKeeper.GetAccount(data.ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, deposit, contractAcct.GetCoins())

	execCmd := MsgExecuteContract{
		Sender:     fred,
		Contract:   contractAddr,
		ExecuteMsg: []byte(`{"release":{}}`),
		Coins:      topUp,
	}
	_, err = h(data.ctx, execCmd)
	require.NoError(t, err)

	// ensure bob now exists and got both payments released
	bobAcct = data.acctKeeper.GetAccount(data.ctx, bob)
	require.NotNil(t, bobAcct)
	balance := bobAcct.GetCoins()
	assert.Equal(t, deposit.Add(topUp...), balance)

	// ensure contract has updated balance
	contractAcct = data.acctKeeper.GetAccount(data.ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins(nil), contractAcct.GetCoins())

	iter := data.keeper.GetContractStoreIterator(data.ctx, contractAddr)
	var models []Model
	for ; iter.Valid(); iter.Next() {
		models = append(models, Model{Key: iter.Key(), Value: iter.Value()})
	}

	expectedConfigStore := state{
		Verifier:    wasmTypes.CanonicalAddress(fred),
		Beneficiary: wasmTypes.CanonicalAddress(bob),
		Funder:      wasmTypes.CanonicalAddress(creator),
	}

	assertContractStore(t, models, expectedConfigStore)
}

func TestHandleExecuteEscrow(t *testing.T) {
	loadContracts()

	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}
	_, err := h(data.ctx, msg)
	require.NoError(t, err)

	bytecode, sdkErr := data.keeper.GetByteCode(data.ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initMsg := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := MsgInstantiateContract{
		Owner:      creator,
		CodeID:     1,
		InitMsg:    initMsgBz,
		InitCoins:  deposit,
		Migratable: true,
	}
	res, err := h(data.ctx, initCmd)
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

	contractInfo, err := data.keeper.GetContractInfo(data.ctx, contractAddr)
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz, true)
	require.Equal(t, expectedContractInfo, contractInfo)

	handleMsg := map[string]interface{}{
		"release": map[string]interface{}{},
	}

	handleMsgBz, err := json.Marshal(handleMsg)
	require.NoError(t, err)

	execCmd := MsgExecuteContract{
		Sender:     fred,
		Contract:   contractAddr,
		ExecuteMsg: handleMsgBz,
		Coins:      topUp,
	}

	res, err = h(data.ctx, execCmd)
	require.NoError(t, err)

	// ensure bob now exists and got both payments released
	bobAcct := data.acctKeeper.GetAccount(data.ctx, bob)
	require.NotNil(t, bobAcct)
	balance := bobAcct.GetCoins()
	assert.Equal(t, deposit.Add(topUp...), balance)

	// ensure contract has updated balance
	contractAcct := data.acctKeeper.GetAccount(data.ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins(nil), contractAcct.GetCoins())
}

func TestHandleMigrate(t *testing.T) {
	loadContracts()

	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	storeMsg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}

	// store two same code
	_, err := h(data.ctx, storeMsg)
	require.NoError(t, err)
	res, err := h(data.ctx, storeMsg)
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

	bytecode, sdkErr := data.keeper.GetByteCode(data.ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, testContract, bytecode)

	_, _, bob := keyPubAddr()
	initData := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initDataBz, err := json.Marshal(initData)
	require.NoError(t, err)

	initMsg := MsgInstantiateContract{
		Owner:      creator,
		CodeID:     1,
		InitMsg:    initDataBz,
		InitCoins:  deposit,
		Migratable: true,
	}
	res, err = h(data.ctx, initMsg)
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

	migrateMsg := NewMsgMigrateContract(creator, contractAddr, newCodeID, migDataBz)
	_, err = h(data.ctx, migrateMsg)
	require.NoError(t, err)

	cInfo, err := data.keeper.GetContractInfo(data.ctx, contractAddr)
	require.NoError(t, err)
	assert.Equal(t, newCodeID, cInfo.CodeID)
}

func TestHandleUpdateOwner(t *testing.T) {
	loadContracts()

	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit...))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	storeMsg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}

	// store two same code
	_, err := h(data.ctx, storeMsg)
	require.NoError(t, err)

	_, _, bob := keyPubAddr()
	initData := map[string]interface{}{
		"verifier":    fred.String(),
		"beneficiary": bob.String(),
	}
	initDataBz, err := json.Marshal(initData)
	require.NoError(t, err)

	initMsg := MsgInstantiateContract{
		Owner:      creator,
		CodeID:     1,
		InitMsg:    initDataBz,
		InitCoins:  deposit,
		Migratable: true,
	}
	res, err := h(data.ctx, initMsg)
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

	updateOwnerMsg := NewMsgUpdateContractOwner(creator, fred, contractAddr)
	_, err = h(data.ctx, updateOwnerMsg)
	require.NoError(t, err)

	cInfo, err := data.keeper.GetContractInfo(data.ctx, contractAddr)
	require.NoError(t, err)
	require.Equal(t, fred, cInfo.Owner)
}
