package ante_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/terra-money/core/app"

	core "github.com/terra-money/core/types"
	"github.com/terra-money/core/x/auth"
	"github.com/terra-money/core/x/auth/ante"
	"github.com/terra-money/core/x/bank"
	"github.com/terra-money/core/x/msgauth"
	oracleexported "github.com/terra-money/core/x/oracle/exported"
	"github.com/terra-money/core/x/treasury"
	"github.com/terra-money/core/x/wasm"
	wasmconfig "github.com/terra-money/core/x/wasm/config"
)

// returns context and app with params set on account keeper
func createTestApp() (*app.TerraApp, sdk.Context) {
	db := dbm.NewMemDB()

	tapp := app.NewTerraApp(log.NewNopLogger(), db, nil, true, 0, map[int64]bool{}, wasmconfig.DefaultConfig())
	ctx := tapp.BaseApp.NewContext(true, abci.Header{})
	tapp.GetTreasuryKeeper().SetParams(ctx, treasury.DefaultParams())

	return tapp, ctx
}

func TestEnsureMempoolFeesGas(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	mtd := ante.NewTaxFeeDecorator(tapp.GetTreasuryKeeper())
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// msg and signatures
	msg1 := types.NewTestMsg(addr1)
	fee := types.NewTestStdFee()

	msgs := []sdk.Msg{msg1}

	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	// Set high gas price so standard test fee fails
	lunaPrice := sdk.NewDecCoinFromDec(core.MicroLunaDenom, sdk.NewDec(200).Quo(sdk.NewDec(100000)))
	highGasPrice := []sdk.DecCoin{lunaPrice}
	ctx = ctx.WithMinGasPrices(highGasPrice)

	// Set IsCheckTx to true
	ctx = ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should have errored on too low fee for local gasPrice")

	// Set IsCheckTx to false
	ctx = ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check minGasPrice in DeliverTx
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "MempoolFeeDecorator returned error in DeliverTx")

	// Set IsCheckTx back to true for testing sufficient mempool fee
	ctx = ctx.WithIsCheckTx(true)

	lunaPrice = sdk.NewDecCoinFromDec(core.MicroLunaDenom, sdk.NewDec(0).Quo(sdk.NewDec(100000)))
	lowGasPrice := []sdk.DecCoin{lunaPrice}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice")
}

func TestEnsureMempoolFeesSend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msgs := []sdk.Msg{bank.NewMsgSend(addr1, addr1, sendCoins)}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should errored on low fee for local gasPrice + tax")

	expectedTax := tk.GetTaxRate(ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	fee.Amount = sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax))
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice + tax")
}

func TestEnsureMempoolFeesMultiSend(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msgs := []sdk.Msg{bank.NewMsgMultiSend(
		[]bank.Input{
			bank.NewInput(addr1, sendCoins),
			bank.NewInput(addr1, sendCoins),
		},
		[]bank.Output{
			bank.NewOutput(addr1, sendCoins.Add(sendCoins...)),
		},
	)}

	expectedTax := tk.GetTaxRate(ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	fee := auth.NewStdFee(100000, sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax)))
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should errored on low fee for local gasPrice + tax")

	fee.Amount = sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax.Add(expectedTax)))
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice + tax")
}

func TestEnsureMempoolFeesInstantiateContract(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msgs := []sdk.Msg{wasm.NewMsgInstantiateContract(addr1, 0, []byte{}, sendCoins, true)}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should errored on low fee for local gasPrice + tax")

	expectedTax := tk.GetTaxRate(ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	fee.Amount = sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax))
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice + tax")
}

func TestEnsureMempoolFeesExecuteContract(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))
	msgs := []sdk.Msg{wasm.NewMsgExecuteContract(addr1, addr1, []byte{}, sendCoins)}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should errored on low fee for local gasPrice + tax")

	expectedTax := tk.GetTaxRate(ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	fee.Amount = sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax))
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice + tax")
}

func TestEnsureMempoolFeesExecAuthorized(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	sendAmount := int64(1000000)
	sendCoins := sdk.NewCoins(sdk.NewInt64Coin(core.MicroSDRDenom, sendAmount))

	msgs := []sdk.Msg{msgauth.NewMsgExecAuthorized(addr1, []sdk.Msg{bank.NewMsgSend(addr1, addr1, sendCoins)})}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Decorator should errored on low fee for local gasPrice + tax")

	expectedTax := tk.GetTaxRate(ctx).MulInt64(sendAmount).TruncateInt()
	if taxCap := tk.GetTaxCap(ctx, core.MicroSDRDenom); expectedTax.GT(taxCap) {
		expectedTax = taxCap
	}

	fee.Amount = sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, expectedTax))
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err, "Decorator should not have errored on fee higher than local gasPrice + tax")
}

func TestEnsureNoMempoolFeesForOracleMessages(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{{Denom: "uusd", Amount: sdk.NewDec(100)}}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	tk := tapp.GetTreasuryKeeper()
	mtd := ante.NewTaxFeeDecorator(tk)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, _ := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	msgs := []sdk.Msg{oracleexported.MsgAggregateExchangeRatePrevote{}, oracleexported.MsgAggregateExchangeRateVote{}}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)
}
