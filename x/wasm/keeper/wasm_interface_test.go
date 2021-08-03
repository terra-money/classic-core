package keeper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/types"
)

func TestEcoding(t *testing.T) {
	pubKeys := []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	addrs := []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
	}

	invalidAddr := "xrnd1d02kd90n38qvr3qb9qof83fn2d2"
	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmvmtypes.CosmosMsg
		// set if valid
		output sdk.Msg
		// set if invalid
		isError bool
	}{
		"simple execute": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					Execute: &wasmvmtypes.ExecuteMsg{
						ContractAddr: addrs[1].String(),
						Msg:          []byte("{}"),
						Funds:        wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)},
					},
				},
			},
			output: &types.MsgExecuteContract{
				Sender:     addrs[0].String(),
				Contract:   addrs[1].String(),
				Coins:      sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1234)),
				ExecuteMsg: []byte("{}"),
			},
		},
		"simple instantiate without admin": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					Instantiate: &wasmvmtypes.InstantiateMsg{
						CodeID: 1,
						Msg:    []byte("{}"),
						Funds:  wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)},
					},
				},
			},
			output: &types.MsgInstantiateContract{
				Sender:    addrs[0].String(),
				CodeID:    1,
				InitCoins: sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1234)),
				InitMsg:   []byte("{}"),
			},
		},
		"simple instantiate with admin": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					Instantiate: &wasmvmtypes.InstantiateMsg{
						CodeID: 1,
						Msg:    []byte("{}"),
						Funds:  wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)},
						Admin:  addrs[0].String(),
					},
				},
			},
			output: &types.MsgInstantiateContract{
				Sender:    addrs[0].String(),
				CodeID:    1,
				InitCoins: sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1234)),
				InitMsg:   []byte("{}"),
				Admin:     addrs[0].String(),
			},
		},
		"simple migrate": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					Migrate: &wasmvmtypes.MigrateMsg{
						ContractAddr: addrs[1].String(),
						Msg:          []byte("{}"),
						NewCodeID:    1,
					},
				},
			},
			output: &types.MsgMigrateContract{
				Admin:      addrs[0].String(),
				NewCodeID:  1,
				Contract:   addrs[1].String(),
				MigrateMsg: []byte("{}"),
			},
		},
		"simple update admin": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					UpdateAdmin: &wasmvmtypes.UpdateAdminMsg{
						ContractAddr: addrs[1].String(),
						Admin:        addrs[2].String(),
					},
				},
			},
			output: &types.MsgUpdateContractAdmin{
				Admin:    addrs[0].String(),
				Contract: addrs[1].String(),
				NewAdmin: addrs[2].String(),
			},
		},
		"simple clear admin": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					ClearAdmin: &wasmvmtypes.ClearAdminMsg{
						ContractAddr: addrs[1].String(),
					},
				},
			},
			output: &types.MsgClearContractAdmin{
				Admin:    addrs[0].String(),
				Contract: addrs[1].String(),
			},
		},
		"invalid address execute": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Wasm: &wasmvmtypes.WasmMsg{
					Execute: &wasmvmtypes.ExecuteMsg{
						ContractAddr: invalidAddr,
						Msg:          []byte("{}"),
						Funds:        wasmvmtypes.Coins{wasmvmtypes.NewCoin(1234, core.MicroLunaDenom)},
					},
				},
			},
			isError: true,
		},
	}

	parser := NewWasmMsgParser()
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := parser.Parse(tc.sender, tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}

func TestQueryRaw(t *testing.T) {
	input := CreateTestInput(t)

	input.WasmKeeper.SetContractStore(input.Ctx, Addrs[0], []types.Model{
		{
			Key:   []byte("key1"),
			Value: []byte("value1"),
		},
		{
			Key:   []byte("key2"),
			Value: []byte("value2"),
		},
		{
			Key:   []byte("key3"),
			Value: []byte("value3"),
		},
		{
			Key:   []byte("key4"),
			Value: []byte("value4"),
		},
	})

	querier := NewWasmQuerier(input.WasmKeeper)
	res, err := querier.Query(input.Ctx, wasmvmtypes.QueryRequest{
		Wasm: &wasmvmtypes.WasmQuery{
			Raw: &wasmvmtypes.RawQuery{
				ContractAddr: Addrs[0].String(),
				Key:          []byte("key1"),
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, res, []byte("value1"))
}

func TestQueryContractInfo(t *testing.T) {
	input := CreateTestInput(t)

	input.WasmKeeper.SetContractInfo(input.Ctx, Addrs[0], types.NewContractInfo(1, Addrs[0], Addrs[1], sdk.AccAddress{}, []byte{}))

	bz, err := json.Marshal(CosmosQuery{
		ContractInfo: &ContractInfoQueryParams{
			ContractAddress: Addrs[0].String(),
		},
	})
	require.NoError(t, err)

	querier := NewWasmQuerier(input.WasmKeeper)
	res, err := querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var contractInfoResponse ContractInfoQueryResponse
	err = json.Unmarshal(res, &contractInfoResponse)
	require.NoError(t, err)

	require.Equal(t, contractInfoResponse, ContractInfoQueryResponse{
		Address: Addrs[0].String(),
		Creator: Addrs[1].String(),
		Admin:   "",
		CodeID:  1,
	})
}
