package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
)

func TestSendTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	sendTxCmd := SendTxCmd(cdc)
	txCmd.AddCommand(sendTxCmd)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`send`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--to=terra12c5s58hnc3c0pjr5x7u68upsgzg2r8fwq5nlsy`,
		`--coins=1000000uluna`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}
