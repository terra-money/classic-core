package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	types1 "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/keeper/wasmtesting"
)

func getPortFn(ctx sdk.Context) string {
	return "wasm.transfer"
}

func TestIBCEcoding(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload ibc reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// upload ibc reflect send code
	ibcReflectCode, err := ioutil.ReadFile("./testdata/ibc_reflect.wasm")
	require.NoError(t, err)
	ibcReflectID, err := keeper.StoreCode(ctx, creator, ibcReflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), ibcReflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// now we set contract as verifier of an ibc-reflect
	initMsg := IbcReflectExampleInitMsg{
		ReflectCodeID: reflectID,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	ibcReflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000))
	ibcReflectAddr, _, err := keeper.InstantiateContract(ctx, ibcReflectID, creator, sdk.AccAddress{}, initMsgBz, ibcReflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, ibcReflectAddr)

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

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmvmtypes.CosmosMsg
		// set if valid
		output sdk.Msg
		// set if invalid
		isError bool
	}{
		"close channel": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				IBC: &wasmvmtypes.IBCMsg{
					CloseChannel: &wasmvmtypes.CloseChannelMsg{ChannelID: "testchannel-1"},
				},
			},
			/* it always fails because there's no opened channel
			   output: &channeltypes.MsgChannelCloseInit{
			       PortId:    "wasm.transfer",
			       ChannelId: "testchannel-1",
			       Signer:    string(addrs[0]),
			   },
			*/
			isError: true, // can't open channel in test
		},
		"transfer": {
			sender: addrs[0],
			input: wasmvmtypes.CosmosMsg{
				IBC: &wasmvmtypes.IBCMsg{
					Transfer: &wasmvmtypes.TransferMsg{
						ChannelID: "testchannel-1",
						ToAddress: "terra1x46rqay4d3cssq8gxxvqz8xt6nwlz4td20k38",
						Amount:    wasmvmtypes.Coin{Denom: "uluna", Amount: "40000"}, //sdk.NewInt64Coin(core.MicroLunaDenom, 40000),
						Timeout:   wasmvmtypes.IBCTimeout{Block: &wasmvmtypes.IBCTimeoutBlock{Revision: 1, Height: 1}, Timestamp: 100000000},
					},
				},
			},
			output: &ibctransfertypes.MsgTransfer{
				SourcePort:       "wasm.transfer",
				SourceChannel:    "testchannel-1",
				Token:            sdk.NewCoin("uluna", sdk.NewInt(40000)),
				Sender:           (ibcReflectAddr.String()),
				Receiver:         `terra1x46rqay4d3cssq8gxxvqz8xt6nwlz4td20k38`,
				TimeoutHeight:    types1.Height{RevisionNumber: 1, RevisionHeight: 1},
				TimeoutTimestamp: 100000000,
			},
		},
	}

	parser := NewIBCMsgParser(wasmtesting.MockIBCTransferKeeper{
		GetPortFn: getPortFn,
	})
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := parser.Parse(ctx, ibcReflectAddr, tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}

func TestIBCQuery(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper, ibcKeeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper, input.IBCKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload ibc reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// upload ibc reflect send code
	ibcReflectCode, err := ioutil.ReadFile("./testdata/ibc_reflect.wasm")
	require.NoError(t, err)
	ibcReflectID, err := keeper.StoreCode(ctx, creator, ibcReflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), ibcReflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// now we set contract as verifier of an ibc-reflect
	initMsg := IbcReflectExampleInitMsg{
		ReflectCodeID: reflectID,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	ibcReflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000))
	ibcReflectAddr, _, err := keeper.InstantiateContract(ctx, ibcReflectID, creator, sdk.AccAddress{}, initMsgBz, ibcReflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, ibcReflectAddr)

	contractInfo, err := keeper.GetContractInfo(ctx, ibcReflectAddr)
	require.NoError(t, err)

	querier := NewIBCQuerier(keeper, ibcKeeper.ChannelKeeper)

	requestPort := wasmvmtypes.QueryRequest{
		IBC: &wasmvmtypes.IBCQuery{
			PortID: &wasmvmtypes.PortIDQuery{},
		},
	}
	responsePortBz, err := querier.Query(ctx, ibcReflectAddr, requestPort)
	require.NoError(t, err)
	responsePort := wasmvmtypes.PortIDResponse{}
	err = json.Unmarshal(responsePortBz, &responsePort)
	require.NoError(t, err)
	require.Equal(t, responsePort.PortID, contractInfo.GetIBCPortID())

	requestChannel := wasmvmtypes.QueryRequest{
		IBC: &wasmvmtypes.IBCQuery{
			Channel: &wasmvmtypes.ChannelQuery{
				PortID:    contractInfo.GetIBCPortID(),
				ChannelID: "testchannel-1",
			},
		},
	}
	responseChannelBz, err := querier.Query(ctx, ibcReflectAddr, requestChannel)
	require.NoError(t, err)
	responseChannel := wasmvmtypes.ChannelResponse{}
	err = json.Unmarshal(responseChannelBz, &responseChannel)
	require.NoError(t, err)

	requestListChannels := wasmvmtypes.QueryRequest{
		IBC: &wasmvmtypes.IBCQuery{
			ListChannels: &wasmvmtypes.ListChannelsQuery{
				PortID: contractInfo.GetIBCPortID(),
			},
		},
	}
	responseListChannelsBz, err := querier.Query(ctx, ibcReflectAddr, requestListChannels)
	require.NoError(t, err)
	responseListChannels := wasmvmtypes.ListChannelsResponse{}
	err = json.Unmarshal(responseListChannelsBz, &responseListChannels)
	require.NoError(t, err)

}
