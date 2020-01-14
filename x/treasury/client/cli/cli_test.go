package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
	"github.com/terra-project/core/x/treasury"
)

func TestQueryTaxRate(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryTaxRate := GetCmdQueryTaxRate(cdc)

	// Name check
	require.Equal(t, treasury.QueryTaxRate, queryTaxRate.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryTaxRate.Args))

	// Check Flags
	epochFlag := queryTaxRate.Flag(flagEpoch)
	require.NotNil(t, epochFlag)
}

func TestQueryTaxCap(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryTaxCap := GetCmdQueryTaxCap(cdc)

	// Name check
	require.Equal(t, treasury.QueryTaxCap, queryTaxCap.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryTaxCap.Args))

	// Check Flags
	denomFlag := queryTaxCap.Flag(flagDenom)
	require.NotNil(t, denomFlag)
	require.Equal(t, []string{"true"}, denomFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryIssuance(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryIssuance := GetCmdQueryIssuance(cdc)

	// Name check
	require.Equal(t, treasury.QueryIssuance, queryIssuance.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryIssuance.Args))

	// Check Flags
	denomFlag := queryIssuance.Flag(flagDenom)
	require.NotNil(t, denomFlag)
	require.Equal(t, []string{"true"}, denomFlag.Annotations[cobra.BashCompOneRequiredFlag])

	dayFlag := queryIssuance.Flag(flagDenom)
	require.NotNil(t, dayFlag)
}

func TestQueryMiningRewardWeight(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryMiningRewardWeight := GetCmdQueryMiningRewardWeight(cdc)

	// Name check
	require.Equal(t, treasury.QueryMiningRewardWeight, queryMiningRewardWeight.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryMiningRewardWeight.Args))

	// Check Flags
	epochFlag := queryMiningRewardWeight.Flag(flagEpoch)
	require.NotNil(t, epochFlag)
}

func TestQueryTaxProceeds(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryTaxProceeds := GetCmdQueryTaxProceeds(cdc)

	// Name check
	require.Equal(t, treasury.QueryTaxProceeds, queryTaxProceeds.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryTaxProceeds.Args))

	// Check Flags
	epochFlag := queryTaxProceeds.Flag(flagEpoch)
	require.NotNil(t, epochFlag)
}

func TestQuerySeigniorageProceeds(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	querySeigniorageProceeds := GetCmdQuerySeigniorageProceeds(cdc)

	// Name check
	require.Equal(t, treasury.QuerySeigniorageProceeds, querySeigniorageProceeds.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(querySeigniorageProceeds.Args))

	// Check Flags
	epochFlag := querySeigniorageProceeds.Flag(flagEpoch)
	require.NotNil(t, epochFlag)
}

func TestQueryCurrentEpoch(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryCurrentEpoch := GetCmdQueryCurrentEpoch(cdc)

	// Name check
	require.Equal(t, treasury.QueryCurrentEpoch, queryCurrentEpoch.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryCurrentEpoch.Args))
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParams := GetCmdQueryParams(cdc)

	// Name check
	require.Equal(t, treasury.QueryParams, queryParams.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParams.Args))
}
