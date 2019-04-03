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
	flagProgram     = "program"
	flagProgramID   = "program-id"
	flagOption      = "option"
)

type program struct {
	Title       string
	Description string
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
  "executor": terra1nk5lsuvy0rcfjcdr8au8za0wq25rat0qa07p6t,
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

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			// Get from address
			from := cliCtx.GetFromAddress()

			// Pull associated account
			submitter, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			submitterCoins := submitter.GetCoins()

			// Query params to get deposit amount
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryParams), nil)
			if err != nil {
				return err
			}

			var params budget.Params
			cdc.MustUnmarshalJSON(res, &params)

			// Check submitter has enough coins to pay a deposit
			if submitterCoins.AmountOf(params.Deposit.Denom).LT(params.Deposit.Amount) {
				return fmt.Errorf(strings.TrimSpace(`
					account %s has insufficient amount of coins to pay a deposit.\n
					Required: %s\n
					Given:    %s\n`),
					from, params.Deposit.String(), submitterCoins.String())
			}

			// Get executor address
			executorAddr, err := sdk.AccAddressFromBech32(program.Executor)
			if err != nil {
				return err
			}

			msg := budget.NewMsgSubmitProgram(program.Title, program.Description, from, executorAddr)
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

	if len(programFile) == 0 {
		program.Title = viper.GetString(flagTitle)
		program.Description = viper.GetString(flagDescription)
		program.Executor = viper.GetString(flagExecutor)

		// Check title existence
		if len(program.Title) == 0 {
			return nil, fmt.Errorf("--%s flag is required", flagTitle)
		}

		// Check executor existence
		if len(program.Executor) == 0 {
			return nil, fmt.Errorf("--%s flag is required", flagExecutor)
		}

		return program, nil
	}

	for _, flag := range programFlags {
		if len(viper.GetString(flag)) > 0 {
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
		Use:   "vote",
		Short: "Vote for an candidate/active program, options: yes or no",
		Long: strings.TrimSpace(`
Submit a vote for an candidate/active program.

You can find the program-id of active programs by running terracli query budget actives
You can find the program-id of candidate programs by running terracli query budget candidates

$ terracli tx budget vote --program-id 1  --option yes --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			// Get voting address
			from := cliCtx.GetFromAddress()

			// Check flag program-id is given
			programStrID := viper.GetString(flagProgramID)
			if len(programStrID) == 0 {
				return fmt.Errorf("--program-id flag is required")
			}

			// Validate that the program id is a uint
			programID, err := strconv.ParseUint(programStrID, 10, 64)
			if err != nil {
				return fmt.Errorf("given program-id {%s} is not a valid format; program-id should be formatted as integer", programStrID)
			}

			// Find out which vote option user chose
			var option bool
			optionStr := viper.GetString(flagOption)
			if optionStr == "yes" || optionStr == "true" {
				option = true
			} else if optionStr == "no" || optionStr == "false" {
				option = false
			} else {
				return fmt.Errorf(`given option {%s} is not valid format;\n option should be formatted as "yes" or "no"`, optionStr)
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

	cmd.Flags().String(flagProgramID, "", "the program ID to vote")
	cmd.Flags().String(flagOption, "", "yes or no")

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdWithdrawProgram(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw a program from consideration",
		Long: strings.TrimSpace(`
Withdraw a program from consideration. The deposit is only refunded if the program is already in the active set. 

$ terracli tx budget withdraw --program-id 1 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			// Get voting address
			from := cliCtx.GetFromAddress()

			// validate that the program id is a uint
			programStrID := viper.GetString(flagProgramID)
			if programStrID == "" {
				return fmt.Errorf("--program-id flag is required")
			}

			programID, err := strconv.ParseUint(programStrID, 10, 64)
			if err != nil {
				return fmt.Errorf("given program-id %s not a valid int, please input a valid program-id", programStrID)
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

	cmd.Flags().String(flagProgramID, "", "the program ID to withdraw")

	return cmd
}
