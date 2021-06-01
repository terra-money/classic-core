package cli

import (
	"context"
	"strings"

	"github.com/terra-money/core/x/treasury/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

const (
	flagDenom = "denom"
	flagEpoch = "epoch"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	oracleQueryCmd := &cobra.Command{
		Use:                        "treasury",
		Short:                      "Querying commands for the treasury module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleQueryCmd.AddCommand(
		GetCmdQueryTaxRate(),
		GetCmdQueryTaxCap(),
		GetCmdQueryTaxCaps(),
		GetCmdQueryRewardWeight(),
		GetCmdQueryTaxProceeds(),
		GetCmdQuerySeigniorageProceeds(),
		GetCmdQueryIndicators(),
		GetCmdQueryParams(),
	)

	return oracleQueryCmd

}

// GetCmdQueryTaxRate implements the query tax-rate command.
func GetCmdQueryTaxRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-rate",
		Args:  cobra.NoArgs,
		Short: "Query the stability tax rate",
		Long: strings.TrimSpace(`
Query the stability tax rate of the current epoch.

$ terrad query treasury tax-rate
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TaxRate(context.Background(), &types.QueryTaxRateRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTaxCap implements the query taxcap command.
func GetCmdQueryTaxCap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-cap [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current stability tax cap of a denom asset",
		Long: strings.TrimSpace(`
Query the current stability tax cap of the denom asset. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terrad query treasury tax-cap ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			denom := args[0]
			res, err := queryClient.TaxCap(context.Background(), &types.QueryTaxCapRequest{
				Denom: denom,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTaxCaps implements the query tax-caps command.
func GetCmdQueryTaxCaps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-caps",
		Args:  cobra.NoArgs,
		Short: "Query the current stability tax caps for all denom assets",
		Long: strings.TrimSpace(`
Query the current stability tax caps of the all denom assets. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terrad query treasury tax-caps
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TaxCaps(context.Background(), &types.QueryTaxCapsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryRewardWeight implements the query reward-weight command.
func GetCmdQueryRewardWeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-weight",
		Args:  cobra.NoArgs,
		Short: "Query the reward weight",
		Long: strings.TrimSpace(`
Query the reward rate of the current epoch.

$ terrad query treasury reward-weight
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RewardWeight(context.Background(), &types.QueryRewardWeightRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTaxProceeds implements the query tax-proceeds command.
func GetCmdQueryTaxProceeds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-proceeds",
		Args:  cobra.NoArgs,
		Short: "Query the tax proceeds for the current epoch",
		Long: strings.TrimSpace(`
Query the tax proceeds corresponding to the current epoch. The return value will be sdk.Coins{} of all the taxes collected. 

$ terrad query treasury tax-proceeds
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TaxProceeds(context.Background(), &types.QueryTaxProceedsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQuerySeigniorageProceeds implements the query seigniorage-proceeds command.
func GetCmdQuerySeigniorageProceeds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seigniorage-proceeds",
		Args:  cobra.NoArgs,
		Short: "Query the seigniorage proceeds for the current epoch",
		Long: strings.TrimSpace(`
Query the seigniorage proceeds corresponding to the current epoch. The return value will be in units of 'uluna' coins. 

$ terrad query treasury seigniorage-proceeds
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.SeigniorageProceeds(context.Background(), &types.QuerySeigniorageProceedsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryIndicators implements the query indicators command.
func GetCmdQueryIndicators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indicators",
		Args:  cobra.NoArgs,
		Short: "Query the current Treasury indicators",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Indicators(context.Background(), &types.QueryIndicatorsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Treasury parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
