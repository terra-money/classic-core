package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/msgauth/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	authorizationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the msg authorization module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authorizationQueryCmd.AddCommand(
		GetCmdQueryGrants(),
	)

	return authorizationQueryCmd
}

// GetCmdQueryGrants implements the query grants command.
func GetCmdQueryGrants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant [granter-addr] [grantee-addr]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query grants of a granter or between a granter-grantee pair",
		Long: strings.TrimSpace(`
Query grants between a granter-grantee pair,

$ terrad query msgauth grant terra... terra...

Or, query all grants of a granter,

$ terrad query msgauth grant terra... 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			granter := args[0]
			_, err = sdk.AccAddressFromBech32(granter)
			if err != nil {
				return err
			}

			// when grantee address was given,
			// query grants
			if len(args) == 2 {
				grantee := args[1]
				_, err = sdk.AccAddressFromBech32(grantee)
				if err != nil {
					return err
				}

				res, err := queryClient.Grants(context.Background(),
					&types.QueryGrantsRequest{Granter: granter, Grantee: grantee},
				)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.AllGrants(context.Background(),
				&types.QueryAllGrantsRequest{Granter: granter},
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
