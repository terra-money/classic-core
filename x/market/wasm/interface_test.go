package wasm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market/internal/keeper"
	"github.com/terra-project/core/x/market/internal/types"
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
		"simple swap": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap": {"trader": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
						addrs[0], core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			output: []sdk.Msg{
				types.MsgSwap{
					Trader:    addrs[0],
					OfferCoin: sdk.NewInt64Coin(core.MicroLunaDenom, 1234),
					AskDenom:  core.MicroSDRDenom,
				},
			},
		},
		"simple swap send": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap_send": {"from_address": "%s", "to_address": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
						addrs[0], addrs[1], core.MicroLunaDenom, core.MicroSDRDenom,
					),
				),
			},
			output: []sdk.Msg{
				types.MsgSwapSend{
					FromAddress: addrs[0],
					ToAddress:   addrs[1],
					OfferCoin:   sdk.NewInt64Coin(core.MicroLunaDenom, 1234),
					AskDenom:    core.MicroSDRDenom,
				},
			},
		},
		"invalid swap amount": {
			sender: addrs[0],
			input: wasmTypes.CosmosMsg{
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
			input: wasmTypes.CosmosMsg{
				Custom: []byte(
					fmt.Sprintf(
						`{"swap": {"trader": "%s", "offer_coin": {"amount": "1234", "denom": "%s"}, "ask_denom": "%s"}}`,
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
		Swap: queryParams,
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// overflow query
	overflowAmt, _ := sdk.NewIntFromString("1000000000000000000000000000000000")
	overflowOfferCoin := sdk.NewCoin(core.MicroLunaDenom, overflowAmt)
	queryParams = types.NewQuerySwapParams(overflowOfferCoin, core.MicroSDRDenom)
	bz, err = json.Marshal(CosmosQuery{
		Swap: queryParams,
	})
	require.NoError(t, err)

	_, err = querier.QueryCustom(input.Ctx, bz)
	require.Error(t, err)

	// valid query
	queryParams = types.NewQuerySwapParams(offerCoin, core.MicroSDRDenom)
	bz, err = json.Marshal(CosmosQuery{
		Swap: queryParams,
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
