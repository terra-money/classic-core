package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/terra-project/core/x/oracle/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	oracleTxCmd := &cobra.Command{
		Use:                        "oracle",
		Short:                      "Oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleTxCmd.AddCommand(client.PostCommands(
		GetCmdPrevote(cdc),
		GetCmdVote(cdc),
		GetCmdDelegateFeederPermission(cdc),
	)...)

	return oracleTxCmd
}

// GetCmdPrevote will create a Prevote tx and sign it with the given key.
func GetCmdPrevote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prevote [salt] [exchangeRate] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle prevote for the exchangeRate of Luna",
		Long: strings.TrimSpace(`
Submit an oracle prevote for the exchangeRate of Luna denominated in the input denom.
The purpose of prevote is to hide vote exchangeRate with hash which is formatted 
as hex string in SHA256("salt:exchangeRate:denom:voter")

# Prevote
$ terracli tx oracle prevote 1234 8888.0ukrw

where "ukrw" is the denominating currency, and "8888.0" is the exchangeRate of micro Luna in micro KRW from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle prevote 1234 8888.0ukrw terravaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			exchangeRate, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return fmt.Errorf("given exchangeRate {%s} is not a valid format; exchangeRate should be formatted as DecCoin", exchangeRate)
			}

			// Get from address
			voter := cliCtx.GetFromAddress()
			denom := exchangeRate.Denom
			amount := exchangeRate.Amount

			// By default the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errors.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			hashBytes, err := types.VoteHash(salt, amount, denom, validator)
			if err != nil {
				return err
			}

			hash := hex.EncodeToString(hashBytes)

			msg := types.NewMsgPrevote(hash, denom, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdVote will create a Vote tx and sign it with the given key.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [salt] [exchangeRate] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle vote for the exchangeRate of Luna",
		Long: strings.TrimSpace(`
Submit a vote for the exchangeRate of Luna denominated in the input denom. Companion to a prevote submitted in the previous vote period. 

$ terracli tx oracle vote 1234 8890.0ukrw

where "ukrw" is the denominating currency, and "8890.0" is the exchangeRate of micro Luna in micro KRW from the voter's point of view.

"salt" should match the salt used to generate the SHA256 hex in the associated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle vote 1234 8890.0ukrw terravaloper1....
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			exchangeRate, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return fmt.Errorf("given exchangeRate {%s} is not a valid format; exchangeRate should be formatted as DecCoin", exchangeRate)
			}

			// Get from address
			voter := cliCtx.GetFromAddress()
			denom := exchangeRate.Denom
			amount := exchangeRate.Amount

			// By default the voter is voting on behalf of itself
			validator := sdk.ValAddress(voter)

			// Override validator if validator is given
			if len(args) == 3 {
				parsedVal, err := sdk.ValAddressFromBech32(args[2])
				if err != nil {
					return errors.Wrap(err, "validator address is invalid")
				}
				validator = parsedVal
			}

			msg := types.NewMsgVote(amount, salt, denom, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDelegateFeederPermission will create a feeder permission delegation tx and sign it with the given key.
func GetCmdDelegateFeederPermission(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-feeder [feeder]",
		Args:  cobra.ExactArgs(1),
		Short: "Delegate the permission to vote for the oracle to an address",
		Long: strings.TrimSpace(`
Delegate the permission to vote for the oracle to an address.

Delegation can keep your validator operator key offline and use a separate replaceable key online.

$ terracli tx oracle set-feeder terra1...

where "terra1..." is the address you want to delegate your voting rights to.
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get from address
			voter := cliCtx.GetFromAddress()

			// The address the right is being delegated from
			validator := sdk.ValAddress(voter)

			feederStr := args[0]
			feeder, err := sdk.AccAddressFromBech32(feederStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegateConsent(validator, feeder)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
