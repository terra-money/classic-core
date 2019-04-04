package cli

import (
	"fmt"
	"strings"
	"github.com/terra-project/core/types/assets"
	"github.com/terra-project/core/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DenomList is array of denom
type DenomList []string

func (dl DenomList) String() (out string) {
	out = strings.Join(dl, "\n")
	return
}

// GetCmdQueryPrice implements the query price command.
func GetCmdQueryPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryPrice,
		Short: "Query the current price of a denom asset",
		Long: strings.TrimSpace(`
Query the current price of a denom asset. You can find the current list of active denoms by running: terracli query oracle active

$ terracli query oracle price --denom mkrw
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

			var price sdk.Dec
			cdc.MustUnmarshalJSON(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	cmd.Flags().String(flagDenom, "", "target denom to get the price")
	return cmd
}

// GetCmdQueryActive implements the query active command.
func GetCmdQueryActive(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryActive,
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

			var actives DenomList
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
		Short: "Query outstanding oracle votes, filtered by denom and voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle votes, filtered by denom and voter address.

$ terracli query oracle votes --denom="musd" --voter="terrad8duyufdshs..."

returns oracle votes submitted by terrad8duyufdshs... for denom musd 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			if len(denom) == 0 {
				return fmt.Errorf("--denom flag is required")
			}

			if !assets.IsValidDenom(denom) {
				return fmt.Errorf("The denom is not known: %s", denom)
			}

			// Check voter address exists, then valids
			var voterAddress sdk.AccAddress

			bechVoterAddr := viper.GetString(flagVoter)
			if len(bechVoterAddr) != 0 {
				var err error

				voterAddress, err = sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			params := oracle.NewQueryVoteParams(voterAddress, denom)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes oracle.PriceBallot
			cdc.MustUnmarshalJSON(res, &matchingVotes)

			return cliCtx.PrintOutput(matchingVotes)
		},
	}

	cmd.Flags().String(flagDenom, "", "filter by votes matching the denom")
	cmd.Flags().String(flagVoter, "", "(optional) filter by votes by voter")

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   oracle.QueryParams,
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
