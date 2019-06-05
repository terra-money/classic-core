package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/terra-project/core/x/market"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCmdQuerySwap implements the query swap amount command.
func GetCmdQuerySwap(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Args:  cobra.NoArgs,
		Short: "Query the real amount of swap operation",
		Long: strings.TrimSpace(`
Query the real amount of swap operation which a user can receive. 

$ terracli query query swap --ask-denom usdr --offer-coin 5000000uluna
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			askDenom := viper.GetString(flagAskDenom)

			// parse offerCoin
			offerCoinStr := viper.GetString(flagOfferCoin)
			offerCoin, err := sdk.ParseCoin(offerCoinStr)
			if err != nil {
				return err
			}

			params := market.NewQuerySwapParams(offerCoin)
			bz := cdc.MustMarshalJSON(params)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", market.QuerierRoute, market.QuerySwap, askDenom), bz)
			if err != nil {
				return err
			}

			var retCoin sdk.Coin
			err = cdc.UnmarshalJSON(res, &retCoin)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(retCoin)
		},
	}

	cmd.Flags().String(flagAskDenom, "", "Denom of the asset to swap to")
	cmd.Flags().String(flagOfferCoin, "", "The asset to swap from e.g. 1000ukrw")

	cmd.MarkFlagRequired(flagAskDenom)
	cmd.MarkFlagRequired(flagOfferCoin)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   market.QueryParams,
		Args:  cobra.NoArgs,
		Short: "Query the current market params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", market.QuerierRoute, market.QueryParams), nil)
			if err != nil {
				return err
			}

			var params market.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
