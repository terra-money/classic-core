package wasmbinding_test

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	customante "github.com/classic-terra/core/custom/auth/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// go test -v -run ^TestWasmTestSuite/TestTax$ github.com/classic-terra/core/wasmbinding/test
func (s *WasmTestSuite) TestTax() {
	s.SetupTest()
	taxRate := sdk.NewDecWithPrec(11, 2)            // 11%
	s.App.TreasuryKeeper.SetTaxRate(s.Ctx, taxRate) // 11%

	payer := s.TestAccs[0]
	toAddress := s.TestAccs[1]
	// fund an account
	s.FundAcc(payer, sdk.NewCoins(sdk.NewInt64Coin("uluna", 1000000000)))
	s.Ctx = s.Ctx.WithBlockHeight(customante.TaxPowerUpgradeHeight + 1)

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(payer, TERRA_BINDINGS_DIR)
	s.Require().NotEmpty(contractAddr)

	// make a bank send message
	coin := sdk.NewInt64Coin("uluna", 10000)
	taxAmount := coin.Amount.ToDec().Mul(taxRate).TruncateInt()
	updateAmt := coin.Amount.Add(taxAmount)
	updatedCoins := sdk.NewCoins(sdk.NewInt64Coin("uluna", updateAmt.Int64()))
	reflectMsg := ReflectExec{
		ReflectMsg: &ReflectMsgs{
			Msgs: []wasmvmtypes.CosmosMsg{{
				Bank: &wasmvmtypes.BankMsg{
					Send: &wasmvmtypes.SendMsg{
						ToAddress: toAddress.String(),
						Amount: []wasmvmtypes.Coin{{
							Denom:  coin.Denom,
							Amount: coin.Amount.String(),
						}},
					},
				},
			}},
		},
	}
	reflectBz, err := json.Marshal(reflectMsg)
	s.Require().NoError(err)

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(s.App.WasmKeeper)
	_, err = contractKeeper.Execute(s.Ctx, contractAddr, payer, reflectBz, updatedCoins)
	s.Require().NoError(err)

	// check balance
	res := s.App.BankKeeper.GetAllBalances(s.Ctx, payer)
	s.Require().Equal(sdk.NewInt(1000000000).Sub(updateAmt), res.AmountOf("uluna"))
}
