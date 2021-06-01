package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"

	feeutils "github.com/terra-money/core/custom/auth/client/utils"
)

// EstimateTxFeeRequestHandlerFn returns estimated tx fee. In particular,
// it takes 'auto' for the gas field, then simulates and computes gas consumption.
func EstimateTxFeeRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req feeutils.EstimateFeeReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		stdFee, err := feeutils.ComputeFeesWithBaseReq(clientCtx, req.BaseReq, req.Msgs...)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		response := feeutils.EstimateFeeResp{Fee: *stdFee}
		rest.PostProcessResponse(w, clientCtx, response)
	}
}
