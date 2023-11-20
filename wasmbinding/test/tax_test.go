package wasmbinding_test

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	core "github.com/classic-terra/core/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// go test -v -run ^TestWasmTestSuite/TestTax$ github.com/classic-terra/core/v2/wasmbinding/test
func (s *WasmTestSuite) TestTax() {
	s.SetupTest()
	taxRate := sdk.NewDecWithPrec(11, 2)            // 11%
	s.App.TreasuryKeeper.SetTaxRate(s.Ctx, taxRate) // 11%

	payer := s.TestAccs[0]
	toAddress := s.TestAccs[1]
	// fund an account
	s.FundAcc(payer, sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, 1000000000)))

	// instantiate reflect contract
	contractAddr := s.InstantiateContract(payer, TerraBindingsPath)
	s.Require().NotEmpty(contractAddr)

	// make a bank send message
	coin := sdk.NewInt64Coin(core.MicroLunaDenom, 10000)
	taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()
	updateAmt := coin.Amount.Add(taxAmount)
	updatedCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, updateAmt.Int64()))
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
	s.Require().Equal(sdk.NewInt(1000000000).Sub(updateAmt), res.AmountOf(core.MicroLunaDenom))
}
