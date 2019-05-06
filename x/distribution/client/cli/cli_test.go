package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
	"github.com/terra-project/core/x/budget"

	"github.com/cosmos/cosmos-sdk/client"
	dt "github.com/cosmos/cosmos-sdk/x/distribution"
	dtk "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
)

func TestWithdrawRewardsTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	distTxCmd := &cobra.Command{
		Use:   "distr",
		Short: "Distribution transactions subcommands",
	}

	txCmd.AddCommand(distTxCmd)

	distTxCmd.AddCommand(client.PostCommands(
		GetCmdWithdrawRewards(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`distr`,
		`withdraw-rewards`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--validator=terravaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldcl4phj`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestSetWithdrawAddrTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	distTxCmd := &cobra.Command{
		Use:   "distr",
		Short: "Distribution transactions subcommands",
	}

	txCmd.AddCommand(distTxCmd)

	distTxCmd.AddCommand(client.PostCommands(
		GetCmdSetWithdrawAddr(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`distr`,
		`set-withdraw-addr`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--withdraw-to=terra13pqzy3n7ekfnpt9gmk9xtulzl49qw7td0hrsgh`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestWithdrawAllRewardsTx(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	withdrawAllRewardsCmd := GetCmdWithdrawAllRewards(cdc, dt.QuerierRoute)

	// Name check
	require.Equal(t, withdrawAllRewardsCmd.Name(), "withdraw-all-rewards")

	// NoArg check
	require.Equal(t, testutil.FS(withdrawAllRewardsCmd.Args), testutil.FS(cobra.PositionalArgs(cobra.NoArgs)))
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParamsCmd := GetCmdQueryParams(dt.QuerierRoute, cdc)

	// Name check
	require.Equal(t, queryParamsCmd.Name(), budget.QueryParams)

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParamsCmd.Args))
}

func TestQueryValidatorOutstandingRewards(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorOutstandingRewardsCmd := GetCmdQueryValidatorOutstandingRewards(dt.QuerierRoute, cdc)

	// Name check
	require.Equal(t, dtk.QueryValidatorOutstandingRewards, strings.ReplaceAll(queryValidatorOutstandingRewardsCmd.Name(), "-", "_"))

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorOutstandingRewardsCmd.Args))

	// Check Flags
	validatorFlag := queryValidatorOutstandingRewardsCmd.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryValidatorCommission(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorCommission := GetCmdQueryValidatorCommission(dt.QuerierRoute, cdc)

	// Name check
	require.Equal(t, dtk.QueryValidatorCommission, strings.ReplaceAll(queryValidatorCommission.Name(), "-", "_"))

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorCommission.Args))

	// Check Flags
	validatorFlag := queryValidatorCommission.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryValidatorSlashes(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorSlashes := GetCmdQueryValidatorSlashes(dt.QuerierRoute, cdc)

	// Name check
	require.Equal(t, dtk.QueryValidatorSlashes, strings.ReplaceAll(queryValidatorSlashes.Name(), "-", "_"))

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorSlashes.Args))

	// Check Flags
	validatorFlag := queryValidatorSlashes.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])

	startHeightFlag := queryValidatorSlashes.Flag(flagStartHeight)
	require.NotNil(t, startHeightFlag)
	require.Equal(t, []string{"true"}, startHeightFlag.Annotations[cobra.BashCompOneRequiredFlag])

	endHeightFlag := queryValidatorSlashes.Flag(flagEndHeight)
	require.NotNil(t, endHeightFlag)
	require.Equal(t, []string{"true"}, endHeightFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryDelegatorRewards(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryDelegatroRewards := GetCmdQueryDelegatorRewards(dt.QuerierRoute, cdc)

	// Name check
	require.Equal(t, dtk.QueryDelegationRewards, "delegation_"+queryDelegatroRewards.Name())
	require.Equal(t, dtk.QueryDelegatorTotalRewards, "delegator_total_"+queryDelegatroRewards.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryDelegatroRewards.Args))

	// Check Flags
	validatorFlag := queryDelegatroRewards.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)

	delegatorFlag := queryDelegatroRewards.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}
