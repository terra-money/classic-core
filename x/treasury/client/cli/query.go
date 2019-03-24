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
		Use:   treasury.QueryTaxRate + " [epoch]",
		Short: "Query the stability tax rate",
		Long: strings.TrimSpace(`
Query the stability tax rate at the specified epoch.

$ terracli query treasury taxrate --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryTaxRate, epoch), nil)
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
		Use:   treasury.QueryTaxCap + " [denom]",
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
		Use:   treasury.QueryIssuance + " [denom]",
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

// GetCmdQueryActiveClaims implements the query active-claims command.
func GetCmdQueryActiveClaims(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryActiveClaims,
		Short: "Query claims that have yet to be redeemed by the treasury",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryActiveClaims), nil)
			if err != nil {
				return err
			}

			var claims types.ClaimPool
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &claims)
			return cliCtx.PrintOutput(claims)
		},
	}

	return cmd
}

// GetCmdQueryMiningRewardWeight implements the query reward-weight command.
func GetCmdQueryMiningRewardWeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryMiningRewardWeight,
		Short: "Query the mining reward weight",
		Long: strings.TrimSpace(`
Query the mining reward rate at the specified epoch.

$ terracli query treasury reward-weight --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryMiningRewardWeight, epoch), nil)
			if err != nil {
				return err
			}

			var rewardWeight sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &rewardWeight)
			return cliCtx.PrintOutput(rewardWeight)
		},
	}

	return cmd
}

// GetCmdQueryTaxProceeds implements the query tax-proceeds command.
func GetCmdQueryTaxProceeds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryTaxProceeds,
		Short: "Query the tax proceeds for the epoch",
		Long: strings.TrimSpace(`
Query the tax proceeds corresponding to the given epoch. The return value will be sdk.Coins{} of all the taxes collected. 

$ terracli query treasury tax-proceeds --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QueryTaxProceeds, epoch), nil)
			if err != nil {
				return err
			}

			var taxProceeds sdk.Coins
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &taxProceeds)
			return cliCtx.PrintOutput(taxProceeds)
		},
	}

	return cmd
}

// GetCmdQuerySeigniorageProceeds implements the query seigniorage-proceeds command.
func GetCmdQuerySeigniorageProceeds(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QuerySeigniorageProceeds,
		Short: "Query the seigniorage proceeds for the epoch",
		Long: strings.TrimSpace(`
Query the seigniorage proceeds corresponding to the given epoch. The return value will be in units of Luna coins. 

$ terracli query treasury seigniorage-proceeds --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			epoch := viper.GetInt(flagEpoch)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, treasury.QuerySeigniorageProceeds, epoch), nil)
			if err != nil {
				return err
			}

			var seigniorageProceeds sdk.Coins
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &seigniorageProceeds)
			return cliCtx.PrintOutput(seigniorageProceeds)
		},
	}

	return cmd
}

// GetCmdQueryCurrentEpoch implements the query seigniorage-proceeds command.
func GetCmdQueryCurrentEpoch(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryCurrentEpoch,
		Short: "Query the current epoch number",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				return err
			}

			var curEpoch sdk.Int
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &curEpoch)
			return cliCtx.PrintOutput(curEpoch)
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
