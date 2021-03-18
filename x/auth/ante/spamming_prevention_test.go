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
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/core/x/auth/ante"
	oracleexported "github.com/terra-project/core/x/oracle/exported"
)

func TestOracleSpamming(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "wasmtest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	viper.Set(flags.FlagHome, tempDir)

	// setup
	tapp, ctx := createTestApp()

	lowGasPrice := []sdk.DecCoin{}
	ctx = ctx.WithMinGasPrices(lowGasPrice)

	ok := tapp.GetOracleKeeper()
	mtd := ante.NewSpammingPreventionDecorator(ok)
	antehandler := sdk.ChainAnteDecorators(mtd)

	// keys and addresses
	priv1, _, _ := types.KeyTestPubAddr()
	privs, accNums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}

	msgs := []sdk.Msg{oracleexported.MsgAggregateExchangeRatePrevote{}, oracleexported.MsgAggregateExchangeRateVote{}}

	fee := auth.NewStdFee(100000, sdk.NewCoins())
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Error(t, err)

	msgs = []sdk.Msg{oracleexported.MsgAggregateExchangeRatePrevote{}, oracleexported.MsgAggregateExchangeRateVote{}}

	fee = auth.NewStdFee(100000, sdk.NewCoins())
	tx = types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)
	_, err = antehandler(ctx, tx, false)
	require.Error(t, err)
}
