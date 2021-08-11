package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	feeutils "github.com/terra-money/core/custom/auth/client/utils"
)

// GetTxFeesEstimateCommand will create a send tx and sign it with the given key.
func GetTxFeesEstimateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-fee [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Estimate required fee (stability + gas) and gas amount",
		Long: strings.TrimSpace(`
Estimate fees for the given stdTx

$ terrad tx estimate-fee [file] --gas-adjustment 1.4 --gas-prices 0.015uluna
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stdTx, err := authclient.ReadTxFromFile(clientCtx, args[0])
			if err != nil {
				return err
			}

			stdFee, err := feeutils.ComputeFeesWithCmd(clientCtx, cmd.Flags(), stdTx.GetMsgs()...)
			if err != nil {
				return err
			}

			response := feeutils.EstimateFeeResp{Fee: *stdFee}
			return clientCtx.PrintObjectLegacy(response)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
