package wasm

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terra-project/core/x/wasm/internal/types"

	wasmTypes "github.com/confio/go-cosmwasm/types"
)

func TestHandleCreate(t *testing.T) {
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
				WASMByteCode: escrowContract,
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

			res := h(data.ctx, tc.msg)
			if !tc.isValid {
				require.False(t, res.IsOK(), "%#v", res)
				_, err := data.keeper.GetCodeInfo(data.ctx, 1)
				require.Error(t, err)
				return
			}
			require.True(t, res.IsOK(), "%#v", res)
			_, err := data.keeper.GetCodeInfo(data.ctx, 1)
			require.NoError(t, err)
		})
	}
}
func TestHandleInstantiate(t *testing.T) {
	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}
	res := h(data.ctx, msg)
	require.True(t, res.IsOK())

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
		Sender:    creator,
		CodeID:    1,
		InitMsg:   initMsgBz,
		InitCoins: nil,
	}
	res = h(data.ctx, initCmd)
	fmt.Print(res.Log)
	require.True(t, res.IsOK())

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
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz)
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
	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: testContract,
	}
	res := h(data.ctx, msg)
	require.True(t, res.IsOK())

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
		Sender:    creator,
		CodeID:    1,
		InitMsg:   initMsgBz,
		InitCoins: deposit,
	}
	res = h(data.ctx, initCmd)
	require.True(t, res.IsOK())

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
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz)
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
		Sender:   fred,
		Contract: contractAddr,
		Msg:      []byte(`{"release":{}}`),
		Coins:    topUp,
	}
	res = h(data.ctx, execCmd)
	require.True(t, res.IsOK())

	// ensure bob now exists and got both payments released
	bobAcct = data.acctKeeper.GetAccount(data.ctx, bob)
	require.NotNil(t, bobAcct)
	balance := bobAcct.GetCoins()
	assert.Equal(t, deposit.Add(topUp), balance)

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
	data, cleanup := setupTest(t)
	defer cleanup()

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := createFakeFundedAccount(data.ctx, data.acctKeeper, deposit.Add(deposit))
	fred := createFakeFundedAccount(data.ctx, data.acctKeeper, topUp)

	h := data.module.NewHandler()

	msg := MsgStoreCode{
		Sender:       creator,
		WASMByteCode: escrowContract,
	}
	res := h(data.ctx, &msg)
	require.True(t, res.IsOK())

	bytecode, sdkErr := data.keeper.GetByteCode(data.ctx, 1)
	require.NoError(t, sdkErr)
	require.Equal(t, escrowContract, bytecode)

	_, _, bob := keyPubAddr()
	initMsg := map[string]interface{}{
		"arbiter":    fred.String(),
		"recipient":  bob.String(),
		"end_time":   0,
		"end_height": 0,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := MsgInstantiateContract{
		Sender:    creator,
		CodeID:    1,
		InitMsg:   initMsgBz,
		InitCoins: deposit,
	}
	res = h(data.ctx, initCmd)
	require.True(t, res.IsOK())

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
	expectedContractInfo := types.NewContractInfo(1, contractAddr, creator, initMsgBz)
	require.Equal(t, expectedContractInfo, contractInfo)

	handleMsg := map[string]interface{}{
		"approve": map[string]interface{}{},
	}

	handleMsgBz, err := json.Marshal(handleMsg)
	require.NoError(t, err)

	execCmd := MsgExecuteContract{
		Sender:   fred,
		Contract: contractAddr,
		Msg:      handleMsgBz,
		Coins:    topUp,
	}

	res = h(data.ctx, execCmd)
	require.True(t, res.IsOK())

	// ensure bob now exists and got both payments released
	bobAcct := data.acctKeeper.GetAccount(data.ctx, bob)
	require.NotNil(t, bobAcct)
	balance := bobAcct.GetCoins()
	assert.Equal(t, deposit.Add(topUp), balance)

	// ensure contract has updated balance
	contractAcct := data.acctKeeper.GetAccount(data.ctx, contractAddr)
	require.NotNil(t, contractAcct)
	assert.Equal(t, sdk.Coins(nil), contractAcct.GetCoins())
}
