package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func TestEncoding(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	valAddr := make(sdk.ValAddress, 1)
	valAddr[0] = 12
	valAddr2 := make(sdk.ValAddress, 2)
	valAddr2[1] = 123

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmvmtypes.CosmosMsg
		// set if valid
		output sdk.Msg
		// set if invalid
		isError bool
	}{
		"distribution withdraw ": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Distribution: &wasmvmtypes.DistributionMsg{
					WithdrawDelegatorReward: &wasmvmtypes.WithdrawDelegatorRewardMsg{
						Validator: valAddr2.String(),
					},
				},
			},
			output: &distrtypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: addrs[0].String(),
				ValidatorAddress: valAddr2.String(),
			},
		},
		"staking withdraw ": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Distribution: &wasmvmtypes.DistributionMsg{
					SetWithdrawAddress: &wasmvmtypes.SetWithdrawAddressMsg{
						Address: addrs[1].String(),
					},
				},
			},
			output: &distrtypes.MsgSetWithdrawAddress{
				DelegatorAddress: addrs[0].String(),
				WithdrawAddress:  addrs[1].String(),
			},
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
