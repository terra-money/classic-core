package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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
		"staking delegate to non-validator": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Staking: &wasmvmtypes.StakingMsg{
					Delegate: &wasmvmtypes.DelegateMsg{
						Validator: addrs[1].String(),
						Amount:    wasmvmtypes.NewCoin(777, "stake"),
					},
				},
			},
			isError: true,
		},
		"staking undelegate": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Staking: &wasmvmtypes.StakingMsg{
					Undelegate: &wasmvmtypes.UndelegateMsg{
						Validator: valAddr.String(),
						Amount:    wasmvmtypes.NewCoin(555, "stake"),
					},
				},
			},
			output: &stakingtypes.MsgUndelegate{
				DelegatorAddress: addrs[0].String(),
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewInt64Coin("stake", 555),
			},
		},
		"staking redelegate": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Staking: &wasmvmtypes.StakingMsg{
					Redelegate: &wasmvmtypes.RedelegateMsg{
						SrcValidator: valAddr.String(),
						DstValidator: valAddr2.String(),
						Amount:       wasmvmtypes.NewCoin(222, "stake"),
					},
				},
			},
			output: &stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    addrs[0].String(),
				ValidatorSrcAddress: valAddr.String(),
				ValidatorDstAddress: valAddr2.String(),
				Amount:              sdk.NewInt64Coin("stake", 222),
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
