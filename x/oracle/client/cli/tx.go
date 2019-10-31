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
		GetCmdExchangeRatePrevote(cdc),
		GetCmdExchangeRateVote(cdc),
		GetCmdDelegateFeederPermission(cdc),
	)...)

	return oracleTxCmd
}

// GetCmdExchangeRatePrevote will create a exchangeRatePrevote tx and sign it with the given key.
func GetCmdExchangeRatePrevote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prevote [salt] [exchange_rate] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle prevote for the exchange rate of Luna",
		Long: strings.TrimSpace(`
Submit an oracle prevote for the exchange rate of Luna denominated in the input denom.
The purpose of prevote is to hide vote exchnage rate with hash which is formatted 
as hex string in SHA256("salt:exchange_rate:denom:voter")

# Prevote
$ terracli tx oracle prevote 1234 8888.0ukrw

where "ukrw" is the denominating currency, and "8888.0" is the exchange rate of micro Luna in micro KRW from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle prevote 1234 8888.0ukrw terravaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			rate, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange_rate should be formatted as DecCoin", rate)
			}

			// Get from address
			voter := cliCtx.GetFromAddress()
			denom := rate.Denom
			amount := rate.Amount

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

			msg := types.NewMsgExchangeRatePrevote(hash, denom, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdExchangeRateVote will create a exchangeRateVote tx and sign it with the given key.
func GetCmdExchangeRateVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [salt] [exchange_rate] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle vote for the exchange_rate of Luna",
		Long: strings.TrimSpace(`
Submit a vote for the exchange_rate of Luna w.r.t the input denom. Companion to a prevote submitted in the previous vote period. 

$ terracli tx oracle vote 1234 8890.0ukrw

where "ukrw" is the denominating currency, and "8890.0" is the exchange rate of micro Luna in micro KRW from the voter's point of view.

"salt" should match the salt used to generate the SHA256 hex in the associated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle vote 1234 8890.0ukrw terravaloper1....
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			rate, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange rate should be formatted as DecCoin", rate)
			}

			// Get from address
			voter := cliCtx.GetFromAddress()
			denom := rate.Denom
			amount := rate.Amount

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

			msg := types.NewMsgExchangeRateVote(amount, salt, denom, voter, validator)
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
Delegate the permission to submit exchange rate votes for the oracle to an address.

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

			msg := types.NewMsgDelegateFeedConsent(validator, feeder)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
