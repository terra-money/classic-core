package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/x/market/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	marketQueryCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Querying commands for the market module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQuerySwap(queryRoute, cdc),
		GetCmdQueryTerraPoolDelta(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)

	return marketQueryCmd
}

// GetCmdQuerySwap implements the query swap simulation result command.
func GetCmdQuerySwap(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Query a quote for a swap operation",
		Long: strings.TrimSpace(`
Query a quote for how many coins can be received in a swap operation. Note; rates are dynamic and can quickly change.

$ terracli query swap 5000000uluna usdr
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// parse offerCoin
			offerCoinStr := args[0]
			offerCoin, err := sdk.ParseCoin(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]

			params := types.NewQuerySwapParams(offerCoin, askDenom)
			bz := cdc.MustMarshalJSON(params)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QuerySwap, askDenom), bz)
			if err != nil {
				return err
			}

			var retCoin sdk.Coin
			cdc.MustUnmarshalJSON(res, &retCoin)
			return cliCtx.PrintOutput(retCoin)
		},
	}

	return cmd
}

// GetCmdQueryTerraPoolDelta implements the query terra pool delta command.
func GetCmdQueryTerraPoolDelta(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terra-pool-delta",
		Args:  cobra.NoArgs,
		Short: "Query terra pool delta",
		Long: `Query terra pool delta, which is usdr amount used for swap operation from the TerraPool.
It can be negative if the market wants more Terra than Luna, and vice versa if the market wants more Luna.

$ terracli query market terra-pool-delta
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTerraPoolDelta), nil)
			if err != nil {
				return err
			}

			var poolDelta sdk.Dec
			cdc.MustUnmarshalJSON(res, &poolDelta)
			return cliCtx.PrintOutput(poolDelta)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
