package keeper

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	core "github.com/terra-money/core/types"
	marketwasm "github.com/terra-money/core/x/market/wasm"
	oraclewasm "github.com/terra-money/core/x/oracle/wasm"
	treasurylegacy "github.com/terra-money/core/x/wasm/legacyqueriers/treasury"
	"github.com/terra-money/core/x/wasm/types"
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
	Coin      wasmvmtypes.Coin `json:"coin"`
	Recipient string           `json:"recipient"`
}

// MakerQueryMsg nolint
type MakerQueryMsg struct {
	Simulate simulateQuery `json:"simulate"`
}

type simulateQuery struct {
	OfferCoin wasmvmtypes.Coin `json:"offer"`
}

type simulateResponse struct {
	SellCoin wasmvmtypes.Coin `json:"sell"`
	BuyCoin  wasmvmtypes.Coin `json:"buy"`
}

// MakerTreasuryQuerymsg nolint
type MakerTreasuryQuerymsg struct {
	Reflect treasuryQueryMsg `json:"reflect,omitempty"`
}

type treasuryQueryMsg struct {
	TerraQueryWrapper treasuryQueryWrapper `json:"query"`
}

type treasuryQueryWrapper struct {
	Route     string                     `json:"route"`
	QueryData treasurylegacy.CosmosQuery `json:"query_data"`
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
type bindingsTesterContractInfoQuerymsg struct {
	ContractInfo contractInfoQueryMsg `json:"contract_info"`
}
type swapQueryMsg struct {
	OfferCoin wasmvmtypes.Coin `json:"offer_coin"`
	AskDenom  string           `json:"ask_denom"`
}
type taxRateQueryMsg struct{}
type taxCapQueryMsg struct {
	Denom string `json:"denom"`
}
type exchangeRatesQueryMsg struct {
	BaseDenom   string   `json:"base_denom"`
	QuoteDenoms []string `json:"quote_denoms"`
}
type contractInfoQueryMsg struct {
	ContractAddress string `json:"contract_address"`
}

func TestInstantiateMaker(t *testing.T) {
	input := CreateTestInput(t)

	ctx, keeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.OracleKeeper
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	_, _, creatorAddr := keyPubAddr()

	// upload staking derivatives code
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
	makerAddr, _, err := keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, sdk.AccAddress{}, initBz, nil)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	// invalid init msg
	_, _, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, sdk.AccAddress{}, []byte{}, nil)
	require.Error(t, err)
}

