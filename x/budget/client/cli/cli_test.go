package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
	"github.com/terra-project/core/x/budget"

	"github.com/cosmos/cosmos-sdk/client"
)

func TestSubmitProgramTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	budgetTxCmd := &cobra.Command{
		Use:   "budget",
		Short: "budget transaction subcommands",
	}

	txCmd.AddCommand(budgetTxCmd)

	budgetTxCmd.AddCommand(client.PostCommands(
		GetCmdSubmitProgram(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`budget`,
		`submit-program`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--title=testprogram`,
		`--description=testprogramtestprogram`,
		`--executor=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestVoteTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	budgetTxCmd := &cobra.Command{
		Use:   "budget",
		Short: "budget transaction subcommands",
	}

	txCmd.AddCommand(budgetTxCmd)

	budgetTxCmd.AddCommand(client.PostCommands(
		GetCmdVote(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`budget`,
		`vote`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--program-id=1`,
		`--option=yes`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestWithdrawProgramTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	budgetTxCmd := &cobra.Command{
		Use:   "budget",
		Short: "budget transaction subcommands",
	}

	txCmd.AddCommand(budgetTxCmd)

	budgetTxCmd.AddCommand(client.PostCommands(
		GetCmdWithdrawProgram(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`budget`,
		`withdraw`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--program-id=1`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestQueryProgram(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryProgramCmd := GetCmdQueryProgram(cdc)

	// Name check
	require.Equal(t, budget.QueryProgram, queryProgramCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryProgramCmd.Args))

	// Check Flags
	programFlag := queryProgramCmd.Flag(flagProgramID)
	require.NotNil(t, programFlag)
	require.Equal(t, []string{"true"}, programFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryActives(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryActivesCmd := GetCmdQueryActives(cdc)

	// Name check
	require.Equal(t, budget.QueryActiveList, queryActivesCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryActivesCmd.Args))
}

func TestQueryCandidates(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryCandidatesCmd := GetCmdQueryCandidates(cdc)

	// Name check
	require.Equal(t, queryCandidatesCmd.Name(), budget.QueryCandidateList)

	// NoArg check
	require.Equal(t, testutil.FS(queryCandidatesCmd.Args), testutil.FS(cobra.PositionalArgs(cobra.NoArgs)))
}

func TestQueryVotes(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryVotesCmd := GetCmdQueryVotes(cdc)

	// Name check
	require.Equal(t, budget.QueryVotes, queryVotesCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryVotesCmd.Args))

	// Check Flags
	programFlag := queryVotesCmd.Flag(flagProgramID)
	require.NotNil(t, programFlag)

	voterFlag := queryVotesCmd.Flag(flagVoter)
	require.NotNil(t, voterFlag)
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParamsCmd := GetCmdQueryParams(cdc)

	// Name check
	require.Equal(t, budget.QueryParams, queryParamsCmd.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParamsCmd.Args))
}
