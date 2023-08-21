package ante_test

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/classic-terra/core/v2/custom/auth/ante"
	core "github.com/classic-terra/core/v2/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// go test -v -run ^TestAnteTestSuite/TestSplitTax$ github.com/classic-terra/core/v2/custom/auth/ante
func (suite *AnteTestSuite) TestSplitTax() {
	suite.runSplitTaxTest(sdk.NewDecWithPrec(1, 0))  // 100%
	suite.runSplitTaxTest(sdk.NewDecWithPrec(1, 1))  // 10%
	suite.runSplitTaxTest(sdk.NewDecWithPrec(1, 2))  // 0.1%
	suite.runSplitTaxTest(sdk.NewDecWithPrec(0, 0))  // 0% burn all taxes (old burn tax behavior)
	suite.runSplitTaxTest(sdk.NewDecWithPrec(-1, 1)) // -10% invalid rate
}

func (suite *AnteTestSuite) runSplitTaxTest(burnSplitRate sdk.Dec) {
	suite.SetupTest(true) // setup
	require := suite.Require()
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	tk := suite.app.TreasuryKeeper
	bk := suite.app.BankKeeper
	dk := suite.app.DistrKeeper
	ak := suite.app.AccountKeeper
	mfd := ante.NewBurnTaxFeeDecorator(ak, tk, bk, dk)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// Set burn split tax
	tk.SetBurnSplitRate(suite.ctx, burnSplitRate)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msg := banktypes.NewMsgSend(addr1, addr1, sendCoins)

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	require.NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	require.NoError(err)

	// Send taxes to fee collector to simulate DeductFeeDecorator antehandler
	taxes := suite.DeductFees(sendAmount)
	feeCollector := ak.GetModuleAccount(suite.ctx, types.FeeCollectorName)

	// expected: fee collector = taxes
	amountFeeBefore := bk.GetAllBalances(suite.ctx, feeCollector.GetAddress())
	require.Equal(amountFeeBefore, taxes)

	totalSupplyBefore, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})
	require.NoError(err)
	fmt.Printf(
		"Before: TotalSupply %v, FeeCollector %v\n",
		totalSupplyBefore,
		amountFeeBefore,
	)

	// send tx to BurnTaxFeeDecorator antehandler
	_, err = antehandler(suite.ctx, tx, false)
	require.NoError(err)

	// burn the burn account
	tk.BurnCoinsFromBurnAccount(suite.ctx)

	feeCollectorAfter := sdk.NewDecCoinsFromCoins(bk.GetAllBalances(suite.ctx, ak.GetModuleAddress(types.FeeCollectorName))...)
	burnTax := sdk.NewDecCoinsFromCoins(taxes...)

	if burnSplitRate.IsPositive() {
		distributionDeltaCoins := burnTax.MulDec(burnSplitRate)

		// expected: community pool 50%
		fmt.Printf("BurnSplitRate %v, DistributionDeltaCoins %v\n", burnSplitRate, distributionDeltaCoins)
		require.Equal(feeCollectorAfter, distributionDeltaCoins)
		burnTax = burnTax.Sub(distributionDeltaCoins)
	}

	totalSupplyAfter, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})
	require.NoError(err)
	if !burnTax.Empty() {
		// expected: total supply = tax - split tax
		require.Equal(
			sdk.NewDecCoinsFromCoins(totalSupplyBefore.Sub(totalSupplyAfter...)...),
			burnTax,
		)
	}

	fmt.Printf(
		"After: TotalSupply %v, FeeCollector %v\n",
		totalSupplyAfter,
		feeCollectorAfter,
	)
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
