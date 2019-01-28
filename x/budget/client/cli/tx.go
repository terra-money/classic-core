package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"terra/x/budget"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTitle       = "title"
	flagDescription = "description"
	flagDeposit     = "deposit"
	flagVoter       = "voter"
	flagOption      = "option"
	flagState       = "state"
	flagNumLimit    = "limit"
	flagPrgram      = "program"
)

type program struct {
	Title       string
	Description string
	Type        string
	Deposit     string
}

var programFlags = []string{
	flagTitle,
	flagDescription,
	flagDeposit,
}

// GetCmdSubmitProgram implements submitting a program transaction command.
func GetCmdSubmitProgram(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-program",
		Short: "Submit a program along with an initial deposit",
		Long: strings.TrimSpace(`
Submit a program along with an initial deposit. program title, description, type and deposit can be given directly or through a program JSON file. For example:

$ terracli budget submit-program --program="path/to/program.json" --from mykey

where program.json contains:

{
  "title": "Test program",
  "description": "My awesome program",
  "type": "Text",
  "deposit": "10terra"
}

is equivalent to

$ terracli budget submit-program --title="Test program" --description="My awesome program" --type="Text" --deposit="10test" --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			program, err := parseSubmitProgramFlags()
			if err != nil {
				return err
			}

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			// Get from address
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// Pull associated account
			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// Find deposit amount
			amount, err := sdk.ParseCoins(program.Deposit)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(amount) {
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			ProgramType, err := budget.ProgramTypeFromString(program.Type)
			if err != nil {
				return err
			}

			msg := budget.NewMsgSubmitProgram(program.Title, program.Description, from, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			// Build and sign the transaction, then broadcast to Tendermint
			// programID must be returned, and it is a part of response.
			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTitle, "", "title of program")
	cmd.Flags().String(flagDescription, "", "description of program")
	cmd.Flags().String(flagDeposit, "", "deposit of program")
	cmd.Flags().String(flagProgram, "", "program file path (if this path is given, other program flags are ignored)")

	return cmd
}

func parseSubmitProgramFlags() (*program, error) {
	program := &program{}
	programFile := viper.GetString(flagProgram)

	if programFile == "" {
		program.Title = viper.GetString(flagTitle)
		program.Description = viper.GetString(flagDescription)
		program.Deposit = viper.GetString(flagDeposit)
		return program, nil
	}

	for _, flag := range programFlags {
		if viper.GetString(flag) != "" {
			return nil, fmt.Errorf("--%s flag provided alongside --program, which is a noop", flag)
		}
	}

	contents, err := ioutil.ReadFile(programFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, program)
	if err != nil {
		return nil, err
	}

	return program, nil
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [program-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active program, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(`
Submit a vote for an acive program. You can find the program-id by running terracli query budget programs:

$ terracli tx budget vote 1 yes --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			// Get voting address
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			// check to see if the program is in the store
			_, err = queryProgram(programID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch program-id %d: %s", programID, err)
			}

			// Find out which vote option user chose
			byteVoteOption, err := budget.VoteOptionFromString(govClientUtils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := budget.NewMsgVote(from, programID, byteVoteOption)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// If generate only print the transaction
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
