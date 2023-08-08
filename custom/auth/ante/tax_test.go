package ante_test

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/classic-terra/core/v2/custom/auth/ante"
	core "github.com/classic-terra/core/v2/types"
	"github.com/classic-terra/core/v2/types/fork"
	markettypes "github.com/classic-terra/core/v2/x/market/types"
)

func (suite *AnteTestSuite) TestEnsureMempoolFeesGas() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set high gas price so standard test fee fails
	atomPrice := sdk.NewDecCoinFromDec("atom", sdk.NewDec(200).Quo(sdk.NewDec(100000)))
	highGasPrice := []sdk.DecCoin{atomPrice}
	suite.ctx = suite.ctx.WithMinGasPrices(highGasPrice)

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should have errored on too low fee for local gasPrice")

	// Set IsCheckTx to false
	suite.ctx = suite.ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check minGasPrice in DeliverTx
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "MempoolFeeDecorator returned error in DeliverTx")

	// Set IsCheckTx back to true for testing sufficient mempool fee
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	atomPrice = sdk.NewDecCoinFromDec("atom", sdk.NewDec(0).Quo(sdk.NewDec(100000)))
	lowGasPrice := []sdk.DecCoin{atomPrice}
	suite.ctx = suite.ctx.WithMinGasPrices(lowGasPrice)

	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
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

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSwapSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoin := sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount)
	msg := markettypes.NewMsgSwapSend(addr1, addr1, sendCoin, core.MicroKRWDenom)

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

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesMultiSend() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgMultiSend(
		[]banktypes.Input{
			banktypes.NewInput(addr1, sendCoins),
			banktypes.NewInput(addr1, sendCoins),
		},
		[]banktypes.Output{
			banktypes.NewOutput(addr1, sendCoins),
			banktypes.NewOutput(addr1, sendCoins),
		},
	)

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

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	// must pass with tax
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax.Add(expectedTax))))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesInstantiateContract() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := &wasmtypes.MsgInstantiateContract{
		Sender: addr1.String(),
		Admin:  addr1.String(),
		CodeID: 0,
		Msg:    []byte{},
		Funds:  sendCoins,
	}

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

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesExecuteContract() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   addr1.String(),
		Contract: addr1.String(),
		Msg:      []byte{},
		Funds:    sendCoins,
	}

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

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesExec() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := authz.NewMsgExec(addr1, []sdk.Msg{banktypes.NewMsgSend(addr1, addr1, sendCoins)})

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(&msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSendLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount))
	msg := banktypes.NewMsgSend(addr1, addr1, sendCoins)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesSwapSendLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoin := sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount)
	msg := markettypes.NewMsgSwapSend(addr1, addr1, sendCoin, core.MicroKRWDenom)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesMultiSendLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount))
	msg := banktypes.NewMsgMultiSend(
		[]banktypes.Input{
			banktypes.NewInput(addr1, sendCoins),
			banktypes.NewInput(addr1, sendCoins),
		},
		[]banktypes.Output{
			banktypes.NewOutput(addr1, sendCoins),
			banktypes.NewOutput(addr1, sendCoins),
		},
	)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	// must pass with tax
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax.Add(expectedTax))))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesInstantiateContractLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount))
	msg := &wasmtypes.MsgInstantiateContract{
		Sender: addr1.String(),
		Admin:  addr1.String(),
		CodeID: 0,
		Msg:    []byte{},
		Funds:  sendCoins,
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesExecuteContractLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount))
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   addr1.String(),
		Contract: addr1.String(),
		Msg:      []byte{},
		Funds:    sendCoins,
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

func (suite *AnteTestSuite) TestEnsureMempoolFeesExecLunaTax() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroLunaDenom, sendAmount))
	msg := authz.NewMsgExec(addr1, []sdk.Msg{banktypes.NewMsgSend(addr1, addr1, sendCoins)})

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(&msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// Set ChainID to columbus-5
	suite.ctx = suite.ctx.WithChainID(core.ColumbusChainID)

	// set zero gas prices
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins())

	// Set IsCheckTx to true
	suite.ctx = suite.ctx.WithIsCheckTx(true)

	// Luna must pass with tax before the specified tax block height
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored when block height is 1")

	// Set the blockheight past the burn tax height block
	suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)

	// antehandler errors with insufficient fees due to tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().Error(err, "Decorator should errored on low fee for local gasPrice + tax")

	tk := suite.app.TreasuryKeeper
	expectedTax := tk.GetTaxRate(suite.ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(suite.ctx, core.MicroLunaDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	// set tax amount
	suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, expectedTax)))
	tx, err = suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	// must pass with tax
	_, err = antehandler(suite.ctx, tx, false)
	suite.Require().NoError(err, "Decorator should not have errored on fee higher than local gasPrice")
}

