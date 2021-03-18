package ante_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/terra-project/core/x/auth/ante"
)

func TestEnsureSoftforkGasCheck(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	_, ctx := createTestApp()

	// setup
	spd := ante.NewSpammingPreventionDecorator()
	antehandler := sdk.ChainAnteDecorators(spd)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// msg and signatures
	msg1 := types.NewTestMsg(addr1)
	fee := types.NewTestStdFee()
	fee.Gas = 100000000

	msgs := []sdk.Msg{msg1}

	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	// Set IsCheckTx to true
	ctx = ctx.WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(ctx, tx, false)
	require.Error(t, err, "Decorator should have errored on too high gas for local gasPrice")

	// Set IsCheckTx to false
	ctx = ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check minGasPrice in DeliverTx
	_, err = antehandler(ctx, tx, false)
	require.Error(t, err, "Decorator should have errored on too high gas for local gasPrice")

	// Set ChainID to columbus-4 and height to before fork
	ctx = ctx.WithChainID("columbus-4")
	ctx = ctx.WithBlockHeight(2379999)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err, "SpammingPreventionDecorator returned error in DeliverTx")

	// Set height to after fork
	ctx = ctx.WithBlockHeight(2380000)

	_, err = antehandler(ctx, tx, false)
	require.Error(t, err, "Decorator should have errored on high gas than hard cap")

	ctx = ctx.WithBlockHeight(2380001)

	_, err = antehandler(ctx, tx, false)
	require.Error(t, err, "Decorator should have errored on high gas than hard cap")
}
