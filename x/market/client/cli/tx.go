package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/terra-money/core/x/market/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Market transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketTxCmd.AddCommand(
		GetSwapCmd(),
	)

	return marketTxCmd
}

// GetSwapCmd will create and send a MsgSwap
func GetSwapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom] [to-address]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offer-coin to the ask-denom currency at the oracle's effective exchange rate. 

$ terrad market swap "1000ukrw" "uusd"

The to-address can be specified. A default to-address is trader.

$ terrad market swap "1000ukrw" "uusd" "terra1..."
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			offerCoinStr := args[0]
			offerCoin, err := sdk.ParseCoinNormalized(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]
			fromAddress := clientCtx.GetFromAddress()

			var msg sdk.Msg
			if len(args) == 3 {
				toAddress, err := sdk.AccAddressFromBech32(args[2])
				if err != nil {
					return err
				}

				msg = types.NewMsgSwapSend(fromAddress, toAddress, offerCoin, askDenom)
				if err = msg.ValidateBasic(); err != nil {
					return err
				}

			} else {
				msg = types.NewMsgSwap(fromAddress, offerCoin, askDenom)
				if err = msg.ValidateBasic(); err != nil {
					return err
				}
			}

			// build and sign the transaction, then broadcast to Tendermint
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
