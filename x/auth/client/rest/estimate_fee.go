package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/terra-money/core/x/auth/client/utils"
)

// EstimateTxFeeRequestHandlerFn returns estimated tx fee. In particular,
// it takes 'auto' for the gas field, then simulates and computes gas consumption.
func EstimateTxFeeRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.EstimateFeeReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		gasAdjustment, err := utils.ParseFloat64(req.GasAdjustment, flags.DefaultGasAdjustment)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fees, gas, err := utils.ComputeFeesWithStdTx(cliCtx, req.Tx, gasAdjustment, req.GasPrices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := utils.EstimateFeeResp{Fees: fees, Gas: gas}
		rest.PostProcessResponse(w, cliCtx, response)
	}
}
