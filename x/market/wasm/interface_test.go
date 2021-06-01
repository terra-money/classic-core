package wasm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/market/keeper"
	"github.com/terra-money/core/x/market/types"
)

func TestEncoding(t *testing.T) {
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
		"simple swap": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap": {"trader": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
						addrs[0], core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			output: &types.MsgSwap{
				Trader:    addrs[0].String(),
				OfferCoin: sdk.NewInt64Coin(core.MicroLunaDenom, 1234),
				AskDenom:  core.MicroSDRDenom,
			},
		},
		"simple swap send": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap_send": {"from_address": "%s", "to_address": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
						addrs[0], addrs[1], core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			output: &types.MsgSwapSend{
				FromAddress: addrs[0].String(),
				ToAddress:   addrs[1].String(),
				OfferCoin:   sdk.NewInt64Coin(core.MicroLunaDenom, 1234),
				AskDenom:    core.MicroSDRDenom,
			},
		},
		"invalid swap amount": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap": {"trader": "%s", "offer_coin": {"amount": "1234.123", "denom": "%s"}, "ask_denom": "%s"}}`,
						addrs[0], core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			isError: true,
		},
		"invalid address": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap_send": {"to_address": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
						invalidAddr, core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			isError: true,
		},
	}

	parser := NewWasmMsgParser()
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := parser.ParseCustom(tc.sender, tc.input.Custom)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}

func TestQuerySwap(t *testing.T) {
	input := keeper.CreateTestInput(t)

	price := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, price)

	querier := NewWasmQuerier(input.MarketKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// recursive query
	offerCoin := sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(10))
	queryParams := types.NewQuerySwapParams(offerCoin, core.MicroLunaDenom)
	bz, err := json.Marshal(CosmosQuery{
		Swap: &queryParams,
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// overflow query
	overflowAmt, _ := sdk.NewIntFromString("1000000000000000000000000000000000")
	overflowOfferCoin := sdk.NewCoin(core.MicroLunaDenom, overflowAmt)
	queryParams = types.NewQuerySwapParams(overflowOfferCoin, core.MicroSDRDenom)
	bz, err = json.Marshal(CosmosQuery{
		Swap: &queryParams,
	})
	require.NoError(t, err)

	_, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query
	queryParams = types.NewQuerySwapParams(offerCoin, core.MicroSDRDenom)
	bz, err = json.Marshal(CosmosQuery{
		Swap: &queryParams,
	})
	require.NoError(t, err)

	res, err = querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var swapResponse SwapQueryResponse
	err = json.Unmarshal(res, &swapResponse)
	require.NoError(t, err)

	swapAmount, ok := sdk.NewIntFromString(swapResponse.Receive.Amount)
	require.True(t, ok)
	require.Equal(t, core.MicroSDRDenom, swapResponse.Receive.Denom)
	require.True(t, sdk.NewInt(17).GTE(swapAmount))
	require.True(t, swapAmount.IsPositive())
}
