package cli

import (
	"fmt"
	"strings"

	"github.com/terra-project/core/x/oracle/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	oracleQueryCmd := &cobra.Command{
		Use:                        "oracle",
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryExchangeRates(cdc),
		GetCmdQueryVotes(cdc),
		GetCmdQueryPrevotes(cdc),
		GetCmdQueryActives(cdc),
		GetCmdQueryParams(cdc),
		GetCmdQueryFeederDelegation(cdc),
		GetCmdQueryMissCounter(cdc),
		GetCmdQueryAggregatePrevote(cdc),
		GetCmdQueryAggregateVote(cdc),
		GetCmdQueryVoteTargets(cdc),
		GetCmdQueryTobinTaxes(cdc),
	)...)

	return oracleQueryCmd

}

// GetCmdQueryExchangeRates implements the query rate command.
func GetCmdQueryExchangeRates(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rates [denom]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the current Luna exchange rate w.r.t an asset",
		Long: strings.TrimSpace(`
Query the current exchange rate of Luna with an asset. 
You can find the current list of active denoms by running

$ terracli query oracle exchange-rates 

Or, can filter with denom

$ terracli query oracle exchange-rates ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryExchangeRates), nil)
				if err != nil {
					return err
				}

				var rate sdk.DecCoins
				cdc.MustUnmarshalJSON(res, &rate)
				return cliCtx.PrintOutput(rate)
			}

			denom := args[0]
			params := types.NewQueryExchangeRateParams(denom)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryExchangeRate), bz)
			if err != nil {
				return err
			}

			var rates sdk.Dec
			cdc.MustUnmarshalJSON(res, &rates)
			return cliCtx.PrintOutput(rates)

		},
	}
	return cmd
}

// Actives receiver struct
type Denoms []string

// String implements fmt.Stringer interface
func (a Denoms) String() string {
	return strings.Join(a, ",")
}

// GetCmdQueryActives implements the query actives command.
func GetCmdQueryActives(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actives",
		Args:  cobra.NoArgs,
		Short: "Query the active list of Terra assets recognized by the oracle",
		Long: strings.TrimSpace(`
Query the active list of Terra assets recognized by the types.

$ terracli query oracle actives
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryActives), nil)
			if err != nil {
				return err
			}

			var actives Denoms
			cdc.MustUnmarshalJSON(res, &actives)
			return cliCtx.PrintOutput(actives)
		},
	}

	return cmd
}

// GetCmdQueryVotes implements the query vote command.
func GetCmdQueryVotes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [denom] [validator]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query outstanding oracle votes, filtered by denom and voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle votes, filtered by denom and voter address.

$ terracli query oracle votes uusd terravaloper...
$ terracli query oracle votes uusd 

returns oracle votes submitted by the validator for the denom uusd 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]

			// Check voter address exists, then valids
			var voterAddress sdk.ValAddress
			if len(args) >= 2 {
				bechVoterAddr := args[1]

				var err error
				voterAddress, err = sdk.ValAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			params := types.NewQueryVotesParams(voterAddress, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes types.ExchangeRateVotes
			cdc.MustUnmarshalJSON(res, &matchingVotes)

			return cliCtx.PrintOutput(matchingVotes)
		},
	}

	return cmd
}

// GetCmdQueryPrevotes implements the query prevote command.
func GetCmdQueryPrevotes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prevotes [denom] [validator]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query outstanding oracle prevotes, filtered by denom and voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle prevotes, filtered by denom and voter address.

$ terracli query oracle prevotes uusd terravaloper...
$ terracli query oracle prevotes uusd

returns oracle prevotes submitted by the validator for denom uusd 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]

			// Check voter address exists, then valids
			var voterAddress sdk.ValAddress
			if len(args) >= 2 {
				bechVoterAddr := args[1]

				var err error
				voterAddress, err = sdk.ValAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			params := types.NewQueryPrevotesParams(voterAddress, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPrevotes), bz)
			if err != nil {
				return err
			}

			var matchingPrevotes types.ExchangeRatePrevotes
			cdc.MustUnmarshalJSON(res, &matchingPrevotes)

			return cliCtx.PrintOutput(matchingPrevotes)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}

