package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/classic-terra/core/v2/x/dyncomm/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	dyncommQueryCmd := &cobra.Command{
		Use:                        "dyncomm",
		Short:                      "Querying commands for the dyncomm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	dyncommQueryCmd.AddCommand(
		GetCmdQueryRate(),
		GetCmdQueryParams(),
	)

	return dyncommQueryCmd
}

// GetCmdQueryRate implements the query min dynamic commission rate result command.
func GetCmdQueryRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the min dynamic commission rate of a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// parse validator address
			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Rate(context.Background(),
				&types.QueryRateRequest{ValidatorAddr: addr.String()},
			)
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
		Short: "Query the current dyncomm params",
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
