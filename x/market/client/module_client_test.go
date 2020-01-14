package client

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/app"
)

const (
	storeKey = string("budget")
)

var (
	txCmdList = map[string]bool{
		"swap": true,
	}
)

func TestTxCmdInvariant(t *testing.T) {

	cdc := app.MakeCodec()
	mc := NewModuleClient(storeKey, cdc)

	for _, cmd := range mc.GetTxCmd().Commands() {
		_, ok := txCmdList[cmd.Name()]
		require.True(t, ok)
	}

	require.Equal(t, len(txCmdList), len(mc.GetTxCmd().Commands()))
}
