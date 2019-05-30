package cli

import (
	"encoding/hex"
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
	flagProofPrice = "proof-price"
	flagProofSalt  = "proof-salt"

	flagSalt  = "salt"
	flagPrice = "price"
	flagHash  = "Hash"

	flagDenom     = "denom"
	flagValidator = "validator"
	flagFeeder    = "feeder"

	flagOffline = "offline"
)

// GetCmdPriceVote will create a send tx and sign it with the given key.
func GetCmdPriceVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Submit an oracle vote for the price of Luna",
		Long: strings.TrimSpace(`
Submit an oracle prevote and vote for the price of Luna denominated in the input denom.
This message has piggybacking structure. Firstly users should submit the hash of real vote (prevote). After then in next vote period,
users should submit vote with proof (= price and slat) to verify prevote and make real vote.

# Prevote and Vote Both
$ terracli oracle vote --denom "ukrw" --proof-price "8890" --hash "72f374291b0428453bf481ec9d4b0b2440299b62" --proof-salt "1234" --from mykey
$ terracli oracle vote --denom "ukrw" --proof-price "8890" --price "8888" --salt "4321" --proof-salt "1234" --from mykey

# Vote Only
$ terracli oracle vote --denom "ukrw" --proof-price "8890" --proof-salt "1234" --from mykey

# Prevote Only
$ terracli oracle vote --denom "ukrw" --price "8888" --salt "4321"--from mykey

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
			denom := viper.GetString(flagDenom)
			priceStr := viper.GetString(flagPrice)
			hash := viper.GetString(flagHash)
			salt := viper.GetString(flagSalt)
			proofPriceStr := viper.GetString(flagProofPrice)
			proofSalt := viper.GetString(flagProofSalt)

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
			var proofPrice sdk.Dec
			if len(proofPriceStr) == 0 {
				proofPrice = sdk.ZeroDec()
			} else {
				var err sdk.Error
				proofPrice, err = sdk.NewDecFromStr(proofPriceStr)
				if err != nil {
					return fmt.Errorf("given price {%s} is not a valid format; price should be formatted as float", proofPriceStr)
				}
			}

			if len(hash) == 0 && (len(priceStr) > 0 && len(salt) > 0) {
				price, err := sdk.NewDecFromStr(priceStr)
				if err != nil {
					return fmt.Errorf("given price {%s} is not a valid format; price should be formatted as float", priceStr)
				}

				hashBytes, err2 := oracle.VoteHash(salt, price, denom, validator)
				if err2 != nil {
					return err2
				}

				hash = hex.EncodeToString(hashBytes)
			}

			msg := oracle.NewMsgPriceFeed(hash, proofSalt, denom, voter, validator, proofPrice)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, offline)
		},
	}

	cmd.Flags().String(flagDenom, "", "denominating currency")
	cmd.Flags().String(flagValidator, "", "validator on behalf of which to vote (for delegated feeders)")
	cmd.Flags().String(flagHash, "", "hex string; hash of next vote; empty == skip prevote")
	cmd.Flags().String(flagPrice, "", "price of Luna in denom currency is to make provte hash; this field is required to submit prevote in case absense of hash")
	cmd.Flags().String(flagSalt, "", "salt is to make prevote hash; this field is required to submit prevote in case  absense of hash")
	cmd.Flags().String(flagProofPrice, "", "proof price of Luna in denom currency was used to make prevote hash; initial prevote does not require this field")
	cmd.Flags().String(flagProofSalt, "", "proof salt was used to make prevote hash; initial prevote does not require this field")
	cmd.Flags().Bool(flagOffline, false, " Offline mode; Do not query a full node")

	cmd.MarkFlagRequired(flagDenom)

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

$ terracli oracle set-feeder --feeder terra1...... --from mykey

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

			feederStr := viper.GetString(flagFeeder)

			feeder, err := sdk.AccAddressFromBech32(feederStr)
			if err != nil {
				return err
			}

			msg := oracle.NewMsgDelegateFeederPermission(validator, feeder)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagFeeder, "", "account the voting right will be delegated to")

	cmd.MarkFlagRequired(flagFeeder)

	return cmd
}
