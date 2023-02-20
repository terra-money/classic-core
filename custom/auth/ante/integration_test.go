package ante_test

import (
	"fmt"

	"github.com/classic-terra/core/custom/auth/ante"
	customante "github.com/classic-terra/core/custom/auth/ante"
	core "github.com/classic-terra/core/types"
	treasurytypes "github.com/classic-terra/core/x/treasury/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// go test -v -run ^TestAnteTestSuite/TestIntegrationTaxExemption$ github.com/classic-terra/core/custom/auth/ante
func (suite *AnteTestSuite) TestIntegrationTaxExemption() {
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
		name              string
		msgSigner         int
		msgCreator        func() []sdk.Msg
		expectedFeeAmount int64
	}{
		{
			name:      "MsgSend(exemption -> exemption)",
			msgSigner: 0,
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			expectedFeeAmount: 0,
		}, {
			name:      "MsgSend(normal -> normal)",
			msgSigner: 2,
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[2], addrs[3], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			expectedFeeAmount: feeAmt,
		}, {
			name:      "MsgSend(exemption -> normal), MsgSend(exemption -> exemption)",
			msgSigner: 0,
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[2], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg2)

				return msgs
			},
			// tax this one hence burn amount is fee amount
			expectedFeeAmount: feeAmt,
		}, {
			name:      "MsgSend(exemption -> exemption), MsgMultiSend(exemption -> normal, exemption)",
			msgSigner: 0,
			msgCreator: func() []sdk.Msg {
				var msgs []sdk.Msg

				msg1 := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sendCoin))
				msgs = append(msgs, msg1)
				msg2 := banktypes.NewMsgMultiSend(
					[]banktypes.Input{
						{
							Address: addrs[0].String(),
							Coins:   sdk.NewCoins(sendCoin.Add(sendCoin)),
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
			expectedFeeAmount: feeAmt * 2,
		},
	}

	for _, c := range cases {
		suite.SetupTest(true) // setup
		tk := suite.app.TreasuryKeeper
		ak := suite.app.AccountKeeper
		bk := suite.app.BankKeeper
		dk := suite.app.DistrKeeper

		// Set burn split rate to 50%
		// fee amount should be 500, 50% of 10000
		tk.SetBurnSplitRate(suite.ctx, sdk.NewDecWithPrec(5, 1)) //50%

		feeCollector := ak.GetModuleAccount(suite.ctx, types.FeeCollectorName)
		burnModule := ak.GetModuleAccount(suite.ctx, treasurytypes.BurnModuleName)

		encodingConfig := suite.SetupEncoding()
		antehandler, err := customante.NewAnteHandler(
			customante.HandlerOptions{
				AccountKeeper:      ak,
				BankKeeper:         bk,
				FeegrantKeeper:     suite.app.FeeGrantKeeper,
				OracleKeeper:       suite.app.OracleKeeper,
				TreasuryKeeper:     suite.app.TreasuryKeeper,
				SigGasConsumer:     ante.DefaultSigVerificationGasConsumer,
				SignModeHandler:    encodingConfig.TxConfig.SignModeHandler(),
				IBCChannelKeeper:   suite.app.IBCKeeper.ChannelKeeper,
				DistributionKeeper: dk,
			},
		)
		suite.Require().NoError(err)

		fmt.Printf("CASE = %s \n", c.name)
		suite.ctx = suite.ctx.WithBlockHeight(ante.TaxPowerUpgradeHeight)
		suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[0].String())
		tk.AddBurnTaxExemptionAddress(suite.ctx, addrs[1].String())

		for i := 0; i < 4; i++ {
			fundCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, 1_000_000_000_000))
			acc := ak.NewAccountWithAddress(suite.ctx, addrs[i])
			suite.Require().NoError(acc.SetAccountNumber(uint64(i)))
			ak.SetAccount(suite.ctx, acc)
			bk.MintCoins(suite.ctx, minttypes.ModuleName, fundCoins)
			bk.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addrs[i], fundCoins)
		}

		// case 1 provides zero fee so not enough fee
		// case 2 provides enough fee
		feeCases := []int64{0, feeAmt}
		for i := 0; i < 1; i++ {
			feeAmount := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, feeCases[i]))
			gasLimit := testdata.NewTestGasLimit()
			suite.Require().NoError(suite.txBuilder.SetMsgs(c.msgCreator()...))
			suite.txBuilder.SetFeeAmount(feeAmount)
			suite.txBuilder.SetGasLimit(gasLimit)

			privs, accNums, accSeqs := []cryptotypes.PrivKey{privs[c.msgSigner]}, []uint64{uint64(c.msgSigner)}, []uint64{uint64(i)}
			tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
			suite.Require().NoError(err)

			feeCollectorBefore := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
			burnBefore := bk.GetBalance(suite.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
			communityBefore := dk.GetFeePool(suite.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
			supplyBefore := bk.GetSupply(suite.ctx, core.MicroSDRDenom)

			_, err = antehandler(suite.ctx, tx, false)
			if i == 0 && c.expectedFeeAmount != 0 {
				suite.Require().EqualError(err, fmt.Sprintf("insufficient fees; got: \"\", required: \"%dusdr\" = \"\"(gas) +\"%dusdr\"(stability): insufficient fee", c.expectedFeeAmount, c.expectedFeeAmount))
			} else {
				suite.Require().NoError(err)
			}

			feeCollectorAfter := bk.GetBalance(suite.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
			burnAfter := bk.GetBalance(suite.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
			communityAfter := dk.GetFeePool(suite.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
			supplyAfter := bk.GetSupply(suite.ctx, core.MicroSDRDenom)

			if i == 0 {
				suite.Require().Equal(feeCollectorBefore, feeCollectorAfter)
				suite.Require().Equal(burnBefore, burnAfter)
				suite.Require().Equal(communityBefore, communityAfter)
				suite.Require().Equal(supplyBefore, supplyAfter)
			}

			if i == 1 {
				suite.Require().Equal(feeCollectorBefore, feeCollectorAfter)
				splitAmount := sdk.NewInt(int64(float64(c.expectedFeeAmount) * 0.5))
				suite.Require().Equal(burnBefore, burnAfter.AddAmount(splitAmount))
				suite.Require().Equal(communityBefore, communityAfter.Add(sdk.NewDecFromInt(splitAmount)))
				suite.Require().Equal(supplyBefore, supplyAfter.SubAmount(splitAmount))
			}
		}
	}
}
