package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestLegacyReflectReflectContractSend(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	_, _, bob := keyPubAddr()

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect_legacy.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// upload hackatom escrow code
	escrowCode, err := ioutil.ReadFile("./testdata/hackatom_legacy.wasm")
	require.NoError(t, err)
	escrowID, err := keeper.StoreCode(ctx, creator, escrowCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), escrowID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// now we set contract as verifier of an escrow
	initMsg := HackatomExampleInitMsg{
		Verifier:    reflectAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	escrowStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000))
	escrowAddr, _, err := keeper.InstantiateContract(ctx, escrowID, creator, sdk.AccAddress{}, initMsgBz, escrowStart)
	require.NoError(t, err)
	require.NotEmpty(t, escrowAddr)

	// let's make sure all balances make sense
	checkAccount(t, ctx, accKeeper, bankKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 35000))) // 100k - 40k - 25k
	checkAccount(t, ctx, accKeeper, bankKeeper, reflectAddr, reflectStart)
	checkAccount(t, ctx, accKeeper, bankKeeper, escrowAddr, escrowStart)
	checkAccount(t, ctx, accKeeper, bankKeeper, bob, nil)

	// now for the trick.... we reflect a message through the reflect to call the escrow
	// we also send an additional 14k tokens there.
	// this should reduce the reflect balance by 14k (to 26k)
	// this 14k is added to the escrow, then the entire balance is sent to bob (total: 39k)
	approveMsg := []byte(`{"release":{}}`)
	msgs := []wasmvmtypes.CosmosMsg{{
		Wasm: &wasmvmtypes.WasmMsg{
			Execute: &wasmvmtypes.ExecuteMsg{
				ContractAddr: escrowAddr.String(),
				Msg:          approveMsg,
				Funds: []wasmvmtypes.Coin{{
					Denom:  core.MicroLunaDenom,
					Amount: "14000",
				}},
			},
		},
	}}
	reflectSend := ReflectHandleMsg{
		Reflect: &reflectPayload{
			Msgs: msgs,
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, reflectAddr, creator, reflectSendBz, nil)
	require.NoError(t, err)

	// did this work???
	checkAccount(t, ctx, accKeeper, bankKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 35000)))     // same as before
	checkAccount(t, ctx, accKeeper, bankKeeper, reflectAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 26000))) // 40k - 14k (from send)
	checkAccount(t, ctx, accKeeper, bankKeeper, escrowAddr, sdk.Coins{})                                                 // emptied reserved
	checkAccount(t, ctx, accKeeper, bankKeeper, bob, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 39000)))         // all escrow of 25k + 14k

}

func TestLegacyReflectStargateQuery(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	funds := sdk.NewCoins(sdk.NewInt64Coin("denom", 320000))
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	expectedBalance := funds.Sub(contractStart)
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, funds)

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect_legacy.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractAddr, _, err := keeper.InstantiateContract(ctx, codeID, creator, sdk.AccAddress{}, []byte("{}"), contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// first, normal query for the bank balance (to make sure our query is proper)
	bankQuery := wasmvmtypes.QueryRequest{
		Bank: &wasmvmtypes.BankQuery{
			AllBalances: &wasmvmtypes.AllBalancesQuery{
				Address: creator.String(),
			},
		},
	}
	simpleQueryBz, err := json.Marshal(ReflectQueryMsg{
		Chain: &ChainQuery{Request: &bankQuery},
	})
	require.NoError(t, err)
	simpleRes, err := keeper.queryToContract(ctx, contractAddr, simpleQueryBz)
	require.NoError(t, err)
	var simpleChain ChainResponse
	mustParse(t, simpleRes, &simpleChain)
	var simpleBalance wasmvmtypes.AllBalancesResponse
	mustParse(t, simpleChain.Data, &simpleBalance)
	require.Equal(t, len(expectedBalance), len(simpleBalance.Amount))
	assert.Equal(t, expectedBalance[0].Amount.String(), simpleBalance.Amount[0].Amount)
	assert.Equal(t, expectedBalance[0].Denom, simpleBalance.Amount[0].Denom)

	// now, try to build a protobuf query
	protoQuery := banktypes.QueryAllBalancesRequest{
		Address: creator.String(),
	}
	protoQueryBin, err := proto.Marshal(&protoQuery)
	protoRequest := wasmvmtypes.QueryRequest{
		Stargate: &wasmvmtypes.StargateQuery{
			Path: "/cosmos.bank.v1beta1.Query/AllBalances",
			Data: protoQueryBin,
		},
	}
	protoQueryBz, err := json.Marshal(ReflectQueryMsg{
		Chain: &ChainQuery{Request: &protoRequest},
	})
	require.NoError(t, err)

	// make a query on the chain
	protoRes, err := keeper.queryToContract(ctx, contractAddr, protoQueryBz)
	require.NoError(t, err)
	var protoChain ChainResponse
	mustParse(t, protoRes, &protoChain)

	// unmarshal raw protobuf response
	var protoResult banktypes.QueryAllBalancesResponse
	err = proto.Unmarshal(protoChain.Data, &protoResult)
	require.NoError(t, err)
	assert.Equal(t, expectedBalance, protoResult.Balances)
}

