package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	tutils "github.com/terra-money/core/x/auth/client/utils"
)

// GetTxFeesEstimateCommand will create a send tx and sign it with the given key.
func GetTxFeesEstimateCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-fee [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Estimate required fee (stability + gas) and gas amount",
		Long: strings.TrimSpace(`
Estimate fees for the given stdTx

$ terracli tx estimate-fee [file] --gas-adjustment 1.4 --gas-prices 0.015uluna
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			stdTx, err := utils.ReadStdTxFromFile(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			gasAdjustment := viper.GetFloat64(flags.FlagGasAdjustment)

			var gasPrices sdk.DecCoins
			gasPricesStr := viper.GetString(flags.FlagGasPrices)
			if len(gasPricesStr) != 0 {
				gasPrices, err = sdk.ParseDecCoins(gasPricesStr)
				if err != nil {
					return err
				}
			}

			fees, gas, err := tutils.ComputeFeesWithStdTx(cliCtx, stdTx, gasAdjustment, gasPrices)

			if err != nil {
				return err
			}

			response := tutils.EstimateFeeResp{Fees: fees, Gas: gas}
			return cliCtx.PrintOutput(response)
		},
	}

	cmd = flags.GetCommands(cmd)[0]

	cmd.Flags().Float64(flags.FlagGasAdjustment, flags.DefaultGasAdjustment, "adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")
	cmd.Flags().String(flags.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10uluna)")
	// cmd.MarkFlagRequired(client.FlagGasAdjustment)

	return cmd
}
