package cli

import (
	"math"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDenom        = "denom"
	flagPrice        = "price"
	flagVoterAddress = "address"
)

// VoteCmd will create a send tx and sign it with the given key.
func VoteCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle",
		Short: "Submit a vote for the price oracle",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

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
			price := sdk.NewDecWithPrec(int64(math.Round(viper.GetFloat64(flagPrice)*100)), 2)

			// build and sign the transaction, then broadcast to Tendermint
			msg := oracle.NewPriceFeedMsg(denom, price, voter)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagVoterAddress, "", "Validator address of the voter")
	cmd.Flags().String(flagDenom, "", "Denom of the asset to vote on")
	cmd.Flags().Float32(flagPrice, 0.0, "Price of the asset denominated in SDR, to 2nd decimal precision")

	cmd.MarkFlagRequired(flagDenom)
	cmd.MarkFlagRequired(flagPrice)
	cmd.MarkFlagRequired(flagVoterAddress)

	return cmd
}
