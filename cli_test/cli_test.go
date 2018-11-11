package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"
)

var (
	terradHome = ""
	terracliHome   = ""
)

func init() {
	terradHome, terracliHome = getTestingHomeDirs()
}

func TestInitStartSequence(t *testing.T) {
	os.RemoveAll(terradHome)
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	executeInit(t)
	executeStart(t, servAddr, port)
}

func executeInit(t *testing.T) {
	var (
		chainID string
		initRes map[string]json.RawMessage
	)
	_, stderr := tests.ExecuteT(t, fmt.Sprintf("terrad --home=%s --home-client=%s init --name=test", terradHome, terracliHome), app.DefaultKeyPass)
	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(t, err)
	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)
}

func executeStart(t *testing.T, servAddr, port string) {
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("terrad start --home=%s --rpc.laddr=%v", terradHome, servAddr))
	defer proc.Stop(false)
	tests.WaitForTMStart(port)
}

func getTestingHomeDirs() (string, string) {
	tmpDir := os.TempDir()
	terradHome := fmt.Sprintf("%s%s.test_terrad", tmpDir, string(os.PathSeparator))
	terracliHome := fmt.Sprintf("%s%s.test_terracli", tmpDir, string(os.PathSeparator))
	return terradHome, terracliHome
}
