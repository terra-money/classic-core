package keys

import (
	"bufio"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"

	"github.com/stretchr/testify/assert"

	"github.com/terra-project/core/testutil"
)

func Test_runAddCmdBasic(t *testing.T) {
	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	// Now enter password
	cleanUp1 := client.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\n")))
	defer cleanUp1()
	err := runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	/// Test Text - Replace? >> FAIL
	viper.Set(cli.OutputFlag, OutputFormatText)
	// Now enter password
	cleanUp2 := client.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\n")))
	defer cleanUp2()
	err = runAddCmd(cmd, []string{"keyname1"})
	assert.Error(t, err)

	/// Test Text - Replace? Answer >> PASS
	viper.Set(cli.OutputFlag, OutputFormatText)
	// Now enter password
	cleanUp3 := client.OverrideStdin(bufio.NewReader(strings.NewReader("y\ntest1234\n")))
	defer cleanUp3()
	err = runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	// Check JSON
	viper.Set(cli.OutputFlag, OutputFormatJSON)
	// Now enter password
	cleanUp4 := client.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\n")))
	defer cleanUp4()
	err = runAddCmd(cmd, []string{"keyname2"})
	assert.NoError(t, err)

	recoverInitialViperState()
}

func Test_runnAddCmdDryRun(t *testing.T) {
	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	viper.Set(flagDryRun, true)

	keyName := "keyname1"

	// Without password
	err := runAddCmd(cmd, []string{keyName})
	assert.NoError(t, err)

	// dry-run will not make any key info
	_, err = GetKeyInfo(keyName)
	assert.Error(t, err)

	recoverInitialViperState()
}

func Test_runAddCmdRecover(t *testing.T) {
	testutil.PrepareCmdTest()

	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	viper.Set(flagRecover, true)
	viper.Set(flagOldHdPath, false)

	keyName := "keyname1"
	password := "test1234\n"
	mnemonic := "candy hint hamster cute inquiry bright industry decide assist wedding carpet fiber arm menu machine lottery type alert fan march argue adapt recycle stomach\n"

	// New HD Path
	cleanUp1 := client.OverrideStdin(bufio.NewReader(strings.NewReader(password + mnemonic)))
	defer cleanUp1()

	err := runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	info, err := GetKeyInfo(keyName)
	assert.NoError(t, err)
	assert.Equal(t, "terra1wxuq9hkt4kes7r9kxh953l7p2cpcw8l73ek5dg", info.GetAddress().String())

	// Old HD Path
	viper.Set(flagOldHdPath, true)
	cleanUp2 := client.OverrideStdin(bufio.NewReader(strings.NewReader("y\n" + password + mnemonic)))
	defer cleanUp2()

	err = runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	info, err = GetKeyInfo(keyName)
	assert.NoError(t, err)
	assert.Equal(t, "terra1gaczd45crhwfa4x05k9747cuxwfmnduvmtyefs", info.GetAddress().String())

	// recover with dry-run flag (default password)
	viper.Set(flagDryRun, true)
	viper.Set(flagOldHdPath, true)

	cleanUp3 := client.OverrideStdin(bufio.NewReader(strings.NewReader(mnemonic)))
	defer cleanUp3()

	err = runAddCmd(cmd, []string{"keyname1"})
	assert.NoError(t, err)

	recoverInitialViperState()
}

