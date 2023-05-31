package wasmbinding_test

import (
	"github.com/classic-terra/core/wasmbinding/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TERRA_BINDINGS_DIR           = "../testdata/terra_reflect.wasm"
	TERRA_RENOVATED_BINDINGS_DIR = "../testdata/old/bindings_tester.wasm"
)

// go test -v -run ^TestWasmTestSuite/TestBindingsAll$ github.com/classic-terra/core/wasmbinding/test
func (s *WasmTestSuite) TestBindingsAll() {
	cases := []struct {
		name        string
		dir         string
		executeFunc func(contract sdk.AccAddress, sender sdk.AccAddress, msg bindings.TerraMsg, funds sdk.Coin) error
		queryFunc   func(contract sdk.AccAddress, request bindings.TerraQuery, response interface{})
	}{
		{
			name:        "Terra",
			dir:         TERRA_BINDINGS_DIR,
			executeFunc: s.executeCustom,
			queryFunc:   s.queryCustom,
		},
		{
			name:        "Old Terra bindings",
			dir:         TERRA_RENOVATED_BINDINGS_DIR,
			executeFunc: s.executeOldBindings,
			queryFunc:   s.queryOldBindings,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// Msg
			s.Run("TestSwap", func() {
				s.Swap(tc.dir, tc.executeFunc)
			})
			s.Run("TestSwapSend", func() {
				s.SwapSend(tc.dir, tc.executeFunc)
			})

			// Query
			s.Run("TestQuerySwap", func() {
				s.QuerySwap(tc.dir, tc.queryFunc)
			})
			s.Run("TestQueryExchangeRates", func() {
				s.QueryExchangeRates(tc.dir, tc.queryFunc)
			})
			s.Run("TestQueryTaxRate", func() {
				s.QueryTaxRate(tc.dir, tc.queryFunc)
			})
			s.Run("TestQueryTaxCap", func() {
				s.QueryTaxCap(tc.dir, tc.queryFunc)
			})
		})
	}
}
