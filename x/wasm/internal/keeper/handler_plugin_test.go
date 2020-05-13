package keeper

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/bank"
// 	"github.com/cosmos/cosmos-sdk/x/distribution"
// 	"github.com/cosmos/cosmos-sdk/x/staking"

// 	"github.com/terra-project/core/x/wasm/internal/types"

// 	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
// )

// func TestEncoding(t *testing.T) {
// 	_, _, addr1 := keyPubAddr()
// 	_, _, addr2 := keyPubAddr()
// 	invalidAddr := "xrnd1d02kd90n38qvr3qb9qof83fn2d2"
// 	valAddr := make(sdk.ValAddress, sdk.AddrLen)
// 	valAddr[0] = 12
// 	valAddr2 := make(sdk.ValAddress, sdk.AddrLen)
// 	valAddr2[1] = 123

// 	jsonMsg := json.RawMessage(`{"foo": 123}`)

// 	cases := map[string]struct {
// 		sender sdk.AccAddress
// 		input  wasmTypes.CosmosMsg
// 		// set if valid
// 		output []sdk.Msg
// 		// set if invalid
// 		isError bool
// 	}{
// 		"simple send": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Bank: &wasmTypes.BankMsg{
// 					Send: &wasmTypes.SendMsg{
// 						FromAddress: addr1.String(),
// 						ToAddress:   addr2.String(),
// 						Amount: []wasmTypes.Coin{
// 							{
// 								Denom:  "uatom",
// 								Amount: "12345",
// 							},
// 							{
// 								Denom:  "usdt",
// 								Amount: "54321",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				bank.MsgSend{
// 					FromAddress: addr1,
// 					ToAddress:   addr2,
// 					Amount: sdk.Coins{
// 						sdk.NewInt64Coin("uatom", 12345),
// 						sdk.NewInt64Coin("usdt", 54321),
// 					},
// 				},
// 			},
// 		},
// 		"invalid send amount": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Bank: &wasmTypes.BankMsg{
// 					Send: &wasmTypes.SendMsg{
// 						FromAddress: addr1.String(),
// 						ToAddress:   addr2.String(),
// 						Amount: []wasmTypes.Coin{
// 							{
// 								Denom:  "uatom",
// 								Amount: "123.456",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			isError: true,
// 		},
// 		"invalid address": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Bank: &wasmTypes.BankMsg{
// 					Send: &wasmTypes.SendMsg{
// 						FromAddress: addr1.String(),
// 						ToAddress:   invalidAddr,
// 						Amount: []wasmTypes.Coin{
// 							{
// 								Denom:  "uatom",
// 								Amount: "7890",
// 							},
// 						},
// 					},
// 				},
// 			},
// 			isError: true,
// 		},
// 		"wasm execute": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Wasm: &wasmTypes.WasmMsg{
// 					Execute: &wasmTypes.ExecuteMsg{
// 						ContractAddr: addr2.String(),
// 						Msg:          jsonMsg,
// 						Send: []wasmTypes.Coin{
// 							wasmTypes.NewCoin(12, "eth"),
// 						},
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				types.MsgExecuteContract{
// 					Sender:   addr1,
// 					Contract: addr2,
// 					Msg:      jsonMsg,
// 					Coins:    sdk.NewCoins(sdk.NewInt64Coin("eth", 12)),
// 				},
// 			},
// 		},
// 		"wasm instantiate": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Wasm: &wasmTypes.WasmMsg{
// 					Instantiate: &wasmTypes.InstantiateMsg{
// 						CodeID: 7,
// 						Msg:    jsonMsg,
// 						Send: []wasmTypes.Coin{
// 							wasmTypes.NewCoin(123, "eth"),
// 						},
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				types.MsgInstantiateContract{
// 					Sender:    addr1,
// 					CodeID:    7,
// 					InitMsg:   jsonMsg,
// 					InitCoins: sdk.NewCoins(sdk.NewInt64Coin("eth", 123)),
// 				},
// 			},
// 		},
// 		"staking delegate": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Delegate: &wasmTypes.DelegateMsg{
// 						Validator: valAddr.String(),
// 						Amount:    wasmTypes.NewCoin(777, "stake"),
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				staking.MsgDelegate{
// 					DelegatorAddress: addr1,
// 					ValidatorAddress: valAddr,
// 					Amount:           sdk.NewInt64Coin("stake", 777),
// 				},
// 			},
// 		},
// 		"staking delegate to non-validator": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Delegate: &wasmTypes.DelegateMsg{
// 						Validator: addr2.String(),
// 						Amount:    wasmTypes.NewCoin(777, "stake"),
// 					},
// 				},
// 			},
// 			isError: true,
// 		},
// 		"staking undelegate": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Undelegate: &wasmTypes.UndelegateMsg{
// 						Validator: valAddr.String(),
// 						Amount:    wasmTypes.NewCoin(555, "stake"),
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				staking.MsgUndelegate{
// 					DelegatorAddress: addr1,
// 					ValidatorAddress: valAddr,
// 					Amount:           sdk.NewInt64Coin("stake", 555),
// 				},
// 			},
// 		},
// 		"staking redelegate": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Redelegate: &wasmTypes.RedelegateMsg{
// 						SrcValidator: valAddr.String(),
// 						DstValidator: valAddr2.String(),
// 						Amount:       wasmTypes.NewCoin(222, "stake"),
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				staking.MsgBeginRedelegate{
// 					DelegatorAddress:    addr1,
// 					ValidatorSrcAddress: valAddr,
// 					ValidatorDstAddress: valAddr2,
// 					Amount:              sdk.NewInt64Coin("stake", 222),
// 				},
// 			},
// 		},
// 		"staking withdraw (implicit recipient)": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Withdraw: &wasmTypes.WithdrawMsg{
// 						Validator: valAddr2.String(),
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				distribution.MsgSetWithdrawAddress{
// 					DelegatorAddress: addr1,
// 					WithdrawAddress:  addr1,
// 				},
// 				distribution.MsgWithdrawDelegatorReward{
// 					DelegatorAddress: addr1,
// 					ValidatorAddress: valAddr2,
// 				},
// 			},
// 		},
// 		"staking withdraw (explicit recipient)": {
// 			sender: addr1,
// 			input: wasmTypes.CosmosMsg{
// 				Staking: &wasmTypes.StakingMsg{
// 					Withdraw: &wasmTypes.WithdrawMsg{
// 						Validator: valAddr2.String(),
// 						Recipient: addr2.String(),
// 					},
// 				},
// 			},
// 			output: []sdk.Msg{
// 				distribution.MsgSetWithdrawAddress{
// 					DelegatorAddress: addr1,
// 					WithdrawAddress:  addr2,
// 				},
// 				distribution.MsgWithdrawDelegatorReward{
// 					DelegatorAddress: addr1,
// 					ValidatorAddress: valAddr2,
// 				},
// 			},
// 		},
// 	}

// 	encoder := defaultEncoders()
// 	for name, tc := range cases {
// 		tc := tc
// 		t.Run(name, func(t *testing.T) {
// 			res, err := encoder.encode(tc.sender, tc.input)
// 			if tc.isError {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 				assert.Equal(t, tc.output, res)
// 			}
// 		})
// 	}

// }
