package cli

import (
	"strings"
	"terra/x/market"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOfferCoin = "offerCoin"
	flagAskDenom  = "askDenom"
)

// GetSwapCmd will create and send a MsgSwap
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offerCoin] [askDenom]",
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offerCoin to the askDenom currency at the oracle's effective exchange rate. 

$ terracli market swap --offerCoin="1000krw" --askDenom="usd"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			askDenom := viper.GetString(flagAskDenom)
			offerCoin, err := sdk.ParseCoin(viper.GetString(flagOfferCoin))
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := market.NewMsgSwap(fromAddress, offerCoin, askDenom)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagOfferCoin, "", "The asset to swap from e.g. 1000krw")
	cmd.Flags().String(flagAskDenom, "", "Denom of the asset to swap to")

	return cmd
}
