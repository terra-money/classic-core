package cli

import (
	"fmt"
	"strconv"

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
	flagExecutor    = "executor"
	flagVoter       = "voter"
	flagOption      = "option"
	flagNumLimit    = "limit"
	flagProgram     = "program"
	flagProgramID   = "program-id"

	queryRoute = "budget"
)

type program struct {
	Title       string
	Description string
	Deposit     string
	Executor    string
}

var programFlags = []string{
	flagTitle,
	flagDescription,
	flagExecutor,
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
  "description": "My awesome program (include a website link for impact)",
  "submitter": terra1nk5lsuvy0rcfjcdr8au8za0wq25rat0qa07p6t,
  "executor": terra1nk5lsuvy0rcfjcdr8au8za0wq25rat0qa07p6t,
  "deposit": "10terra"
}

is equivalent to

$ terracli budget submit-program --title="Test program" --description="My awesome program" ... --from mykey
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
			from := cliCtx.GetFromAddress()

			// Pull associated account
			submitter, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			executor, err := cliCtx.GetAccount([]byte(program.Executor))
			if err != nil {
				return err
			}

			msg := budget.NewMsgSubmitProgram(program.Title, program.Description, submitter.GetAddress(), executor.GetAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagTitle, "", "title of program")
	cmd.Flags().String(flagDescription, "", "(optional) description of program")
	cmd.Flags().String(flagExecutor, "", "executor of program")
	cmd.Flags().String(flagProgram, "", "program file path (if this path is given, other program flags are ignored)")

	return cmd
}

func parseSubmitProgramFlags() (*program, error) {
	program := &program{}
	programFile := viper.GetString(flagProgram)

	if programFile == "" {
		program.Title = viper.GetString(flagTitle)
		program.Description = viper.GetString(flagDescription)
		program.Executor = viper.GetString(flagExecutor)
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
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [program-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active program, options: yes or no",
		Long: strings.TrimSpace(`
Submit a vote for an active program. 

You can find the program-id of active programs by running terracli query budget actives
You can find the program-id of inactive (candidate) programs by running terracli query budget candidates

$ terracli tx budget vote 1 yes --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			// Get voting address
			from := cliCtx.GetFromAddress()

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			// Find out which vote option user chose
			var option bool
			if args[1] == "yes" {
				option = true
			}

			// Build vote message and run basic validation
			msg := budget.NewMsgVoteProgram(programID, option, from)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdWithdrawProgram(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "withdraw a program from consideration",
		Long: strings.TrimSpace(`
Withdraw a program from consideration. The deposit is only refunded if the program is already in the active set. 

$ terracli tx budget withdraw 1 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			// Get voting address
			from := cliCtx.GetFromAddress()

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			// Build vote message and run basic validation
			msg := budget.NewMsgWithdrawProgram(programID, from)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	return cmd
}
