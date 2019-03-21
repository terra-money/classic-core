package cli

import (
	"fmt"
	"strings"
	"terra/types"
	"terra/x/treasury"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagDenom = "denom"
	flagEpoch = "epoch"
)

// GetCmdQueryTaxRate implements the query taxrate command.
func GetCmdQueryTaxRate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "taxrate",
		Short: "Query the stability tax rate",
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

// GetCmdQueryMRL implements the query mrl command.
func GetCmdQueryMRL(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mrl",
		Short: "Query the unit mining rewards per luna for the epoch",
		Long: strings.TrimSpace(`
Query the unit mining reward for luna at the given epoch. 
mining rewards are a sum of transaction tax and seigniorage rewards. 

$ terracli query treasury mrl --epoch=5
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryMRL, epoch), nil)
			if err != nil {
				return err
			}

			var mrl sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &mrl)
			return cliCtx.PrintOutput(mrl)
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

$ terracli query treasury taxcap --denom="krw"
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

// GetCmdQueryIssuance implements the query issuance command.
func GetCmdQueryIssuance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issuance [denom]",
		Short: "Query the current issuance of the [denom] asset",
		Long: strings.TrimSpace(`
Query the current issuance of the [denom] asset. 

$ terracli query treasury issuance --denom="krw"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryIssuance, denom), nil)
			if err != nil {
				return err
			}

			var issuance sdk.Int
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &issuance)
			return cliCtx.PrintOutput(issuance)
		},
	}

	return cmd
}

// GetCmdQueryActiveClaims implements the query activeclaims command.
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

			var claims types.Claims
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

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryMiningRewardWeight, epoch), nil)
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

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryBalance, epoch), nil)
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
