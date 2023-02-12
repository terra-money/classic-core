package ante_test

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/classic-terra/core/custom/auth/ante"
	core "github.com/classic-terra/core/types"
	treasury "github.com/classic-terra/core/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *AnteTestSuite) TestEnsureBurnTaxModule() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewBurnTaxFeeDecorator(suite.app.TreasuryKeeper, suite.app.BankKeeper, suite.app.DistrKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgSend(addr1, addr1, sendCoins)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass without burn before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the tax height block
	suite.ctx = suite.ctx.WithBlockHeight(10000000)
	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	taxes := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, expectedTax.Int64()))

	bk := suite.app.BankKeeper
	bk.MintCoins(suite.ctx, minttypes.ModuleName, sendCoins)

	// Populate the FeeCollector module with taxes
	bk.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.FeeCollectorName, taxes)
	feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.FeeCollectorName)

	amountFee := bk.GetAllBalances(suite.ctx, feeCollector.GetAddress())
	suite.Require().Equal(amountFee, taxes)
	totalSupply, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})

	// must pass with tax and burn
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")

	// Burn the taxes
	tk.BurnCoinsFromBurnAccount(suite.ctx)
	suite.Require().NoError(err)

	supplyAfterBurn, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})

	// Total supply should have decreased by the tax amount
	suite.Require().Equal(taxes, totalSupply.Sub(supplyAfterBurn))

}

// go test -v -run ^TestAnteTestSuite/TestSplitTax$ github.com/classic-terra/core/custom/auth/ante
func (suite *AnteTestSuite) TestSplitTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewBurnTaxFeeDecorator(suite.app.TreasuryKeeper, suite.app.BankKeeper, suite.app.DistrKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgSend(addr1, addr1, sendCoins)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// ===tax should not be split before TaxPowerSplitHeight===
	suite.ctx = suite.ctx.WithBlockHeight(ante.TaxPowerUpgradeHeight)
	burnModule := suite.app.AccountKeeper.GetModuleAddress(treasury.BurnModuleName)

	feePoolBefore := suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool
	burnBefore := suite.app.BankKeeper.GetAllBalances(suite.ctx, burnModule)

	// send taxes to fee collector to simulate DeductFeeDecorator antehandler
	taxes := suite.DeductFees(sendAmount)

	// send tx to BurnTaxFeeDecorator antehandler
	_, err = antehandler(suite.ctx, tx, false)

	feePoolAfter := suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool
	burnAfter := suite.app.BankKeeper.GetAllBalances(suite.ctx, burnModule)

	// expected: 1000 tax goes to burn
	suite.Require().Equal(feePoolAfter, feePoolBefore)
	suite.Require().Equal(burnBefore.Add(taxes...), burnAfter)

	// ===tax should be split after TaxPowerSplitHeight===
	suite.ctx = suite.ctx.WithBlockHeight(ante.TaxPowerSplitHeight)

	feePoolBefore = suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool
	burnBefore = suite.app.BankKeeper.GetAllBalances(suite.ctx, burnModule)

	// send taxes to fee collector to simulate DeductFeeDecorator antehandler
	taxes = suite.DeductFees(sendAmount)
	splitTaxRate := suite.app.TreasuryKeeper.GetBurnSplitRate(suite.ctx)
	splitTaxesDec := splitTaxRate.MulInt(taxes.AmountOf(core.MicroSDRDenom))
	splitTaxesCoin := sdk.NewCoin(core.MicroSDRDenom, splitTaxesDec.RoundInt())

	// send tx to BurnTaxFeeDecorator antehandler
	_, err = antehandler(suite.ctx, tx, false)

	feePoolAfter = suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool
	burnAfter = suite.app.BankKeeper.GetAllBalances(suite.ctx, burnModule)
	taxesAfter := taxes[0].Sub(splitTaxesCoin)

	// expected: 500 tax goes to fee pool, 500 tax goes to burn
	suite.Require().Equal(splitTaxesDec, feePoolAfter[0].Amount)
	suite.Require().Equal(burnBefore.Add(taxesAfter), burnAfter)
}

func (suite *AnteTestSuite) DeductFees(sendAmount int64) sdk.Coins {
	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}
	taxes := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, expectedTax.Int64()))
	bk := suite.app.BankKeeper
	bk.MintCoins(suite.ctx, minttypes.ModuleName, taxes)
	// populate the FeeCollector module with taxes
	bk.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.FeeCollectorName, taxes)

	return taxes
}
