package keys

import (
	"bufio"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"
)

func Test_runDeleteCmd(t *testing.T) {
	deleteKeyCommand := deleteKeyCommand()

	yesF, _ := deleteKeyCommand.Flags().GetBool(flagYes)
	forceF, _ := deleteKeyCommand.Flags().GetBool(flagForce)
	assert.False(t, yesF)
	assert.False(t, forceF)

	fakeKeyName1 := "runDeleteCmd_Key1"
	fakeKeyName2 := "runDeleteCmd_Key2"
	fakeKeyName3 := "runDeleteCmd_Key3"
	fakeKeyName4 := "runDeleteCmd_Key4"

	// Now add a temporary keybase
	kbHome, cleanUp := tests.NewTestCaseDir(t)
	defer cleanUp()
	viper.Set(cli.HomeFlag, kbHome)

	// Now
	kb, err := NewKeyBaseFromHomeFlag()
	assert.NoError(t, err)
	_, err = kb.CreateAccount(fakeKeyName1, tests.TestMnemonic, "", "", 0, 0)
	assert.NoError(t, err)
	_, err = kb.CreateAccount(fakeKeyName2, tests.TestMnemonic, "", "", 0, 1)
	assert.NoError(t, err)
	_, err = kb.CreateAccount(fakeKeyName3, tests.TestMnemonic, "", "test1234", 0, 0)
	assert.NoError(t, err)
	_, err = kb.CreateAccount(fakeKeyName4, tests.TestMnemonic, "", "test1234", 0, 0)
	assert.NoError(t, err)

	err = runDeleteCmd(deleteKeyCommand, []string{"blah"})
	require.Error(t, err)
	require.Equal(t, "Key blah not found", err.Error())

	// User confirmation missing
	err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName1})
	require.Error(t, err)
	require.Equal(t, "EOF", err.Error())

	{
		_, err = kb.Get(fakeKeyName1)
		require.NoError(t, err)

		// Now there is a confirmation
		cleanUp := client.OverrideStdin(bufio.NewReader(strings.NewReader("y\n")))
		defer cleanUp()
		err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName1})
		require.NoError(t, err)

		_, err = kb.Get(fakeKeyName1)
		require.Error(t, err) // Key1 is gone
	}

	viper.Set(flagYes, true)
	_, err = kb.Get(fakeKeyName2)
	require.NoError(t, err)
	err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName2})
	require.NoError(t, err)
	_, err = kb.Get(fakeKeyName2)
	require.Error(t, err) // Key2 is gone

	// Invalid Password
	viper.Set(flagYes, false)
	_, err = kb.Get(fakeKeyName3)
	require.NoError(t, err)
	cleanUp2 := client.OverrideStdin(bufio.NewReader(strings.NewReader("invalid\n")))
	defer cleanUp2()
	err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName3})
	require.Error(t, err)
	_, err = kb.Get(fakeKeyName3)
	require.NoError(t, err) // Key3 is not gone

	// Valid Password
	cleanUp3 := client.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\n")))
	defer cleanUp3()
	err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName3})
	require.NoError(t, err)
	_, err = kb.Get(fakeKeyName3)
	require.Error(t, err) // Key3 is gone

	// Force Delete
	viper.Set(flagForce, true)
	_, err = kb.Get(fakeKeyName4)
	require.NoError(t, err)
	err = runDeleteCmd(deleteKeyCommand, []string{fakeKeyName4})
	require.NoError(t, err)
	_, err = kb.Get(fakeKeyName4)
	require.Error(t, err) // Key4 is gone

	// TODO: Write another case for !keys.Local
}

func Test_confirmDeletion(t *testing.T) {
	type args struct {
		buf *bufio.Reader
	}

	answerYes := bufio.NewReader(strings.NewReader("y\n"))
	answerYes2 := bufio.NewReader(strings.NewReader("Y\n"))
	answerNo := bufio.NewReader(strings.NewReader("n\n"))
	answerInvalid := bufio.NewReader(strings.NewReader("245\n"))

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Y", args{answerYes}, false},
		{"y", args{answerYes2}, false},
		{"N", args{answerNo}, true},
		{"BAD", args{answerInvalid}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := confirmDeletion(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("confirmDeletion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
