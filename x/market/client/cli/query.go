package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/market/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	marketQueryCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Querying commands for the market module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketQueryCmd.AddCommand(
		GetCmdQuerySwap(),
		GetCmdQueryMintPoolDelta(),
		GetCmdQueryBurnPoolDelta(),
		GetCmdQueryParams(),
	)

	return marketQueryCmd
}

// GetCmdQuerySwap implements the query swap simulation result command.
func GetCmdQuerySwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Query a quote for a swap operation",
		Long: strings.TrimSpace(`
Query a quote for how many coins can be received in a swap operation. Note; rates are dynamic and can quickly change.

$ terrad query swap 5000000uluna usdr
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// parse offerCoin
			offerCoinStr := args[0]
			_, err = sdk.ParseCoinNormalized(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]

			res, err := queryClient.Swap(context.Background(),
				&types.QuerySwapRequest{OfferCoin: offerCoinStr, AskDenom: askDenom},
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

// GetCmdQueryMintPoolDelta implements the query mint pool delta command.
func GetCmdQueryMintPoolDelta() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-pool-delta",
		Args:  cobra.NoArgs,
		Short: "Query mint pool delta",
		Long: `Query mint pool delta, which is usdr amount used for mint operation from the MintBasePool.

$ terrad query market mint-pool-delta
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MintPoolDelta(context.Background(),
				&types.QueryMintPoolDeltaRequest{},
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

// GetCmdQueryBurn poolDelta implements the query mint pool delta command.
func GetCmdQueryBurnPoolDelta() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-pool-delta",
		Args:  cobra.NoArgs,
		Short: "Query burn pool delta",
		Long: `Query burn pool delta, which is usdr amount used for mint operation from the BurnBasePool.

$ terrad query market burn-pool-delta
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BurnPoolDelta(context.Background(),
				&types.QueryBurnPoolDeltaRequest{},
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
		Short: "Query the current market params",
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
