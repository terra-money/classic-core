package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

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
	flagDenom   = "denom"
	flagPrice   = "price"
	flagVoter   = "voter"
	flagOffline = "offline"

	flagValidator = "validator"
	flagDelegatee = "delegatee"
)

// GetCmdPriceVote will create a send tx and sign it with the given key.
func GetCmdPriceVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Submit an oracle vote for the price of Luna",
		Long: strings.TrimSpace(`
Submit an oracle vote for the price of Luna denominated in the input denom.

$ terracli oracle vote --denom "ukrw" --price "8890" --from mykey

where "ukrw" is the denominating currency, and "8890" is the price of micro Luna in micro KRW from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli oracle vote --denom "ukrw" --price "8890" --from mykey --validator terravaloper1.......
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			offline := viper.GetBool(flagOffline)

			if !offline {
				if err := cliCtx.EnsureAccountExists(); err != nil {
					return err
				}
			}

			// Get from address
			voter := cliCtx.GetFromAddress()

			// Check the denom exists and valid
			denom := viper.GetString(flagDenom)

			// Check the price flag exists
			priceStr := viper.GetString(flagPrice)

			// By default the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if flag is set
			valStr := viper.GetString(flagValidator)
			if len(valStr) != 0 {
				parsedVal, err := sdk.ValAddressFromBech32(valStr)
				if err != nil {
					return errors.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			// Parse the price to Dec
			price, err := sdk.NewDecFromStr(priceStr)
			if err != nil {
				return fmt.Errorf("given price {%s} is not a valid format; price should be formatted as float", priceStr)
			}

			msg := oracle.NewMsgPriceFeed(denom, price, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().String(flagDenom, "", "denominating currency")
	cmd.Flags().String(flagPrice, "", "price of Luna in denom currency")
	cmd.Flags().String(flagValidator, "", "(optional) validator on behalf of which to vote (for delegated feeders)")
	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	cmd.MarkFlagRequired(flagDenom)
	cmd.MarkFlagRequired(flagPrice)

	return cmd
}

// GetCmdDelegateFeederPermission will create a feeder permission delegation tx and sign it with the given key.
func GetCmdDelegateFeederPermission(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-feeder",
		Short: "Delegate the permission to vote for the oracle to an address",
		Long: strings.TrimSpace(`
Delegate the permission to vote for the oracle to an address.
That way you can keep your validator operator key offline and use a separate replaceable key online.

$ terracli oracle set-feeder --delegatee terra1...... --from mykey

where "terra1abceuihfu93fud" is the address you want to delegate your voting rights to.
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

			// The address the right is being delegated from
			validator := sdk.ValAddress(voter)

			delegateeStr := viper.GetString(flagDelegatee)
			if len(delegateeStr) == 0 {
				return fmt.Errorf("--delegate flag is required")
			}
			delegatee, err := sdk.AccAddressFromBech32(delegateeStr)
			if err != nil {
				return errors.Wrap(err, "delegate is not a valid account address")
			}

			msg := oracle.NewMsgDelegateFeederPermission(validator, delegatee)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagDelegatee, "", "account the voting right will be delegated to")

	cmd.MarkFlagRequired(flagDelegatee)

	return cmd
}