func Test_runAddCmdPubkeyAndMultisig(t *testing.T) {
	testutil.PrepareCmdTest()

	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	/// Public Key Test
	viper.Set(cli.OutputFlag, OutputFormatText)
	viper.Set(flagRecover, true)

	pubkey1 := "terrapub1addwnpepqtmg9m7jy8xxqwnq05xh2rymsfph0mrfhzuz2lae3k09sn7qqwew7cgk76c"
	pubkey2 := "terrapub1addwnpepqdn2knqsda3zxq4uv24yg5wp97e48sxdhuqyplmpya5eeujlm5zk5chrdt8"
	pubkey3 := "terrapub1addwnpepqtycrza0rc9lk288gk9epwhdmft95t737vrctu75vp7h39l9rh24vxag39p"

	keyName1 := "keyname1"
	keyName2 := "keyname2"
	keyName3 := "keyname3"

	// Invalid Public Key
	viper.Set(FlagPublicKey, "invalid")

	err := runAddCmd(cmd, []string{keyName1})
	assert.Error(t, err)

	// Valid
	viper.Set(FlagPublicKey, pubkey1)

	err = runAddCmd(cmd, []string{keyName1})
	assert.NoError(t, err)

	info, err := GetKeyInfo(keyName1)
	assert.NoError(t, err)
	assert.Equal(t, "terra18smrf782hvjjeu3am06flc7nge2xvf8f2426q4", info.GetAddress().String())

	viper.Set(FlagPublicKey, pubkey2)

	err = runAddCmd(cmd, []string{keyName2})
	assert.NoError(t, err)

	info, err = GetKeyInfo(keyName2)
	assert.NoError(t, err)
	assert.Equal(t, "terra1mzm0v94uchdufn806hxxzu6q4m3xclx2yzpdv8", info.GetAddress().String())

	viper.Set(FlagPublicKey, pubkey3)

	err = runAddCmd(cmd, []string{keyName3})
	assert.NoError(t, err)

	info, err = GetKeyInfo(keyName3)
	assert.NoError(t, err)
	assert.Equal(t, "terra1yr0sqzfraffdwv9c33gx2dqhcl525muheuefaf", info.GetAddress().String())

	// Multisig Test
	keyNameMultisig := "keyNameMultisig"
	viper.Set(FlagPublicKey, "")
	viper.Set(flagMultisig, keyName1+" "+keyName2+" "+keyName3)

	// Invalid Threashold
	viper.Set(flagMultiSigThreshold, 0)
	err = runAddCmd(cmd, []string{keyNameMultisig})
	assert.Error(t, err)

	// Invalid Key Name
	viper.Set(flagMultiSigThreshold, 2)
	viper.Set(flagMultisig, keyName1+" "+keyName2+" hihi")
	err = runAddCmd(cmd, []string{keyNameMultisig})
	assert.Error(t, err)

	// Valid
	viper.Set(flagMultisig, keyName1+" "+keyName2+" "+keyName3)
	err = runAddCmd(cmd, []string{keyNameMultisig})
	assert.NoError(t, err)

	info, err = GetKeyInfo(keyNameMultisig)
	assert.NoError(t, err)
	assert.Equal(t, "terra1tswgrqcdauaw06dxeycj8ctr5etlah6aqg7elm", info.GetAddress().String())

	recoverInitialViperState()
}

func Test_runAddCmdInteractive(t *testing.T) {
	testutil.PrepareCmdTest()

	cmd := addKeyCommand()
	assert.NotNil(t, cmd)

	// Prepare a keybase
	kbHome, kbCleanUp := tests.NewTestCaseDir(t)
	assert.NotNil(t, kbHome)
	defer kbCleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	/// Test Text
	viper.Set(cli.OutputFlag, OutputFormatText)
	viper.Set(flagInteractive, true)
	viper.Set(flagOldHdPath, false)

	keyName := "keyname1"
	password := "test1234\n"
	mnemonic := "candy hint hamster cute inquiry bright industry decide assist wedding carpet fiber arm menu machine lottery type alert fan march argue adapt recycle stomach\n"
	bip39Passphrase := "hihi\nhihi\n"

	// New HD path
	cleanUp1 := client.OverrideStdin(bufio.NewReader(strings.NewReader(password + mnemonic + bip39Passphrase)))
	defer cleanUp1()

	err := runAddCmd(cmd, []string{keyName})
	assert.NoError(t, err)

	info, err := GetKeyInfo(keyName)
	assert.NoError(t, err)
	assert.Equal(t, "terra1smea3fuwun5ggfjep25gd7yv8kvw3mvx2hw3zm", info.GetAddress().String())

	viper.Set(flagOldHdPath, true)
	// Old HD path
	cleanUp2 := client.OverrideStdin(bufio.NewReader(strings.NewReader("y\n" + password + mnemonic + bip39Passphrase)))
	defer cleanUp2()

	err = runAddCmd(cmd, []string{keyName})
	assert.NoError(t, err)

	info, err = GetKeyInfo(keyName)
	assert.NoError(t, err)
	assert.Equal(t, "terra1nv4nsd7tfl8xc2dm7rry5exwcf350wjguk0x2c", info.GetAddress().String())

	recoverInitialViperState()
}

func recoverInitialViperState() {
	viper.Set(flagInteractive, false)
	viper.Set(flagRecover, false)
	viper.Set(FlagPublicKey, "")
	viper.Set(flagMultisig, "")
	viper.Set(flagMultiSigThreshold, 0)
	viper.Set(flagDryRun, false)
}
