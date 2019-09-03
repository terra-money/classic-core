package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/terra-project/core/x/treasury/internal/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
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
	oracleQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryTaxRate(cdc),
		GetCmdQueryTaxCap(cdc),
		GetCmdQueryHistoricalIssuance(cdc),
		GetCmdQueryRewardWeight(cdc),
		GetCmdQueryParams(cdc),
		GetCmdQueryTaxProceeds(cdc),
		GetCmdQuerySeigniorageProceeds(cdc),
		GetCmdQueryCurrentEpoch(cdc),
		GetCmdQueryParams(cdc),
	)...)

	return oracleQueryCmd

}

// GetCmdQueryTaxRate implements the query taxrate command.
func GetCmdQueryTaxRate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tax-rate [epoch]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the stability tax rate",
		Long: strings.TrimSpace(`
Query the stability tax rate at the specified epoch.

$ terracli query treasury tax-rate 14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch int64
			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var err error
				epoch, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return errors.New(sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				}
			}

			params := types.NewQueryTaxRateParams(epoch)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxRate), bz)
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

// GetCmdQueryHistoricalIssuance implements the query issuance command.
func GetCmdQueryHistoricalIssuance(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "historical-issuance [epoch]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the epoch historical issuance",
		Long: strings.TrimSpace(`
Query the epoch issuance

$ terracli query treasury historical-issuance 0"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch int64
			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var err error
				epoch, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return errors.New(sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				}
			}

			params := types.NewQueryHistoricalIssuanceParams(epoch)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryHistoricalIssuance), bz)
			if err != nil {
				return err
			}

			var issuance sdk.Coins
			cdc.MustUnmarshalJSON(res, &issuance)
			return cliCtx.PrintOutput(issuance)
		},
	}

	return cmd
}

// GetCmdQueryRewardWeight implements the query reward-weight command.
func GetCmdQueryRewardWeight(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-weight [epoch]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the reward weight",
		Long: strings.TrimSpace(`
Query the reward rate at the specified epoch.

$ terracli query treasury reward-weight 14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch int64
			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var err error
				epoch, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return errors.New(sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				}
			}

			params := types.NewQueryRewardWeightParams(epoch)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRewardWeight), bz)
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
		Use:   "tax-proceeds [epoch]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the tax proceeds for the epoch",
		Long: strings.TrimSpace(`
Query the tax proceeds corresponding to the given epoch. The return value will be sdk.Coins{} of all the taxes collected. 

$ terracli query treasury tax-proceeds 14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch int64
			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var err error
				epoch, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return errors.New(sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				}
			}

			params := types.NewQueryTaxProceedsParams(epoch)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTaxProceeds), bz)
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
		Use:   "seigniorage-proceeds [epoch]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the seigniorage proceeds for the epoch",
		Long: strings.TrimSpace(`
Query the seigniorage proceeds corresponding to the given epoch. The return value will be in units of 'uluna' coins. 

$ terracli query treasury seigniorage-proceeds 14
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var epoch int64
			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
				if err != nil {
					return err
				}

				cdc.MustUnmarshalJSON(res, &epoch)
			} else {
				var err error
				epoch, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return errors.New(sdk.AppendMsgToErr("Falied to parse epoch", err.Error()))
				}
			}

			params := types.NewQuerySeigniorageParams(epoch)
			bz := cdc.MustMarshalJSON(params)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySeigniorageProceeds), bz)
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

// GetCmdQueryCurrentEpoch implements the query seigniorage-proceeds command.
func GetCmdQueryCurrentEpoch(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-epoch",
		Args:  cobra.NoArgs,
		Short: "Query the current epoch number",
		Long: strings.TrimSpace(`
Query the current epoch, starting at 0.

$ terracli query treasury current-epoch
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentEpoch), nil)
			if err != nil {
				return err
			}

			var curEpoch int64
			cdc.MustUnmarshalJSON(res, &curEpoch)
			return cliCtx.PrintOutput(sdk.NewInt(curEpoch))
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
