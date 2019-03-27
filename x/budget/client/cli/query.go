package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"terra/x/budget"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetCmdQueryProgram implements the query program command.
func GetCmdQueryProgram(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program [program-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single program",
		Long: strings.TrimSpace(`
Query details for a program. 

You can find the program-id of active programs by running terracli query budget actives
You can find the program-id of inactive (candidate) programs by running terracli query budget candidates

$ terracli query budget program 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid uint, please input a valid program-id", args[0])
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", queryRoute, budget.QueryProgram, programID), nil)
			if err != nil {
				return err
			}

			var program budget.Program
			err = cdc.UnmarshalJSON(res, &program)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(program)
		},
	}

	return cmd
}

// GetCmdQueryActives implements a query actives command.
func GetCmdQueryActives(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryActiveList,
		Short: "Query active programs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryActiveList), nil)
			if err != nil {
				return err
			}

			var actives []budget.Program
			cdc.MustUnmarshalJSON(res, &actives)

			if len(actives) == 0 {
				fmt.Println("No active Programs found")
				return nil
			}

			for _, program := range actives {
				fmt.Println(program.String())
			}

			return nil
		},
	}

	return cmd
}

// GetCmdQueryCandidates implements the query program candidates command.
func GetCmdQueryCandidates(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryCandidateList,
		Short: "Query candidate programs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryActiveList), nil)
			if err != nil {
				return err
			}

			var candidates []budget.Program
			cdc.MustUnmarshalJSON(res, &candidates)

			if len(candidates) == 0 {
				fmt.Println("No candidates Programs found")
				return nil
			}

			for _, program := range candidates {
				fmt.Println(program.String())
			}

			return nil
		},
	}

	return cmd
}

// GetCmdQueryVotes implements the command to query for program votes.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryVotes,
		Short: "Query votes, filtered by voterAddress ",
		Long: strings.TrimSpace(`
Query vote details for a single program by its identifier.

Example:
$ terracli query budget votes 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := budget.QueryVotesParams{}

			// Get voting address
			voterAddrStr := viper.GetString(flagVoter)
			if len(voterAddrStr) > 0 {
				acc, err := cliCtx.GetAccount([]byte(voterAddrStr))
				if err != nil {
					return err
				}

				params.Voter = acc.GetAddress()
			}

			programIDStr := viper.GetString(flagProgramID)
			if len(programIDStr) > 0 {
				// validate that the program id is a uint
				programID, err := strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", args[0])
				}

				params.ProgramID = programID
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes []budget.MsgVoteProgram
			cdc.MustUnmarshalJSON(res, &matchingVotes)

			if len(matchingVotes) == 0 {
				fmt.Println("No matching votes found")
				return nil
			}

			for _, vote := range matchingVotes {
				fmt.Println(vote.String())
			}

			return nil
		},
	}

	cmd.Flags().String(flagVoter, "", "voter for the program")

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryParams,
		Short: "Query the current budget params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, budget.QueryParams), nil)
			if err != nil {
				return err
			}

			var params budget.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}
