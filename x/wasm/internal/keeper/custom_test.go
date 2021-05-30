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

	core "github.com/terra-money/core/types"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	oraclewasm "github.com/terra-money/core/x/oracle/wasm"
	treasurywasm "github.com/terra-money/core/x/treasury/wasm"
	"github.com/terra-money/core/x/wasm/internal/types"
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

// Binding query messages
type bindingsTesterSwapQueryMsg struct {
	Swap swapQueryMsg `json:"swap"`
}
type bindingsTesterTaxRateQueryMsg struct {
	TaxRate taxRateQueryMsg `json:"tax_rate"`
}
type bindingsTesterTaxCapQueryMsg struct {
	TaxCap taxCapQueryMsg `json:"tax_cap"`
}
type bindingsTesterExchangeRatesQueryMsg struct {
	ExchangeRates exchangeRatesQueryMsg `json:"exchange_rates"`
}
type swapQueryMsg struct {
	OfferCoin wasmTypes.Coin `json:"offer_coin"`
	AskDenom  string         `json:"ask_denom"`
}
type taxRateQueryMsg struct{}
type taxCapQueryMsg struct {
	Denom string `json:"denom"`
}
type exchangeRatesQueryMsg struct {
	BaseDenom   string   `json:"base_denom"`
	QuoteDenoms []string `json:"quote_denoms"`
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

	input, _, testerAddr, offerCoin := setupBindingsTesterContract(t)

	ctx, keeper, marketKeeper := input.Ctx, input.WasmKeeper, input.MarketKeeper

	swapQueryMsg := bindingsTesterSwapQueryMsg{
		Swap: swapQueryMsg{
			OfferCoin: wasmTypes.Coin{
				Denom:  core.MicroSDRDenom,
				Amount: offerCoin.Amount.String(),
			},
			AskDenom: core.MicroLunaDenom,
		},
	}

	retCoin, spread, err := marketKeeper.ComputeSwap(input.Ctx, offerCoin, core.MicroLunaDenom)
	retAmount := retCoin.Amount.Mul(sdk.OneDec().Sub(spread)).TruncateInt()

	bz, err := json.Marshal(swapQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var swapResponse marketwasm.SwapQueryResponse
	err = json.Unmarshal(res, &swapResponse)
	require.NoError(t, err)
	require.Equal(t, wasmTypes.Coin{
		Denom:  core.MicroLunaDenom,
		Amount: retAmount.String(),
	}, swapResponse.Receive)
}

func TestTreasuryQuerier(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, _, testerAddr, _ := setupBindingsTesterContract(t)
	ctx, keeper, treasuryKeeper := input.Ctx, input.WasmKeeper, input.TreasuryKeeper

	taxRate := treasuryKeeper.GetTaxRate(ctx)
	taxRateQueryMsg := bindingsTesterTaxRateQueryMsg{
		TaxRate: taxRateQueryMsg{},
	}

	bz, err := json.Marshal(taxRateQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var taxRateResponse treasurywasm.TaxRateQueryResponse
	err = json.Unmarshal(res, &taxRateResponse)
	require.NoError(t, err)

	taxRateDec, err := sdk.NewDecFromStr(taxRateResponse.Rate)
	require.NoError(t, err)
	require.Equal(t, taxRate, taxRateDec)

	taxCap := treasuryKeeper.GetTaxCap(ctx, core.MicroSDRDenom)
	taxCapQueryMsg := bindingsTesterTaxCapQueryMsg{
		TaxCap: taxCapQueryMsg{
			Denom: core.MicroSDRDenom,
		},
	}

	bz, err = json.Marshal(taxCapQueryMsg)
	require.NoError(t, err)

	res, err = keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var taxCapResponse treasurywasm.TaxCapQueryResponse
	err = json.Unmarshal(res, &taxCapResponse)
	require.NoError(t, err)
	require.Equal(t, taxCap.String(), taxCapResponse.Cap)
}

func TestExchangeRatesQuerier(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	input, _, testerAddr, _ := setupBindingsTesterContract(t)

	ctx, keeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.OracleKeeper

	exchangeRateQueryMsg := bindingsTesterExchangeRatesQueryMsg{
		ExchangeRates: exchangeRatesQueryMsg{
			BaseDenom: core.MicroKRWDenom,
			QuoteDenoms: []string{
				core.MicroLunaDenom,
			},
		},
	}

	KRWExchangeRate, err := oracleKeeper.GetLunaExchangeRate(ctx, core.MicroKRWDenom)
	require.NoError(t, err)

	bz, err := json.Marshal(exchangeRateQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var exchangeRateResponse oraclewasm.ExchangeRatesQueryResponse
	err = json.Unmarshal(res, &exchangeRateResponse)
	require.NoError(t, err)

	exchangeRateDec, err := sdk.NewDecFromStr(exchangeRateResponse.ExchangeRates[0].ExchangeRate)
	require.NoError(t, err)
	require.Equal(t, KRWExchangeRate, exchangeRateDec)
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

func setupBindingsTesterContract(t *testing.T) (input TestInput, creatorAddr, bindingsTesterAddr sdk.AccAddress, initCoin sdk.Coin) {
	input = CreateTestInput(t)

	ctx, keeper, accKeeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.OracleKeeper

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	lunaPriceInUSD := sdk.NewDecWithPrec(15, 1)
	lunaPriceInKRW := sdk.NewDec(1300)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroUSDDenom, lunaPriceInUSD)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, lunaPriceInKRW)

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	initCoin = sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)

	creatorAddr = createFakeFundedAccount(ctx, accKeeper, sdk.NewCoins(initCoin))

	// upload staking derivates code
	bindingsTCode, err := ioutil.ReadFile("./testdata/bindings_tester.wasm")
	require.NoError(t, err)
	bindingsTesterID, err := keeper.StoreCode(ctx, creatorAddr, bindingsTCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), bindingsTesterID)

	type EmptyStruct struct{}
	initBz, err := json.Marshal(&EmptyStruct{})
	bindingsTesterAddr, err = keeper.InstantiateContract(input.Ctx, bindingsTesterID, creatorAddr, initBz, nil, true)
	require.NoError(t, err)
	require.NotEmpty(t, bindingsTesterAddr)

	return
}
