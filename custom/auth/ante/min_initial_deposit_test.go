package ante_test

import (
	// "fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/classic-terra/core/v2/custom/auth/ante"
	// core "github.com/terra-money/core/types"
	// treasury "github.com/terra-money/core/x/treasury/types"

	// "github.com/cosmos/cosmos-sdk/types/query"
	// cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	// "github.com/cosmos/cosmos-sdk/x/auth/types"
	// minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *AnteTestSuite) TestMinInitialDepositRatioDefault() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	midd := ante.NewMinInitialDepositDecorator(suite.app.GovKeeper, suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(midd)

	// set required deposit to uluna
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govtypes.DefaultDepositParams())
	govparams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	govparams.MinDeposit = sdk.NewCoins(
		sdk.NewCoin("uluna", sdk.NewInt(1_000_000)),
	)
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govparams)

	// set initial deposit ratio to 0.0
	ratio := sdk.ZeroDec()
	suite.app.TreasuryKeeper.SetMinInitialDepositRatio(suite.ctx, ratio)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	prop1 := govtypes.NewTextProposal("prop1", "prop1")
	depositCoins1 := sdk.NewCoins()

	// create prop tx
	msg, _ := govtypes.NewMsgSubmitProposal(prop1, depositCoins1, addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// antehandler should not error
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "error: Proposal whithout initial deposit should have gone through")
}

func (suite *AnteTestSuite) TestMinInitialDepositRatioWithSufficientDeposit() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	midd := ante.NewMinInitialDepositDecorator(suite.app.GovKeeper, suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(midd)

	// set required deposit to uluna
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govtypes.DefaultDepositParams())
	govparams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	govparams.MinDeposit = sdk.NewCoins(
		sdk.NewCoin("uluna", sdk.NewInt(1_000_000)),
	)
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govparams)

	// set initial deposit ratio to 0.2
	ratio := sdk.NewDecWithPrec(2, 1)
	suite.app.TreasuryKeeper.SetMinInitialDepositRatio(suite.ctx, ratio)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	prop1 := govtypes.NewTextProposal("prop1", "prop1")
	depositCoins1 := sdk.NewCoins(
		sdk.NewCoin("uluna", sdk.NewInt(200_000)),
	)

	// create prop tx
	msg, _ := govtypes.NewMsgSubmitProposal(prop1, depositCoins1, addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// antehandler should not error
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "error: Proposal with sufficient initial deposit should have gone through")
}

func (suite *AnteTestSuite) TestMinInitialDepositRatioWithInsufficientDeposit() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	midd := ante.NewMinInitialDepositDecorator(suite.app.GovKeeper, suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(midd)

	// set required deposit to uluna
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govtypes.DefaultDepositParams())
	govparams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	govparams.MinDeposit = sdk.NewCoins(
		sdk.NewCoin("uluna", sdk.NewInt(1_000_000)),
	)
	suite.app.GovKeeper.SetDepositParams(suite.ctx, govparams)

	// set initial deposit ratio to 0.2
	ratio := sdk.NewDecWithPrec(2, 1)
	suite.app.TreasuryKeeper.SetMinInitialDepositRatio(suite.ctx, ratio)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	prop1 := govtypes.NewTextProposal("prop1", "prop1")
	depositCoins1 := sdk.NewCoins(
		sdk.NewCoin("uluna", sdk.NewInt(100_000)),
	)

	// create prop tx
	msg, _ := govtypes.NewMsgSubmitProposal(prop1, depositCoins1, addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// antehandler should not error
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "error: Proposal with insufficient initial deposit should have failed")
}