func TestMarketQuerier(t *testing.T) {
	input, _, testerAddr, offerCoin := setupBindingsTesterContract(t)

	ctx, keeper, marketKeeper := input.Ctx, input.WasmKeeper, input.MarketKeeper

	swapQueryMsg := bindingsTesterSwapQueryMsg{
		Swap: swapQueryMsg{
			OfferCoin: wasmvmtypes.Coin{
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
	require.Equal(t, wasmvmtypes.Coin{
		Denom:  core.MicroLunaDenom,
		Amount: retAmount.String(),
	}, swapResponse.Receive)
}

func TestTreasuryQuerier(t *testing.T) {
	input, _, testerAddr, _ := setupBindingsTesterContract(t)
	ctx, keeper := input.Ctx, input.WasmKeeper

	taxRate := sdk.ZeroDec()
	taxRateQueryMsg := bindingsTesterTaxRateQueryMsg{
		TaxRate: taxRateQueryMsg{},
	}

	bz, err := json.Marshal(taxRateQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var taxRateResponse treasurylegacy.TaxRateQueryResponse
	err = json.Unmarshal(res, &taxRateResponse)
	require.NoError(t, err)

	taxRateDec, err := sdk.NewDecFromStr(taxRateResponse.Rate)
	require.NoError(t, err)
	require.Equal(t, taxRate, taxRateDec)

	taxCap := sdk.ZeroInt()
	taxCapQueryMsg := bindingsTesterTaxCapQueryMsg{
		TaxCap: taxCapQueryMsg{
			Denom: core.MicroSDRDenom,
		},
	}

	bz, err = json.Marshal(taxCapQueryMsg)
	require.NoError(t, err)

	res, err = keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var taxCapResponse treasurylegacy.TaxCapQueryResponse
	err = json.Unmarshal(res, &taxCapResponse)
	require.NoError(t, err)
	require.Equal(t, taxCap.String(), taxCapResponse.Cap)
}

func TestExchangeRatesQuerier(t *testing.T) {
	input, _, testerAddr, _ := setupBindingsTesterContract(t)

	ctx, keeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.OracleKeeper

	exchangeRateQueryMsg := bindingsTesterExchangeRatesQueryMsg{
		ExchangeRates: exchangeRatesQueryMsg{
			BaseDenom: core.MicroLunaDenom,
			QuoteDenoms: []string{
				core.MicroKRWDenom,
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

func TestContractInfoQuerier(t *testing.T) {
	input, _, testerAddr, _ := setupBindingsTesterContract(t)

	ctx, keeper := input.Ctx, input.WasmKeeper

	contractInfoQueryMsg := bindingsTesterContractInfoQuerymsg{
		ContractInfo: contractInfoQueryMsg{
			ContractAddress: testerAddr.String(),
		},
	}

	bz, err := json.Marshal(contractInfoQueryMsg)
	require.NoError(t, err)

	res, err := keeper.queryToContract(ctx, testerAddr, bz)
	require.NoError(t, err)

	var contractInfoResponse ContractInfoQueryResponse
	err = json.Unmarshal(res, &contractInfoResponse)
	require.NoError(t, err)

	contractInfo, err := keeper.GetContractInfo(ctx, testerAddr)
	require.NoError(t, err)
	require.Equal(t, contractInfoResponse, ContractInfoQueryResponse{
		CodeID:  contractInfo.CodeID,
		Address: contractInfo.Address,
		Creator: contractInfo.Creator,
		Admin:   contractInfo.Admin,
	})
}

func TestBuyMsg(t *testing.T) {
	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper, bankKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper

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

	checkAccount(t, ctx, accKeeper, bankKeeper, creatorAddr, sdk.Coins{})
	checkAccount(t, ctx, accKeeper, bankKeeper, makerAddr, expectedRetCoins)

	// unauthorized
	bob := createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(offerCoin))
	_, err = keeper.ExecuteContract(ctx, makerAddr, bob, bz, sdk.NewCoins(offerCoin))
	require.Error(t, err)
}

func TestBuyAndSendMsg(t *testing.T) {
	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper, bankKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper

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

	checkAccount(t, ctx, accKeeper, bankKeeper, creatorAddr, expectedRetCoins)
	checkAccount(t, ctx, accKeeper, bankKeeper, makerAddr, sdk.Coins{})
}

func TestSellMsg(t *testing.T) {
	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	ctx, keeper, accKeeper, bankKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper

	sellAmount := sdk.NewInt(rand.Int63()%10000 + 2)
	sellCoin := sdk.NewCoin(core.MicroLunaDenom, sellAmount)
	err := FundAccount(input, creatorAddr, sdk.NewCoins(sellCoin))
	require.NoError(t, err)

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

	checkAccount(t, ctx, accKeeper, bankKeeper, creatorAddr, sdk.NewCoins(offerCoin))
	checkAccount(t, ctx, accKeeper, bankKeeper, makerAddr, expectedRetCoins)

	// unauthorized
	bob := createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(sellCoin))
	_, err = keeper.ExecuteContract(ctx, makerAddr, bob, bz, sdk.NewCoins(sellCoin))
	require.Error(t, err)
}

func TestSendMsg(t *testing.T) {
	input, creatorAddr, makerAddr, offerCoin := setupMakerContract(t)

	// Check tax charging
	ctx, keeper, accKeeper, bankKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper

	sendMsg := MakerHandleMsg{
		Send: &sendPayload{
			Coin:      types.EncodeSdkCoin(offerCoin),
			Recipient: creatorAddr.String(),
		},
	}

	bz, err := json.Marshal(&sendMsg)

	_, err = keeper.ExecuteContract(ctx, makerAddr, creatorAddr, bz, sdk.NewCoins(offerCoin))
	require.NoError(t, err)

	checkAccount(t, ctx, accKeeper, bankKeeper, creatorAddr, sdk.NewCoins(offerCoin))
}

func setupMakerContract(t *testing.T) (input TestInput, creatorAddr, makerAddr sdk.AccAddress, initCoin sdk.Coin) {
	input = CreateTestInput(t)

	ctx, keeper, accKeeper, bankKeeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper, input.OracleKeeper

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	initCoin = sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)

	creatorAddr = createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(initCoin))

	// upload staking derivatives code
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
	makerAddr, _, err = keeper.InstantiateContract(input.Ctx, makerID, creatorAddr, sdk.AccAddress{}, initBz, nil)
	require.NoError(t, err)
	require.NotEmpty(t, makerAddr)

	return
}

func setupBindingsTesterContract(t *testing.T) (input TestInput, creatorAddr, bindingsTesterAddr sdk.AccAddress, initCoin sdk.Coin) {
	input = CreateTestInput(t)

	ctx, keeper, accKeeper, bankKeeper, oracleKeeper := input.Ctx, input.WasmKeeper, input.AccKeeper, input.BankKeeper, input.OracleKeeper

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	lunaPriceInUSD := sdk.NewDecWithPrec(15, 1)
	lunaPriceInKRW := sdk.NewDec(1300)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroSDRDenom, lunaPriceInSDR)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroUSDDenom, lunaPriceInUSD)
	oracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, lunaPriceInKRW)

	swapAmountInSDR := lunaPriceInSDR.MulInt64(rand.Int63()%10000 + 2).TruncateInt()
	initCoin = sdk.NewCoin(core.MicroSDRDenom, swapAmountInSDR)

	creatorAddr = createFakeFundedAccount(ctx, accKeeper, bankKeeper, sdk.NewCoins(initCoin))

	// upload binding_tester contract codes
	bindingsTCode, err := ioutil.ReadFile("./testdata/bindings_tester.wasm")
	require.NoError(t, err)
	bindingsTesterID, err := keeper.StoreCode(ctx, creatorAddr, bindingsTCode)
	require.NoError(t, err)
	require.Equal(t, uint64(1), bindingsTesterID)

	type EmptyStruct struct{}
	initBz, err := json.Marshal(&EmptyStruct{})
	bindingsTesterAddr, _, err = keeper.InstantiateContract(input.Ctx, bindingsTesterID, creatorAddr, sdk.AccAddress{}, initBz, nil)
	require.NoError(t, err)
	require.NotEmpty(t, bindingsTesterAddr)

	return
}
