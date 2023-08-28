package wasmbinding_test

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	core "github.com/classic-terra/core/v2/types"
	"github.com/classic-terra/core/v2/wasmbinding/bindings"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
	treasurytypes "github.com/classic-terra/core/v2/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// go test -v -run ^TestQuerySwap$ github.com/classic-terra/core/v2/wasmbinding/test
// oracle rate: 1 uluna = 1.7 usdr
// 1000 uluna from trader goes to contract
// 1666 usdr (after 2% tax) is swapped into
func (s *WasmTestSuite) QuerySwap(contractPath string, queryFunc func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})) {
	s.SetupTest()
	actor := s.RandomAccountAddresses(1)[0]

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	// setup swap environment
	// Set Oracle Price
	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	s.App.OracleKeeper.SetLunaExchangeRate(s.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	// Calculate expected swapped SDR
	expectedSwappedSDR := sdk.NewDec(1000).Mul(lunaPriceInSDR)
	tax := markettypes.DefaultMinStabilitySpread.Mul(expectedSwappedSDR)
	expectedSwappedSDR = expectedSwappedSDR.Sub(tax)

	// query swap
	query := bindings.TerraQuery{
		Swap: &markettypes.QuerySwapParams{
			OfferCoin: sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(1000)),
			AskDenom:  core.MicroSDRDenom,
		},
	}

	resp := bindings.SwapQueryResponse{}
	queryFunc(contractAddr, query, &resp)

	s.Require().Equal(expectedSwappedSDR.TruncateInt().String(), resp.Receive.Amount)
}

// go test -v -run ^TestQueryExchangeRates$ github.com/classic-terra/core/v2/wasmbinding/test
func (s *WasmTestSuite) QueryExchangeRates(contractPath string, queryFunc func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})) {
	s.SetupTest()
	actor := s.RandomAccountAddresses(1)[0]

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	lunaPriceInSDR := sdk.NewDecWithPrec(17, 1)
	s.App.OracleKeeper.SetLunaExchangeRate(s.Ctx, core.MicroSDRDenom, lunaPriceInSDR)

	query := bindings.TerraQuery{
		ExchangeRates: &bindings.ExchangeRateQueryParams{
			BaseDenom:   core.MicroLunaDenom,
			QuoteDenoms: []string{core.MicroSDRDenom},
		},
	}

	resp := bindings.ExchangeRatesQueryResponse{}
	queryFunc(contractAddr, query, &resp)

	s.Require().Equal(lunaPriceInSDR, sdk.MustNewDecFromStr(resp.ExchangeRates[0].ExchangeRate))
}

// go test -v -run ^TestQueryTaxRate$ github.com/classic-terra/core/v2/wasmbinding/test
func (s *WasmTestSuite) QueryTaxRate(contractPath string, queryFunc func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})) {
	s.SetupTest()
	actor := s.RandomAccountAddresses(1)[0]

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	query := bindings.TerraQuery{
		TaxRate: &struct{}{},
	}

	resp := bindings.TaxRateQueryResponse{}
	queryFunc(contractAddr, query, &resp)

	s.Require().Equal(treasurytypes.DefaultTaxRate, sdk.MustNewDecFromStr(resp.Rate))
}

