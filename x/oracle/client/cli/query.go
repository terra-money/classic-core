package cli

import (
	"fmt"
	"strings"

	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdQueryPrice implements the query price command.
func GetCmdQueryPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryPrice,
		Args:  cobra.NoArgs,
		Short: "Query the current Luna exchange rate w.r.t an asset",
		Long: strings.TrimSpace(`
Query the current exchange rate of Luna with an asset. You can find the current list of active denoms by running: terracli query oracle active

$ terracli query oracle price --denom ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			if denom == "" {
				return fmt.Errorf("--denom flag is required")
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, oracle.QueryPrice, denom), nil)
			if err != nil {
				return err
			}

			var price oracle.QueryPriceResponse
			cdc.MustUnmarshalJSON(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	cmd.Flags().String(flagDenom, "", "target denom to get the price")

	cmd.MarkFlagRequired(flagDenom)
	return cmd
}

// GetCmdQueryActive implements the query active command.
func GetCmdQueryActive(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryActive,
		Args:  cobra.NoArgs,
		Short: "Query the active list of Terra assets recognized by the oracle",
		Long: strings.TrimSpace(`
Query the active list of Terra assets recognized by the oracle.

$ terracli query oracle active
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryActive), nil)
			if err != nil {
				return err
			}

			var actives oracle.QueryActiveResponse
			cdc.MustUnmarshalJSON(res, &actives)
			return cliCtx.PrintOutput(actives)
		},
	}

	return cmd
}

// GetCmdQueryVotes implements the query vote command.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryVotes,
		Args:  cobra.NoArgs,
		Short: "Query outstanding oracle votes, filtered by denom and voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle votes, filtered by denom and voter address.

$ terracli query oracle votes --denom="uusd" --validator="terravaloper..."

returns oracle votes submitted by the validator for the denom uusd 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)

			// Check voter address exists, then valids
			var voterAddress sdk.ValAddress

			bechVoterAddr := viper.GetString(flagValidator)
			if len(bechVoterAddr) != 0 {
				var err error

				voterAddress, err = sdk.ValAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			params := oracle.NewQueryVotesParams(voterAddress, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes oracle.QueryVotesResponse
			cdc.MustUnmarshalJSON(res, &matchingVotes)

			return cliCtx.PrintOutput(matchingVotes)
		},
	}

	cmd.Flags().String(flagDenom, "", "filter by votes matching the denom")
	cmd.Flags().String(flagValidator, "", "(optional) filter by votes by validator")

	cmd.MarkFlagRequired(flagDenom)

	return cmd
}

// GetCmdQueryPrevotes implements the query prevote command.
func GetCmdQueryPrevotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryPrevotes,
		Args:  cobra.NoArgs,
		Short: "Query outstanding oracle prevotes, filtered by denom and voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle prevotes, filtered by denom and voter address.

$ terracli query oracle prevotes --denom="uusd" --validator="terravaloper..."

returns oracle prevotes submitted by the validator for denom uusd 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)

			// Check voter address exists, then valids
			var voterAddress sdk.ValAddress

			bechVoterAddr := viper.GetString(flagValidator)
			if len(bechVoterAddr) != 0 {
				var err error

				voterAddress, err = sdk.ValAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			params := oracle.NewQueryPrevotesParams(voterAddress, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryPrevotes), bz)
			if err != nil {
				return err
			}

			var matchingPrevotes oracle.QueryPrevotesResponse
			cdc.MustUnmarshalJSON(res, &matchingPrevotes)

			return cliCtx.PrintOutput(matchingPrevotes)
		},
	}

	cmd.Flags().String(flagDenom, "", "filter by prevotes matching the denom")
	cmd.Flags().String(flagValidator, "", "(optional) filter by prevotes by validator")

	cmd.MarkFlagRequired(flagDenom)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryParams,
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle params",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryParams), nil)
			if err != nil {
				return err
			}

			var params oracle.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}

// GetCmdQueryFeederDelegation implements the query feeder delegation command
func GetCmdQueryFeederDelegation(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryFeederDelegation,
		Short: "Query the oracle feeder delegate account",
		Long: strings.TrimSpace(`
Query the account the validator's oracle voting right is delegated to.

$ terracli query oracle feeder --validator terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valString := viper.GetString(flagValidator)
			if len(valString) == 0 {
				return fmt.Errorf("--validator flag is required")
			}
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := oracle.NewQueryFeederDelegationParams(validator)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryFeederDelegation), bz)
			if err != nil {
				return err
			}

			var delegatee oracle.QueryFeederDelegationResponse
			cdc.MustUnmarshalJSON(res, &delegatee)
			return cliCtx.PrintOutput(delegatee)
		},
	}

	cmd.Flags().String(flagValidator, "", "validator which owns the oracle voting rights")

	cmd.MarkFlagRequired(flagValidator)

	return cmd
}
