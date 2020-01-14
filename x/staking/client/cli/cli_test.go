package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/terra-project/core/testutil"
)

func TestCreateValidatorTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	stakingTxCmd := &cobra.Command{
		Use:   "staking",
		Short: "Querying commands for the staking module",
	}

	txCmd.AddCommand(stakingTxCmd)

	stakingTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateValidator(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`staking`,
		`create-validator`,
		`--amount=1000000uluna`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--pubkey=terravalconspub1zcjduepqp4hglzr4m4u9lnxxp7f4tv9r0gt86jpcfk84m8fzhke4r3u7ctzsgxgpj5`,
		`--moniker=validator`,
		`--identity=identity-url`,
		`--details=details-description`,
		`--website=website-url`,
		`--commission-rate=0.10`,
		`--commission-max-rate=0.20`,
		`--commission-max-change-rate=0.01`,
		`--min-self-delegation=1`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestEditValidatorTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	stakingTxCmd := &cobra.Command{
		Use:   "staking",
		Short: "Querying commands for the staking module",
	}

	txCmd.AddCommand(stakingTxCmd)

	stakingTxCmd.AddCommand(client.PostCommands(
		GetCmdEditValidator(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`staking`,
		`edit-validator`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--moniker=validator`,
		`--identity=identity-url`,
		`--details=details-description`,
		`--website=website-url`,
		`--commission-rate=0.10`,
		`--min-self-delegation=1`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestDelegateTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	stakingTxCmd := &cobra.Command{
		Use:   "staking",
		Short: "Querying commands for the staking module",
	}

	txCmd.AddCommand(stakingTxCmd)

	stakingTxCmd.AddCommand(client.PostCommands(
		GetCmdDelegate(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`staking`,
		`delegate`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--validator=terravaloper1wg2mlrxdmnnkkykgqg4znky86nyrtc45q7a85l`,
		`--amount=100000000uluna`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestRedelegateTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	stakingTxCmd := &cobra.Command{
		Use:   "staking",
		Short: "Querying commands for the staking module",
	}

	txCmd.AddCommand(stakingTxCmd)

	stakingTxCmd.AddCommand(client.PostCommands(
		GetCmdRedelegate(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`staking`,
		`redelegate`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--addr-validator-source=terravaloper1wg2mlrxdmnnkkykgqg4znky86nyrtc45q7a85l`,
		`--addr-validator-dest=terravaloper1yg5gcx9krjhq0at036uyjg2h0cpxwt5g746drv`,
		`--amount=100000000uluna`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestUnbondTx(t *testing.T) {
	cdc, rootCmd, txCmd, _ := testutil.PrepareCmdTest()

	stakingTxCmd := &cobra.Command{
		Use:   "staking",
		Short: "Querying commands for the staking module",
	}

	txCmd.AddCommand(stakingTxCmd)

	stakingTxCmd.AddCommand(client.PostCommands(
		GetCmdUnbond(cdc),
	)...)

	// normal case all parameter given
	_, err := testutil.ExecuteCommand(
		rootCmd,
		`tx`,
		`staking`,
		`unbond`,
		`--from=terra1wg2mlrxdmnnkkykgqg4znky86nyrtc45q336yv`,
		`--validator=terravaloper1wg2mlrxdmnnkkykgqg4znky86nyrtc45q7a85l`,
		`--amount=100000000uluna`,
		`--generate-only`,
		`--offline`,
		`--chain-id=columbus`,
	)

	require.Nil(t, err)
}

func TestQueryValidator(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidator := GetCmdQueryValidator(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, staking.QueryValidator, queryValidator.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidator.Args))

	// Check Flags
	validatorFlag := queryValidator.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryValidators(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidators := GetCmdQueryValidators(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, staking.QueryValidators, queryValidators.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidators.Args))
}

func TestQueryValidatorUnbondingDelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorUnbondingDelegations := GetCmdQueryValidatorUnbondingDelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "unbonding-delegations-from", queryValidatorUnbondingDelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorUnbondingDelegations.Args))

	// Check Flags
	validatorFlag := queryValidatorUnbondingDelegations.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryValidatorRedelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorRedelegations := GetCmdQueryValidatorRedelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "redelegations-from", queryValidatorRedelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorRedelegations.Args))

	// Check Flags
	validatorFlag := queryValidatorRedelegations.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryDelegation(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryDelegation := GetCmdQueryDelegation(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "delegation", queryDelegation.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryDelegation.Args))

	// Check Flags
	validatorFlag := queryDelegation.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])

	delegatorFlag := queryDelegation.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryDelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryDelegations := GetCmdQueryDelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "delegations", queryDelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryDelegations.Args))

	// Check Flags
	delegatorFlag := queryDelegations.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryValidatorDelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryValidatorDelegations := GetCmdQueryValidatorDelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "delegations-to", queryValidatorDelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryValidatorDelegations.Args))

	// Check Flags
	validatorFlag := queryValidatorDelegations.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryUnbondingDelegation(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryUnbondingDelegation := GetCmdQueryUnbondingDelegation(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "unbonding-delegation", queryUnbondingDelegation.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryUnbondingDelegation.Args))

	// Check Flags
	validatorFlag := queryUnbondingDelegation.Flag(flagAddressValidator)
	require.NotNil(t, validatorFlag)
	require.Equal(t, []string{"true"}, validatorFlag.Annotations[cobra.BashCompOneRequiredFlag])

	delegatorFlag := queryUnbondingDelegation.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryUnbondingDelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryUnbondingDelegations := GetCmdQueryUnbondingDelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "unbonding-delegations", queryUnbondingDelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryUnbondingDelegations.Args))

	// Check Flags
	delegatorFlag := queryUnbondingDelegations.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryRedelegation(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryRedelegation := GetCmdQueryRedelegation(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "redelegation", queryRedelegation.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryRedelegation.Args))

	// Check Flags
	delegatorFlag := queryRedelegation.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])

	validatorSrcFlag := queryRedelegation.Flag(flagAddressValidatorSrc)
	require.NotNil(t, validatorSrcFlag)
	require.Equal(t, []string{"true"}, validatorSrcFlag.Annotations[cobra.BashCompOneRequiredFlag])

	validatorDstFlag := queryRedelegation.Flag(flagAddressValidatorDst)
	require.NotNil(t, validatorDstFlag)
	require.Equal(t, []string{"true"}, validatorDstFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryRedelegations(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryRedelegations := GetCmdQueryRedelegations(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "redelegations", queryRedelegations.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryRedelegations.Args))

	// Check Flags
	delegatorFlag := queryRedelegations.Flag(flagAddressDelegator)
	require.NotNil(t, delegatorFlag)
	require.Equal(t, []string{"true"}, delegatorFlag.Annotations[cobra.BashCompOneRequiredFlag])
}

func TestQueryPool(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryPool := GetCmdQueryPool(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "pool", queryPool.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryPool.Args))
}

func TestQueryParams(t *testing.T) {
	cdc, _, _, _ := testutil.PrepareCmdTest()

	queryParams := GetCmdQueryParams(staking.StoreKey, cdc)

	// Name check
	require.Equal(t, "params", queryParams.Name())

	// NoArg check
	require.Equal(t, testutil.FS(cobra.PositionalArgs(cobra.NoArgs)), testutil.FS(queryParams.Args))
}
