package cli

import (
	"fmt"
	"strings"

	"github.com/terra-money/core/x/treasury/internal/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagDenom = "denom"
	flagEpoch = "epoch"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	oracleQueryCmd := &cobra.Command{
		Use:                        "treasury",
		Short:                      "Querying commands for the treasury module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryTaxRate(cdc),
		GetCmdQueryTaxCap(cdc),
		GetCmdQueryTaxCaps(cdc),
		GetCmdQueryRewardWeight(cdc),
		GetCmdQueryParams(cdc),
		GetCmdQueryTaxProceeds(cdc),
		GetCmdQuerySeigniorageProceeds(cdc),
		GetCmdQueryParams(cdc),
		GetCmdQueryIndicators(cdc),
	)...)

	return oracleQueryCmd

}

// GetCmdQueryTaxRate implements the query tax-rate command.
func GetCmdQueryTaxRate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-rate",
		Args:  cobra.NoArgs,
		Short: "Query the stability tax rate",
		Long: strings.TrimSpace(`
Query the stability tax rate of the current epoch.

$ terracli query treasury tax-rate
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxRate), nil)
			if err != nil {
				return err
			}

			var taxRate sdk.Dec
			cdc.MustUnmarshalJSON(res, &taxRate)
			return cliCtx.PrintOutput(taxRate)
		},
	}

	return cmd
}

// GetCmdQueryTaxCap implements the query taxcap command.
func GetCmdQueryTaxCap(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-cap [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current stability tax cap of a denom asset",
		Long: strings.TrimSpace(`
Query the current stability tax cap of the denom asset. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terracli query treasury tax-cap ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]

			params := types.NewQueryTaxCapParams(denom)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxCap), bz)
			if err != nil {
				return err
			}

			var taxCap sdk.Dec
			cdc.MustUnmarshalJSON(res, &taxCap)
			return cliCtx.PrintOutput(taxCap)
		},
	}

	return cmd
}

// GetCmdQueryTaxCaps implements the query tax-caps command.
func GetCmdQueryTaxCaps(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-caps",
		Args:  cobra.NoArgs,
		Short: "Query the current stability tax caps for all denom assets",
		Long: strings.TrimSpace(`
Query the current stability tax caps of the all denom assets. 
The stability tax levied on a tx is at most tax cap, regardless of the size of the transaction. 

$ terracli query treasury tax-caps
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxCaps), nil)
			if err != nil {
				return err
			}

			var taxCaps types.TaxCapsQueryResponse
			cdc.MustUnmarshalJSON(res, &taxCaps)
			return cliCtx.PrintOutput(taxCaps)
		},
	}

	return cmd
}

// GetCmdQueryRewardWeight implements the query reward-weight command.
func GetCmdQueryRewardWeight(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-weight",
		Args:  cobra.NoArgs,
		Short: "Query the reward weight",
		Long: strings.TrimSpace(`
Query the reward rate of the current epoch.

$ terracli query treasury reward-weight
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardWeight), nil)
			if err != nil {
				return err
			}

			var rewardWeight sdk.Dec
			cdc.MustUnmarshalJSON(res, &rewardWeight)
			return cliCtx.PrintOutput(rewardWeight)
		},
	}

	return cmd
}

// GetCmdQueryTaxProceeds implements the query tax-proceeds command.
func GetCmdQueryTaxProceeds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-proceeds",
		Args:  cobra.NoArgs,
		Short: "Query the tax proceeds for the current epoch",
		Long: strings.TrimSpace(`
Query the tax proceeds corresponding to the current epoch. The return value will be sdk.Coins{} of all the taxes collected. 

$ terracli query treasury tax-proceeds
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxProceeds), nil)
			if err != nil {
				return err
			}

			var taxProceeds sdk.Coins
			cdc.MustUnmarshalJSON(res, &taxProceeds)
			return cliCtx.PrintOutput(taxProceeds)
		},
	}

	return cmd
}

// GetCmdQuerySeigniorageProceeds implements the query seigniorage-proceeds command.
func GetCmdQuerySeigniorageProceeds(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seigniorage-proceeds",
		Args:  cobra.NoArgs,
		Short: "Query the seigniorage proceeds for the current epoch",
		Long: strings.TrimSpace(`
Query the seigniorage proceeds corresponding to the current epoch. The return value will be in units of 'uluna' coins. 

$ terracli query treasury seigniorage-proceeds
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySeigniorageProceeds), nil)
			if err != nil {
				return err
			}

			var seigniorageProceeds sdk.Int
			cdc.MustUnmarshalJSON(res, &seigniorageProceeds)
			return cliCtx.PrintOutput(seigniorageProceeds)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Treasury parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
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

// GetCmdQueryIndicators implements the query params command.
func GetCmdQueryIndicators(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "indicators",
		Args:  cobra.NoArgs,
		Short: "Query the current Treasury indicators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryIndicators), nil)
			if err != nil {
				return err
			}

			var response types.IndicatorQueryResonse
			cdc.MustUnmarshalJSON(res, &response)
			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}
