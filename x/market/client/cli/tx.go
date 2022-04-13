package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/spf13/cobra"

	marketcutils "github.com/terra-money/core/x/market/client/utils"
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

// NewSubmitSeigniorageRouteChangeTxCmd returns a CLI command handler for creating
// a seigniorage route change proposal governance transaction.
func NewSubmitSeigniorageRouteChangeTxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seigniorage-route-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a seigniorage route change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a seigniorage route change proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

Example:
$ %s tx gov submit-proposal seigniorage-route-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "routes": [
    {
      "address": "terra1...",
      "weight": "0.1",
    }
  ],
  "deposit": "1000uluna"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := marketcutils.ParseSeigniorageRouteChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			routes, err := proposal.Routes.ToSeigniorageRoutes()
			if err != nil {
				return err
			}

			content := types.NewSeigniorageRouteChangeProposal(
				proposal.Title, proposal.Description, routes,
			)

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
