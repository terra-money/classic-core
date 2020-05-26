package keeper

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/wasm/internal/types"
)

// MakerInitMsg nolint
type MakerInitMsg struct {
	OfferDenom string `json:"offer"`
	AskDenom   string `json:"ask"`
}

// MakerHandleMsg nolint
type MakerHandleMsg struct {
	Buy  *buyPayload  `json:"buy,omitempty"`
	Sell *sellPayload `json:"sell,omitempty"`
}

// MakerQueryMsg nolint
type MakerQueryMsg struct {
	Simulate simulateQuery `json:"simulate"`
}

type simulateQuery struct {
	OfferCoin wasmTypes.Coin `json:"offer"`
}

type simulateResponse struct {
	SellCoin wasmTypes.Coin `json:"sell"`
	BuyCoin  wasmTypes.Coin `json:"buy"`
}

type buyPayload struct {
	Limit uint64 `json:"limit,omitempty"`
}

type sellPayload struct {
	Limit uint64 `json:"limit,omitempty"`
}

func TestInstantiateMaker(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input := CreateTestInput(t)

	ctx, keeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.OracleKeeper
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	_, _, creatorAddr := keyPubAddr()

	// upload staking derivates code
	makingCode, err := ioutil.ReadFile("./testdata/maker.wasm")
	require.NoError(t, err)
	makerID, err := keeper.StoreCode(ctx, creatorAddr, makingCode, true)
	require.NoError(t, err)
	require.Equal(t, uint64(1), makerID)

	// valid instantiate
	initMsg := MakerInitMsg{
		OfferDenom: core.MicroSDRDenom,
		AskDenom:   core.MicroLunaDenom,
	}

	initBz, err := json.Marshal(&initMsg)
	makerAddr, err := keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, initBz, nil)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	// invalid init msg
	_, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, []byte{}, nil)
	require.Error(t, err)
}

func TestQuerier(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, _, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper := input.Ctx, input.WasmKeeper

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	expectedRetCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()))

	// querier test
	swapQueryMsg := MakerQueryMsg{
		Simulate: simulateQuery{
			OfferCoin: types.EncodeSdkCoin(offerCoin),
		},
	}

	bz, err := json.Marshal(swapQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, makerAddr, bz)
	require.NoError(t, err)

	var simulRes simulateResponse
	err = json.Unmarshal(res, &simulRes)
	require.NoError(t, err)

	sellCoin, err := types.ParseToCoin(simulRes.SellCoin)
	require.NoError(t, err)
	require.Equal(t, offerCoin, sellCoin)

	buyCoin, err := types.ParseToCoin(simulRes.BuyCoin)
	require.NoError(t, err)
	require.Equal(t, expectedRetCoins[0], buyCoin)
}

func TestBuyMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	expectedRetCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()))

	// buy without limit
	buyMsg := MakerHandleMsg{
		Buy: &buyPayload{},
	}

	bz, err := json.Marshal(&buyMsg)

	// normal buy
	_, err = keeper.ExecuteContract(ctx, makerAddr, creatorAddr, bz, sdk.NewCoins(offerCoin))
	require.NoError(t, err)

	checkAccount(t, ctx, accKeeper, creatorAddr, sdk.Coins{})
	checkAccount(t, ctx, accKeeper, makerAddr, expectedRetCoins)

	// unauthorized
	bob := createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(offerCoin))
	_, err = keeper.ExecuteContract(ctx, makerAddr, bob, bz, sdk.NewCoins(offerCoin))
	require.Error(t, err)
}

func TestSellMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper

	sellAmount := sdk.NewInt(rand.Int63()%10000 + 2)
	sellCoin := sdk.NewCoin(core.MicroLunaDenom, sellAmount)
	creatorAcc := accKeeper.GetAccount(ctx, creatorAddr)
	creatorAcc.SetCoins(creatorAcc.GetCoins().Add(sellCoin))
	accKeeper.SetAccount(ctx, creatorAcc)

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, sellCoin, core.MicroSDRDenom)
	expectedRetCoins := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()))

	// sell without limit
	sellMsg := MakerHandleMsg{
		Sell: &sellPayload{},
	}

	bz, err := json.Marshal(&sellMsg)

	// normal sell
	_, err = keeper.ExecuteContract(ctx, makerAddr, creatorAddr, bz, sdk.NewCoins(sellCoin))
	require.NoError(t, err)

	checkAccount(t, ctx, accKeeper, creatorAddr, sdk.NewCoins(offerCoin))
	checkAccount(t, ctx, accKeeper, makerAddr, expectedRetCoins)

	// unauthorized
	bob := createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(sellCoin))
	_, err = keeper.ExecuteContract(ctx, makerAddr, bob, bz, sdk.NewCoins(sellCoin))
	require.Error(t, err)
}

func setupMakerContract(t *testing.T) (input TestInput, creatorAddr, makerAddr sdk.AccAddress, initCoin sdk.Coin) {
	input = CreateTestInput(t)

	ctx, keeper, accKeeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.OracleKeeper

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	initCoin = sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)

	creatorAddr = createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(initCoin))

	// upload staking derivates code
	makingCode, err := ioutil.ReadFile("./testdata/maker.wasm")
	require.NoError(t, err)
	makerID, err := keeper.StoreCode(ctx, creatorAddr, makingCode, true)
	require.NoError(t, err)
	require.Equal(t, uint64(1), makerID)

	initMsg := MakerInitMsg{
		OfferDenom: core.MicroSDRDenom,
		AskDenom:   core.MicroLunaDenom,
	}

	initBz, err := json.Marshal(&initMsg)
	makerAddr, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, initBz, nil)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	return
}
