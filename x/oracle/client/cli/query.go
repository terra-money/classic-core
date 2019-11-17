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
		GetCmdQueryExchangeRate(cdc),
		GetCmdQueryVotes(cdc),
		GetCmdQueryPrevotes(cdc),
		GetCmdQueryActive(cdc),
		GetCmdQueryParams(cdc),
		GetCmdQueryFeederDelegation(cdc),
		GetCmdQueryMissCounter(cdc),
	)...)

	return oracleQueryCmd

}

// GetCmdQueryExchangeRate implements the query rate command.
func GetCmdQueryExchangeRate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rate [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current Luna exchange rate w.r.t an asset",
		Long: strings.TrimSpace(`
Query the current exchange rate of Luna with an asset. You can find the current list of active denoms by running: terracli query oracle active

$ terracli query oracle exchange-rate ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

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

			var rate sdk.Dec
			cdc.MustUnmarshalJSON(res, &rate)
			return cliCtx.PrintOutput(rate)
		},
	}
	return cmd
}

// GetCmdQueryActive implements the query actives command.
func GetCmdQueryActive(cdc *codec.Codec) *cobra.Command {
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

			var actives types.DenomList
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

			var delegatee sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &delegatee)
			return cliCtx.PrintOutput(delegatee)
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
