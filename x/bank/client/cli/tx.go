package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        bank.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SendTxCmd(cdc),
	)
	return txCmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [to_address] [coins]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.MsgSend{FromAddress: cliCtx.GetFromAddress(), ToAddress: to, Amount: coins}

			if !cliCtx.GenerateOnly && txBldr.Fees().IsZero() {
				// extimate tax and gas
				fees, gas, err := feeutils.ComputeFees(cliCtx, feeutils.ComputeReqParams{
					Memo:          txBldr.Memo(),
					ChainID:       txBldr.ChainID(),
					AccountNumber: txBldr.AccountNumber(),
					Sequence:      txBldr.Sequence(),
					GasPrices:     txBldr.GasPrices(),
					Gas:           fmt.Sprintf("%d", txBldr.Gas()),
					GasAdjustment: fmt.Sprintf("%f", txBldr.GasAdjustment()),
					Msgs:          []sdk.Msg{msg},
				})

				if err != nil {
					return err
				}

				// override gas and fees
				txBldr = auth.NewTxBuilder(txBldr.TxEncoder(), txBldr.AccountNumber(), txBldr.Sequence(),
					gas, txBldr.GasAdjustment(), false, txBldr.ChainID(), txBldr.Memo(), fees, sdk.DecCoins{})
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
