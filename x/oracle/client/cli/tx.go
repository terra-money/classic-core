package cli

import (
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
		GetCmdAggregateExchangeRatePrevote(cdc),
		GetCmdAggregateExchangeRateVote(cdc),
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
The purpose of prevote is to hide exchange rate vote with hash which is formatted 
as hex string in SHA256("{salt}:{exchange_rate}:{denom}:{voter}")

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
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange_rate should be formatted as DecCoin; %s", rate, err.Error())
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

			hash := types.GetVoteHash(salt, amount, denom, validator)

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

"salt" should match the salt used to generate the SHA256 hex in the aggregated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle vote 1234 8890.0ukrw terravaloper1....
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			rate, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange_rate should be formatted as DecCoin; %s", rate, err.Error())
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

// GetCmdAggregateExchangeRatePrevote will create a aggregateExchangeRatePrevote tx and sign it with the given key.
func GetCmdAggregateExchangeRatePrevote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [salt] [exchange_rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate prevote for the exchange rates of Luna",
		Long: strings.TrimSpace(`
Submit an oracle aggregate prevote for the exchange rates of Luna denominated in multiple denoms.
The purpose of aggregate prevote is to hide aggregate exchange rate vote with hash which is formatted 
as hex string in SHA256("{salt}:{exchange_rate}{denom},...,{exchange_rate}{denom}:{voter}")

# Aggregate Prevote
$ terracli tx oracle aggregate-prevote 1234 8888.0ukrw,1.243uusd,0.99usdr 

where "ukrw,uusd,usdr" is the denominating currencies, and "8888.0,1.243,0.99" is the exchange rates of micro Luna in micro denoms from the voter's point of view.

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle aggregate-prevote 1234 8888.0ukrw,1.243uusd,0.99usdr terravaloper1...
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err := types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rates {%s} is not a valid format; exchange_rate should be formatted as DecCoins; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := cliCtx.GetFromAddress()

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

			hash := types.GetAggregateVoteHash(salt, exchangeRatesStr, validator)

			msg := types.NewMsgAggregateExchangeRatePrevote(hash, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdAggregateExchangeRateVote will create a aggregateExchangeRateVote tx and sign it with the given key.
func GetCmdAggregateExchangeRateVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [salt] [exchange_rates] [validator]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Submit an oracle aggregate vote for the exchange_rates of Luna",
		Long: strings.TrimSpace(`
Submit a aggregate vote for the exchange_rates of Luna w.r.t the input denom. Companion to a prevote submitted in the previous vote period. 

$ terracli tx oracle aggregate-vote 1234 8888.0ukrw,1.243uusd,0.99usdr 

where "ukrw,uusd,usdr" is the denominating currencies, and "8888.0,1.243,0.99" is the exchange rates of micro Luna in micro denoms from the voter's point of view.

"salt" should match the salt used to generate the SHA256 hex in the aggregated pre-vote. 

If voting from a voting delegate, set "validator" to the address of the validator to vote on behalf of:
$ terracli tx oracle aggregate-vote 1234 8888.0ukrw,1.243uusd,0.99usdr terravaloper1....
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			salt := args[0]
			exchangeRatesStr := args[1]
			_, err := types.ParseExchangeRateTuples(exchangeRatesStr)
			if err != nil {
				return fmt.Errorf("given exchange_rate {%s} is not a valid format; exchange rate should be formatted as DecCoin; %s", exchangeRatesStr, err.Error())
			}

			// Get from address
			voter := cliCtx.GetFromAddress()

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

			msg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRatesStr, voter, validator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
