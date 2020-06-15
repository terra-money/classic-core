package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
)

func TestEncoding(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(2)
	invalidAddr := "xrnd1d02kd90n38qvr3qb9qof83fn2d2"

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmTypes.CosmosMsg
		// set if valid
		output []sdk.Msg
		// set if invalid
		isError bool
	}{
		"simple send": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addrs[0].String(),
						ToAddress:   addrs[1].String(),
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
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
			output: []sdk.Msg{
				bank.MsgSend{
					FromAddress: addrs[0],
					ToAddress:   addrs[1],
					Amount: sdk.Coins{
						sdk.NewInt64Coin("uatom", 12345),
						sdk.NewInt64Coin("usdt", 54321),
					},
				},
			},
		},
		"invalid send amount": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addrs[0].String(),
						ToAddress:   addrs[1].String(),
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
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
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addrs[0].String(),
						ToAddress:   invalidAddr,
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
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
