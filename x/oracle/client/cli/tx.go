package cli

import (
	"math"
	"os"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdPriceVote will create a send tx and sign it with the given key.
func GetCmdPriceVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle",
		Short: "Submit a vote for the price oracle",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			voterAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// parse denom of the coin to be voted on
			denom := viper.GetString(flagDenom)
			target := sdk.NewDecWithPrec(int64(math.Round(viper.GetFloat64(flagTargetPrice)*100)), 2)
			current := sdk.NewDecWithPrec(int64(math.Round(viper.GetFloat64(flagCurrentPrice)*100)), 2)

			// build and sign the transaction, then broadcast to Tendermint
			msg := oracle.NewPriceFeedMsg(denom, target, current, voterAddress)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsDenom)
	cmd.Flags().AddFlagSet(fsCurrentPrice)
	cmd.Flags().AddFlagSet(fsTargetPrice)

	cmd.MarkFlagRequired(flagCurrentPrice)
	cmd.MarkFlagRequired(flagTargetPrice)
	cmd.MarkFlagRequired(flagDenom)

	return cmd
}
