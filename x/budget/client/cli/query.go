package cli

import (
	"fmt"
	"strconv"
	"strings"

	"terra/x/budget"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCmdQueryProgram implements the query program command.
func GetCmdQueryProgram(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single program",
		Long: strings.TrimSpace(`
Query details for a program. You can find the program-id by running terracli query budget program:

$ terracli query budget program 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid uint, please input a valid program-id", args[0])
			}

			// Query the program
			res, err := queryProgram(programID, cliCtx, cdc, queryRoute)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

func queryProgram(ProgramID uint64, cliCtx context.CLIContext, cdc *codec.Codec, queryRoute string) ([]byte, error) {
	// Construct query
	params := budget.NewQueryProgramParams(ProgramID)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}

	// Query store
	res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/program", queryRoute), bz)
	if err != nil {
		return nil, err
	}
	return res, err
}

// GetCmdQueryPrograms implements a query Programs command.
func GetCmdQueryPrograms(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "programs",
		Short: "Query programs with optional filters",
		Long: strings.TrimSpace(`
Query for a all programs. You can filter the returns with the following flags:

$ terracli query budget programs --voter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ terracli query budget programs --status (inactive|rejected|active|legacy)
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositorAddr := viper.GetString(flagDepositor)
			bechVoterAddr := viper.GetString(flagVoter)
			strProgramState := viper.GetString(flagState)
			numLimit := uint64(viper.GetInt64(flagNumLimit))

			var voterAddr sdk.AccAddress
			var ProgramStatus budget.ProgramStatus

			params := budget.NewQueryProgramsParams(ProgramStatus, numLimit, voterAddr)

			if len(bechVoterAddr) != 0 {
				voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
				params.Voter = voterAddr
			}

			if len(strProgramState) != 0 {
				params.ProgramState = strProgramState
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/programs", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingPrograms []budget.program
			err = cdc.UnmarshalJSON(res, &matchingPrograms)
			if err != nil {
				return err
			}

			if len(matchingPrograms) == 0 {
				fmt.Println("No matching Programs found")
				return nil
			}

			for _, program := range matchingPrograms {
				fmt.Printf("  %d - %s\n", program.GetProgramID(), program.GetTitle())
			}

			return nil
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] Programs. Defaults to all Programs")
	cmd.Flags().String(flagVoter, "", "(optional) filter by Programs voted on by voted")
	cmd.Flags().String(flagState, "", "(optional) filter Programs by program state, state: inactive/active/legacied/rejected")

	return cmd
}

// Command to Get a program Information
// GetCmdQueryVote implements the query program vote command.
func GetCmdQueryVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [program-id] [voter-address]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(`
Query details for a single vote on a program given its identifier.

Example:
$ terracli query budget vote 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			ProgramID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			// check to see if the program is in the store
			_, err = queryProgram(ProgramID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch program-id %d: %s", ProgramID, err)
			}

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := budget.NewQueryVoteParams(ProgramID, voterAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/vote", queryRoute), bz)
			if err != nil {
				return err
			}

			var vote budget.Vote
			cdc.UnmarshalJSON(res, &vote)

			if vote.Empty() {
				res, err = gcutils.QueryVoteByTxQuery(cdc, cliCtx, params)
				if err != nil {
					return err
				}
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetCmdQueryVotes implements the command to query for program votes.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a program",
		Long: strings.TrimSpace(`
Query vote details for a single program by its identifier.

Example:
$ terracli query budget votes 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			ProgramID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			params := budget.NewQueryProgramParams(ProgramID)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// check to see if the program is in the store
			res, err := queryProgram(ProgramID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch program-id %d: %s", ProgramID, err)
			}

			var program budget.program
			if err := cdc.UnmarshalJSON(res, &program); err != nil {
				return err
			}

			propStatus := program.GetStatus()
			if !(propStatus == budget.StatusVotingPeriod || propStatus == budget.StatusDepositPeriod) {
				res, err = gcutils.QueryVotesByTxQuery(cdc, cliCtx, params)
			} else {
				res, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", queryRoute), bz)
			}

			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetCmdQueryTally implements the command to query for program tally result.
func GetCmdQueryTally(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tally [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a program vote",
		Long: strings.TrimSpace(`
Query tally of votes on a program. You can find the program-id by running terracli query budget programs:

$ terracli query budget tally 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			ProgramID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
			}

			// check to see if the program is in the store
			_, err = queryProgram(ProgramID, cliCtx, cdc, queryRoute)
			if err != nil {
				return fmt.Errorf("Failed to fetch program-id %d: %s", ProgramID, err)
			}

			// Construct query
			params := budget.NewQueryProgramParams(ProgramID)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// Query store
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tally", queryRoute), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetCmdQueryParams queries the params of the budget process.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters of the budget process",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Query store
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}
