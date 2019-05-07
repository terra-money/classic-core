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
		"params":               true,
		"tax-rate":             true,
		"tax-cap":              true,
		"reward-weight":        true,
		"seigniorage-proceeds": true,
		"active-claims":        true,
		"current-epoch":        true,
		"issuance":             true,
		"tax-proceeds":         true,
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
