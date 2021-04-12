package keeper

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/wasm/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// ReflectInitMsg is {}

// ReflectHandleMsg is used to encode handle messages
type ReflectHandleMsg struct {
	Reflect        *reflectPayload    `json:"reflect_msg,omitempty"`
	ReflectSubCall *reflectSubPayload `json:"reflect_sub_call,omitempty"`
	Change         *ownerPayload      `json:"change_owner,omitempty"`
}

type ownerPayload struct {
	Owner sdk.Address `json:"owner"`
}

type reflectPayload struct {
	Msgs []wasmvmtypes.CosmosMsg `json:"msgs"`
}

type reflectSubPayload struct {
	Msgs []wasmvmtypes.SubMsg `json:"msgs"`
}

// ReflectQueryMsg is used to encode query messages
type ReflectQueryMsg struct {
	Owner         *struct{}   `json:"owner,omitempty"`
	Capitalized   *Text       `json:"capitalized,omitempty"`
	Chain         *ChainQuery `json:"chain,omitempty"`
	SubCallResult *SubCall    `json:"sub_call_result,omitempty"`
}

type ChainQuery struct {
	Request *wasmvmtypes.QueryRequest `json:"request,omitempty"`
}

type Text struct {
	Text string `json:"text"`
}

type SubCall struct {
	ID uint64 `json:"id"`
}

type OwnerResponse struct {
	Owner string `json:"owner,omitempty"`
}

type ChainResponse struct {
	Data []byte `json:"data,omitempty"`
}

func buildReflectQuery(t *testing.T, query *ReflectQueryMsg) []byte {
	bz, err := json.Marshal(query)
	require.NoError(t, err)
	return bz
}

func mustParse(t *testing.T, data []byte, res interface{}) {
	err := json.Unmarshal(data, res)
	require.NoError(t, err)
}

const ReflectFeatures = "staking,reflect,stargate"

func TestReflectReflectContractSend(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	_, _, bob := keyPubAddr()

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// upload hackatom escrow code
	escrowCode, err := ioutil.ReadFile("./testdata/hackatom.wasm")
	require.NoError(t, err)
	escrowID, err := keeper.StoreCode(ctx, creator, escrowCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), escrowID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, []byte("{}"), reflectStart, true)
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
	escrowAddr, _, err := keeper.InstantiateContract(ctx, escrowID, creator, initMsgBz, escrowStart, true)
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
				Send: []wasmvmtypes.Coin{{
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

func TestReflectReflectCustomMsg(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		"reflect": reflectQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		"reflect": reflectRawMsgParser{
			cdc: MakeTestCodec(t),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	bob := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)
	_, _, fred := keyPubAddr()

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	contractAddr, _, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart, true)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// set owner to bob
	transfer := ReflectHandleMsg{
		Change: &ownerPayload{
			Owner: bob,
		},
	}
	transferBz, err := json.Marshal(transfer)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, creator, transferBz, nil)
	require.NoError(t, err)

	// check some account values
	checkAccount(t, ctx, accKeeper, bankKeeper, contractAddr, contractStart)
	checkAccount(t, ctx, accKeeper, bankKeeper, bob, deposit)
	checkAccount(t, ctx, accKeeper, bankKeeper, fred, nil)

	// bob can send contract's tokens to fred (using SendMsg)
	msgs := []wasmvmtypes.CosmosMsg{{
		Bank: &wasmvmtypes.BankMsg{
			Send: &wasmvmtypes.SendMsg{
				ToAddress: fred.String(),
				Amount: []wasmvmtypes.Coin{{
					Denom:  core.MicroLunaDenom,
					Amount: "15000",
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
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, reflectSendBz, nil)
	require.NoError(t, err)

	// fred got coins
	checkAccount(t, ctx, accKeeper, bankKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 15000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, bankKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000)))
	checkAccount(t, ctx, accKeeper, bankKeeper, bob, deposit)

	// construct an opaque message
	var sdkSendMsg sdk.Msg = &banktypes.MsgSend{
		FromAddress: contractAddr.String(),
		ToAddress:   fred.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 23000)),
	}

	opaque, err := toMaskRawMsg(MakeTestCodec(t), sdkSendMsg)
	require.NoError(t, err)
	reflectOpaque := ReflectHandleMsg{
		Reflect: &reflectPayload{
			Msgs: []wasmvmtypes.CosmosMsg{opaque},
		},
	}
	reflectOpaqueBz, err := json.Marshal(reflectOpaque)
	require.NoError(t, err)

	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, reflectOpaqueBz, nil)
	require.NoError(t, err)

	// fred got more coins
	checkAccount(t, ctx, accKeeper, bankKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 38000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, bankKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 2000)))
	checkAccount(t, ctx, accKeeper, bankKeeper, bob, deposit)
}

