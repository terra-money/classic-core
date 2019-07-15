package tx

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type (
	// EstimateReq defines a tx encoding request.
	EstimateFeeReq struct {
		Tx            auth.StdTx   `json:"tx"`
		GasAdjustment string       `json:"gas_adjustment"`
		GasPrices     sdk.DecCoins `json:"gas_prices"`
	}

	// EstimateResp defines a tx encoding response.
	EstimateFeeResp struct {
		Fees sdk.Coins `json:"fees"`
		Gas  uint64    `json:"gas"`
	}
)

func (r EstimateFeeResp) String() string {
	return fmt.Sprintf(`EstimateFeeResp
	fees: %s,
	gas:  %d`,
		r.Fees, r.Gas)
}

// EstimateTxFeeRequestHandlerFn returns estimated tx fee. In particular,
// it takes 'auto' for the gas field, then simulates and computes gas consumption.
func EstimateTxFeeRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EstimateFeeReq

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		gasAdjustment, err := parseFloat64(req.GasAdjustment, client.DefaultGasAdjustment)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fees, gas, err := ComputeFeesWithStdTx(cliCtx, cdc, req.Tx, gasAdjustment, req.GasPrices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := EstimateFeeResp{Fees: fees, Gas: gas}
		rest.PostProcessResponse(w, cdc, response, cliCtx.Indent)
	}
}

// GetExstimateTxFeesCommand will create a send tx and sign it with the given key.
func GetExstimateTxFeesCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-fee [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create and sign a send tx",
		Long: strings.TrimSpace(`
Estimate fees for the given stdTx

$ terracli tx estimate-fee [file] --gas-adjustment 1.4 --gas-prices 0.015uluna
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			stdTx, err := utils.ReadStdTxFromFile(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			gasAdjustment := viper.GetFloat64(client.FlagGasAdjustment)

			var gasPrices sdk.DecCoins
			gasPricesStr := viper.GetString(client.FlagGasPrices)
			if len(gasPricesStr) != 0 {
				gasPrices, err = sdk.ParseDecCoins(gasPricesStr)
				if err != nil {
					return err
				}
			}

			fees, gas, err := ComputeFeesWithStdTx(cliCtx, cdc, stdTx, gasAdjustment, gasPrices)

			if err != nil {
				return err
			}

			response := EstimateFeeResp{Fees: fees, Gas: gas}
			return cliCtx.PrintOutput(response)
		},
	}

	cmd = client.GetCommands(cmd)[0]

	cmd.Flags().Float64(client.FlagGasAdjustment, client.DefaultGasAdjustment, "adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")
	cmd.Flags().String(client.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10uluna)")
	// cmd.MarkFlagRequired(client.FlagGasAdjustment)

	return cmd
}
