package cli

import (
	"fmt"
	"strings"
	"terra/x/treasury"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagDenom = "denom"
)

// GetCmdQueryTaxRate implements the query taxrate command.
func GetCmdQueryTaxRate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "taxrate",
		Short: "Query the current stability tax rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryTaxRate), nil)
			if err != nil {
				return err
			}

			var taxRate sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &taxRate)
			return cliCtx.PrintOutput(taxRate)
		},
	}

	return cmd
}

// GetCmdQueryTaxCap implements the query taxcap command.
func GetCmdQueryTaxCap(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "taxcap [denom]",
		Short: "Query the current stability tax cap of the [denom] asset",
		Long: strings.TrimSpace(`
Query the current stability tax cap of the [denom] asset. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terracli query treasury taxcap krw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryTaxCap, denom), nil)
			if err != nil {
				return err
			}

			var price sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	return cmd
}

// GetCmdQueryActive implements the query active command.
func GetCmdQueryActiveClaims(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activeclaims",
		Short: "Query claims that have yet to be redeemed by the treasury",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryActiveClaims), nil)
			if err != nil {
				return err
			}

			var claims treasury.Claims
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &claims)
			return cliCtx.PrintOutput(claims)
		},
	}

	return cmd
}

// GetCmdQueryMiningWeight implements the query miningweight command.
func GetCmdQueryMiningWeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miningweight",
		Short: "Query the current mining weight",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryMiningRewardWeight), nil)
			if err != nil {
				return err
			}

			var miningWeight sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &miningWeight)
			return cliCtx.PrintOutput(miningWeight)
		},
	}

	return cmd
}

// GetCmdQueryBalance implements the query balance command.
func GetCmdQueryBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Query the current Treasury balance",
		Long: strings.TrimSpace(`
Query the current Treasury balance, denominated in TerraSDR. 
Balance clears periodically to satisfy claims registered with the Treasury. 

$ terracli query treasury balance
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryBalance), nil)
			if err != nil {
				return err
			}

			var balance sdk.Int
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &balance)
			return cliCtx.PrintOutput(balance)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current Treasury params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryParams), nil)
			if err != nil {
				return err
			}

			var params treasury.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
