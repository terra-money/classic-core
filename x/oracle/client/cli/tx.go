package cli

import (
	"fmt"
	"strings"

	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDenom = "denom"
	flagPrice = "price"
	flagVoter = "voter"
)

// GetCmdPriceVote will create a send tx and sign it with the given key.
func GetCmdPriceVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Submit an oracle vote for the price of Luna",
		Long: strings.TrimSpace(`
Submit an oracle vote for the price of Luna denominated in the input denom.

$ terracli oracle vote --denom "mkrw" --price "8890" --from mykey

where "mkrw" is the denominating currency, and "8890" is the price of micro Luna in micro KRW from the voter's point of view. 
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			// Get from address
			voter := cliCtx.GetFromAddress()

			// Check the denom exists and valid
			denom := viper.GetString(flagDenom)
			if len(denom) == 0 {
				return fmt.Errorf("--denom flag is required")
			}

			// Check the price flag exists
			priceStr := viper.GetString(flagPrice)
			if len(priceStr) == 0 {
				return fmt.Errorf("--price flag is required")
			}

			// Parse the price to Dec
			price, err := sdk.NewDecFromStr(priceStr)
			if err != nil {
				return fmt.Errorf("given price {%s} is not a valid format; price should be formatted as float", priceStr)
			}

			msg := oracle.NewMsgPriceFeed(denom, price, voter)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagDenom, "", "denominating currency")
	cmd.Flags().String(flagPrice, "", "price of Luna in denom currency")

	return cmd
}
