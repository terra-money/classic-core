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
	"github.com/terra-project/core/x/treasury"
	treasurywasm "github.com/terra-project/core/x/treasury/wasm"
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
	Send *sendPayload `json:"send,omitempty"`
}

type buyPayload struct {
	Limit     uint64 `json:"limit,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}

type sellPayload struct {
	Limit     uint64 `json:"limit,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}

type sendPayload struct {
	Coin      wasmTypes.Coin `json:"coin"`
	Recipient string         `json:"recipient"`
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

// MakerTreasuryQuerymsg nolint
type MakerTreasuryQuerymsg struct {
	Reflect treasuryQueryMsg `json:"reflect,omitempty"`
}

type treasuryQueryMsg struct {
	TerraQueryWrapper treasuryQueryWrapper `json:"query"`
}

type treasuryQueryWrapper struct {
	Route     string                   `json:"route"`
	QueryData treasurywasm.CosmosQuery `json:"query_data"`
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
	makerID, err := keeper.StoreCode(ctx, creatorAddr, makingCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), makerID)

	// valid instantiate
	initMsg := MakerInitMsg{
		OfferDenom: core.MicroSDRDenom,
		AskDenom:   core.MicroLunaDenom,
	}

	initBz, err := json.Marshal(&initMsg)
	makerAddr, err := keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, initBz, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	// invalid init msg
	_, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, []byte{}, nil, true)
	require.Error(t, err)
}

func TestMarketQuerier(t *testing.T) {
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

func TestTreasuryQuerier(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, _, makerAddr, _ := setupMakerContract(t)

	ctx, keeper, treasuryKeeper := input.Ctx, input.WasmKeeper, input.TreasuryKeeper

	expectedTaxRate := treasuryKeeper.GetTaxRate(ctx)
	expectedTaxCap := treasuryKeeper.GetTaxCap(ctx, core.MicroSDRDenom)

	// querier test
	taxRateQueryMsg := MakerTreasuryQuerymsg{
		Reflect: treasuryQueryMsg{
			TerraQueryWrapper: treasuryQueryWrapper{
				Route: types.WasmQueryRouteTreasury,
				QueryData: treasurywasm.CosmosQuery{
					TaxRate: &struct{}{},
				},
			},
		},
	}

	bz, err := json.Marshal(taxRateQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, makerAddr, bz)
	require.NoError(t, err)

	var taxRateRes treasurywasm.TaxRateQueryResponse
	err = json.Unmarshal(res, &taxRateRes)
	require.NoError(t, err)

	taxRate, err := sdk.NewDecFromStr(taxRateRes.Rate)
	require.NoError(t, err)
	require.Equal(t, expectedTaxRate, taxRate)

	taxCapQueryMsg := MakerTreasuryQuerymsg{
		Reflect: treasuryQueryMsg{
			TerraQueryWrapper: treasuryQueryWrapper{
				Route: types.WasmQueryRouteTreasury,
				QueryData: treasurywasm.CosmosQuery{
					TaxCap: &treasury.QueryTaxCapParams{
						Denom: core.MicroSDRDenom,
					},
				},
			},
		},
	}

	bz, err = json.Marshal(taxCapQueryMsg)
	require.NoError(t, err)

	res, err = keeper.queryToContract(ctx, makerAddr, bz)
	require.NoError(t, err)

	var taxCapRes treasurywasm.TaxCapQueryResponse
	err = json.Unmarshal(res, &taxCapRes)
	require.NoError(t, err)

	taxCap, ok := sdk.NewIntFromString(taxCapRes.Cap)
	require.True(t, ok)
	require.Equal(t, expectedTaxCap, taxCap)
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

func TestBuyAndSendMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper, treasuryKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.TreasuryKeeper
	treasuryKeeper.SetTaxRate(ctx, sdk.ZeroDec())

	retCoin, spread, err := input.MarketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	expectedRetCoins := sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()))

	// buy without limit
	buyMsg := MakerHandleMsg{
		Buy: &buyPayload{
			Recipient: creatorAddr.String(),
		},
	}

	bz, err := json.Marshal(&buyMsg)

	// normal buy
	_, err = keeper.ExecuteContract(ctx, makerAddr, creatorAddr, bz, sdk.NewCoins(offerCoin))
	require.NoError(t, err)

	checkAccount(t, ctx, accKeeper, creatorAddr, expectedRetCoins)
	checkAccount(t, ctx, accKeeper, makerAddr, sdk.Coins{})
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

func TestSendMsg(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	// Check tax charging
	ctx, keeper, accKeeper, treasuryKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.TreasuryKeeper
	taxRate := treasuryKeeper.GetTaxRate(ctx)
	taxCap := treasuryKeeper.GetTaxCap(ctx, core.MicroSDRDenom)

	sendMsg := MakerHandleMsg{
		Send: &sendPayload{
			Coin:      types.EncodeSdkCoin(offerCoin),
			Recipient: creatorAddr.String(),
		},
	}

	bz, err := json.Marshal(&sendMsg)

	expectedTaxAmount := taxRate.MulInt(offerCoin.Amount).TruncateInt()
	if expectedTaxAmount.GT(taxCap) {
		expectedTaxAmount = taxCap
	}

	_, err = keeper.ExecuteContract(ctx, makerAddr, creatorAddr, bz, sdk.NewCoins(offerCoin))
	require.NoError(t, err)

	checkAccount(t, ctx, accKeeper, creatorAddr, sdk.NewCoins(offerCoin.Sub(sdk.NewCoin(offerCoin.Denom, expectedTaxAmount))))
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
	makerID, err := keeper.StoreCode(ctx, creatorAddr, makingCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), makerID)

	initMsg := MakerInitMsg{
		OfferDenom: core.MicroSDRDenom,
		AskDenom:   core.MicroLunaDenom,
	}

	initBz, err := json.Marshal(&initMsg)
	makerAddr, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, initBz, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	return
}