// go test -v -run ^TestAnteTestSuite/TestTaxExemption$ github.com/classic-terra/core/v2/custom/auth/ante
func (suite *AnteTestSuite) TestTaxExemption() {
	// keys and addresses
	var privs []cryptotypes.PrivKey
	var addrs []sdk.AccAddress

	// 0, 1: exemption
	// 2, 3: normal
	for i := 0; i < 4; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		privs = append(privs, priv)
		addrs = append(addrs, addr)
	}

	// set send amount
	sendAmt := int64(1000000)
	sendCoin := sdk.NewInt64Coin(core.MicroSDRDenom, sendAmt)
	feeAmt := int64(1000)

	cases := []struct {
		name           string
		msgSigner      cryptotypes.PrivKey
		msgCreator     func() []sdk.Msg
		minFeeAmount   int64
		expectProceeds int64
	}{
		{
			name:      "MsgSend(exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			minFeeAmount:   0,
			expectProceeds: 0,
		}, {
			name:      "MsgSend(normal -> normal)",
			msgSigner: privs[2],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[2], addrs[3], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		}, {
			name:      "MsgSend(exemption -> normal), MsgSend(exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[2], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg2)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		}, {
			name:      "MsgSend(exemption -> exemption), MsgMultiSend(exemption -> normal, exemption -> exemption)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgMultiSend(
					[]banktypes.Input{
						{
							Address: addrs[0].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
						{
							Address: addrs[0].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
					},
					[]banktypes.Output{
						{
							Address: addrs[2].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
						{
							Address: addrs[1].String(),
							Coins:   sdk.NewCoins(sendCoin),
						},
					},
				)
				msgs = append(msgs, msg2)

				return msgs
			},
			minFeeAmount:   feeAmt * 2,
			expectProceeds: feeAmt * 2,
		}, {
			name:      "MsgExecuteContract(exemption), MsgExecuteContract(normal)",
			msgSigner: privs[3],
			msgCreator: func() []sdk.Msg {
				sendAmount := int64(1000000)
				sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
				// get wasm code for wasm contract create and instantiate
				wasmCode, err := os.ReadFile("./testdata/hackatom.wasm")
				suite.Require().NoError(err)
				per := wasmkeeper.NewDefaultPermissionKeeper(suite.app.WasmKeeper)
				// set wasm default params
				suite.app.WasmKeeper.SetParams(suite.ctx, wasmtypes.DefaultParams())
				// wasm create
				CodeID, _, err := per.Create(suite.ctx, addrs[0], wasmCode, nil)
				suite.Require().NoError(err)
				// params for contract init
				r := wasmkeeper.HackatomExampleInitMsg{Verifier: addrs[0], Beneficiary: addrs[0]}
				bz, err := json.Marshal(r)
				suite.Require().NoError(err)
				// change block time for contract instantiate
				suite.ctx = suite.ctx.WithBlockTime(time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC))
				// instantiate contract then set the contract address to tax exemption
				addr, _, err := per.Instantiate(suite.ctx, CodeID, addrs[0], nil, bz, "my label", nil)
				suite.Require().NoError(err)
				suite.app.TreasuryKeeper.AddBurnTaxExemptionAddress(suite.ctx, addr.String())
				// instantiate contract then not set to tax exemption
				addr1, _, err := per.Instantiate(suite.ctx, CodeID, addrs[0], nil, bz, "my label", nil)
				suite.Require().NoError(err)

				var msgs []sdk.Msg
				// msg and signatures
				msg1 := &wasmtypes.MsgExecuteContract{
					Sender:   addrs[0].String(),
					Contract: addr.String(),
					Msg:      []byte{},
					Funds:    sendCoins,
				}
				msgs = append(msgs, msg1)

				msg2 := &wasmtypes.MsgExecuteContract{
					Sender:   addrs[3].String(),
					Contract: addr1.String(),
					Msg:      []byte{},
					Funds:    sendCoins,
				}
				msgs = append(msgs, msg2)
				return msgs
			},
			minFeeAmount:   feeAmt,
			expectProceeds: feeAmt,
		},
	}

	// there should be no coin in burn module
	for _, c := range cases {
		suite.SetupTest(true) // setup
		require := suite.Require()
		tk := suite.app.TreasuryKeeper
		ak := suite.app.AccountKeeper
		bk := suite.app.BankKeeper

		// Set burn split rate to 50%
		tk.SetBurnSplitRate(suite.ctx, sdk.NewDecWithPrec(5, 1))

		fmt.Printf("CASE = %s \n", c.name)
		suite.ctx = suite.ctx.WithBlockHeight(fork.BurnTaxUpgradeHeight)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[0].String())
		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[1].String())

		mfd := ante.NewTaxFeeDecorator(suite.app.TreasuryKeeper)
		antehandler := sdk.ChainAnteDecorators(
			mfd,
			cosmosante.NewDeductFeeDecorator(ak, bk, suite.app.FeeGrantKeeper),
		)

		for i := 0; i < 4; i++ {
			fundCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, 1000000000))
			acc := ak.NewAccountWithAddress(suite.ctx, addrs[i])
			ak.SetAccount(suite.ctx, acc)
			bk.MintCoins(suite.ctx, minttypes.ModuleName, fundCoins)
			bk.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addrs[i], fundCoins)
		}

		// msg and signatures
		feeAmount := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, c.minFeeAmount))
		gasLimit := testdata.NewTestGasLimit()
		require.NoError(suite.txBuilder.SetMsgs(c.msgCreator()...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{c.msgSigner}, []uint64{0}, []uint64{0}
		tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
		require.NoError(err)

		_, err = antehandler(suite.ctx, tx, false)
		require.NoError(err)

		// check fee collector
		feeCollector := ak.GetModuleAccount(suite.ctx, types.FeeCollectorName)
		amountFee := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
		require.Equal(amountFee, sdk.NewCoin("usdr", sdk.NewInt(c.minFeeAmount)))

		// check tax proceeds
		taxProceeds := suite.app.TreasuryKeeper.PeekEpochTaxProceeds(suite.ctx)
		require.Equal(taxProceeds, sdk.NewCoins(sdk.NewCoin("usdr", sdk.NewInt(c.expectProceeds))))
	}
}
