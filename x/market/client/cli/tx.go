package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"

	feeutils "github.com/terra-money/core/x/auth/client/utils"
	"github.com/terra-money/core/x/market/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Market transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketTxCmd.AddCommand(flags.PostCommands(
		GetSwapCmd(cdc),
	)...)

	return marketTxCmd
}

// GetSwapCmd will create and send a MsgSwap
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom] [to-address]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offer-coin to the ask-denom currency at the oracle's effective exchange rate. 

$ terracli market swap "1000ukrw" "uusd"

The to-address can be specfied. A default to-address is trader.

$ terracli market swap "1000ukrw" "uusd" "terra1..."
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			offerCoinStr := args[0]
			offerCoin, err := sdk.ParseCoin(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]
			fromAddress := cliCtx.GetFromAddress()

			var msg sdk.Msg
			if len(args) == 3 {
				toAddress, err := sdk.AccAddressFromBech32(args[2])
				if err != nil {
					return err
				}

				msg = types.NewMsgSwapSend(fromAddress, toAddress, offerCoin, askDenom)
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
			} else {
				msg = types.NewMsgSwap(fromAddress, offerCoin, askDenom)
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
