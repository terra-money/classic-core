package wasmbinding_test

import (
	"github.com/classic-terra/core/v2/wasmbinding/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TerraBindingsPath          = "../testdata/terra_reflect.wasm"
	TerraRenovatedBindingsPath = "../testdata/old/bindings_tester.wasm"
)

// go test -v -run ^TestWasmTestSuite/TestBindingsAll$ github.com/classic-terra/core/v2/wasmbinding/test
func (s *WasmTestSuite) TestBindingsAll() {
	cases := []struct {
		name        string
		path        string
		executeFunc func(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error
		queryFunc   func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})
	}{
		{
			name:        "Terra",
			path:        TerraBindingsPath,
			executeFunc: s.executeCustom,
			queryFunc:   s.queryCustom,
		},
		{
			name:        "Old Terra bindings",
			path:        TerraRenovatedBindingsPath,
			executeFunc: s.executeOldBindings,
			queryFunc:   s.queryOldBindings,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// Msg
			s.Run("TestSwap", func() {
				s.Swap(tc.path, tc.executeFunc)
			})
			s.Run("TestSwapSend", func() {
				s.SwapSend(tc.path, tc.executeFunc)
			})

			// Query
			s.Run("TestQuerySwap", func() {
				s.QuerySwap(tc.path, tc.queryFunc)
			})
			s.Run("TestQueryExchangeRates", func() {
				s.QueryExchangeRates(tc.path, tc.queryFunc)
			})
			s.Run("TestQueryTaxRate", func() {
				s.QueryTaxRate(tc.path, tc.queryFunc)
			})
			s.Run("TestQueryTaxCap", func() {
				s.QueryTaxCap(tc.path, tc.queryFunc)
			})
		})
	}
}
