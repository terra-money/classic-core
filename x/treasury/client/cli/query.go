package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/terra-project/core/x/treasury"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagDenom = "denom"
	flagDay   = "day"
	flagEpoch = "epoch"
)

// GetCmdQueryTaxRate implements the query taxrate command.
func GetCmdQueryTaxRate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryTaxRate,
		Args:  cobra.NoArgs,
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
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxRate, epoch.String()), nil)
			if err != nil {
				return err
			}

			var taxRate sdk.Dec
			cdc.MustUnmarshalJSON(res, &taxRate)
			return cliCtx.PrintOutput(taxRate)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) an epoch number which you wants to get tax rate of; default is current epoch")
	return cmd
}

// GetCmdQueryTaxCap implements the query taxcap command.
func GetCmdQueryTaxCap(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryTaxCap,
		Args:  cobra.NoArgs,
		Short: "Query the current stability tax cap of a denom asset",
		Long: strings.TrimSpace(`
Query the current stability tax cap of the denom asset. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terracli query treasury taxcap --denom="ukrw"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxCap, denom), nil)
			if err != nil {
				return err
			}

			var price sdk.Dec
			cdc.MustUnmarshalJSON(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	cmd.Flags().String(flagDenom, "", "the denom for which you want to know the taxcap of")

	cmd.MarkFlagRequired(flagDenom)

	return cmd
}

// GetCmdQueryIssuance implements the query issuance command.
func GetCmdQueryIssuance(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryIssuance,
		Args:  cobra.NoArgs,
		Short: "Query the current issuance of a denom asset",
		Long: strings.TrimSpace(`
Query the current issuance of a denom asset. 

$ terracli query treasury issuance --denom="ukrw --day 0"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			dayStr := viper.GetString(flagDay)

			if len(dayStr) != 0 {
				_, err := strconv.ParseInt(dayStr, 10, 64)
				if err != nil {
					return err
				}
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", treasury.QuerierRoute, treasury.QueryIssuance, denom, dayStr), nil)
			if err != nil {
				return err
			}

			var issuance sdk.Int
			cdc.MustUnmarshalJSON(res, &issuance)
			return cliCtx.PrintOutput(issuance)
		},
	}

	cmd.Flags().String(flagDenom, "", "the denom which you want to know the issueance of")
	cmd.Flags().String(flagDay, "", "the # of date after genesis time, a user want to query")

	cmd.MarkFlagRequired(flagDenom)

	return cmd
}

// GetCmdQueryMiningRewardWeight implements the query reward-weight command.
func GetCmdQueryMiningRewardWeight(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryMiningRewardWeight,
		Args:  cobra.NoArgs,
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
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryMiningRewardWeight, epoch.String()), nil)
			if err != nil {
				return err
			}

			var rewardWeight sdk.Dec
			cdc.MustUnmarshalJSON(res, &rewardWeight)
			return cliCtx.PrintOutput(rewardWeight)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) an epoch number which you wants to get reward weight of; default is current epoch")

	return cmd
}

// GetCmdQueryTaxProceeds implements the query tax-proceeds command.
func GetCmdQueryTaxProceeds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryTaxProceeds,
		Args:  cobra.NoArgs,
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
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QueryTaxProceeds, epoch.String()), nil)
			if err != nil {
				return err
			}

			var taxProceeds sdk.Coins
			cdc.MustUnmarshalJSON(res, &taxProceeds)
			return cliCtx.PrintOutput(taxProceeds)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) an epoch number which you wants to get tax proceeds of; default is current epoch")

	return cmd
}

// GetCmdQuerySeigniorageProceeds implements the query seigniorage-proceeds command.
func GetCmdQuerySeigniorageProceeds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QuerySeigniorageProceeds,
		Args:  cobra.NoArgs,
		Short: "Query the seigniorage proceeds for the epoch",
		Long: strings.TrimSpace(`
Query the seigniorage proceeds corresponding to the given epoch. The return value will be in units of 'uluna' coins. 

$ terracli query treasury seigniorage-proceeds --epoch=14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch sdk.Int
			epochStr := viper.GetString(flagEpoch)
			if len(epochStr) == 0 {
				res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
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

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", treasury.QuerierRoute, treasury.QuerySeigniorageProceeds, epoch.String()), nil)
			if err != nil {
				return err
			}

			var seigniorageProceeds sdk.Int
			cdc.MustUnmarshalJSON(res, &seigniorageProceeds)
			return cliCtx.PrintOutput(seigniorageProceeds)
		},
	}

	cmd.Flags().String(flagEpoch, "", "(optional) an epoch number which you wants to get seigniorage proceeds of; default is current epoch")

	return cmd
}

// GetCmdQueryCurrentEpoch implements the query seigniorage-proceeds command.
func GetCmdQueryCurrentEpoch(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryCurrentEpoch,
		Args:  cobra.NoArgs,
		Short: "Query the current epoch number",
		Long: strings.TrimSpace(`
Query the current epoch, starting at 0.

$ terracli query treasury current-epoch
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryCurrentEpoch), nil)
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
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   treasury.QueryParams,
		Args:  cobra.NoArgs,
		Short: "Query the current Treasury params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", treasury.QuerierRoute, treasury.QueryParams), nil)
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
