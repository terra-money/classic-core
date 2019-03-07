package cli

import (
	"fmt"
	"math"
	"strings"
	"terra/x/oracle"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDenom = "denom"
	flagPrice = "price"
	flagVoter = "voter"
)

// GetCmdPriceVote will create a send tx and sign it with the given key.
func GetCmdPriceVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [denom] [price]",
		Short: "Submit an oracle vote for the price of Luna",
		Long: strings.TrimSpace(`
Submit an oracle vote for the price of Luna. Reference currency denom and price should be given as input. 

$ terracli oracle vote --denom="krw" --price="8890.12"

where "krw" is the denominating currency, and "8890.12" is the price of Luna in KRW from the voter's point of view. 
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			voterAddress := cliCtx.GetFromAddress()

			// parse denom of the coin to be voted on
			denom := viper.GetString(flagDenom)
			price := viper.GetFloat64(flagPrice)
			cleanPrice := sdk.NewDecWithPrec(int64(math.Round(price*100)), 2)

			msg := oracle.NewPriceFeedMsg(denom, cleanPrice, voterAddress)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagDenom, "", "denominating currency")
	cmd.Flags().String(flagPrice, "", "price of Luna in denom currency")

	return cmd
}

// GetCmdQueryPrice implements the query price command.
func GetCmdQueryPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price [denom]",
		Short: "Query the current price of [denom] asset",
		Long: strings.TrimSpace(`
Query the current price of [denom] asset. You can find the current list of active denoms by running: terracli query oracle active

$ terracli query oracle price krw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, oracle.QueryVotes, denom), nil)
			if err != nil {
				return err
			}

			var price sdk.Dec
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &price)
			return cliCtx.PrintOutput(price)
		},
	}

	return cmd
}

type DenomList []string

func (dl DenomList) String() (out string) {
	for _, denom := range dl {
		out += fmt.Sprintf("\n %s", denom)
	}
	return
}

// GetCmdQueryActive implements the query active command.
func GetCmdQueryActive(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active",
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
			cdc.MustUnmarshalBinaryLengthPrefixed(res, &actives)
			return cliCtx.PrintOutput(actives)
		},
	}
	return cmd
}

// GetCmdQueryVote implements the query vote command.
func GetCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [denom] [voterAddress]",
		Short: "Query outstanding oracle votes, filtered by [denom] and [voterAddress].",
		Long: strings.TrimSpace(`
Query outstanding oracle votes, filtered by [denom] and [voterAddress].

$ terracli query oracle votes --denom="usd" --voterAddress="terrad8duyufdshs..."

returns oracle votes submitted by terrad8duyufdshs... for denom [usd] 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := viper.GetString(flagDenom)
			var voterAddress sdk.AccAddress

			params := oracle.NewQueryVoteParams(voterAddress, denom)

			bechVoterAddr := viper.GetString(flagVoter)
			if len(bechVoterAddr) != 0 {
				voterAddress, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
				params.Voter = voterAddress
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, oracle.QueryVotes), bz)
			if err != nil {
				return err
			}

			var matchingVotes oracle.PriceBallot
			err = cdc.UnmarshalJSON(res, &matchingVotes)
			if err != nil {
				return err
			}

			if len(matchingVotes) == 0 {
				return fmt.Errorf("No matching votes found")
			}

			return cliCtx.PrintOutput(matchingVotes)
		},
	}

	cmd.Flags().String(flagDenom, "", "(optional) filter by votes matching the denom")
	cmd.Flags().String(flagVoter, "", "(optional) filter by votes by voter")

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current Treasury params",
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
