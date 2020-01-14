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
	queryCmdList = map[string]bool{
		"params":                     true,
		"validator":                  true,
		"validators":                 true,
		"unbonding-delegations-from": true,
		"redelegations-from":         true,
		"delegation":                 true,
		"delegations":                true,
		"delegations-to":             true,
		"unbonding-delegation":       true,
		"unbonding-delegations":      true,
		"redelegation":               true,
		"redelegations":              true,
		"pool":                       true,
	}

	txCmdList = map[string]bool{
		"create-validator": true,
		"edit-validator":   true,
		"delegate":         true,
		"redelegate":       true,
		"unbond":           true,
	}
)

func TestQueryCmdInvariant(t *testing.T) {

	cdc := app.MakeCodec()
	mc := NewModuleClient(storeKey, cdc)

	for _, cmd := range mc.GetQueryCmd().Commands() {
		_, ok := queryCmdList[cmd.Name()]
		require.True(t, ok)
	}

	require.Equal(t, len(queryCmdList), len(mc.GetQueryCmd().Commands()))
}

func TestTxCmdInvariant(t *testing.T) {

	cdc := app.MakeCodec()
	mc := NewModuleClient(storeKey, cdc)

	for _, cmd := range mc.GetTxCmd().Commands() {
		_, ok := txCmdList[cmd.Name()]
		require.True(t, ok)
	}

	require.Equal(t, len(txCmdList), len(mc.GetTxCmd().Commands()))
}
