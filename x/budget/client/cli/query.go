package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/terra-project/core/x/budget"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetCmdQueryProgram implements the query program command.
func GetCmdQueryProgram(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program",
		Args:  cobra.NoArgs,
		Short: "Query details of a single program",
		Long: strings.TrimSpace(`
Query details for a program. 

You can find the program-id of active programs by running terracli query budget actives
You can find the program-id of inactive (candidate) programs by running terracli query budget candidates

$ terracli query budget program --program-id 1
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the program id is a uint
			programIDStr := viper.GetString(flagProgramID)
			if len(programIDStr) == 0 {
				return fmt.Errorf("--program-id flag is required")
			}

			programID, err := strconv.ParseUint(programIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("given program-id %s not a valid format\n, program-id should be formatted as integer", programIDStr)
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d", budget.QuerierRoute, budget.QueryProgram, programID), nil)
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

	cmd.Flags().String(flagProgramID, "", "the program ID to query")

	cmd.MarkFlagRequired(flagProgramID)

	return cmd
}

// GetCmdQueryActives implements a query actives command.
func GetCmdQueryActives(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryActiveList,
		Args:  cobra.NoArgs,
		Short: "Query active programs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryActiveList), nil)
			if err != nil {
				return err
			}

			var actives budget.QueryActiveListResponse
			cdc.MustUnmarshalJSON(res, &actives)

			return cliCtx.PrintOutput(actives)
		},
	}

	return cmd
}

// GetCmdQueryCandidates implements the query program candidates command.
func GetCmdQueryCandidates(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryCandidateList,
		Args:  cobra.NoArgs,
		Short: "Query candidate programs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryCandidateList), nil)
			if err != nil {
				return err
			}

			var candidates budget.QueryCandidateListResponse
			cdc.MustUnmarshalJSON(res, &candidates)

			return cliCtx.PrintOutput(candidates)
		},
	}

	return cmd
}

// GetCmdQueryVotes implements the command to query for program votes.
func GetCmdQueryVotes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryVotes,
		Args:  cobra.NoArgs,
		Short: "Query votes, filtered by voter and program id ",
		Long: strings.TrimSpace(`
Query vote details filtered by voter address and program id.

Example:
$ terracli query budget votes --program-id 1 --voter terra1nk5lsuvy0rcfjcdr8au8za0wq25rat0qa07p6t
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := budget.QueryVotesParams{}

			// Get voting address
			voterAddrStr := viper.GetString(flagVoter)
			if len(voterAddrStr) > 0 {
				voterAddress, err := sdk.AccAddressFromBech32(voterAddrStr)
				if err != nil {
					return err
				}

				params.Voter = voterAddress
			}

			programIDStr := viper.GetString(flagProgramID)

			// validate that the program id is a uint
			programID, err := strconv.ParseUint(programIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("program-id %s not a valid int, please input a valid program-id", programIDStr)
			}

			params.ProgramID = programID

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes budget.QueryVotesResponse
			cdc.MustUnmarshalJSON(res, &matchingVotes)

			return cliCtx.PrintOutput(matchingVotes)
		},
	}

	cmd.Flags().String(flagProgramID, "0", "(optional) the program ID to query; defalut 0 for all programs")
	cmd.Flags().String(flagVoter, "", "(optional) voter for the program")

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   budget.QueryParams,
		Args:  cobra.NoArgs,
		Short: "Query the current budget params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", budget.QuerierRoute, budget.QueryParams), nil)
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
