package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/classic-terra/core/v2/x/feeshare/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	feesQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	feesQueryCmd.AddCommand(
		GetCmdQueryFeeShares(),
		GetCmdQueryFeeShare(),
		GetCmdQueryParams(),
		GetCmdQueryDeployerFeeShares(),
		GetCmdQueryWithdrawerFeeShares(),
	)

	return feesQueryCmd
}

// GetCmdQueryFeeShares implements a command to return all registered contracts
// for fee distribution
func GetCmdQueryFeeShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contracts",
		Short: "Query all FeeShares",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryFeeSharesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.FeeShares(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryFeeShare implements a command to return a registered contract for fee
// distribution
func GetCmdQueryFeeShare() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contract [contract_address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a registered contract for fee distribution by its bech32 address",
		Long:    "Query a registered contract for fee distribution by its bech32 address",
		Example: fmt.Sprintf("%s query feeshare contract <contract-address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryFeeShareRequest{ContractAddress: args[0]}
			if err := req.ValidateBasic(); err != nil {
				return err
			}

			// Query store
			res, err := queryClient.FeeShare(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryParams implements a command to return the current FeeShare
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current feeshare module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}

			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDeployerFeeShares implements a command that returns all contracts
// that a deployer has registered for fee distribution
func GetCmdQueryDeployerFeeShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployer-contracts [deployer_address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query all contracts that a given deployer has registered for feeshare distribution",
		Long:    "Query all contracts that a given deployer has registered for feeshare distribution",
		Example: fmt.Sprintf("%s query feeshare deployer-contracts <deployer-address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			deployerFeeShareReq := &types.QueryDeployerFeeSharesRequest{
				DeployerAddress: args[0],
				Pagination:      pageReq,
			}
			if deployerFeeShareReq.ValidateBasic() != nil {
				return err
			}

			// Query store
			res, err := queryClient.DeployerFeeShares(context.Background(), deployerFeeShareReq)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryWithdrawerFeeShares implements a command that returns all
// contracts that have registered for fee distribution with a given withdraw
// address
func GetCmdQueryWithdrawerFeeShares() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdrawer-contracts [withdraw_address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query all contracts that have been registered for feeshare distribution with a given withdrawer address",
		Long:    "Query all contracts that have been registered for feeshare distribution with a given withdrawer address",
		Example: fmt.Sprintf("%s query feeshare withdrawer-contracts <withdrawer-address>", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			withdrawReq := &types.QueryWithdrawerFeeSharesRequest{
				WithdrawerAddress: args[0],
				Pagination:        pageReq,
			}

			if err := withdrawReq.ValidateBasic(); err != nil {
				return err
			}

			// Query store
			res, err := queryClient.WithdrawerFeeShares(context.Background(), withdrawReq)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
