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

const (
	flagDenom        = "denom"
	flagTargetPrice  = "targetprice"
	flagCurrentPrice = "currentprice"
	flagVoterAddress = "address"
)

// GetPriceFeedCmd will create a send tx and sign it with the given key.
func GetPriceFeedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle",
		Short: "Submit a vote for the price oracle",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			voterStr := viper.GetString(flagVoterAddress)

			voter, err := sdk.AccAddressFromBech32(voterStr)
			if err != nil {
				return err
			}

			// parse denom of the coin to be voted on
			denom := viper.GetString(flagDenom)
			target := sdk.NewDecWithPrec(int64(math.Round(viper.GetFloat64(flagTargetPrice)*100)), 2)
			current := sdk.NewDecWithPrec(int64(math.Round(viper.GetFloat64(flagCurrentPrice)*100)), 2)

			// build and sign the transaction, then broadcast to Tendermint
			msg := oracle.NewPriceFeedMsg(denom, target, current, voter)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagVoterAddress, "", "Validator address of the voter")
	cmd.Flags().String(flagDenom, "", "Denom of the asset to vote on")
	cmd.Flags().Float32(flagTargetPrice, 0.0, "Price of the asset denominated in Luna, to 2nd decimal precision")
	cmd.Flags().Float32(flagCurrentPrice, 0.0, "Price of the asset denominated in Luna, to 2nd decimal precision")

	cmd.MarkFlagRequired(flagDenom)
	cmd.MarkFlagRequired(flagTargetPrice)
	cmd.MarkFlagRequired(flagCurrentPrice)
	cmd.MarkFlagRequired(flagVoterAddress)

	return cmd
}