// go test -v -run ^TestQueryTaxCap$ github.com/classic-terra/core/v2/wasmbinding/test
func (s *WasmTestSuite) QueryTaxCap(contractPath string, queryFunc func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})) {
	s.SetupTest()
	actor := s.RandomAccountAddresses(1)[0]

	// fund
	s.FundAcc(actor, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(actor, contractPath)
	s.Require().NotEmpty(contractAddr)

	query := bindings.TerraQuery{
		TaxCap: &treasurytypes.QueryTaxCapParams{
			Denom: core.MicroSDRDenom,
		},
	}

	resp := bindings.TaxCapQueryResponse{}
	queryFunc(contractAddr, query, &resp)

	s.Require().Equal(treasurytypes.DefaultTaxPolicy.Cap.Amount.String(), resp.Cap)
}

type ReflectQuery struct {
	Chain *ChainRequest `json:"chain,omitempty"`
}

type ChainRequest struct {
	Request wasmvmtypes.QueryRequest `json:"request"`
}

type ChainResponse struct {
	Data []byte `json:"data"`
}

func (s *WasmTestSuite) queryCustom(contract sdk.AccAddress, request bindings.TerraQuery, response interface{}) {
	msgBz, err := json.Marshal(request)
	s.Require().NoError(err)

	query := ReflectQuery{
		Chain: &ChainRequest{
			Request: wasmvmtypes.QueryRequest{Custom: msgBz},
		},
	}
	queryBz, err := json.Marshal(query)
	s.Require().NoError(err)

	resBz, err := s.App.WasmKeeper.QuerySmart(s.Ctx, contract, queryBz)
	s.Require().NoError(err)
	var resp ChainResponse
	err = json.Unmarshal(resBz, &resp)
	s.Require().NoError(err)
	err = json.Unmarshal(resp.Data, response)
	s.Require().NoError(err)
}

// old bindings contract query
// Binding query messages
type bindingsTesterSwapQueryMsg struct {
	Swap swapQueryMsg `json:"swap"`
}
type bindingsTesterTaxRateQueryMsg struct {
	TaxRate struct{} `json:"tax_rate"`
}
type bindingsTesterTaxCapQueryMsg struct {
	TaxCap *treasurytypes.QueryTaxCapParams `json:"tax_cap"`
}
type bindingsTesterExchangeRatesQueryMsg struct {
	ExchangeRates *bindings.ExchangeRateQueryParams `json:"exchange_rates"`
}
type swapQueryMsg struct {
	OfferCoin wasmvmtypes.Coin `json:"offer_coin"`
	AskDenom  string           `json:"ask_denom"`
}

func (s *WasmTestSuite) queryOldBindings(contract sdk.AccAddress, request bindings.TerraQuery, response interface{}) {
	var msgBz []byte
	switch {
	case request.Swap != nil:
		query := bindingsTesterSwapQueryMsg{
			Swap: swapQueryMsg{
				OfferCoin: wasmvmtypes.Coin{
					Denom:  request.Swap.OfferCoin.Denom,
					Amount: request.Swap.OfferCoin.Amount.String(),
				},
				AskDenom: request.Swap.AskDenom,
			},
		}
		var err error
		msgBz, err = json.Marshal(query)
		s.Require().NoError(err)
	case request.ExchangeRates != nil:
		query := bindingsTesterExchangeRatesQueryMsg{
			ExchangeRates: request.ExchangeRates,
		}
		var err error
		msgBz, err = json.Marshal(query)
		s.Require().NoError(err)
	case request.TaxRate != nil:
		query := bindingsTesterTaxRateQueryMsg{
			TaxRate: struct{}{},
		}
		var err error
		msgBz, err = json.Marshal(query)
		s.Require().NoError(err)
	case request.TaxCap != nil:
		query := bindingsTesterTaxCapQueryMsg{
			TaxCap: request.TaxCap,
		}
		var err error
		msgBz, err = json.Marshal(query)
		s.Require().NoError(err)
	}

	resBz, err := s.App.WasmKeeper.QuerySmart(s.Ctx, contract, msgBz)
	s.Require().NoError(err)
	err = json.Unmarshal(resBz, response)
	s.Require().NoError(err)
}

func (s *WasmTestSuite) queryStargate(contract sdk.AccAddress, request bindings.TerraQuery, response interface{}) {
	queryBz, err := json.Marshal(request)
	s.Require().NoError(err)

	resBz, err := s.App.WasmKeeper.QuerySmart(s.Ctx, contract, queryBz)
	s.Require().NoError(err)
	err = json.Unmarshal(resBz, response)
	s.Require().NoError(err)
}
