package cli

import (
	"os"
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
	flagOfferDenom     = "offerDenom"
	flagOfferAmount    = "offerAmount"
	flagAskDenom       = "askDenom"
	flagSwapperAddress = "address"
)

// GetSwapCmd will create and send a SwapMsg
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap offerCoin [offerCoin] askDenom [askDenom]",
		Short: "Atomically swap [offerCoin] asset with [askDenom] asset",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			askDenom := viper.GetString(flagAskDenom)
			offerCoin := sdk.NewInt64Coin(viper.GetString(flagOfferDenom), viper.GetInt64(flagOfferAmount))

			swapperStr := viper.GetString(flagSwapperAddress)
			swapper, err := sdk.AccAddressFromBech32(swapperStr)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := market.NewSwapMsg(swapper, offerCoin, askDenom)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagSwapperAddress, "", "Account address of the swapper")
	cmd.Flags().String(flagAskDenom, "luna", "Denom of the asset to swap to")
	cmd.Flags().String(flagOfferAmount, "", "Amount of the asset to swap from")
	cmd.Flags().String(flagOfferDenom, "", "Denom of the asset to swap from")

	cmd.MarkFlagRequired(flagOfferAmount)
	cmd.MarkFlagRequired(flagOfferDenom)
	cmd.MarkFlagRequired(flagSwapperAddress)

	return cmd
}
