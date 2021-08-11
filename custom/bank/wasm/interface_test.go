package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	core "github.com/terra-money/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestEncoding(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
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
		"simple send": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{
					Send: &wasmvmtypes.SendMsg{
						ToAddress: addrs[1].String(),
						Amount: []wasmvmtypes.Coin{
							{
								Denom:  core.MicroLunaDenom,
								Amount: "12345",
							},
							{
								Denom:  "usdt",
								Amount: "54321",
							},
						},
					},
				},
			},
			output: &banktypes.MsgSend{
				FromAddress: addrs[0].String(),
				ToAddress:   addrs[1].String(),
				Amount: sdk.Coins{
					sdk.NewInt64Coin(core.MicroLunaDenom, 12345),
					sdk.NewInt64Coin("usdt", 54321),
				},
			},
		},
		"invalid send amount": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{
					Send: &wasmvmtypes.SendMsg{
						ToAddress: addrs[1].String(),
						Amount: []wasmvmtypes.Coin{
							{
								Denom:  core.MicroLunaDenom,
								Amount: "123.456",
							},
						},
					},
				},
			},
			isError: true,
		},
		"invalid address": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{
					Send: &wasmvmtypes.SendMsg{
						ToAddress: invalidAddr,
						Amount: []wasmvmtypes.Coin{
							{
								Denom:  core.MicroLunaDenom,
								Amount: "7890",
							},
						},
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
