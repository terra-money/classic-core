package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/terra-project/core/x/oracle"
)

func TestPricePrevoteTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	oracleTxCmd := &cobra.Command{
		Use:   "oracle",
		Short: "Oracle transaction subcommands",
	}

	txCmd.AddCommand(oracleTxCmd)

	oracleTxCmd.AddCommand(client.PostCommands(
		GetCmdPricePrevote(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`oracle`,
		`prevote`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--denom=ukrw`,
		`--price=5555.55`,
		`--salt=1234`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestPriceVoteTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	oracleTxCmd := &cobra.Command{
		Use:   "oracle",
		Short: "Oracle transaction subcommands",
	}

	txCmd.AddCommand(oracleTxCmd)

	oracleTxCmd.AddCommand(client.PostCommands(
		GetCmdPriceVote(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`oracle`,
		`vote`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--denom=ukrw`,
		`--price=5555.55`,
		`--salt=1234`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestDelegateFeederPermissionTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	oracleTxCmd := &cobra.Command{
		Use:   "oracle",
		Short: "Oracle transaction subcommands",
	}

	txCmd.AddCommand(oracleTxCmd)

	oracleTxCmd.AddCommand(client.PostCommands(
		GetCmdDelegateFeederPermission(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`oracle`,
		`set-feeder`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--feeder=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestGetCmdQueryPrice(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryPriceCmd := GetCmdQueryPrice(oracle.QuerierRoute, cdc)

	// Name check
	require.Equal(t, oracle.QueryPrice, queryPriceCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryPriceCmd.Args))

	// Check Flags
	denomFlag := queryPriceCmd.Flag(flagDenom)
	require.NotNil(t, denomFlag)
	require.Equal(t, []string{"true"}, denomFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestGetCmdQueryActive(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryActiveCmd := GetCmdQueryActive(oracle.QuerierRoute, cdc)

	// Name check
	require.Equal(t, oracle.QueryActive, queryActiveCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryActiveCmd.Args))
}

func TestGetCmdQueryVotes(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryVotesCmd := GetCmdQueryVotes(oracle.QuerierRoute, cdc)

	// Name check
	require.Equal(t, oracle.QueryVotes, queryVotesCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryVotesCmd.Args))

	// Check Flags
	denomFlag := queryVotesCmd.Flag(flagDenom)
	require.NotNil(t, denomFlag)
	require.Equal(t, []string{"true"}, denomFlag.Annotations[cobra.BashCompOneRequiredFlag])

	voterFlag := queryVotesCmd.Flag(flagValidator)
	require.NotNil(t, voterFlag)
}

func TestGetCmdQueryPrevotes(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryPrevotesCmd := GetCmdQueryPrevotes(oracle.QuerierRoute, cdc)

	// Name check
	require.Equal(t, oracle.QueryPrevotes, queryPrevotesCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryPrevotesCmd.Args))

	// Check Flags
	denomFlag := queryPrevotesCmd.Flag(flagDenom)
	require.NotNil(t, denomFlag)
	require.Equal(t, []string{"true"}, denomFlag.Annotations[cobra.BashCompOneRequiredFlag])

	voterFlag := queryPrevotesCmd.Flag(flagValidator)
	require.NotNil(t, voterFlag)
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParamsCmd := GetCmdQueryParams(oracle.QuerierRoute, cdc)

	// Name check
	require.Equal(t, queryParamsCmd.Name(), oracle.QueryParams)

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParamsCmd.Args))
}
