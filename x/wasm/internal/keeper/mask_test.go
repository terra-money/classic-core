package keeper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/wasm/internal/types"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
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

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
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
	maskStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	maskAddr, err := keeper.InstantiateContract(ctx, maskID, creator, []byte("{}"), maskStart, true)
	require.NoError(t, err)
	require.NotEmpty(t, maskAddr)

	// now we set contract as verifier of an escrow
	initMsg := InitMsg{
		Verifier:    maskAddr,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)
	escrowStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000))
	escrowAddr, err := keeper.InstantiateContract(ctx, escrowID, creator, initMsgBz, escrowStart, true)
	require.NoError(t, err)
	require.NotEmpty(t, escrowAddr)

	// let's make sure all balances make sense
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 35000))) // 100k - 40k - 25k
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
					Denom:  core.MicroLunaDenom,
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
	_, err = keeper.ExecuteContract(ctx, maskAddr, creator, reflectSendBz, nil)
	require.NoError(t, err)

	// did this work???
	checkAccount(t, ctx, accKeeper, creator, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 35000)))  // same as before
	checkAccount(t, ctx, accKeeper, maskAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 26000))) // 40k - 14k (from send)
	checkAccount(t, ctx, accKeeper, escrowAddr, sdk.Coins{})                                              // emptied reserved
	checkAccount(t, ctx, accKeeper, bob, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 39000)))      // all escrow of 25k + 14k

}

func TestMaskReflectCustomMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		"mask": maskQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		"mask": maskRawMsgParser{
			cdc: makeTestCodec(),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
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
	contractStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	contractAddr, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart, true)
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
	_, err = keeper.ExecuteContract(ctx, contractAddr, creator, transferBz, nil)
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
					Denom:  core.MicroLunaDenom,
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
	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, reflectSendBz, nil)
	require.NoError(t, err)

	// fred got coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 15000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 25000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)

	// construct an opaque message
	var sdkSendMsg sdk.Msg = &bank.MsgSend{
		FromAddress: contractAddr,
		ToAddress:   fred,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 23000)),
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

	_, err = keeper.ExecuteContract(ctx, contractAddr, bob, reflectOpaqueBz, nil)
	require.NoError(t, err)

	// fred got more coins
	checkAccount(t, ctx, accKeeper, fred, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 38000)))
	// contract lost them
	checkAccount(t, ctx, accKeeper, contractAddr, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 2000)))
	checkAccount(t, ctx, accKeeper, bob, deposit)
}

func TestMaskReflectCustomQuery(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	input := CreateTestInput(t)
	ctx, accKeeper, keeper := input.Ctx, input.AccKeeper, input.WasmKeeper

	keeper.RegisterQueriers(map[string]types.WasmQuerierInterface{
		"mask": maskQuerier{},
	})
	keeper.RegisterMsgParsers(map[string]types.WasmMsgParserInterface{
		"mask": maskRawMsgParser{
			cdc: makeTestCodec(),
		},
	})

	deposit := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 100000))
	creator := createFakeFundedAccount(ctx, accKeeper, deposit)

	// upload code
	maskCode, err := ioutil.ReadFile("./testdata/mask.wasm")
	require.NoError(t, err)
	codeID, err := keeper.StoreCode(ctx, creator, maskCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), codeID)

	// creator instantiates a contract and gives it tokens
	contractStart := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 40000))
	contractAddr, err := keeper.InstantiateContract(ctx, codeID, creator, []byte("{}"), contractStart, true)
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
			assert.Equal(t, expected, acct.GetCoins())
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
		return wasmTypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	customMsgData, err := json.Marshal(maskCustomMsg{Raw: rawBz})
	if err != nil {
		return wasmTypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	customMsg, err := json.Marshal(types.WasmCustomMsg{
		Route:   "mask",
		MsgData: customMsgData,
	})
	if err != nil {
		return wasmTypes.CosmosMsg{}, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

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

var _ types.WasmMsgParserInterface = maskRawMsgParser{}

func (p maskRawMsgParser) Parse(_sender sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	return nil, nil
}

func (p maskRawMsgParser) ParseCustom(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	var custom maskCustomMsg
	err := json.Unmarshal(msg, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Raw != nil {
		var sdkMsg sdk.Msg
		err := p.cdc.UnmarshalJSON(custom.Raw, &sdkMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if custom.Debug != "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Custom Debug: %s", custom.Debug)
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom message variant")
}

type maskQuerier struct{}

type maskCustomQuery struct {
	Ping    *struct{} `json:"ping,omitempty"`
	Capital *Text     `json:"capital,omitempty"`
}

var _ types.WasmQuerierInterface = maskQuerier{}

type customQueryResponse struct {
	Msg string `json:"msg"`
}

func (maskQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) {
	return nil, nil
}

func (maskQuerier) QueryCustom(_ sdk.Context, request json.RawMessage) ([]byte, error) {
	var custom maskCustomQuery
	err := json.Unmarshal(request, &custom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if custom.Capital != nil {
		msg := strings.ToUpper(custom.Capital.Text)
		bz, err := json.Marshal(customQueryResponse{Msg: msg})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
	if custom.Ping != nil {
		bz, err := json.Marshal(customQueryResponse{Msg: "pong"})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}

	return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
}
