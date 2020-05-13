package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/terra-project/core/x/bank"
	"github.com/terra-project/core/x/wasm/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// MaskInitMsg is {}

type MaskHandleMsg struct {
	Reflect *reflectPayload `json:"reflect_msg,omitempty"`
	Change  *ownerPayload   `json:"change_owner,omitempty"`
}

type ownerPayload struct {
	Owner sdk.Address `json:"owner"`
}

type reflectPayload struct {
	Msgs []wasmTypes.CosmosMsg `json:"msgs"`
}

// MaskQueryMsg is used to encode query messages
type MaskQueryMsg struct {
	Owner         *struct{} `json:"owner,omitempty"`
	ReflectCustom *Text     `json:"reflect_custom,omitempty"`
}

type Text struct {
	Text string `json:"text"`
}

type OwnerResponse struct {
	Owner string `json:"owner,omitempty"`
}

func TestMaskReflectContractSend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	ctx, accKeeper, keeper := CreateTestInput(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)
	_, _, bob := keyPubAddr()

	// upload mask code
	maskCode, err := ioutil.ReadFile("./testdata/mask.wasm")
	require.NoError(t, err)
	maskID, err := keeper.StoreCode(ctx, creator, maskCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), maskID)

	// upload hackatom escrow code
	escrowCode, err := ioutil.ReadFile("./testdata/contract.wasm")
	require.NoError(t, err)
	escrowID, err := keeper.StoreCode(ctx, creator, escrowCode)
	require.NoError(t, err)
	require.Equal(t, uint64(2), escrowID)

	// creator instantiates a contract and gives it tokens
	maskStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	maskAddr, err := keeper.InstantiateContract(ctx, maskID, creator, []byte("{}"), maskStart)
	require.NoError(t, err)
	require.NotEmpty(t, maskAddr)

	// now we set contract as verifier of an escrow
	initMsg := InitMsg{
		Verifier:    maskAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	escrowStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 25000))
	escrowAddr, err := keeper.InstantiateContract(ctx, escrowID, creator, initMsgBz, escrowStart)
	require.NoError(t, err)
	require.NotEmpty(t, escrowAddr)

	// let's make sure all balances make sense
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 35000))) // 100k - 40k - 25k
	checkAccount(t, ctx, accKeeper, maskAddr, maskStart)
	checkAccount(t, ctx, accKeeper, escrowAddr, escrowStart)
	checkAccount(t, ctx, accKeeper, bob, nil)

	// now for the trick.... we reflect a message through the mask to call the escrow
	// we also send an additional 14k tokens there.
	// this should reduce the mask balance by 14k (to 26k)
	// this 14k is added to the escrow, then the entire balance is sent to bob (total: 39k)
	approveMsg := []byte(`{"release":{}}`)
	msgs := []wasmTypes.CosmosMsg{{
		Wasm: &wasmTypes.WasmMsg{
			Execute: &wasmTypes.ExecuteMsg{
				ContractAddr: escrowAddr.String(),
				Msg:          approveMsg,
				Send: []wasmTypes.Coin{{
					Denom:  "denom",
					Amount: "14000",
				}},
			},
		},
	}}
	reflectSend := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: msgs,
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, maskAddr, creator, nil, reflectSendBz)
	require.NoError(t, err)

	// did this work???
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 35000)))  // same as before
	checkAccount(t, ctx, accKeeper, maskAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 26000))) // 40k - 14k (from send)
	checkAccount(t, ctx, accKeeper, escrowAddr, sdk.Coins{})                                  // emptied reserved
	checkAccount(t, ctx, accKeeper, bob, sdk.NewCoins(sdk.NewInt64Coin("denom", 39000)))      // all escrow of 25k + 14k

}

func TestMaskReflectCustomMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, accKeeper, keeper := CreateTestInput(t)
	keeper.RegisterQueriers(ctx, map[string]types.WasmQuerier{
		"mask": maskQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParser{
		"mask": maskRawMsgParser{
			cdc: makeTestCodec(),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)
	bob := createFakeFundedAccount(ctx, accKeeper, deposit)
	_, _, fred := keyPubAddr()

	// upload code
	maskCode, err := ioutil.ReadFile("./testdata/mask.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, maskCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	contractAddr, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// set owner to bob
	transfer := MaskHandleMsg{
		Change: &ownerPayload{
			Owner: bob,
		},
	}
	transferBz, err := json.Marshal(transfer)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, creator, nil, transferBz)
	require.NoError(t, err)

	// check some account values
	checkAccount(t, ctx, accKeeper, contractAddr, contractStart)
	checkAccount(t, ctx, accKeeper, bob, deposit)
	checkAccount(t, ctx, accKeeper, fred, nil)

	// bob can send contract's tokens to fred (using SendMsg)
	msgs := []wasmTypes.CosmosMsg{{
		Bank: &wasmTypes.BankMsg{
			Send: &wasmTypes.SendMsg{
				FromAddress: contractAddr.String(),
				ToAddress:   fred.String(),
				Amount: []wasmTypes.Coin{{
					Denom:  "denom",
					Amount: "15000",
				}},
			},
		},
	}}
	reflectSend := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: msgs,
		},
	}
	reflectSendBz, err := json.Marshal(reflectSend)
	require.NoError(t, err)
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, nil, reflectSendBz)
	require.NoError(t, err)

	// fred got coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin("denom", 15000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 25000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)

	// construct an opaque message
	var sdkSendMsg sdk.Msg = &bank.MsgSend{
		FromAddress: contractAddr,
		ToAddress:   fred,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("denom", 23000)),
	}
	opaque, err := toMaskRawMsg(keeper.cdc, sdkSendMsg)
	require.NoError(t, err)
	reflectOpaque := MaskHandleMsg{
		Reflect: &reflectPayload{
			Msgs: []wasmTypes.CosmosMsg{opaque},
		},
	}
	reflectOpaqueBz, err := json.Marshal(reflectOpaque)
	require.NoError(t, err)

	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, nil, reflectOpaqueBz)
	require.NoError(t, err)

	// fred got more coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin("denom", 38000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin("denom", 2000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)
}

