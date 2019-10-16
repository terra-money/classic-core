package cli

import (
	"strings"

	"github.com/terra-project/core/x/market/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for market module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	marketTxCmd := &cobra.Command{
		Use:                        "market",
		Short:                      "Market transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketTxCmd.AddCommand(client.PostCommands(
		GetSwapCmd(cdc),
	)...)

	return marketTxCmd
}

// GetSwapCmd will create and send a MsgSwap
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offer-coin] [ask-denom]",
		Args:  cobra.ExactArgs(2),
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offer-coin to the ask-denom currency at the oracle's effective exchange rate. 

$ terracli market swap "1000ukrw" "uusd"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			offerCoinStr := args[0]
			offerCoin, err := sdk.ParseCoin(offerCoinStr)
			if err != nil {
				return err
			}

			askDenom := args[1]

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSwap(fromAddress, offerCoin, askDenom)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
