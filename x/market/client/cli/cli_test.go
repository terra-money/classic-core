package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"

	"github.com/cosmos/cosmos-sdk/client"
)

func TestSwapTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	marketTxCmd := &cobra.Command{
		Use:   "market",
		Short: "Market transaction subcommands",
	}

	txCmd.AddCommand(marketTxCmd)

	marketTxCmd.AddCommand(client.PostCommands(
		GetSwapCmd(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`market`,
		`swap`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--offer-coin=1000uluna`,
		`--ask-denom=ukrw`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}