func TestMaskReflectCustomQuery(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	ctx, accKeeper, keeper := CreateTestInput(t)
	keeper.RegisterQueriers(ctx, map[string]types.WasmQuerier{
		"mask": maskQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParser{
		"mask": maskRawMsgParser{
			cdc: makeTestCodec(),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	// upload code
	maskCode, err := ioutil.ReadFile("./testdata/mask.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, maskCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin("denom", 40000))
	contractAddr, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart)
	require.NoError(t, err)
	require.NotEmpty(t, contractAddr)

	// let's perform a normal query of state
	ownerQuery := MaskQueryMsg{
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
	customQuery := MaskQueryMsg{
		ReflectCustom: &Text{
			Text: "all Caps noW",
		},
	}
	customQueryBz, err := json.Marshal(customQuery)
	require.NoError(t, err)
	custom, err := keeper.queryToContract(ctx, contractAddr, customQueryBz)
	require.NoError(t, err)
	var resp customQueryResponse
	err = json.Unmarshal(custom, &resp)
	require.NoError(t, err)
	assert.Equal(t, resp.Msg, "ALL CAPS NOW")
}

func checkAccount(t *testing.T, ctx sdk.Context, accKeeper auth.AccountKeeper, addr sdk.AccAddress, expected sdk.Coins) {
	acct := accKeeper.GetAccount(ctx, addr)
	if expected == nil {
		assert.Nil(t, acct)
	} else {
		assert.NotNil(t, acct)
		if expected.Empty() {
			// there is confusion between nil and empty slice... let's just treat them the same
			assert.True(t, acct.GetCoins().Empty())
		} else {
			assert.Equal(t, acct.GetCoins(), expected)
		}
	}
}

/**** Code to support custom messages *****/

type maskCustomMsg struct {
	Debug string `json:"debug,omitempty"`
	Raw   []byte `json:"raw,omitempty"`
}

// toMaskRawMsg encodes an sdk msg using amino json encoding.
// Then wraps it as an opaque message
func toMaskRawMsg(cdc *codec.Codec, msg sdk.Msg) (wasmTypes.CosmosMsg, error) {
	rawBz, err := cdc.MarshalJSON(msg)
	if err != nil {
		return wasmTypes.CosmosMsg{}, sdk.ErrInternal(err.Error())
	}
	customMsg, err := json.Marshal(maskCustomMsg{
		Raw: rawBz,
	})
	res := wasmTypes.CosmosMsg{
		Custom: customMsg,
	}
	return res, nil
}

// maskRawMsgParser decodes msg.Data to an sdk.Msg using amino json encoding.
// this needs to be registered on the Encoders
type maskRawMsgParser struct {
	cdc *codec.Codec
}

var _ types.WasmMsgParser = maskRawMsgParser{}

func (p maskRawMsgParser) Parse(_sender sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, sdk.Error) {
	return nil, nil
}

func (p maskRawMsgParser) ParseCustom(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, sdk.Error) {
	var custom maskCustomMsg
	err := json.Unmarshal(msg, &custom)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	if custom.Raw != nil {
		var sdkMsg sdk.Msg
		err := p.cdc.UnmarshalJSON(custom.Raw, &sdkMsg)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if custom.Debug != "" {
		return nil, types.ErrInvalidMsg(fmt.Sprintf("Custom Debug: %s", custom.Debug))
	}
	return nil, types.ErrInvalidMsg("Unknown Custom message variant")
}

type maskQuerier struct{}

type maskCustomQuery struct {
	Ping    *struct{} `json:"ping,omitempty"`
	Capital *Text     `json:"capital,omitempty"`
}

var _ types.WasmQuerier = maskQuerier{}

type customQueryResponse struct {
	Msg string `json:"msg"`
}

func (maskQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, sdk.Error) {
	return nil, nil
}

func (maskQuerier) QueryCustom(_ sdk.Context, request json.RawMessage) ([]byte, sdk.Error) {
	var custom maskCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	if custom.Capital != nil {
		msg := strings.ToUpper(custom.Capital.Text)
		bz, err := json.Marshal(customQueryResponse{Msg: msg})
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		return bz, nil
	}
	if custom.Ping != nil {
		bz, err := json.Marshal(customQueryResponse{Msg: "pong"})
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		return bz, nil
	}

	return nil, types.ErrInvalidMsg("Unknown Custom query variant")
}