// GetCmdQueryFeederDelegation implements the query feeder delegation command
func GetCmdQueryFeederDelegation(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the oracle feeder delegate account",
		Long: strings.TrimSpace(`
Query the account the validator's oracle voting right is delegated to.

$ terracli query oracle feeder terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryFeederDelegationParams(validator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeederDelegation), bz)
			if err != nil {
				return err
			}

			var delegate sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &delegate)
			return cliCtx.PrintOutput(delegate)
		},
	}

	return cmd
}

// GetCmdQueryMissCounter implements the query miss counter of the validator command
func GetCmdQueryMissCounter(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miss [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the # of the miss count",
		Long: strings.TrimSpace(`
Query the # of vote periods missed in this oracle slash window.

$ terracli query oracle miss terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryMissCounterParams(validator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMissCounter), bz)
			if err != nil {
				return err
			}

			var missCounter int64
			cdc.MustUnmarshalJSON(res, &missCounter)
			return cliCtx.PrintOutput(sdk.NewInt(missCounter))
		},
	}

	return cmd
}

// GetCmdQueryAggregatePrevote implements the query aggregate prevote of the validator command
func GetCmdQueryAggregatePrevote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query outstanding oracle aggregate prevote, filtered by voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate prevote, filtered by voter address.

$ terracli query oracle aggregate-prevote terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryAggregatePrevoteParams(validator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregatePrevote), bz)
			if err != nil {
				return err
			}

			var aggregatePrevote types.AggregateExchangeRatePrevote
			cdc.MustUnmarshalJSON(res, &aggregatePrevote)
			return cliCtx.PrintOutput(aggregatePrevote)
		},
	}

	return cmd
}

// GetCmdQueryAggregateVote implements the query aggregate prevote of the validator command
func GetCmdQueryAggregateVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query outstanding oracle aggregate vote, filtered by voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate vote, filtered by voter address.

$ terracli query oracle aggregate-vote terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryAggregateVoteParams(validator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregateVote), bz)
			if err != nil {
				return err
			}

			var aggregateVote types.AggregateExchangeRateVote
			cdc.MustUnmarshalJSON(res, &aggregateVote)
			return cliCtx.PrintOutput(aggregateVote)
		},
	}

	return cmd
}

// GetCmdQueryVoteTargets implements the query params command.
func GetCmdQueryVoteTargets(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-targets",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle vote targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVoteTargets), nil)
			if err != nil {
				return err
			}

			var voteTargets Denoms
			cdc.MustUnmarshalJSON(res, &voteTargets)
			return cliCtx.PrintOutput(voteTargets)
		},
	}

	return cmd
}

// GetCmdQueryTobinTaxes implements the query params command.
func GetCmdQueryTobinTaxes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tobin-taxes [denom]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the current Oracle tobin taxes.",
		Long: strings.TrimSpace(`
Query the current Oracle tobin taxes.

$ terracli query oracle tobin-taxes

Or, can filter with denom

$ terracli query oracle tobin-taxes ukrw

Or, can 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if len(args) == 0 {
				res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTaxes), nil)
				if err != nil {
					return err
				}

				var tobinTaxes types.DenomList
				cdc.MustUnmarshalJSON(res, &tobinTaxes)
				return cliCtx.PrintOutput(tobinTaxes)
			}

			denom := args[0]
			params := types.NewQueryTobinTaxParams(denom)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTax), bz)
			if err != nil {
				return err
			}

			var tobinTax sdk.Dec
			cdc.MustUnmarshalJSON(res, &tobinTax)
			return cliCtx.PrintOutput(tobinTax)
		},
	}

	return cmd
}
