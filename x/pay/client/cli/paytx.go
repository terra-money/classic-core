package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/spf13/cobra"
)

const (
	flagTo    = "to"
	flagCoins = "coins"
)

// PayTxCmd will create a pay tx and sign it with the given key.
func PayTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pay --to [to_address] --coins [amount] --from [from_address or key_name]",
		Args:  cobra.NoArgs,
		Short: "Create and sign a pay tx",
		Long: strings.TrimSpace(`
Create, sign and broadcast pay tx.

In case generate-only, --from should be specified as address not key name.
$ terracli tx pay --to [to_address] --coins [amount] --from [from_address or key_name]
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			coinsStr := viper.GetString(flagCoins)

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(coinsStr)
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			if !cliCtx.GenerateOnly {

				// ensure account has enough coins
				if !account.GetCoins().IsAllGTE(coins) {
					return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
				}

			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.NewMsgSend(from, to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	cmd.Flags().String(flagTo, "", "the address which a user wants to pay")
	cmd.Flags().String(flagCoins, "", "the amount a user wants to transfer")

	cmd.MarkFlagRequired(client.FlagFrom)
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagCoins)

	return cmd
}
