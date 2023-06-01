package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/suite"

	terraapp "github.com/classic-terra/core/v2/app"
	appparams "github.com/classic-terra/core/v2/app/params"
	ante "github.com/classic-terra/core/v2/x/feeshare/ante"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	feesharetypes "github.com/classic-terra/core/v2/x/feeshare/types"
	treasurytypes "github.com/classic-terra/core/v2/x/treasury/types"
)

var emptyWasmOpts []wasm.Option

type AnteTestSuite struct {
	suite.Suite

	app       *terraapp.TerraApp
	ctx       sdk.Context
	clientCtx client.Context
	txBuilder client.TxBuilder
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, tempDir string) (*terraapp.TerraApp, sdk.Context) {
	app := terraapp.NewTerraApp(
		log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		tempDir, simapp.FlagPeriodValue, terraapp.MakeEncodingConfig(),
		simapp.EmptyAppOptions{}, emptyWasmOpts,
	)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.TreasuryKeeper.SetParams(ctx, treasurytypes.DefaultParams())
	app.DistrKeeper.SetParams(ctx, distributiontypes.DefaultParams())
	app.DistrKeeper.SetFeePool(ctx, distributiontypes.InitialFeePool())
	app.FeeShareKeeper.SetParams(ctx, feesharetypes.DefaultParams())

	return app, ctx
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func (suite *AnteTestSuite) SetupTest(isCheckTx bool) {
	tempDir := suite.T().TempDir()
	suite.app, suite.ctx = createTestApp(isCheckTx, tempDir)
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Set up TxConfig.
	encodingConfig := suite.SetupEncoding()

	suite.clientCtx = client.Context{}.
		WithTxConfig(encodingConfig.TxConfig)
}

func (suite *AnteTestSuite) SetupEncoding() simappparams.EncodingConfig {
	encodingConfig := simapp.MakeTestEncodingConfig()
	// We're using TestMsg encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64, chainID string) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			suite.txBuilder, priv, suite.clientCtx.TxConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.txBuilder.GetTx(), nil
}

func TestAnteSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (suite *AnteTestSuite) TestFeeLogic() {
	// We expect all to pass
	feeCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250)))

	testCases := []struct {
		name               string
		incomingFee        sdk.Coins
		govPercent         sdk.Dec
		numContracts       int
		expectedFeePayment sdk.Coins
	}{
		{
			"100% fee / 1 contract",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			1,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(500)), sdk.NewCoin("utoken", sdk.NewInt(250))),
		},
		{
			"100% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			2,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"100% fee / 10 contracts",
			feeCoins,
			sdk.NewDecWithPrec(100, 2),
			10,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(50)), sdk.NewCoin("utoken", sdk.NewInt(25))),
		},
		{
			"67% fee / 7 contracts",
			feeCoins,
			sdk.NewDecWithPrec(67, 2),
			7,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(48)), sdk.NewCoin("utoken", sdk.NewInt(24))),
		},
		{
			"50% fee / 1 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			1,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(250)), sdk.NewCoin("utoken", sdk.NewInt(125))),
		},
		{
			"50% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			2,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(125)), sdk.NewCoin("utoken", sdk.NewInt(62))),
		},
		{
			"50% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(50, 2),
			3,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(83)), sdk.NewCoin("utoken", sdk.NewInt(42))),
		},
		{
			"25% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(25, 2),
			2,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(62)), sdk.NewCoin("utoken", sdk.NewInt(31))),
		},
		{
			"15% fee / 3 contracts",
			feeCoins,
			sdk.NewDecWithPrec(15, 2),
			3,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(25)), sdk.NewCoin("utoken", sdk.NewInt(12))),
		},
		{
			"1% fee / 2 contracts",
			feeCoins,
			sdk.NewDecWithPrec(1, 2),
			2,
			sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, sdk.NewInt(2)), sdk.NewCoin("utoken", sdk.NewInt(1))),
		},
	}

	for _, tc := range testCases {
		coins := ante.FeePayLogic(tc.incomingFee, tc.govPercent, tc.numContracts)

		for _, coin := range coins {
			for _, expectedCoin := range tc.expectedFeePayment {
				if coin.Denom == expectedCoin.Denom {
					suite.Require().Equal(expectedCoin.Amount.Int64(), coin.Amount.Int64(), tc.name)
				}
			}
		}
	}
}

func (suite *AnteTestSuite) TestFeeSharePayout() {
	suite.SetupTest(true) // setup
	require := suite.Require()
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	fsk := suite.app.FeeShareKeeper
	bk := suite.app.BankKeeper
	mfd := ante.NewFeeSharePayoutDecorator(bk, fsk)
	antehandler := sdk.ChainAnteDecorators(mfd)

	// Set the blockheight past the tax height block
	suite.ctx = suite.ctx.WithBlockHeight(10000000)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}

	// Create msg for tx
	msg := &wasmtypes.MsgExecuteContract{
		Sender:   addr1.String(),
		Contract: addr1.String(),
		Msg:      []byte(`{"send":{}}`),
		Funds:    sdk.NewCoins(sdk.NewCoin("utoken", sdk.NewInt(100))),
	}

	require.NoError(suite.txBuilder.SetMsgs(msg))
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	require.NoError(err)

	// send tx to FeeSharePayoutDecorator antehandler
	_, err = antehandler(suite.ctx, tx, false)
	require.NoError(err)
}
