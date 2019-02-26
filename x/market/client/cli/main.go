package cli

import (
	"net/http"
	"strings"
	"terra/x/market"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOfferCoin     = "offerCoin"
	flagOfferDenom    = "offerDenom"
	flagAskDenom      = "askDenom"
	flagTraderAddress = "traderAddress"
)

// GetSwapCmd will create and send a SwapMsg
func GetSwapCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [offerCoin] [askDenom]",
		Short: "Atomically swap currencies at their target exchange rate",
		Long: strings.TrimSpace(`
Swap the offerCoin to the askDenom currency at the oracle's effective exchange rate. 

$ terracli market swap --offerCoin="1000krw" --askDenom="usd"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			askDenom := viper.GetString(flagAskDenom)
			offerCoin, err := sdk.ParseCoin(viper.GetString(flagOfferCoin))
			if err != nil {
				return err
			}

			fromAddress := cliCtx.GetFromAddress()

			// build and sign the transaction, then broadcast to Tendermint
			msg := market.NewSwapMsg(fromAddress, offerCoin, askDenom)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(flagOfferCoin, "", "The asset to swap from e.g. 1000krw")
	cmd.Flags().String(flagAskDenom, "", "Denom of the asset to swap to")

	return cmd
}

// GetCmdQueryActive implements the query active command.
func GetCmdQueryHistory(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [traderAddress] [offerDenom] [askDenom]",
		Short: "Query history of atomic swaps filtered by three optional variables.",
		Long: strings.TrimSpace(`
Query history of atomic swaps filtered by three optional variables.

$ terracli query market history --offerDenom="usd" --askDenom="krw"

Return item count paginated by units of 30 values. 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var traderAddress sdk.AccAddress
			offerDenom := viper.GetString(flagOfferDenom)
			askDenom := viper.GetString(flagAskDenom)

			params := client.QueryHistoryParams{
				TraderAddress: traderAddress,
				AskDenom:      askDenom,
				OfferDenom:    offerDenom,
			}

			traderAddrStr := viper.GetString(flagTraderAddress)
			traderAddress, err := cliCtx.GetAccount([]byte(traderAddrStr))
			if err == nil {
				params.TraderAddress = traderAddress
			}

			res, err := client.QueryHistoryByTxQuery(cdc, cliCtx, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			return cliCtx.PrintOutput(res)
		},
	}
	return cmd
}
