package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
	"github.com/terra-project/core/x/market"

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

func TestQuerySwap(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	querySwapCmd := GetCmdQuerySwap(cdc)

	// Name check
	require.Equal(t, market.QuerySwap, querySwapCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(querySwapCmd.Args))

	// Check Flags
	askDenomFlag := querySwapCmd.Flag(flagAskDenom)
	require.NotNil(t, askDenomFlag)
	require.Equal(t, []string{"true"}, askDenomFlag.Annotations[cobra.BashCompOneRequiredFlag])

	offerCoinFlag := querySwapCmd.Flag(flagOfferCoin)
	require.NotNil(t, offerCoinFlag)
	require.Equal(t, []string{"true"}, offerCoinFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParamsCmd := GetCmdQueryParams(cdc)

	// Name check
	require.Equal(t, market.QueryParams, queryParamsCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParamsCmd.Args))
}
