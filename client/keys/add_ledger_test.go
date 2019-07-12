//+build ledger,test_ledger_mock

package keys

import (
	"bufio"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/testutil"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/tests"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/stretchr/testify/assert"
)

func Test_runAddCmdLedger(t *testing.T) {
	testutil.PrepareCmdTest()

	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)
	viper.Set(client.FlagUseLedger, true)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	// Now enter password
	cleanUp1 := client.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\ntest1234\n")))
	defer cleanUp1()
	err := runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	// Now check that it has been stored properly
	kb, err := NewKeyBaseFromHomeFlag()
	assert.NoError(t, err)
	assert.NotNil(t, kb)
	key1, err := kb.Get("keyname1")
	assert.NoError(t, err)
	assert.NotNil(t, key1)

	assert.Equal(t, "keyname1", key1.GetName())
	assert.Equal(t, keys.TypeLedger, key1.GetType())
	assert.Equal(t,
		"terrapub1addwnpepqvpg7r26nl2pvqqern00m6s9uaax3hauu2rzg8qpjzq9hy6xve7sw0d84m6",
		sdk.MustBech32ifyAccPub(key1.GetPubKey()))

	viper.Set(client.FlagUseLedger, false)
}
