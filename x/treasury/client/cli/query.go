package cli

import (
	"fmt"
	"strings"
	"terra/types"
	"terra/types/assets"
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
		Use:   treasury.QueryTaxRate,
		Short: "Query the stability tax rate",
		Long: strings.TrimSpace(`
Query the stability tax rate at the specified epoch.

$ terracli query treasury taxrate --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch sdk.Int
			epochStr := viper.GetString(flagEpoch)
			if len(epochStr) == 0 {
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var ok bool
				epoch, ok = sdk.NewIntFromString(epochStr)
				if !ok {
					return fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				}
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryTaxRate, epoch.String()), nil)
			if err != nil {
				return err
			}

			var taxRate sdk.Dec
			cdc.MustUnmarshalJSON(res, &taxRate)
			return cliCtx.PrintOutput(taxRate)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) a epoch number which you wants to get tax rate of; default is current epoch")
	return cmd
}

// GetCmdQueryTaxCap implements the query taxcap command.
func GetCmdQueryTaxCap(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryTaxCap,
		Short: "Query the current stability tax cap of a denom asset",
		Long: strings.TrimSpace(`
Query the current stability tax cap of the denom asset. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terracli query treasury taxcap --denom="mkrw"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			if len(denom) == 0 {
				return fmt.Errorf("--denom flag is required")
			}

			if !assets.IsValidDenom(denom) {
				return fmt.Errorf("given denom {%s} is not a valid one", denom)
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryTaxCap, denom), nil)
			if err != nil {
				return err
			}

			var price sdk.Dec
			cdc.MustUnmarshalJSON(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	cmd.Flags().String(flagDenom, "", "the denom which you want to know the taxcap of")

	return cmd
}

// GetCmdQueryIssuance implements the query issuance command.
func GetCmdQueryIssuance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryIssuance,
		Short: "Query the current issuance of a denom asset",
		Long: strings.TrimSpace(`
Query the current issuance of a denom asset. 

$ terracli query treasury issuance --denom="mkrw"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			if len(denom) == 0 {
				return fmt.Errorf("--denom flag is required")
			}

			if !assets.IsValidDenom(denom) {
				return fmt.Errorf("given denom {%s} is not a valid one", denom)
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryIssuance, denom), nil)
			if err != nil {
				return err
			}

			var issuance sdk.Int
			cdc.MustUnmarshalJSON(res, &issuance)
			return cliCtx.PrintOutput(issuance)
		},
	}

	cmd.Flags().String(flagDenom, "", "the denom which you want to know the issueance of")

	return cmd
}

// GetCmdQueryActiveClaims implements the query active-claims command.
func GetCmdQueryActiveClaims(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryActiveClaims,
		Short: "Query claims that have yet to be redeemed by the treasury",
		Long: strings.TrimSpace(`
Query the current active claims from oracle votes and program votes . 

$ terracli query treasury active-claims
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryActiveClaims), nil)
			if err != nil {
				return err
			}

			var claims types.ClaimPool
			cdc.MustUnmarshalJSON(res, &claims)
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

			var epoch sdk.Int
			epochStr := viper.GetString(flagEpoch)
			if len(epochStr) == 0 {
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var ok bool
				epoch, ok = sdk.NewIntFromString(epochStr)
				if !ok {
					return fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				}
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryMiningRewardWeight, epoch.String()), nil)
			if err != nil {
				return err
			}

			var rewardWeight sdk.Dec
			cdc.MustUnmarshalJSON(res, &rewardWeight)
			return cliCtx.PrintOutput(rewardWeight)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) a epoch number which you wants to get reward weight of; default is current epoch")

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

			var epoch sdk.Int
			epochStr := viper.GetString(flagEpoch)
			if len(epochStr) == 0 {
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var ok bool
				epoch, ok = sdk.NewIntFromString(epochStr)
				if !ok {
					return fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				}
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QueryTaxProceeds, epoch.String()), nil)
			if err != nil {
				return err
			}

			var taxProceeds sdk.Coins
			cdc.MustUnmarshalJSON(res, &taxProceeds)
			return cliCtx.PrintOutput(taxProceeds)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) a epoch number which you wants to get tax proceeds of; default is current epoch")

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

			var epoch sdk.Int
			epochStr := viper.GetString(flagEpoch)
			if len(epochStr) == 0 {
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var ok bool
				epoch, ok = sdk.NewIntFromString(epochStr)
				if !ok {
					return fmt.Errorf("the given epoch {%s} is not a valid format; epoch should be formatted as an integer", epochStr)
				}
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, treasury.QuerySeigniorageProceeds, epoch.String()), nil)
			if err != nil {
				return err
			}

			var seigniorageProceeds sdk.Int
			cdc.MustUnmarshalJSON(res, &seigniorageProceeds)
			return cliCtx.PrintOutput(seigniorageProceeds)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) a epoch number which you wants to get seigniorage proceeds of; default is current epoch")

	return cmd
}

// GetCmdQueryCurrentEpoch implements the query seigniorage-proceeds command.
func GetCmdQueryCurrentEpoch(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryCurrentEpoch,
		Short: "Query the current epoch number",
		Long: strings.TrimSpace(`
Query the current epoch.

$ terracli query treasury current-epoch
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, treasury.QueryCurrentEpoch), nil)
			if err != nil {
				return err
			}

			var curEpoch sdk.Int
			cdc.MustUnmarshalJSON(res, &curEpoch)
			return cliCtx.PrintOutput(curEpoch)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryParams,
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