func TestLegacyMaskReflectWasmQueries(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect_legacy.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// for control, let's make some queries directly on the reflect
	ownerQuery := buildReflectQuery(t, &ReflectQueryMsg{Owner: &struct{}{}})
	res, err := keeper.queryToContract(ctx, reflectAddr, ownerQuery)
	require.NoError(t, err)
	var ownerRes OwnerResponse
	mustParse(t, res, &ownerRes)
	require.Equal(t, ownerRes.Owner, creator.String())

	// and a raw query: cosmwasm_storage::Singleton uses 2 byte big-endian length-prefixed to store data
	configKey := append([]byte{0, 6}, []byte("config")...)
	raw := keeper.queryToStore(ctx, reflectAddr, configKey)
	var stateRes reflectState
	mustParse(t, raw, &stateRes)
	require.Equal(t, stateRes.Owner, creator.String())

	// now, let's reflect a smart query into the x/wasm handlers and see if we get the same result
	reflectOwnerQuery := ReflectQueryMsg{Chain: &ChainQuery{Request: &wasmvmtypes.QueryRequest{Wasm: &wasmvmtypes.WasmQuery{
		Smart: &wasmvmtypes.SmartQuery{
			ContractAddr: reflectAddr.String(),
			Msg:          ownerQuery,
		},
	}}}}
	reflectOwnerBin := buildReflectQuery(t, &reflectOwnerQuery)
	res, err = keeper.queryToContract(ctx, reflectAddr, reflectOwnerBin)
	require.NoError(t, err)
	// first we pull out the data from chain response, before parsing the original response
	var reflectRes ChainResponse
	mustParse(t, res, &reflectRes)
	var reflectOwnerRes OwnerResponse
	mustParse(t, reflectRes.Data, &reflectOwnerRes)
	require.Equal(t, reflectOwnerRes.Owner, creator.String())

	// and with queryRaw
	reflectStateQuery := ReflectQueryMsg{Chain: &ChainQuery{Request: &wasmvmtypes.QueryRequest{Wasm: &wasmvmtypes.WasmQuery{
		Raw: &wasmvmtypes.RawQuery{
			ContractAddr: reflectAddr.String(),
			Key:          configKey,
		},
	}}}}
	reflectStateBin := buildReflectQuery(t, &reflectStateQuery)
	res, err = keeper.queryToContract(ctx, reflectAddr, reflectStateBin)
	require.NoError(t, err)
	// first we pull out the data from chain response, before parsing the original response
	var reflectRawRes ChainResponse
	mustParse(t, res, &reflectRawRes)

	// now, with the raw data, we can parse it into state
	var reflectStateRes reflectState
	mustParse(t, reflectRawRes.Data, &reflectStateRes)
	require.Equal(t, reflectStateRes.Owner, creator.String())
}

func TestLegacyWasmRawQueryWithNil(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect_legacy.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, sdk.AccAddress{}, []byte("{}"), reflectStart)
	require.NoError(t, err)
	require.NotEmpty(t, reflectAddr)

	// control: query directly
	missingKey := []byte{0, 1, 2, 3, 4}
	raw := keeper.queryToStore(ctx, reflectAddr, missingKey)
	require.Nil(t, raw)

	// and with queryRaw
	reflectQuery := ReflectQueryMsg{Chain: &ChainQuery{Request: &wasmvmtypes.QueryRequest{Wasm: &wasmvmtypes.WasmQuery{
		Raw: &wasmvmtypes.RawQuery{
			ContractAddr: reflectAddr.String(),
			Key:          missingKey,
		},
	}}}}
	reflectStateBin := buildReflectQuery(t, &reflectQuery)
	res, err := keeper.queryToContract(ctx, reflectAddr, reflectStateBin)
	require.NoError(t, err)

	// first we pull out the data from chain response, before parsing the original response
	var reflectRawRes ChainResponse
	mustParse(t, res, &reflectRawRes)
	// and make sure there is no data
	require.Empty(t, reflectRawRes.Data)
	// we get an empty byte slice not nil (if anyone care in go-land)
	require.Equal(t, []byte{}, reflectRawRes.Data)
}
