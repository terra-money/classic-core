package types_test

import (
	"fmt"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/golang/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bankwasm "github.com/terra-money/core/custom/bank/wasm"
	distrwasm "github.com/terra-money/core/custom/distribution/wasm"
	govwasm "github.com/terra-money/core/custom/gov/wasm"
	stakingwasm "github.com/terra-money/core/custom/staking/wasm"
	core "github.com/terra-money/core/types"
	markettypes "github.com/terra-money/core/x/market/types"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	test_util "github.com/terra-money/core/x/wasm/keeper"
	"github.com/terra-money/core/x/wasm/keeper/wasmtesting"
	"github.com/terra-money/core/x/wasm/types"
)

func TestParse(t *testing.T) {
	input := test_util.CreateTestInput(t)
	ctx := input.Ctx
	encodingConfig := test_util.MakeEncodingConfig(t)

	addrs := []sdk.AccAddress{
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
		sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}

	invalidAddr := "xrnd1d02kd90n38qvr3qb9qof83fn2d2"

	bankMsg := banktypes.NewMsgSend(test_util.Addrs[0], test_util.Addrs[1], sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 10000)))
	bankMsgBin, err := proto.Marshal(bankMsg)
	require.NoError(t, err)

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmvmtypes.CosmosMsg
		// set if valid
		output sdk.Msg
		// set if invalid
		isError bool
	}{
		"BankMsg : invalid address": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{
					Send: &wasmvmtypes.SendMsg{
						ToAddress: invalidAddr,
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
		"BankMsg : simple send": {
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
		"BankMsg : burn coins": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{
					Burn: &wasmvmtypes.BurnMsg{
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
			isError: true,
		},
		"BankMsg : invalid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Bank: &wasmvmtypes.BankMsg{},
			},
			isError: true,
		},
		"GovMsg : vote yes": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Gov: &wasmvmtypes.GovMsg{
					Vote: &wasmvmtypes.VoteMsg{
						ProposalId: 1,
						Vote:       wasmvmtypes.Yes,
					},
				},
			},
			output: &govtypes.MsgVote{
				Voter:      addrs[0].String(),
				ProposalId: 1,
				Option:     govtypes.OptionYes,
			},
		},
		"DistribuitionMsg : invalid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Distribution: &wasmvmtypes.DistributionMsg{
					WithdrawDelegatorReward: &wasmvmtypes.WithdrawDelegatorRewardMsg{
						Validator: invalidAddr,
					},
				},
			},
			isError: true,
		},
		"StakingMsg : invalid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Staking: &wasmvmtypes.StakingMsg{
					Delegate: &wasmvmtypes.DelegateMsg{
						Validator: invalidAddr,
						Amount:    wasmvmtypes.NewCoin(777, "stake"),
					},
				},
			},
			isError: true,
		},
		"WasmMsg : invalid msg": {
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
		"IBCMsg : no open channel": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				IBC: &wasmvmtypes.IBCMsg{
					CloseChannel: &wasmvmtypes.CloseChannelMsg{ChannelID: "testchannel-1"},
				},
			},
			isError: true, // can't open channel in test
		},
		"empty msg": {
			sender:  addrs[0],
			input:   wasmvmtypes.CosmosMsg{},
			isError: true, // can't open channel in test
		},
		"CustomMsg : invalid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(""),
			},
			isError: true,
		},
		"CustomMsg : no registered route": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(`{"route":"treasury", "msg_data":{"swap":{"trader":"cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5","offer_coin":{"denom":"uluna","amount":"10000"},"ask_denom":"uusd"}}}`),
			},
			isError: true,
		},
		"CustomMsg : valid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(fmt.Sprintf(`{"route":"market", "msg_data":{"swap":{"trader":"%s","offer_coin":{"denom":"uluna","amount":"10000"},"ask_denom":"uusd"}}}`, addrs[0])),
			},
			output: &markettypes.MsgSwap{
				Trader:    addrs[0].String(),
				OfferCoin: sdk.NewInt64Coin(core.MicroLunaDenom, 10000),
				AskDenom:  core.MicroUSDDenom,
			},
		},
		"StargateMsg : valid msg": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Stargate: &wasmvmtypes.StargateMsg{
					TypeURL: "/cosmos.bank.v1beta1.MsgSend",
					Value:   bankMsgBin,
				},
			},
			output: bankMsg,
		},
	}

	parser := types.NewWasmMsgParser()

	// test before setting MsgParser
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			_, err := parser.Parse(ctx, tc.sender, tc.input)
			if tc.isError {
				require.Error(t, err)
			}
		})
	}

	parser.Parsers[types.WasmMsgParserRouteBank] = bankwasm.NewWasmMsgParser()
	parser.Parsers[types.WasmMsgParserRouteStaking] = stakingwasm.NewWasmMsgParser()
	parser.Parsers[types.WasmMsgParserRouteMarket] = marketwasm.NewWasmMsgParser()
	parser.Parsers[types.WasmMsgParserRouteDistribution] = distrwasm.NewWasmMsgParser()
	parser.Parsers[types.WasmMsgParserRouteGov] = govwasm.NewWasmMsgParser()
	parser.Parsers[types.WasmMsgParserRouteWasm] = test_util.NewWasmMsgParser()
	parser.StargateParser = test_util.NewStargateWasmMsgParser(encodingConfig.Marshaler)
	parser.IBCParser = test_util.NewIBCMsgParser(wasmtesting.MockIBCTransferKeeper{})

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := parser.Parse(ctx, tc.sender, tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}