func TestReflectReflectCustomQuery(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, bankKeeper, keeper := input.Ctx, input.AccKeeper, input.BankKeeper, input.WasmKeeper

	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		"reflect": reflectQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		"reflect": reflectRawMsgParser{
			cdc: MakeTestCodec(t),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	contractAddr, _, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart, true)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// let's perform a normal query of state
	ownerQuery := ReflectQueryMsg{
		Owner: &struct{}{},
	}
	ownerQueryBz, err := json.Marshal(ownerQuery)
	require.NoError(t, err)
	ownerRes, err := keeper.queryToContract(ctx, contractAddr, ownerQueryBz)
	require.NoError(t, err)
	var res OwnerResponse
	err = json.Unmarshal(ownerRes, &res)
	require.NoError(t, err)
	assert.Equal(t, res.Owner, creator.String())

	// and now making use of the custom querier callbacks
	customQuery := ReflectQueryMsg{
		Capitalized: &Text{
			Text: "all Caps noW",
		},
	}
	customQueryBz, err := json.Marshal(customQuery)
	require.NoError(t, err)
	custom, err := keeper.queryToContract(ctx, contractAddr, customQueryBz)
	require.NoError(t, err)
	var resp capitalizedResponse
	err = json.Unmarshal(custom, &resp)
	require.NoError(t, err)
	assert.Equal(t, resp.Text, "ALL CAPS NOW")
}

func TestReflectStargateQuery(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	funds := sdk.NewCoins(sdk.NewInt64Coin("denom", 320000))
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	expectedBalance := funds.Sub(contractStart)
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, funds)

	// upload code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractAddr, _, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart, false)
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

type reflectState struct {
	Owner string `json:"owner"`
}

func TestMaskReflectWasmQueries(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, []byte("{}"), reflectStart, true)
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

func TestWasmRawQueryWithNil(t *testing.T) {
	input := CreateTestInput(t)
	ctx, accKeeper, keeper, bankKeeper := input.Ctx, input.AccKeeper, input.WasmKeeper, input.BankKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, bankKeeper, deposit)

	// upload reflect code
	reflectCode, err := ioutil.ReadFile("./testdata/reflect.wasm")
	require.NoError(t, err)
	reflectID, err := keeper.StoreCode(ctx, creator, reflectCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), reflectID)

	// creator instantiates a contract and gives it tokens
	reflectStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	reflectAddr, _, err := keeper.InstantiateContract(ctx, reflectID, creator, []byte("{}"), reflectStart, true)
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

func checkAccount(t *testing.T, ctx sdk.Context, accKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, addr sdk.AccAddress, expected sdk.Coins) {
	acct := accKeeper.GetAccount(ctx, addr)
	if expected == nil {
		assert.Nil(t, acct)
	} else {
		assert.NotNil(t, acct)

		balance := bankKeeper.GetAllBalances(ctx, addr)
		if expected.Empty() {
			// there is confusion between nil and empty slice... let's just treat them the same
			assert.True(t, balance.Empty())
		} else {
			assert.Equal(t, expected, balance)
		}
	}
}

/**** Code to support custom messages *****/
type reflectCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toMaskRawMsg encodes an sdk msg using amino json encoding.
// Then wraps it as an opaque message
func toMaskRawMsg(cdc codec.Marshaler, msg sdk.Msg) (wasmvmtypes.CosmosMsg, error) {
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, err
	}
	rawBz, err := cdc.MarshalJSON(any)
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	customMsgData, err := json.Marshal(reflectCustomMsg{Raw: rawBz})
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	customMsg, err := json.Marshal(types.WasmCustomMsg{
		Route:   "reflect",
		MsgData: customMsgData,
	})
	if err != nil {
		return wasmvmtypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	res := wasmvmtypes.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}

// reflectRawMsgParser decodes msg.Data to an sdk.Msg using amino json encoding.
// this needs to be registered on the Encoders
type reflectRawMsgParser struct {
	cdc codec.Marshaler
}

var _ types.WasmMsgParserInterface = reflectRawMsgParser{}

func (p reflectRawMsgParser) Parse(_sender sdk.AccAddress, msg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	return nil, nil
}

func (p reflectRawMsgParser) ParseCustom(_sender sdk.AccAddress, msg json.RawMessage) (sdk.Msg, error) {
	var custom reflectCustomMsg
	err := json.Unmarshal(msg, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Raw != nil {
		var any codectypes.Any
		if err := p.cdc.UnmarshalJSON(custom.Raw, &any); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		var msg sdk.Msg
		if err := p.cdc.UnpackAny(&any, &msg); err != nil {
			return nil, err
		}
		return msg, nil
	}
	if custom.Debug != "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Custom Debug: %s", custom.Debug)
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom message variant")
}

type reflectQuerier struct{}

type reflectCustomQuery struct {
	Ping        *struct{} `json:"ping,omitempty"`
	Capitalized *Text     `json:"capitalized,omitempty"`
}

var _ types.WasmQuerierInterface = reflectQuerier{}

// this is from the go code back to the contract (capitalized or ping)
type customQueryResponse struct {
	Msg string `json:"msg"`
}

// these are the return values from contract -> go depending on type of query
type ownerResponse struct {
	Owner string `json:"owner"`
}

type capitalizedResponse struct {
	Text string `json:"text"`
}

type chainResponse struct {
	Data []byte `json:"data"`
}

func (reflectQuerier) Query(_ sdk.Context, _ wasmvmtypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

func (reflectQuerier) QueryCustom(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom reflectCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Capitalized != nil {
		msg := strings.ToUpper(custom.Capitalized.Text)
		return json.Marshal(customQueryResponse{Msg: msg})
	}
	if custom.Ping != nil {
		return json.Marshal(customQueryResponse{Msg: "pong"})
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
}
