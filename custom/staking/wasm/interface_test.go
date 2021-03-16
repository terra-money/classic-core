package wasm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
)

func TestEncoding(t *testing.T) {
	addrs := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	valAddr := make(sdk.ValAddress, sdk.AddrLen)
	valAddr[0] = 12
	valAddr2 := make(sdk.ValAddress, sdk.AddrLen)
	valAddr2[1] = 123

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmTypes.CosmosMsg
		// set if valid
		output []sdk.Msg
		// set if invalid
		isError bool
	}{
		"staking delegate to non-validator": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Delegate: &wasmTypes.DelegateMsg{
						Validator: addrs[1].String(),
						Amount:    wasmTypes.NewCoin(777, "stake"),
					},
				},
			},
			isError: true,
		},
		"staking undelegate": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Undelegate: &wasmTypes.UndelegateMsg{
						Validator: valAddr.String(),
						Amount:    wasmTypes.NewCoin(555, "stake"),
					},
				},
			},
			output: []sdk.Msg{
				&stakingtypes.MsgUndelegate{
					DelegatorAddress: addrs[0].String(),
					ValidatorAddress: valAddr.String(),
					Amount:           sdk.NewInt64Coin("stake", 555),
				},
			},
		},
		"staking redelegate": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Redelegate: &wasmTypes.RedelegateMsg{
						SrcValidator: valAddr.String(),
						DstValidator: valAddr2.String(),
						Amount:       wasmTypes.NewCoin(222, "stake"),
					},
				},
			},
			output: []sdk.Msg{
				&stakingtypes.MsgBeginRedelegate{
					DelegatorAddress:    addrs[0].String(),
					ValidatorSrcAddress: valAddr.String(),
					ValidatorDstAddress: valAddr2.String(),
					Amount:              sdk.NewInt64Coin("stake", 222),
				},
			},
		},
		"staking withdraw (implicit recipient)": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Withdraw: &wasmTypes.WithdrawMsg{
						Validator: valAddr2.String(),
					},
				},
			},
			output: []sdk.Msg{
				&distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: addrs[0].String(),
					ValidatorAddress: valAddr2.String(),
				},
			},
		},
		"staking withdraw (explicit recipient)": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Withdraw: &wasmTypes.WithdrawMsg{
						Validator: valAddr2.String(),
						Recipient: addrs[1].String(),
					},
				},
			},
			output: []sdk.Msg{
				&distrtypes.MsgSetWithdrawAddress{
					DelegatorAddress: addrs[0].String(),
					WithdrawAddress:  addrs[1].String(),
				},
				&distrtypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: addrs[0].String(),
					ValidatorAddress: valAddr2.String(),
				},
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
