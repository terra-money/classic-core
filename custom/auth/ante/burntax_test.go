package ante_test

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/classic-terra/core/custom/auth/ante"
	core "github.com/classic-terra/core/types"
	treasury "github.com/classic-terra/core/x/treasury/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// go test -v -run ^TestAnteTestSuite/TestSplitTax$ github.com/classic-terra/core/custom/auth/ante
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
	mfd := ante.NewBurnTaxFeeDecorator(suite.app.AccountKeeper, tk, bk, dk)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// Set the blockheight past the tax height block
	suite.ctx = suite.ctx.WithBlockHeight(10000000)

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
	feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.FeeCollectorName)

	// expected: fee collector = taxes
	amountFeeBefore := bk.GetAllBalances(suite.ctx, feeCollector.GetAddress())
	require.Equal(amountFeeBefore, taxes)

	totalSupplyBefore, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})
	fmt.Printf(
		"Before: TotalSupply %v, Community %v, FeeCollector %v\n",
		totalSupplyBefore,
		dk.GetFeePool(suite.ctx).CommunityPool,
		amountFeeBefore,
	)

	// send tx to BurnTaxFeeDecorator antehandler
	_, err = antehandler(suite.ctx, tx, false)
	require.NoError(err)

	communityPoolAfter := dk.GetFeePool(suite.ctx).CommunityPool
	burnTax := sdk.NewDecCoinsFromCoins(taxes...)

	if burnSplitRate.IsPositive() {
		splitTaxesDecCoins := burnTax.MulDec(burnSplitRate)

		// expected: community pool 50%
		require.Equal(communityPoolAfter, splitTaxesDecCoins)

		fmt.Printf("BurnSplitRate %v, splitTaxes %v\n", burnSplitRate, splitTaxesDecCoins)
		burnTax = burnTax.Sub(splitTaxesDecCoins)
	}

	// burn the burn account
	tk.BurnCoinsFromBurnAccount(suite.ctx)

	totalSupplyAfter, _, err := bk.GetPaginatedTotalSupply(suite.ctx, &query.PageRequest{})

	if !burnTax.Empty() {
		// expected: total supply = tax - split tax
		require.Equal(
			sdk.NewDecCoinsFromCoins(totalSupplyBefore.Sub(totalSupplyAfter)...),
			burnTax,
		)
	}

	amountFeeAfter := bk.GetAllBalances(suite.ctx, feeCollector.GetAddress())
	// expected: fee collector = 0
	require.True(amountFeeAfter.Empty())

	fmt.Printf(
		"After: TotalSupply %v, Community %v, FeeCollector %v\n",
		totalSupplyAfter,
		communityPoolAfter,
		amountFeeAfter,
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

// the following binance addresses should not be applied tax
// go test -v -run ^TestAnteTestSuite/TestFilterRecipient$ github.com/classic-terra/core/custom/auth/ante
func (suite *AnteTestSuite) TestFilterRecipient() {
	// keys and addresses
	var privs []cryptotypes.PrivKey
	var addrs []sdk.AccAddress

	// 0, 1: binance
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
		name       string
		msgSigner  cryptotypes.PrivKey
		msgCreator func() []sdk.Msg
		burnAmount int64
		feeAmount  int64
	}{
		{
			name:      "MsgSend(binance -> binance)",
			msgSigner: privs[0],
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			// skip this one hence burn amount is 0
			burnAmount: 0,
			feeAmount:  feeAmt,
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
			burnAmount: feeAmt / 2,
			feeAmount:  feeAmt,
		}, {
			name:      "MsgSend(binance -> normal), MsgSend(binance -> binance)",
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
			burnAmount: (feeAmt * 2) / 2,
			feeAmount:  feeAmt * 2,
		}, {
			name:      "MsgSend(binance -> binance), MsgMultiSend(binance -> normal, binance -> binance)",
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
			// tax this one hence burn amount is fee amount
			burnAmount: (feeAmt * 3) / 2,
			feeAmount:  feeAmt * 3,
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
		suite.ctx = suite.ctx.WithBlockHeight(ante.TaxPowerUpgradeHeight)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[0].String())
		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[1].String())

		mfd := ante.NewBurnTaxFeeDecorator(ak, tk, bk, suite.app.DistrKeeper)
		antehandler := sdk.ChainAnteDecorators(
			cosmosante.NewDeductFeeDecorator(ak, bk, suite.app.FeeGrantKeeper),
			mfd,
		)

		for i := 0; i < 4; i++ {
			fundCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, 1000000000))
			acc := ak.NewAccountWithAddress(suite.ctx, addrs[i])
			ak.SetAccount(suite.ctx, acc)
			bk.MintCoins(suite.ctx, minttypes.ModuleName, fundCoins)
			bk.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addrs[i], fundCoins)
		}

		// msg and signatures
		feeAmount := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, c.feeAmount))
		gasLimit := testdata.NewTestGasLimit()
		require.NoError(suite.txBuilder.SetMsgs(c.msgCreator()...))
		suite.txBuilder.SetFeeAmount(feeAmount)
		suite.txBuilder.SetGasLimit(gasLimit)

		privs, accNums, accSeqs := []cryptotypes.PrivKey{c.msgSigner}, []uint64{0}, []uint64{0}
		tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
		require.NoError(err)

		// check fee decorator and burn module amount before ante handler
		feeCollector := ak.GetModuleAccount(suite.ctx, types.FeeCollectorName)
		burnModule := ak.GetModuleAccount(suite.ctx, treasury.BurnModuleName)

		amountFeeBefore := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
		amountBurnBefore := bk.GetBalance(suite.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
		amountCommunityBefore := suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
		fmt.Printf("before: fee = %v, burn = %v, community = %v\n", amountFeeBefore, amountFeeBefore, amountCommunityBefore)

		_, err = antehandler(suite.ctx, tx, false)
		require.NoError(err)

		// check fee decorator
		amountFee := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
		amountBurn := bk.GetBalance(suite.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
		amountCommunity := suite.app.DistrKeeper.GetFeePool(suite.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
		fmt.Printf("after : fee = %v, burn = %v, community = %v\n", amountFee, amountBurn, amountCommunity)

		if c.burnAmount > 0 {
			require.Equal(amountBurnBefore.Amount.Add(sdk.NewInt(c.burnAmount)), amountBurn.Amount)
			require.Equal(amountFeeBefore, amountFee)
			require.Equal(amountCommunity, amountBurn.Amount.ToDec())
		} else {
			require.Equal(amountBurnBefore, amountBurn)
			require.Equal(amountFeeBefore.Amount.Add(sdk.NewInt(c.feeAmount)), amountFee.Amount)
		}
	}
}
