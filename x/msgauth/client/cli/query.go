package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/terra-money/core/x/msgauth/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	authorizationQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the msg authorization module",
		Long:                       "",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authorizationQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryGrant(queryRoute, cdc),
		GetCmdQueryGrants(queryRoute, cdc),
	)...)

	return authorizationQueryCmd
}

// GetCmdQueryGrant implements the query grant command.
func GetCmdQueryGrant(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "grant [granter-addr] [grantee-addr] [msg-type]",
		Args:  cobra.ExactArgs(3),
		Short: "Query grant entry about a specific msg type between a granter-grantee pair",
		Long: strings.TrimSpace(`
Query grant entry about a specific msg type between a granter-grantee pair,

$ terracli query msgauth grant terra... terra... send
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			granterAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			granteeAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msgAuthorized := args[2]

			params := types.NewQueryGrantParams(granterAddr, granteeAddr, msgAuthorized)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGrant), bz)
			if err != nil {
				return err
			}

			var grant types.AuthorizationGrant
			err = cdc.UnmarshalJSON(res, &grant)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(grant)
		},
	}
}

// GetCmdQueryGrants implements the query grants command.
func GetCmdQueryGrants(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "grants [granter-addr] [grantee-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query grant entries between a granter-grantee pair",
		Long: strings.TrimSpace(`
Query grant entries between a granter-grantee pair,

$ terracli query msgauth grants terra... terra...
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			granterAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			granteeAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := types.NewQueryGrantsParams(granterAddr, granteeAddr)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGrants), bz)
			if err != nil {
				return err
			}

			var grants []types.AuthorizationGrant
			err = cdc.UnmarshalJSON(res, &grants)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(grants)
		},
	}
}
