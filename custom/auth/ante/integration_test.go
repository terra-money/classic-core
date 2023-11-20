package ante_test

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	customante "github.com/classic-terra/core/v2/custom/auth/ante"
	core "github.com/classic-terra/core/v2/types"
	treasurytypes "github.com/classic-terra/core/v2/x/treasury/types"
)

// go test -v -run ^TestAnteTestSuite/TestIntegrationTaxExemption$ github.com/classic-terra/core/v2/custom/auth/ante
func (s *AnteTestSuite) TestIntegrationTaxExemption() {
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
	sendAmt := int64(1_000_000)
	sendCoin := sdk.NewInt64Coin(core.MicroSDRDenom, sendAmt)
	feeAmt := int64(1000)

	cases := []struct {
		name              string
		msgSigner         cryptotypes.PrivKey
		msgCreator        func() []sdk.Msg
		expectedFeeAmount int64
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
			expectedFeeAmount: 0,
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
			expectedFeeAmount: feeAmt,
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
			expectedFeeAmount: feeAmt,
		}, {
			name:      "MsgSend(exemption -> exemption), MsgMultiSend(exemption -> normal, exemption)",
			msgSigner: privs[0],
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
		s.SetupTest(true) // setup
		tk := s.app.TreasuryKeeper
		ak := s.app.AccountKeeper
		bk := s.app.BankKeeper
		dk := s.app.DistrKeeper

		// Set burn split rate to 50%
		// fee amount should be 500, 50% of 10000
		burnSplitRate := sdk.NewDecWithPrec(5, 1)
		tk.SetBurnSplitRate(s.ctx, burnSplitRate) // 50%

		feeCollector := ak.GetModuleAccount(s.ctx, types.FeeCollectorName)
		burnModule := ak.GetModuleAccount(s.ctx, treasurytypes.BurnModuleName)

		encodingConfig := s.SetupEncoding()
		wasmConfig := wasmtypes.DefaultWasmConfig()
		antehandler, err := customante.NewAnteHandler(
			customante.HandlerOptions{
				AccountKeeper:      ak,
				BankKeeper:         bk,
				FeegrantKeeper:     s.app.FeeGrantKeeper,
				OracleKeeper:       s.app.OracleKeeper,
				TreasuryKeeper:     s.app.TreasuryKeeper,
				SigGasConsumer:     ante.DefaultSigVerificationGasConsumer,
				SignModeHandler:    encodingConfig.TxConfig.SignModeHandler(),
				IBCKeeper:          *s.app.IBCKeeper,
				DistributionKeeper: dk,
				WasmConfig:         &wasmConfig,
				TXCounterStoreKey:  s.app.GetKey(wasmtypes.StoreKey),
			},
		)
		s.Require().NoError(err)

		for i := 0; i < 4; i++ {
			coins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, 1_000_000))
			testutil.FundAccount(s.app.BankKeeper, s.ctx, addrs[i], coins)
		}

		s.txBuilder = s.clientCtx.TxConfig.NewTxBuilder()

		tk.AddBurnTaxExemptionAddress(s.ctx, addrs[0].String())
		tk.AddBurnTaxExemptionAddress(s.ctx, addrs[1].String())

		s.Run(c.name, func() {
			// case 1 provides zero fee so not enough fee
			// case 2 provides enough fee
			feeCases := []int64{0, feeAmt}
			for i := 0; i < 1; i++ {
				feeAmount := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, feeCases[i]))
				gasLimit := testdata.NewTestGasLimit()
				s.Require().NoError(s.txBuilder.SetMsgs(c.msgCreator()...))
				s.txBuilder.SetFeeAmount(feeAmount)
				s.txBuilder.SetGasLimit(gasLimit)

				privs, accNums, accSeqs := []cryptotypes.PrivKey{c.msgSigner}, []uint64{3}, []uint64{0}
				tx, err := s.CreateTestTx(privs, accNums, accSeqs, s.ctx.ChainID())
				s.Require().NoError(err)

				// set zero gas prices
				s.ctx = s.ctx.WithMinGasPrices(sdk.NewDecCoins())

				feeCollectorBefore := bk.GetBalance(s.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
				burnBefore := bk.GetBalance(s.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
				communityBefore := dk.GetFeePool(s.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
				supplyBefore := bk.GetSupply(s.ctx, core.MicroSDRDenom)

				_, err = antehandler(s.ctx, tx, false)
				if i == 0 && c.expectedFeeAmount != 0 {
					s.Require().EqualError(err, fmt.Sprintf(
						"insufficient fees; got: \"\", required: \"%dusdr\" = \"\"(gas) + \"%dusdr\"(stability): insufficient fee",
						c.expectedFeeAmount, c.expectedFeeAmount))
				} else {
					s.Require().NoError(err)
				}

				feeCollectorAfter := bk.GetBalance(s.ctx, feeCollector.GetAddress(), core.MicroSDRDenom)
				burnAfter := bk.GetBalance(s.ctx, burnModule.GetAddress(), core.MicroSDRDenom)
				communityAfter := dk.GetFeePool(s.ctx).CommunityPool.AmountOf(core.MicroSDRDenom)
				supplyAfter := bk.GetSupply(s.ctx, core.MicroSDRDenom)

				if i == 0 {
					s.Require().Equal(feeCollectorBefore, feeCollectorAfter)
					s.Require().Equal(burnBefore, burnAfter)
					s.Require().Equal(communityBefore, communityAfter)
					s.Require().Equal(supplyBefore, supplyAfter)
				}

				if i == 1 {
					s.Require().Equal(feeCollectorBefore, feeCollectorAfter)
					splitAmount := burnSplitRate.MulInt64(c.expectedFeeAmount).TruncateInt()
					s.Require().Equal(burnBefore, burnAfter.AddAmount(splitAmount))
					s.Require().Equal(communityBefore, communityAfter.Add(sdk.NewDecFromInt(splitAmount)))
					s.Require().Equal(supplyBefore, supplyAfter.SubAmount(splitAmount))
				}
			}
		})
	}
}
