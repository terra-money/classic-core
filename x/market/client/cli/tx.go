package cli

import (
	"fmt"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/market"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOfferCoin = "offer-coin"
	flagAskDenom  = "ask-denom"
)

// GetSwapCmd will create and send a MsgSwap
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap",
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offer-coin to the ask-denom currency at the oracle's effective exchange rate. 

$ terracli market swap --offer-coin="1000krw" --ask-denom="usd"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			askDenom := viper.GetString(flagAskDenom)
			if len(askDenom) == 0 {
				return fmt.Errorf("--ask-denom flag is required")
			}

			if !assets.IsValidDenom(askDenom) {
				return fmt.Errorf("The denom is not known: %s", askDenom)
			}

			offerCoinStr := viper.GetString(flagOfferCoin)
			if len(offerCoinStr) == 0 {
				return fmt.Errorf("--offset-coin flag is required")
			}

			offerCoin, err := sdk.ParseCoin(offerCoinStr)
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()
			fromAccount, err := cliCtx.GetAccount(fromAddress)
			if err != nil {
				return err
			}

			if fromAccount.GetCoins().AmountOf(offerCoin.Denom).LT(offerCoin.Amount) {
				return fmt.Errorf(strings.TrimSpace(`
					account %s has insufficient amount of coins to pay the offered coins.\n
					Required: %s\n
					Given:    %s\n`),
					fromAddress, offerCoin, fromAccount.GetCoins())
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := market.NewMsgSwap(fromAddress, offerCoin, askDenom)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagOfferCoin, "", "The asset to swap from e.g. 1000ukrw")
	cmd.Flags().String(flagAskDenom, "", "Denom of the asset to swap to")

	return cmd
}
