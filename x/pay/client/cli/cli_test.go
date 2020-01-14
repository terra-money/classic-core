package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
)

func TestSendTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	txCmd.AddCommand(SendTxCmd(cdc))

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`send`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--coins=1000uluna`,
		`--to=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}
